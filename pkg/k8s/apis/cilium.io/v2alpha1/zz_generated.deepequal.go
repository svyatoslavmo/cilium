//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

// Code generated by deepequal-gen. DO NOT EDIT.

package v2alpha1

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumBGPFamily) DeepEqual(other *CiliumBGPFamily) bool {
	if other == nil {
		return false
	}

	if in.Afi != other.Afi {
		return false
	}
	if in.Safi != other.Safi {
		return false
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumBGPNeighbor) DeepEqual(other *CiliumBGPNeighbor) bool {
	if other == nil {
		return false
	}

	if in.PeerAddress != other.PeerAddress {
		return false
	}
	if (in.PeerPort == nil) != (other.PeerPort == nil) {
		return false
	} else if in.PeerPort != nil {
		if *in.PeerPort != *other.PeerPort {
			return false
		}
	}

	if in.PeerASN != other.PeerASN {
		return false
	}
	if (in.EBGPMultihopTTL == nil) != (other.EBGPMultihopTTL == nil) {
		return false
	} else if in.EBGPMultihopTTL != nil {
		if *in.EBGPMultihopTTL != *other.EBGPMultihopTTL {
			return false
		}
	}

	if (in.ConnectRetryTimeSeconds == nil) != (other.ConnectRetryTimeSeconds == nil) {
		return false
	} else if in.ConnectRetryTimeSeconds != nil {
		if *in.ConnectRetryTimeSeconds != *other.ConnectRetryTimeSeconds {
			return false
		}
	}

	if (in.HoldTimeSeconds == nil) != (other.HoldTimeSeconds == nil) {
		return false
	} else if in.HoldTimeSeconds != nil {
		if *in.HoldTimeSeconds != *other.HoldTimeSeconds {
			return false
		}
	}

	if (in.KeepAliveTimeSeconds == nil) != (other.KeepAliveTimeSeconds == nil) {
		return false
	} else if in.KeepAliveTimeSeconds != nil {
		if *in.KeepAliveTimeSeconds != *other.KeepAliveTimeSeconds {
			return false
		}
	}

	if (in.GracefulRestart == nil) != (other.GracefulRestart == nil) {
		return false
	} else if in.GracefulRestart != nil {
		if !in.GracefulRestart.DeepEqual(other.GracefulRestart) {
			return false
		}
	}

	if ((in.Families != nil) && (other.Families != nil)) || ((in.Families == nil) != (other.Families == nil)) {
		in, other := &in.Families, &other.Families
		if other == nil {
			return false
		}

		if len(*in) != len(*other) {
			return false
		} else {
			for i, inElement := range *in {
				if !inElement.DeepEqual(&(*other)[i]) {
					return false
				}
			}
		}
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumBGPNeighborGracefulRestart) DeepEqual(other *CiliumBGPNeighborGracefulRestart) bool {
	if other == nil {
		return false
	}

	if in.Enabled != other.Enabled {
		return false
	}
	if (in.RestartTimeSeconds == nil) != (other.RestartTimeSeconds == nil) {
		return false
	} else if in.RestartTimeSeconds != nil {
		if *in.RestartTimeSeconds != *other.RestartTimeSeconds {
			return false
		}
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumBGPPeeringPolicy) DeepEqual(other *CiliumBGPPeeringPolicy) bool {
	if other == nil {
		return false
	}

	if !in.Spec.DeepEqual(&other.Spec) {
		return false
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumBGPPeeringPolicySpec) DeepEqual(other *CiliumBGPPeeringPolicySpec) bool {
	if other == nil {
		return false
	}

	if (in.NodeSelector == nil) != (other.NodeSelector == nil) {
		return false
	} else if in.NodeSelector != nil {
		if !in.NodeSelector.DeepEqual(other.NodeSelector) {
			return false
		}
	}

	if ((in.VirtualRouters != nil) && (other.VirtualRouters != nil)) || ((in.VirtualRouters == nil) != (other.VirtualRouters == nil)) {
		in, other := &in.VirtualRouters, &other.VirtualRouters
		if other == nil {
			return false
		}

		if len(*in) != len(*other) {
			return false
		} else {
			for i, inElement := range *in {
				if !inElement.DeepEqual(&(*other)[i]) {
					return false
				}
			}
		}
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumBGPVirtualRouter) DeepEqual(other *CiliumBGPVirtualRouter) bool {
	if other == nil {
		return false
	}

	if in.LocalASN != other.LocalASN {
		return false
	}
	if (in.ExportPodCIDR == nil) != (other.ExportPodCIDR == nil) {
		return false
	} else if in.ExportPodCIDR != nil {
		if *in.ExportPodCIDR != *other.ExportPodCIDR {
			return false
		}
	}

	if (in.ServiceSelector == nil) != (other.ServiceSelector == nil) {
		return false
	} else if in.ServiceSelector != nil {
		if !in.ServiceSelector.DeepEqual(other.ServiceSelector) {
			return false
		}
	}

	if ((in.Neighbors != nil) && (other.Neighbors != nil)) || ((in.Neighbors == nil) != (other.Neighbors == nil)) {
		in, other := &in.Neighbors, &other.Neighbors
		if other == nil {
			return false
		}

		if len(*in) != len(*other) {
			return false
		} else {
			for i, inElement := range *in {
				if !inElement.DeepEqual(&(*other)[i]) {
					return false
				}
			}
		}
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumCIDRGroupSpec) DeepEqual(other *CiliumCIDRGroupSpec) bool {
	if other == nil {
		return false
	}

	if ((in.ExternalCIDRs != nil) && (other.ExternalCIDRs != nil)) || ((in.ExternalCIDRs == nil) != (other.ExternalCIDRs == nil)) {
		in, other := &in.ExternalCIDRs, &other.ExternalCIDRs
		if other == nil {
			return false
		}

		if len(*in) != len(*other) {
			return false
		} else {
			for i, inElement := range *in {
				if inElement != (*other)[i] {
					return false
				}
			}
		}
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumEndpointSlice) DeepEqual(other *CiliumEndpointSlice) bool {
	if other == nil {
		return false
	}

	if in.Namespace != other.Namespace {
		return false
	}
	if ((in.Endpoints != nil) && (other.Endpoints != nil)) || ((in.Endpoints == nil) != (other.Endpoints == nil)) {
		in, other := &in.Endpoints, &other.Endpoints
		if other == nil {
			return false
		}

		if len(*in) != len(*other) {
			return false
		} else {
			for i, inElement := range *in {
				if !inElement.DeepEqual(&(*other)[i]) {
					return false
				}
			}
		}
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumL2AnnouncementPolicy) DeepEqual(other *CiliumL2AnnouncementPolicy) bool {
	if other == nil {
		return false
	}

	if !in.Spec.DeepEqual(&other.Spec) {
		return false
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumL2AnnouncementPolicySpec) DeepEqual(other *CiliumL2AnnouncementPolicySpec) bool {
	if other == nil {
		return false
	}

	if (in.NodeSelector == nil) != (other.NodeSelector == nil) {
		return false
	} else if in.NodeSelector != nil {
		if !in.NodeSelector.DeepEqual(other.NodeSelector) {
			return false
		}
	}

	if (in.ServiceSelector == nil) != (other.ServiceSelector == nil) {
		return false
	} else if in.ServiceSelector != nil {
		if !in.ServiceSelector.DeepEqual(other.ServiceSelector) {
			return false
		}
	}

	if in.LoadBalancerIPs != other.LoadBalancerIPs {
		return false
	}
	if in.ExternalIPs != other.ExternalIPs {
		return false
	}
	if ((in.Interfaces != nil) && (other.Interfaces != nil)) || ((in.Interfaces == nil) != (other.Interfaces == nil)) {
		in, other := &in.Interfaces, &other.Interfaces
		if other == nil {
			return false
		}

		if len(*in) != len(*other) {
			return false
		} else {
			for i, inElement := range *in {
				if inElement != (*other)[i] {
					return false
				}
			}
		}
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumLoadBalancerIPPool) DeepEqual(other *CiliumLoadBalancerIPPool) bool {
	if other == nil {
		return false
	}

	if !in.Spec.DeepEqual(&other.Spec) {
		return false
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumLoadBalancerIPPoolCIDRBlock) DeepEqual(other *CiliumLoadBalancerIPPoolCIDRBlock) bool {
	if other == nil {
		return false
	}

	if in.Cidr != other.Cidr {
		return false
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumLoadBalancerIPPoolSpec) DeepEqual(other *CiliumLoadBalancerIPPoolSpec) bool {
	if other == nil {
		return false
	}

	if (in.ServiceSelector == nil) != (other.ServiceSelector == nil) {
		return false
	} else if in.ServiceSelector != nil {
		if !in.ServiceSelector.DeepEqual(other.ServiceSelector) {
			return false
		}
	}

	if ((in.Cidrs != nil) && (other.Cidrs != nil)) || ((in.Cidrs == nil) != (other.Cidrs == nil)) {
		in, other := &in.Cidrs, &other.Cidrs
		if other == nil {
			return false
		}

		if len(*in) != len(*other) {
			return false
		} else {
			for i, inElement := range *in {
				if !inElement.DeepEqual(&(*other)[i]) {
					return false
				}
			}
		}
	}

	if in.Disabled != other.Disabled {
		return false
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CiliumPodIPPool) DeepEqual(other *CiliumPodIPPool) bool {
	if other == nil {
		return false
	}

	if !in.Spec.DeepEqual(&other.Spec) {
		return false
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *CoreCiliumEndpoint) DeepEqual(other *CoreCiliumEndpoint) bool {
	if other == nil {
		return false
	}

	if in.Name != other.Name {
		return false
	}
	if in.IdentityID != other.IdentityID {
		return false
	}
	if (in.Networking == nil) != (other.Networking == nil) {
		return false
	} else if in.Networking != nil {
		if !in.Networking.DeepEqual(other.Networking) {
			return false
		}
	}

	if in.Encryption != other.Encryption {
		return false
	}

	if ((in.NamedPorts != nil) && (other.NamedPorts != nil)) || ((in.NamedPorts == nil) != (other.NamedPorts == nil)) {
		in, other := &in.NamedPorts, &other.NamedPorts
		if other == nil || !in.DeepEqual(other) {
			return false
		}
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *EgressRule) DeepEqual(other *EgressRule) bool {
	if other == nil {
		return false
	}

	if (in.NamespaceSelector == nil) != (other.NamespaceSelector == nil) {
		return false
	} else if in.NamespaceSelector != nil {
		if !in.NamespaceSelector.DeepEqual(other.NamespaceSelector) {
			return false
		}
	}

	if (in.PodSelector == nil) != (other.PodSelector == nil) {
		return false
	} else if in.PodSelector != nil {
		if !in.PodSelector.DeepEqual(other.PodSelector) {
			return false
		}
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *IPPoolSpec) DeepEqual(other *IPPoolSpec) bool {
	if other == nil {
		return false
	}

	if (in.IPv4 == nil) != (other.IPv4 == nil) {
		return false
	} else if in.IPv4 != nil {
		if !in.IPv4.DeepEqual(other.IPv4) {
			return false
		}
	}

	if (in.IPv6 == nil) != (other.IPv6 == nil) {
		return false
	} else if in.IPv6 != nil {
		if !in.IPv6.DeepEqual(other.IPv6) {
			return false
		}
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *IPv4PoolSpec) DeepEqual(other *IPv4PoolSpec) bool {
	if other == nil {
		return false
	}

	if ((in.CIDRs != nil) && (other.CIDRs != nil)) || ((in.CIDRs == nil) != (other.CIDRs == nil)) {
		in, other := &in.CIDRs, &other.CIDRs
		if other == nil {
			return false
		}

		if len(*in) != len(*other) {
			return false
		} else {
			for i, inElement := range *in {
				if inElement != (*other)[i] {
					return false
				}
			}
		}
	}

	if in.MaskSize != other.MaskSize {
		return false
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *IPv6PoolSpec) DeepEqual(other *IPv6PoolSpec) bool {
	if other == nil {
		return false
	}

	if ((in.CIDRs != nil) && (other.CIDRs != nil)) || ((in.CIDRs == nil) != (other.CIDRs == nil)) {
		in, other := &in.CIDRs, &other.CIDRs
		if other == nil {
			return false
		}

		if len(*in) != len(*other) {
			return false
		} else {
			for i, inElement := range *in {
				if inElement != (*other)[i] {
					return false
				}
			}
		}
	}

	if in.MaskSize != other.MaskSize {
		return false
	}

	return true
}
