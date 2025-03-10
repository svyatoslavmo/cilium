// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"

	apiserverK8s "github.com/cilium/cilium/clustermesh-apiserver/k8s"
	cmmetrics "github.com/cilium/cilium/clustermesh-apiserver/metrics"
	apiserverOption "github.com/cilium/cilium/clustermesh-apiserver/option"
	operatorWatchers "github.com/cilium/cilium/operator/watchers"
	cmtypes "github.com/cilium/cilium/pkg/clustermesh/types"
	cmutils "github.com/cilium/cilium/pkg/clustermesh/utils"
	"github.com/cilium/cilium/pkg/controller"
	"github.com/cilium/cilium/pkg/defaults"
	"github.com/cilium/cilium/pkg/gops"
	"github.com/cilium/cilium/pkg/hive"
	"github.com/cilium/cilium/pkg/hive/cell"
	"github.com/cilium/cilium/pkg/identity"
	identityCache "github.com/cilium/cilium/pkg/identity/cache"
	"github.com/cilium/cilium/pkg/ipcache"
	"github.com/cilium/cilium/pkg/k8s"
	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	k8sClient "github.com/cilium/cilium/pkg/k8s/client"
	"github.com/cilium/cilium/pkg/k8s/resource"
	slim_corev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/core/v1"
	"github.com/cilium/cilium/pkg/k8s/synced"
	"github.com/cilium/cilium/pkg/k8s/types"
	"github.com/cilium/cilium/pkg/kvstore"
	"github.com/cilium/cilium/pkg/kvstore/heartbeat"
	"github.com/cilium/cilium/pkg/kvstore/store"
	"github.com/cilium/cilium/pkg/labels"
	"github.com/cilium/cilium/pkg/logging"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/metrics"
	nodeStore "github.com/cilium/cilium/pkg/node/store"
	nodeTypes "github.com/cilium/cilium/pkg/node/types"
	"github.com/cilium/cilium/pkg/option"
	"github.com/cilium/cilium/pkg/pprof"
	"github.com/cilium/cilium/pkg/promise"
)

type configuration struct {
	clusterName             string
	clusterID               uint32
	serviceProxyName        string
	enableExternalWorkloads bool
}

func (c configuration) LocalClusterName() string {
	return c.clusterName
}

func (c configuration) LocalClusterID() uint32 {
	return c.clusterID
}

func (c configuration) K8sServiceProxyNameValue() string {
	return c.serviceProxyName
}

var (
	log = logging.DefaultLogger.WithField(logfields.LogSubsys, "clustermesh-apiserver")

	vp       *viper.Viper
	rootHive *hive.Hive

	rootCmd = &cobra.Command{
		Use:   "clustermesh-apiserver",
		Short: "Run the ClusterMesh apiserver",
		Run: func(cmd *cobra.Command, args []string) {
			if err := rootHive.Run(); err != nil {
				log.Fatal(err)
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			// Overwrite the metrics namespace with the one specific for the ClusterMesh API Server
			metrics.Namespace = metrics.CiliumClusterMeshAPIServerNamespace
			option.Config.Populate(vp)
			if option.Config.Debug {
				log.Logger.SetLevel(logrus.DebugLevel)
			}
			option.LogRegisteredOptions(vp, log)
		},
	}

	mockFile string
	cfg      configuration
)

func init() {
	rootHive = hive.New(
		pprof.Cell,
		cell.Config(pprof.Config{
			PprofAddress: apiserverOption.PprofAddressAPIServer,
			PprofPort:    apiserverOption.PprofPortAPIServer,
		}),
		controller.Cell,

		gops.Cell(defaults.GopsPortApiserver),

		k8sClient.Cell,
		apiserverK8s.ResourcesCell,

		cell.Provide(func() *option.DaemonConfig {
			return option.Config
		}),

		kvstore.Cell(kvstore.EtcdBackendName),
		cell.Provide(func() *kvstore.ExtraOptions { return nil }),
		heartbeat.Cell,

		healthAPIServerCell,
		cmmetrics.Cell,
		usersManagementCell,
		cell.Invoke(registerHooks),
	)
	rootHive.RegisterFlags(rootCmd.Flags())
	rootCmd.AddCommand(rootHive.Command())
	vp = rootHive.Viper()
}

type parameters struct {
	cell.In

	Clientset      k8sClient.Clientset
	Resources      apiserverK8s.Resources
	BackendPromise promise.Promise[kvstore.BackendOperations]
}

func registerHooks(lc hive.Lifecycle, params parameters) error {
	lc.Append(hive.Hook{
		OnStart: func(ctx hive.HookContext) error {
			if !params.Clientset.IsEnabled() {
				return errors.New("Kubernetes client not configured, cannot continue.")
			}

			backend, err := params.BackendPromise.Await(ctx)
			if err != nil {
				return err
			}

			startServer(ctx, params.Clientset, backend, params.Resources)
			return nil
		},
	})
	return nil
}

func readMockFile(ctx context.Context, path string, backend kvstore.BackendOperations) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open file %s: %s", path, err)
	}
	defer f.Close()

	identities := newIdentitySynchronizer(ctx, backend)
	nodes := newNodeSynchronizer(ctx, backend)
	endpoints := newEndpointSynchronizer(ctx, backend)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.Contains(line, "\"CiliumIdentity\""):
			var identity ciliumv2.CiliumIdentity
			err := json.Unmarshal([]byte(line), &identity)
			if err != nil {
				log.WithError(err).WithField("line", line).Warning("Unable to unmarshal CiliumIdentity")
			} else {
				identities.upsert(ctx, resource.NewKey(&identity), &identity)
			}
		case strings.Contains(line, "\"CiliumNode\""):
			var node ciliumv2.CiliumNode
			err = json.Unmarshal([]byte(line), &node)
			if err != nil {
				log.WithError(err).WithField("line", line).Warning("Unable to unmarshal CiliumNode")
			} else {
				nodes.upsert(ctx, resource.NewKey(&node), &node)
			}
		case strings.Contains(line, "\"CiliumEndpoint\""):
			var endpoint types.CiliumEndpoint
			err = json.Unmarshal([]byte(line), &endpoint)
			if err != nil {
				log.WithError(err).WithField("line", line).Warning("Unable to unmarshal CiliumEndpoint")
			} else {
				endpoints.upsert(ctx, resource.NewKey(&endpoint), &endpoint)
			}
		case strings.Contains(line, "\"Service\""):
			var service slim_corev1.Service
			err = json.Unmarshal([]byte(line), &service)
			if err != nil {
				log.WithError(err).WithField("line", line).Warning("Unable to unmarshal Service")
			} else {
				operatorWatchers.K8sSvcCache.UpdateService(&service, nil)
			}
		case strings.Contains(line, "\"Endpoints\""):
			var endpoints slim_corev1.Endpoints
			err = json.Unmarshal([]byte(line), &endpoints)
			if err != nil {
				log.WithError(err).WithField("line", line).Warning("Unable to unmarshal Endpoints")
			} else {
				operatorWatchers.K8sSvcCache.UpdateEndpoints(k8s.ParseEndpoints(&endpoints), nil)
			}
		default:
			log.Warningf("Unknown line in mockfile %s: %s", path, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	identities.synced(ctx)
	nodes.synced(ctx)
	endpoints.synced(ctx)

	return nil
}

func runApiserver() error {
	flags := rootCmd.Flags()
	flags.BoolP(option.DebugArg, "D", false, "Enable debugging mode")
	option.BindEnv(vp, option.DebugArg)

	flags.Duration(option.CRDWaitTimeout, 5*time.Minute, "Cilium will exit if CRDs are not available within this duration upon startup")
	option.BindEnv(vp, option.CRDWaitTimeout)

	flags.String(option.IdentityAllocationMode, option.IdentityAllocationModeCRD, "Method to use for identity allocation")
	option.BindEnv(vp, option.IdentityAllocationMode)

	flags.Uint32Var(&cfg.clusterID, option.ClusterIDName, 0, "Cluster ID")
	option.BindEnv(vp, option.ClusterIDName)

	flags.StringVar(&cfg.clusterName, option.ClusterName, "default", "Cluster name")
	option.BindEnv(vp, option.ClusterName)

	flags.StringVar(&mockFile, "mock-file", "", "Read from mock file")

	flags.StringVar(&cfg.serviceProxyName, option.K8sServiceProxyName, "", "Value of K8s service-proxy-name label for which Cilium handles the services (empty = all services without service.kubernetes.io/service-proxy-name label)")
	option.BindEnv(vp, option.K8sServiceProxyName)

	flags.Duration(option.AllocatorListTimeoutName, defaults.AllocatorListTimeout, "Timeout for listing allocator state before exiting")
	option.BindEnv(vp, option.AllocatorListTimeoutName)

	flags.Bool(option.EnableWellKnownIdentities, defaults.EnableWellKnownIdentities, "Enable well-known identities for known Kubernetes components")
	option.BindEnv(vp, option.EnableWellKnownIdentities)

	flags.Bool(option.K8sEnableEndpointSlice, defaults.K8sEnableEndpointSlice, "Enable support of Kubernetes EndpointSlice")
	option.BindEnv(vp, option.K8sEnableEndpointSlice)

	// The default values is set to true to match the existing behavior in case
	// the flag is not configured (for instance by the legacy cilium CLI).
	flags.BoolVar(&cfg.enableExternalWorkloads, option.EnableExternalWorkloads, true, "Enable support for external workloads")
	option.BindEnv(vp, option.EnableExternalWorkloads)

	vp.BindPFlags(flags)

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := runApiserver(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type identitySynchronizer struct {
	store   store.SyncStore
	encoder func([]byte) string
}

func newIdentitySynchronizer(ctx context.Context, backend kvstore.BackendOperations) synchronizer {
	identitiesStore := store.NewWorkqueueSyncStore(cfg.LocalClusterName(), backend,
		path.Join(identityCache.IdentitiesPath, "id"),
		store.WSSWithSyncedKeyOverride(identityCache.IdentitiesPath))
	go identitiesStore.Run(ctx)

	return &identitySynchronizer{store: identitiesStore, encoder: backend.Encode}
}

func parseLabelArrayFromMap(base map[string]string) labels.LabelArray {
	array := make(labels.LabelArray, 0, len(base))
	for sourceAndKey, value := range base {
		array = append(array, labels.NewLabel(sourceAndKey, value, ""))
	}
	return array.Sort()
}

func (is *identitySynchronizer) upsert(ctx context.Context, _ resource.Key, obj runtime.Object) error {
	identity := obj.(*ciliumv2.CiliumIdentity)
	scopedLog := log.WithField(logfields.Identity, identity.Name)
	if len(identity.SecurityLabels) == 0 {
		scopedLog.WithError(errors.New("missing security labels")).Warning("Ignoring invalid identity")
		// Do not return an error, since it is pointless to retry.
		// We will receive a new update event if the security labels change.
		return nil
	}

	labelArray := parseLabelArrayFromMap(identity.SecurityLabels)

	var labels []byte
	for _, l := range labelArray {
		labels = append(labels, l.FormatForKVStore()...)
	}

	scopedLog.Info("Upserting identity in etcd")
	kv := store.NewKVPair(identity.Name, is.encoder(labels))
	if err := is.store.UpsertKey(ctx, kv); err != nil {
		// The only errors surfaced by WorkqueueSyncStore are the unrecoverable ones.
		log.WithError(err).Warning("Unable to upsert identity in etcd")
	}

	return nil
}

func (is *identitySynchronizer) delete(ctx context.Context, key resource.Key) error {
	scopedLog := log.WithField(logfields.Identity, key.Name)
	scopedLog.Info("Deleting identity from etcd")

	if err := is.store.DeleteKey(ctx, store.NewKVPair(key.Name, "")); err != nil {
		// The only errors surfaced by WorkqueueSyncStore are the unrecoverable ones.
		scopedLog.WithError(err).Warning("Unable to delete node from etcd")
	}

	return nil
}

func (is *identitySynchronizer) synced(ctx context.Context) error {
	log.Info("Initial list of identities successfully received from Kubernetes")
	return is.store.Synced(ctx)
}

type nodeStub struct {
	cluster string
	name    string
}

func (n *nodeStub) GetKeyName() string {
	return nodeTypes.GetKeyNodeName(n.cluster, n.name)
}

type nodeSynchronizer struct {
	store store.SyncStore
}

func newNodeSynchronizer(ctx context.Context, backend kvstore.BackendOperations) synchronizer {
	nodesStore := store.NewWorkqueueSyncStore(cfg.LocalClusterName(), backend, nodeStore.NodeStorePrefix)
	go nodesStore.Run(ctx)

	return &nodeSynchronizer{store: nodesStore}
}

func (ns *nodeSynchronizer) upsert(ctx context.Context, _ resource.Key, obj runtime.Object) error {
	n := nodeTypes.ParseCiliumNode(obj.(*ciliumv2.CiliumNode))
	n.Cluster = cfg.clusterName
	n.ClusterID = cfg.clusterID

	scopedLog := log.WithField(logfields.Node, n.Name)
	scopedLog.Info("Upserting node in etcd")

	if err := ns.store.UpsertKey(ctx, &n); err != nil {
		// The only errors surfaced by WorkqueueSyncStore are the unrecoverable ones.
		log.WithError(err).Warning("Unable to upsert node in etcd")
	}

	return nil
}

func (ns *nodeSynchronizer) delete(ctx context.Context, key resource.Key) error {
	n := nodeStub{
		cluster: cfg.clusterName,
		name:    key.Name,
	}

	scopedLog := log.WithFields(logrus.Fields{logfields.Node: key.Name})
	scopedLog.Info("Deleting node from etcd")

	if err := ns.store.DeleteKey(ctx, &n); err != nil {
		// The only errors surfaced by WorkqueueSyncStore are the unrecoverable ones.
		scopedLog.WithError(err).Warning("Unable to delete node from etcd")
	}

	return nil
}

func (ns *nodeSynchronizer) synced(ctx context.Context) error {
	log.Info("Initial list of nodes successfully received from Kubernetes")
	return ns.store.Synced(ctx)
}

type ipmap map[string]struct{}

type endpointSynchronizer struct {
	store store.SyncStore
	cache map[string]ipmap
}

func newEndpointSynchronizer(ctx context.Context, backend kvstore.BackendOperations) synchronizer {
	endpointsStore := store.NewWorkqueueSyncStore(cfg.LocalClusterName(), backend,
		path.Join(ipcache.IPIdentitiesPath, ipcache.DefaultAddressSpace),
		store.WSSWithSyncedKeyOverride(ipcache.IPIdentitiesPath))
	go endpointsStore.Run(ctx)

	return &endpointSynchronizer{
		store: endpointsStore,
		cache: make(map[string]ipmap),
	}
}

func (es *endpointSynchronizer) upsert(ctx context.Context, key resource.Key, obj runtime.Object) error {
	endpoint := obj.(*types.CiliumEndpoint)
	ips := make(ipmap)
	stale := es.cache[key.String()]

	if n := endpoint.Networking; n != nil {
		for _, address := range n.Addressing {
			for _, ip := range []string{address.IPV4, address.IPV6} {
				if ip == "" {
					continue
				}

				scopedLog := log.WithFields(logrus.Fields{logfields.Endpoint: key.String(), logfields.IPAddr: ip})
				entry := identity.IPIdentityPair{
					IP:           net.ParseIP(ip),
					HostIP:       net.ParseIP(n.NodeIP),
					K8sNamespace: endpoint.Namespace,
					K8sPodName:   endpoint.Name,
				}

				if endpoint.Identity != nil {
					entry.ID = identity.NumericIdentity(endpoint.Identity.ID)
				}

				if endpoint.Encryption != nil {
					entry.Key = uint8(endpoint.Encryption.Key)
				}

				scopedLog.Info("Upserting endpoint in etcd")
				if err := es.store.UpsertKey(ctx, &entry); err != nil {
					// The only errors surfaced by WorkqueueSyncStore are the unrecoverable ones.
					scopedLog.WithError(err).Warning("Unable to upsert endpoint in etcd")
					continue
				}

				ips[ip] = struct{}{}
				delete(stale, ip)
			}
		}
	}

	// Delete the stale endpoint IPs from the KVStore.
	es.deleteEndpoints(ctx, key, stale)
	es.cache[key.String()] = ips

	return nil
}

func (es *endpointSynchronizer) delete(ctx context.Context, key resource.Key) error {
	es.deleteEndpoints(ctx, key, es.cache[key.String()])
	delete(es.cache, key.String())
	return nil
}

func (es *endpointSynchronizer) synced(ctx context.Context) error {
	log.Info("Initial list of endpoints successfully received from Kubernetes")
	return es.store.Synced(ctx)
}

func (es *endpointSynchronizer) deleteEndpoints(ctx context.Context, key resource.Key, ips ipmap) {
	for ip := range ips {
		scopedLog := log.WithFields(logrus.Fields{logfields.Endpoint: key.String(), logfields.IPAddr: ip})
		scopedLog.Info("Deleting endpoint from etcd")

		entry := identity.IPIdentityPair{IP: net.ParseIP(ip)}
		if err := es.store.DeleteKey(ctx, &entry); err != nil {
			// The only errors surfaced by WorkqueueSyncStore are the unrecoverable ones.
			scopedLog.WithError(err).Warning("Unable to delete endpoint from etcd")
		}
	}
}

type synchronizer interface {
	upsert(ctx context.Context, key resource.Key, obj runtime.Object) error
	delete(ctx context.Context, key resource.Key) error
	synced(ctx context.Context) error
}

func synchronize[T runtime.Object](ctx context.Context, r resource.Resource[T], sync synchronizer) {
	for event := range r.Events(ctx) {
		switch event.Kind {
		case resource.Upsert:
			event.Done(sync.upsert(ctx, event.Key, event.Object))
		case resource.Delete:
			event.Done(sync.delete(ctx, event.Key))
		case resource.Sync:
			event.Done(sync.synced(ctx))
		}
	}
}

func startServer(startCtx hive.HookContext, clientset k8sClient.Clientset, backend kvstore.BackendOperations, resources apiserverK8s.Resources) {
	log.WithFields(logrus.Fields{
		"cluster-name": cfg.clusterName,
		"cluster-id":   cfg.clusterID,
	}).Info("Starting clustermesh-apiserver...")

	if mockFile == "" {
		synced.SyncCRDs(startCtx, clientset, synced.ClusterMeshAPIServerResourceNames(), &synced.Resources{}, &synced.APIGroups{})
	}

	var err error

	config := cmtypes.CiliumClusterConfig{
		ID: cfg.clusterID,
		Capabilities: cmtypes.CiliumClusterConfigCapabilities{
			SyncedCanaries: true,
		},
	}

	if err := cmutils.SetClusterConfig(context.Background(), cfg.clusterName, &config, backend); err != nil {
		log.WithError(err).Fatal("Unable to set local cluster config on kvstore")
	}

	if cfg.enableExternalWorkloads {
		mgr := NewVMManager(clientset, backend)
		_, err = store.JoinSharedStore(store.Configuration{
			Backend:              backend,
			Prefix:               nodeStore.NodeRegisterStorePrefix,
			KeyCreator:           nodeStore.RegisterKeyCreator,
			SharedKeyDeleteDelay: defaults.NodeDeleteDelay,
			Observer:             mgr,
		})
		if err != nil {
			log.WithError(err).Fatal("Unable to set up node register store in etcd")
		}
	}

	ctx := context.Background()
	if mockFile != "" {
		if err := readMockFile(ctx, mockFile, backend); err != nil {
			log.WithError(err).Fatal("Unable to read mock file")
		}
	} else {
		go synchronize(ctx, resources.CiliumIdentities, newIdentitySynchronizer(ctx, backend))
		go synchronize(ctx, resources.CiliumNodes, newNodeSynchronizer(ctx, backend))
		go synchronize(ctx, resources.CiliumSlimEndpoints, newEndpointSynchronizer(ctx, backend))
		operatorWatchers.StartSynchronizingServices(ctx, &sync.WaitGroup{}, operatorWatchers.ServiceSyncParameters{
			ServiceSyncConfiguration: cfg,
			Clientset:                clientset,
			Services:                 resources.Services,
			Endpoints:                resources.Endpoints,
			Backend:                  backend,
			SharedOnly:               !cfg.enableExternalWorkloads,
		})
	}

	log.Info("Initialization complete")
}
