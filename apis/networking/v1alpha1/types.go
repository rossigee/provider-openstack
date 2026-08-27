/*
Copyright 2025 The Crossplane Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Network parameters define the desired state of an OpenStack Neutron network.
type NetworkParameters struct {
	// Name is the human-readable name for the network. If not specified,
	// the metadata.name is used.
	// +optional
	Name string `json:"name,omitempty"`

	// Description of the network.
	// +optional
	Description string `json:"description,omitempty"`

	// AdminStateUp is the administrative state of the network (true = up).
	// +kubebuilder:default=true
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// TenantID is the project owner of the network.
	// +optional
	TenantID string `json:"tenantId,omitempty"`

	// Shared indicates whether the network is shared across all projects.
	// +kubebuilder:default=false
	// +optional
	Shared *bool `json:"shared,omitempty"`

	// DNSDomain is the DNS domain for the network.
	// +optional
	DNSDomain string `json:"dnsDomain,omitempty"`

	// PortSecurityEnabled controls port security on the network.
	// +kubebuilder:default=true
	// +optional
	PortSecurityEnabled *bool `json:"portSecurityEnabled,omitempty"`

	// Tags is a set of arbitrary key/value pairs for network tagging.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// NetworkStatus defines the observed state of an OpenStack network.
type NetworkStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	// AtProvider contains observations of the network state at the provider.
	AtProvider NetworkProviderStatus `json:"atProvider,omitempty"`

	// ConnectionDetails contains the connection details for this resource.
	ConnectionDetails []byte `json:"connectionDetails,omitempty"`
}

// NetworkProviderStatus defines the observed state of the network at the provider.
type NetworkProviderStatus struct {
	// NetworkID is the unique identifier of the network in OpenStack.
	NetworkID string `json:"networkId,omitempty"`

	// Status indicates the current status of the network (ACTIVE, DOWN, BUILD, etc.).
	Status string `json:"status,omitempty"`

	// Subnets is a list of subnet IDs associated with this network.
	Subnets []string `json:"subnets,omitempty"`

	// TenantID is the project owner of the network.
	TenantID string `json:"tenantId,omitempty"`

	// AdminStateUp is the administrative state of the network.
	AdminStateUp bool `json:"adminStateUp,omitempty"`

	// Shared indicates if the network is shared.
	Shared bool `json:"shared,omitempty"`

	// ProviderNetworkType is the physical network type (e.g., "flat", "vlan", "vxlan").
	ProviderNetworkType string `json:"providerNetworkType,omitempty"`

	// ProviderPhysicalNetwork is the physical network on which the network is created.
	ProviderPhysicalNetwork string `json:"providerPhysicalNetwork,omitempty"`

	// ProviderSegmentationID is the segmentation ID (VLAN ID, VXLAN VNI, etc.).
	ProviderSegmentationID int `json:"providerSegmentationId,omitempty"`

	// PortSecurityEnabled indicates whether port security is enabled.
	PortSecurityEnabled bool `json:"portSecurityEnabled,omitempty"`

	// DNSDomain is the DNS domain.
	DNSDomain string `json:"dnsDomain,omitempty"`

	// RevisionNumber is the network revision number for optimistic locking.
	RevisionNumber int `json:"revisionNumber,omitempty"`

	// Tags is the set of tags on the network.
	Tags []string `json:"tags,omitempty"`

	// CreatedAt is the creation timestamp in OpenStack.
	CreatedAt string `json:"createdAt,omitempty"`

	// UpdatedAt is the last update timestamp in OpenStack.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// Network is a managed resource that represents an OpenStack Neutron network.
type Network struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkSpec   `json:"spec"`
	Status NetworkStatus `json:"status,omitempty"`
}

// NetworkSpec defines the desired state of a Network.
type NetworkSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     NetworkParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// NetworkList contains a list of Network resources.
type NetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Network `json:"items"`
}

// Subnet parameters define the desired state of an OpenStack Neutron subnet.
type SubnetParameters struct {
	// NetworkID is the ID of the network to which the subnet belongs.
	// +kubebuilder:validation:Required
	NetworkID string `json:"networkId"`

	// Name is the human-readable name for the subnet.
	// +optional
	Name string `json:"name,omitempty"`

	// Description of the subnet.
	// +optional
	Description string `json:"description,omitempty"`

	// CIDR is the CIDR notation of the subnet (e.g., "10.0.0.0/24").
	// +kubebuilder:validation:Required
	CIDR string `json:"cidr"`

	// IPVersion is the IP version (4 or 6).
	// +kubebuilder:default=4
	// +optional
	IPVersion int `json:"ipVersion,omitempty"`

	// Gateway is the default gateway IP address. If not specified,
	// OpenStack allocates one automatically.
	// +optional
	Gateway string `json:"gateway,omitempty"`

	// DNSNameservers is a list of DNS nameserver IP addresses.
	// +optional
	DNSNameservers []string `json:"dnsNameservers,omitempty"`

	// AllocationPools define the IP address ranges for allocation.
	// If not specified, the entire CIDR range is used.
	// +optional
	AllocationPools []AllocationPool `json:"allocationPools,omitempty"`

	// HostRoutes define additional static routes for the subnet.
	// +optional
	HostRoutes []HostRoute `json:"hostRoutes,omitempty"`

	// TenantID is the project owner of the subnet.
	// +optional
	TenantID string `json:"tenantId,omitempty"`

	// EnableDHCP controls whether DHCP is enabled for the subnet.
	// +kubebuilder:default=true
	// +optional
	EnableDHCP *bool `json:"enableDHCP,omitempty"`

	// IPv6AddressMode is the IPv6 address mode (e.g., "slaac", "dhcpv6-stateful").
	// +optional
	IPv6AddressMode string `json:"ipv6AddressMode,omitempty"`

	// IPv6RAMode is the IPv6 RA mode.
	// +optional
	IPv6RAMode string `json:"ipv6RaMode,omitempty"`

	// SubnetPoolID references a subnet pool for address allocation.
	// +optional
	SubnetPoolID string `json:"subnetPoolId,omitempty"`

	// PrefixLength is the prefix length when using a subnet pool.
	// +optional
	PrefixLength *int `json:"prefixLength,omitempty"`

	// Tags is a set of arbitrary key/value pairs for subnet tagging.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// AllocationPool defines a range of IP addresses.
type AllocationPool struct {
	// Start is the starting IP address of the allocation pool.
	// +kubebuilder:validation:Required
	Start string `json:"start"`

	// End is the ending IP address of the allocation pool.
	// +kubebuilder:validation:Required
	End string `json:"end"`
}

// HostRoute defines a static route for the subnet.
type HostRoute struct {
	// Destination is the destination CIDR of the route.
	// +kubebuilder:validation:Required
	Destination string `json:"destination"`

	// Nexthop is the next hop IP address.
	// +kubebuilder:validation:Required
	Nexthop string `json:"nexthop"`
}

// SubnetStatus defines the observed state of an OpenStack subnet.
type SubnetStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider SubnetProviderStatus `json:"atProvider,omitempty"`
}

// SubnetProviderStatus defines the observed state of the subnet at the provider.
type SubnetProviderStatus struct {
	// SubnetID is the unique identifier of the subnet in OpenStack.
	SubnetID string `json:"subnetId,omitempty"`

	// NetworkID is the network this subnet belongs to.
	NetworkID string `json:"networkId,omitempty"`

	// Status is the current status of the subnet.
	Status string `json:"status,omitempty"`

	// CIDR is the CIDR of the subnet.
	CIDR string `json:"cidr,omitempty"`

	// Gateway is the default gateway.
	Gateway string `json:"gateway,omitempty"`

	// DNSNameservers is the list of DNS nameservers.
	DNSNameservers []string `json:"dnsNameservers,omitempty"`

	// AllocationPools is the list of allocation pools.
	AllocationPools []AllocationPool `json:"allocationPools,omitempty"`

	// HostRoutes is the list of host routes.
	HostRoutes []HostRoute `json:"hostRoutes,omitempty"`

	// IPVersion is the IP version.
	IPVersion int `json:"ipVersion,omitempty"`

	// EnableDHCP indicates if DHCP is enabled.
	EnableDHCP bool `json:"enableDHCP,omitempty"`

	// TenantID is the project owner.
	TenantID string `json:"tenantId,omitempty"`

	// RevisionNumber for optimistic locking.
	RevisionNumber int `json:"revisionNumber,omitempty"`

	// Tags on the subnet.
	Tags []string `json:"tags,omitempty"`

	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`

	// UpdatedAt is the last update timestamp.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="CIDR",type="string",JSONPath=".status.atProvider.cidr"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// Subnet is a managed resource that represents an OpenStack Neutron subnet.
type Subnet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubnetSpec   `json:"spec"`
	Status SubnetStatus `json:"status,omitempty"`
}

// SubnetSpec defines the desired state of a Subnet.
type SubnetSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     SubnetParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// SubnetList contains a list of Subnet resources.
type SubnetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Subnet `json:"items"`
}

// Router parameters define the desired state of an OpenStack Neutron router.
type RouterParameters struct {
	// Name is the human-readable name for the router.
	// +optional
	Name string `json:"name,omitempty"`

	// Description of the router.
	// +optional
	Description string `json:"description,omitempty"`

	// AdminStateUp is the administrative state of the router.
	// +kubebuilder:default=true
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// TenantID is the project owner of the router.
	// +optional
	TenantID string `json:"tenantId,omitempty"`

	// Distributed indicates whether the router is distributed.
	// +kubebuilder:default=false
	// +optional
	Distributed *bool `json:"distributed,omitempty"`

	// Tags is a set of arbitrary key/value pairs for router tagging.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// RouterStatus defines the observed state of an OpenStack router.
type RouterStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider RouterProviderStatus `json:"atProvider,omitempty"`
}

// RouterProviderStatus defines the observed state of the router at the provider.
type RouterProviderStatus struct {
	// RouterID is the unique identifier of the router in OpenStack.
	RouterID string `json:"routerId,omitempty"`

	// Status is the current status of the router.
	Status string `json:"status,omitempty"`

	// TenantID is the project owner.
	TenantID string `json:"tenantId,omitempty"`

	// AdminStateUp is the administrative state.
	AdminStateUp bool `json:"adminStateUp,omitempty"`

	// Distributed indicates if the router is distributed.
	Distributed bool `json:"distributed,omitempty"`

	// GatewayNetworkID is the external network ID connected to the router gateway.
	GatewayNetworkID string `json:"gatewayNetworkId,omitempty"`

	// GatewayExternalIP is the external IP allocated to the router.
	GatewayExternalIP string `json:"gatewayExternalIp,omitempty"`

	// Routes are additional static routes configured on the router.
	Routes []RouterRoute `json:"routes,omitempty"`

	// RevisionNumber for optimistic locking.
	RevisionNumber int `json:"revisionNumber,omitempty"`

	// Tags on the router.
	Tags []string `json:"tags,omitempty"`

	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`

	// UpdatedAt is the last update timestamp.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// RouterRoute defines a static route on a router.
type RouterRoute struct {
	// Destination is the destination CIDR.
	Destination string `json:"destination"`

	// Nexthop is the next hop IP address.
	Nexthop string `json:"nexthop"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// Router is a managed resource that represents an OpenStack Neutron router.
type Router struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouterSpec   `json:"spec"`
	Status RouterStatus `json:"status,omitempty"`
}

// RouterSpec defines the desired state of a Router.
type RouterSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     RouterParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// RouterList contains a list of Router resources.
type RouterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Router `json:"items"`
}

// RouterInterface parameters define the desired state of an OpenStack router interface.
type RouterInterfaceParameters struct {
	// RouterID is the ID of the router to attach the interface to.
	// +kubebuilder:validation:Required
	RouterID string `json:"routerId"`

	// Type is the type of interface (subnet, network, or acp). Typically "subnet".
	// +kubebuilder:default=subnet
	// +kubebuilder:validation:Enum=subnet;network;acp
	// +optional
	Type string `json:"type,omitempty"`

	// SubnetID is the ID of the subnet to attach to the router.
	// Required when Type is "subnet".
	// +optional
	SubnetID string `json:"subnetId,omitempty"`

	// SubnetCIDR is the CIDR of the subnet (used to look up the subnet by CIDR when SubnetID is not specified).
	// +optional
	SubnetCIDR string `json:"subnetCidr,omitempty"`

	// ProjectID is the project ID for the router interface.
	// +optional
	ProjectID string `json:"projectId,omitempty"`
}

// RouterInterfaceStatus defines the observed state of a router interface.
type RouterInterfaceStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider RouterInterfaceProviderStatus `json:"atProvider,omitempty"`
}

// RouterInterfaceProviderStatus defines the observed state of the router interface at the provider.
type RouterInterfaceProviderStatus struct {
	// InterfaceID is the unique identifier of the router interface.
	InterfaceID string `json:"interfaceId,omitempty"`

	// RouterID is the router this interface is attached to.
	RouterID string `json:"routerId,omitempty"`

	// SubnetID is the subnet this interface is attached to.
	SubnetID string `json:"subnetId,omitempty"`

	// TenantID is the project owner.
	TenantID string `json:"tenantId,omitempty"`

	// PortID is the port created for this interface.
	PortID string `json:"portId,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// RouterInterface is a managed resource that represents an OpenStack router interface.
type RouterInterface struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouterInterfaceSpec   `json:"spec"`
	Status RouterInterfaceStatus `json:"status,omitempty"`
}

// RouterInterfaceSpec defines the desired state of a RouterInterface.
type RouterInterfaceSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     RouterInterfaceParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// RouterInterfaceList contains a list of RouterInterface resources.
type RouterInterfaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RouterInterface `json:"items"`
}

// SecurityGroupParameters define the desired state of an OpenStack Neutron security group.
type SecurityGroupParameters struct {
	// Name is the human-readable name for the security group.
	// +optional
	Name string `json:"name,omitempty"`

	// Description of the security group.
	// +optional
	Description string `json:"description,omitempty"`

	// TenantID is the project owner of the security group.
	// +optional
	TenantID string `json:"tenantId,omitempty"`

	// Tags is a set of arbitrary key/value pairs for security group tagging.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// SecurityGroupStatus defines the observed state of an OpenStack security group.
type SecurityGroupStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider SecurityGroupProviderStatus `json:"atProvider,omitempty"`
}

// SecurityGroupProviderStatus defines the observed state of the security group at the provider.
type SecurityGroupProviderStatus struct {
	// SecurityGroupID is the unique identifier of the security group.
	SecurityGroupID string `json:"securityGroupId,omitempty"`

	// Name is the name of the security group.
	Name string `json:"name,omitempty"`

	// Description of the security group.
	Description string `json:"description,omitempty"`

	// Rules is the list of security group rules.
	Rules []SecurityGroupRuleInfo `json:"rules,omitempty"`

	// TenantID is the project owner.
	TenantID string `json:"tenantId,omitempty"`

	// Tags on the security group.
	Tags []string `json:"tags,omitempty"`

	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`

	// UpdatedAt is the last update timestamp.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// SecurityGroupRuleInfo describes a security group rule observed in OpenStack.
type SecurityGroupRuleInfo struct {
	// ID is the rule identifier.
	ID string `json:"id"`

	// Direction is the traffic direction (ingress or egress).
	Direction string `json:"direction"`

	// Ethertype is the ethertype (IPv4 or IPv6).
	Ethertype string `json:"ethertype"`

	// Protocol is the IP protocol (tcp, udp, icmp, etc.).
	Protocol string `json:"protocol,omitempty"`

	// PortRangeMin is the minimum port number in the range.
	PortRangeMin *int `json:"portRangeMin,omitempty"`

	// PortRangeMax is the maximum port number in the range.
	PortRangeMax *int `json:"portRangeMax,omitempty"`

	// RemoteIPPrefix is the CIDR of the remote source/destination.
	RemoteIPPrefix string `json:"remoteIPPrefix,omitempty"`

	// RemoteGroupID references another security group.
	RemoteGroupID string `json:"remoteGroupId,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Rules",type="integer",JSONPath=".status.atProvider.rules"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// SecurityGroup is a managed resource that represents an OpenStack Neutron security group.
type SecurityGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecurityGroupSpec   `json:"spec"`
	Status SecurityGroupStatus `json:"status,omitempty"`
}

// SecurityGroupSpec defines the desired state of a SecurityGroup.
type SecurityGroupSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     SecurityGroupParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// SecurityGroupList contains a list of SecurityGroup resources.
type SecurityGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecurityGroup `json:"items"`
}

// SecurityGroupRuleParameters define the desired state of an OpenStack security group rule.
type SecurityGroupRuleParameters struct {
	// SecurityGroupID is the security group this rule belongs to.
	// +kubebuilder:validation:Required
	SecurityGroupID string `json:"securityGroupId"`

	// Direction is the traffic direction (ingress or egress).
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=ingress;egress
	Direction string `json:"direction"`

	// Ethertype is the ethertype (IPv4 or IPv6).
	// +kubebuilder:default="IPv4"
	// +kubebuilder:validation:Enum=IPv4;IPv6
	// +optional
	Ethertype string `json:"ethertype,omitempty"`

	// Protocol is the IP protocol (tcp, udp, icmp, or empty for any).
	// +optional
	Protocol string `json:"protocol,omitempty"`

	// PortRangeMin is the minimum port number in the range.
	// +optional
	PortRangeMin *int `json:"portRangeMin,omitempty"`

	// PortRangeMax is the maximum port number in the range.
	// +optional
	PortRangeMax *int `json:"portRangeMax,omitempty"`

	// RemoteIPPrefix is the CIDR of the remote source/destination.
	// Mutually exclusive with RemoteGroupID.
	// +optional
	RemoteIPPrefix string `json:"remoteIPPrefix,omitempty"`

	// RemoteGroupID references another security group as the source/destination.
	// Mutually exclusive with RemoteIPPrefix.
	// +optional
	RemoteGroupID string `json:"remoteGroupId,omitempty"`

	// TenantID is the project owner.
	// +optional
	TenantID string `json:"tenantId,omitempty"`
}

// SecurityGroupRuleStatus defines the observed state of a security group rule.
type SecurityGroupStatusRuleStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider SecurityGroupRuleProviderStatus `json:"atProvider,omitempty"`
}

// SecurityGroupRuleProviderStatus defines the observed state of the rule at the provider.
type SecurityGroupRuleProviderStatus struct {
	// RuleID is the unique identifier of the rule.
	RuleID string `json:"ruleId,omitempty"`

	// SecurityGroupID is the parent security group.
	SecurityGroupID string `json:"securityGroupId,omitempty"`

	// Direction is the traffic direction.
	Direction string `json:"direction,omitempty"`

	// Ethertype is the ethertype.
	Ethertype string `json:"ethertype,omitempty"`

	// Protocol is the IP protocol.
	Protocol string `json:"protocol,omitempty"`

	// PortRangeMin is the minimum port.
	PortRangeMin *int `json:"portRangeMin,omitempty"`

	// PortRangeMax is the maximum port.
	PortRangeMax *int `json:"portRangeMax,omitempty"`

	// RemoteIPPrefix is the CIDR.
	RemoteIPPrefix string `json:"remoteIPPrefix,omitempty"`

	// RemoteGroupID references another security group.
	RemoteGroupID string `json:"remoteGroupId,omitempty"`

	// TenantID is the project owner.
	TenantID string `json:"tenantId,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Direction",type="string",JSONPath=".spec.forProvider.direction"
// +kubebuilder:printcolumn:name="Protocol",type="string",JSONPath=".spec.forProvider.protocol"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// SecurityGroupRule is a managed resource that represents an OpenStack security group rule.
type SecurityGroupRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecurityGroupRuleSpec         `json:"spec"`
	Status SecurityGroupStatusRuleStatus `json:"status,omitempty"`
}

// SecurityGroupRuleSpec defines the desired state of a SecurityGroupRule.
type SecurityGroupRuleSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     SecurityGroupRuleParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// SecurityGroupRuleList contains a list of SecurityGroupRule resources.
type SecurityGroupRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecurityGroupRule `json:"items"`
}

// FloatingIPParameters define the desired state of an OpenStack Neutron floating IP.
type FloatingIPParameters struct {
	// FloatingNetworkID is the ID of the external network from which to allocate the floating IP.
	// +kubebuilder:validation:Required
	FloatingNetworkID string `json:"floatingNetworkId"`

	// FloatingIP is a specific floating IP address to use. If not specified,
	// OpenStack allocates one from the pool.
	// +optional
	FloatingIP string `json:"floatingIP,omitempty"`

	// PortID is the ID of the port to associate with the floating IP.
	// +optional
	PortID string `json:"portId,omitempty"`

	// FixedIP is the specific IP address of the port's fixed IP to associate with.
	// Used when the port has multiple fixed IPs.
	// +optional
	FixedIP string `json:"fixedIP,omitempty"`

	// TenantID is the project owner of the floating IP.
	// +optional
	TenantID string `json:"tenantId,omitempty"`

	// Description of the floating IP.
	// +optional
	Description string `json:"description,omitempty"`

	// Tags is a set of arbitrary key/value pairs for floating IP tagging.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// FloatingIPStatus defines the observed state of an OpenStack floating IP.
type FloatingIPStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider FloatingIPProviderStatus `json:"atProvider,omitempty"`

	// ConnectionDetails stores the connection details for the floating IP.
	ConnectionDetails ConnectionDetails `json:"connectionDetails,omitempty"`
}

// ConnectionDetails stores connection details.
type ConnectionDetails struct {
	// Endpoint is the floating IP address.
	Endpoint string `json:"endpoint,omitempty"`
}

// FloatingIPProviderStatus defines the observed state of the floating IP at the provider.
type FloatingIPProviderStatus struct {
	// FloatingIPID is the unique identifier of the floating IP.
	FloatingIPID string `json:"floatingIPId,omitempty"`

	// FloatingIP is the allocated floating IP address.
	FloatingIP string `json:"floatingIP,omitempty"`

	// FixedIP is the associated fixed IP address.
	FixedIP string `json:"fixedIP,omitempty"`

	// PortID is the associated port ID.
	PortID string `json:"portId,omitempty"`

	// Status is the current status (DOWN, ACTIVE, ERROR).
	Status string `json:"status,omitempty"`

	// TenantID is the project owner.
	TenantID string `json:"tenantId,omitempty"`

	// RouterID is the ID of the router used for the floating IP.
	RouterID string `json:"routerId,omitempty"`

	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`

	// UpdatedAt is the last update timestamp.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Floating IP",type="string",JSONPath=".status.atProvider.floatingIP"
// +kubebuilder:printcolumn:name="Fixed IP",type="string",JSONPath=".status.atProvider.fixedIP"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// FloatingIP is a managed resource that represents an OpenStack Neutron floating IP.
type FloatingIP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FloatingIPSpec   `json:"spec"`
	Status FloatingIPStatus `json:"status,omitempty"`
}

// FloatingIPSpec defines the desired state of a FloatingIP.
type FloatingIPSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     FloatingIPParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// FloatingIPList contains a list of FloatingIP resources.
type FloatingIPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FloatingIP `json:"items"`
}
