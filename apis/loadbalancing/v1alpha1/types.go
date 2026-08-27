package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LoadBalancerParameters define the desired state of an OpenStack Octavia load balancer.
type LoadBalancerParameters struct {
	// Name is the human-readable name for the load balancer.
	// +optional
	Name string `json:"name,omitempty"`

	// Description of the load balancer.
	// +optional
	Description string `json:"description,omitempty"`

	// VipSubnetID is the UUID of the subnet on which to allocate the virtual IP.
	// +optional
	VipSubnetID string `json:"vipSubnetId,omitempty"`

	// VipNetworkID is the UUID of the network on which to allocate the virtual IP.
	// +optional
	VipNetworkID string `json:"vipNetworkId,omitempty"`

	// VipPortID is the UUID of the port associated with the VIP address.
	// +optional
	VipPortID string `json:"vipPortId,omitempty"`

	// VipAddress is the IP address of the load balancer.
	// +optional
	VipAddress string `json:"vipAddress,omitempty"`

	// VipQosPolicyID is the ID of the QoS Policy which will apply to the Virtual IP.
	// +optional
	VipQosPolicyID string `json:"vipQosPolicyId,omitempty"`

	// AdminStateUp is the administrative state of the load balancer (true = UP).
	// +kubebuilder:default=true
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// ProjectID is the UUID of the project who owns the load balancer.
	// +optional
	ProjectID string `json:"projectId,omitempty"`

	// FlavorID is the UUID of a flavor if set.
	// +optional
	FlavorID string `json:"flavorId,omitempty"`

	// AvailabilityZone is the name of an Octavia availability zone.
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// Provider is the name of the provider.
	// +optional
	Provider string `json:"provider,omitempty"`

	// Tags is a set of resource tags.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// LoadBalancerStatus defines the observed state of an OpenStack load balancer.
type LoadBalancerStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider LoadBalancerProviderStatus `json:"atProvider,omitempty"`
}

// LoadBalancerProviderStatus defines the observed state of the load balancer at the provider.
type LoadBalancerProviderStatus struct {
	// LoadBalancerID is the unique identifier of the load balancer in OpenStack.
	LoadBalancerID string `json:"loadBalancerId,omitempty"`

	// Name is the human-readable name of the load balancer.
	Name string `json:"name,omitempty"`

	// Description of the load balancer.
	Description string `json:"description,omitempty"`

	// ProvisioningStatus is the provisioning status of the load balancer (ACTIVE, PENDING_CREATE, ERROR, etc.).
	ProvisioningStatus string `json:"provisioningStatus,omitempty"`

	// OperatingStatus is the operating status (ONLINE, OFFLINE, etc.).
	OperatingStatus string `json:"operatingStatus,omitempty"`

	// VipAddress is the IP address of the load balancer.
	VipAddress string `json:"vipAddress,omitempty"`

	// VipPortID is the UUID of the port associated with the VIP.
	VipPortID string `json:"vipPortId,omitempty"`

	// VipSubnetID is the UUID of the subnet for the VIP.
	VipSubnetID string `json:"vipSubnetId,omitempty"`

	// VipNetworkID is the UUID of the network for the VIP.
	VipNetworkID string `json:"vipNetworkId,omitempty"`

	// ProjectID is the project owner.
	ProjectID string `json:"projectId,omitempty"`

	// AdminStateUp is the administrative state.
	AdminStateUp bool `json:"adminStateUp,omitempty"`

	// FlavorID is the UUID of the flavor.
	FlavorID string `json:"flavorId,omitempty"`

	// Provider is the provider name.
	Provider string `json:"provider,omitempty"`

	// Tags on the load balancer.
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
// +kubebuilder:printcolumn:name="VIP Address",type="string",JSONPath=".status.atProvider.vipAddress"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.provisioningStatus"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// LoadBalancer is a managed resource that represents an OpenStack Octavia load balancer.
type LoadBalancer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoadBalancerSpec   `json:"spec"`
	Status LoadBalancerStatus `json:"status,omitempty"`
}

// LoadBalancerSpec defines the desired state of a LoadBalancer.
type LoadBalancerSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     LoadBalancerParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// LoadBalancerList contains a list of LoadBalancer resources.
type LoadBalancerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoadBalancer `json:"items"`
}

// ListenerParameters define the desired state of an OpenStack Octavia listener.
type ListenerParameters struct {
	// LoadbalancerID is the UUID of the load balancer on which to provision this listener.
	// +kubebuilder:validation:Required
	LoadbalancerID string `json:"loadbalancerId"`

	// Protocol is the protocol - TCP, UDP, HTTP, HTTPS, or TERMINATED_HTTPS.
	// +kubebuilder:validation:Required
	Protocol string `json:"protocol"`

	// ProtocolPort is the port on which to listen for client traffic.
	// +kubebuilder:validation:Required
	ProtocolPort int `json:"protocolPort"`

	// Name is the human-readable name for the listener.
	// +optional
	Name string `json:"name,omitempty"`

	// Description of the listener.
	// +optional
	Description string `json:"description,omitempty"`

	// DefaultPoolID is the ID of the default pool with which the listener is associated.
	// +optional
	DefaultPoolID string `json:"defaultPoolId,omitempty"`

	// ConnLimit is the maximum number of connections allowed. Default is -1 (no limit).
	// +optional
	ConnLimit *int `json:"connLimit,omitempty"`

	// DefaultTlsContainerRef is a reference to a Barbican container of TLS secrets.
	// +optional
	DefaultTlsContainerRef string `json:"defaultTlsContainerRef,omitempty"`

	// SniContainerRefs is a list of references to TLS secrets for SNI.
	// +optional
	SniContainerRefs []string `json:"sniContainerRefs,omitempty"`

	// AdminStateUp is the administrative state of the listener (true = UP).
	// +kubebuilder:default=true
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// ProjectID is the UUID of the project who owns the listener.
	// +optional
	ProjectID string `json:"projectId,omitempty"`

	// TimeoutClientData is the frontend client inactivity timeout in milliseconds.
	// +optional
	TimeoutClientData *int `json:"timeoutClientData,omitempty"`

	// TimeoutMemberData is the backend member inactivity timeout in milliseconds.
	// +optional
	TimeoutMemberData *int `json:"timeoutMemberData,omitempty"`

	// TimeoutMemberConnect is the backend member connection timeout in milliseconds.
	// +optional
	TimeoutMemberConnect *int `json:"timeoutMemberConnect,omitempty"`

	// TimeoutTCPInspect is the time in milliseconds to wait for additional TCP packets.
	// +optional
	TimeoutTCPInspect *int `json:"timeoutTcpInspect,omitempty"`

	// InsertHeaders is a dictionary of optional headers to insert into the request.
	// +optional
	InsertHeaders map[string]string `json:"insertHeaders,omitempty"`

	// AllowedCIDRs is a list of IPv4/IPv6 CIDRs.
	// +optional
	AllowedCIDRs []string `json:"allowedCidrs,omitempty"`

	// Tags is a set of resource tags.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// ListenerStatus defines the observed state of an OpenStack listener.
type ListenerStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider ListenerProviderStatus `json:"atProvider,omitempty"`
}

// ListenerProviderStatus defines the observed state of the listener at the provider.
type ListenerProviderStatus struct {
	// ListenerID is the unique identifier of the listener in OpenStack.
	ListenerID string `json:"listenerId,omitempty"`

	// Name is the human-readable name of the listener.
	Name string `json:"name,omitempty"`

	// Description of the listener.
	Description string `json:"description,omitempty"`

	// LoadbalancerID is the load balancer this listener belongs to.
	LoadbalancerID string `json:"loadbalancerId,omitempty"`

	// Protocol is the listener protocol.
	Protocol string `json:"protocol,omitempty"`

	// ProtocolPort is the listener port.
	ProtocolPort int `json:"protocolPort,omitempty"`

	// DefaultPoolID is the default pool ID.
	DefaultPoolID string `json:"defaultPoolId,omitempty"`

	// ConnLimit is the connection limit.
	ConnLimit int `json:"connLimit,omitempty"`

	// AdminStateUp is the administrative state.
	AdminStateUp bool `json:"adminStateUp,omitempty"`

	// ProvisioningStatus is the provisioning status.
	ProvisioningStatus string `json:"provisioningStatus,omitempty"`

	// OperatingStatus is the operating status.
	OperatingStatus string `json:"operatingStatus,omitempty"`

	// ProjectID is the project owner.
	ProjectID string `json:"projectId,omitempty"`

	// Tags on the listener.
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
// +kubebuilder:printcolumn:name="Protocol",type="string",JSONPath=".spec.forProvider.protocol"
// +kubebuilder:printcolumn:name="Port",type="integer",JSONPath=".spec.forProvider.protocolPort"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.provisioningStatus"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// Listener is a managed resource that represents an OpenStack Octavia listener.
type Listener struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ListenerSpec   `json:"spec"`
	Status ListenerStatus `json:"status,omitempty"`
}

// ListenerSpec defines the desired state of a Listener.
type ListenerSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     ListenerParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// ListenerList contains a list of Listener resources.
type ListenerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Listener `json:"items"`
}

// SessionPersistence defines the session persistence configuration.
type SessionPersistence struct {
	// Type is the persistence mode (SOURCE_IP, HTTP_COOKIE, APP_COOKIE).
	// +kubebuilder:validation:Enum=SOURCE_IP;HTTP_COOKIE;APP_COOKIE
	// +optional
	Type string `json:"type,omitempty"`

	// CookieName is the cookie name for APP_COOKIE persistence.
	// +optional
	CookieName string `json:"cookieName,omitempty"`
}

// PoolParameters define the desired state of an OpenStack Octavia pool.
type PoolParameters struct {
	// LBMethod is the load balancing algorithm (ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP, SOURCE_IP_PORT).
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=ROUND_ROBIN;LEAST_CONNECTIONS;SOURCE_IP;SOURCE_IP_PORT
	LBMethod string `json:"lbMethod"`

	// Protocol is the protocol of the pool (TCP, UDP, HTTP, HTTPS).
	// +kubebuilder:validation:Required
	Protocol string `json:"protocol"`

	// LoadbalancerID is the UUID of the load balancer. One of LoadbalancerID or ListenerID must be provided.
	// +optional
	LoadbalancerID string `json:"loadbalancerId,omitempty"`

	// ListenerID is the UUID of the listener. One of LoadbalancerID or ListenerID must be provided.
	// +optional
	ListenerID string `json:"listenerId,omitempty"`

	// Name is the human-readable name for the pool.
	// +optional
	Name string `json:"name,omitempty"`

	// Description of the pool.
	// +optional
	Description string `json:"description,omitempty"`

	// Persistence is the session persistence configuration.
	// +optional
	Persistence *SessionPersistence `json:"persistence,omitempty"`

	// AdminStateUp is the administrative state of the pool (true = UP).
	// +kubebuilder:default=true
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// ProjectID is the UUID of the project who owns the pool.
	// +optional
	ProjectID string `json:"projectId,omitempty"`

	// Tags is a set of resource tags.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// PoolStatus defines the observed state of an OpenStack pool.
type PoolStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider PoolProviderStatus `json:"atProvider,omitempty"`
}

// PoolProviderStatus defines the observed state of the pool at the provider.
type PoolProviderStatus struct {
	// PoolID is the unique identifier of the pool in OpenStack.
	PoolID string `json:"poolId,omitempty"`

	// Name is the human-readable name of the pool.
	Name string `json:"name,omitempty"`

	// Description of the pool.
	Description string `json:"description,omitempty"`

	// LBMethod is the load balancing algorithm.
	LBMethod string `json:"lbMethod,omitempty"`

	// Protocol is the pool protocol.
	Protocol string `json:"protocol,omitempty"`

	// AdminStateUp is the administrative state.
	AdminStateUp bool `json:"adminStateUp,omitempty"`

	// ProvisioningStatus is the provisioning status.
	ProvisioningStatus string `json:"provisioningStatus,omitempty"`

	// OperatingStatus is the operating status.
	OperatingStatus string `json:"operatingStatus,omitempty"`

	// ProjectID is the project owner.
	ProjectID string `json:"projectId,omitempty"`

	// MonitorID is the associated health monitor ID.
	MonitorID string `json:"monitorId,omitempty"`

	// Tags on the pool.
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
// +kubebuilder:printcolumn:name="Algorithm",type="string",JSONPath=".spec.forProvider.lbMethod"
// +kubebuilder:printcolumn:name="Protocol",type="string",JSONPath=".spec.forProvider.protocol"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.provisioningStatus"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// Pool is a managed resource that represents an OpenStack Octavia pool.
type Pool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PoolSpec   `json:"spec"`
	Status PoolStatus `json:"status,omitempty"`
}

// PoolSpec defines the desired state of a Pool.
type PoolSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     PoolParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// PoolList contains a list of Pool resources.
type PoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pool `json:"items"`
}

// MemberParameters define the desired state of an OpenStack Octavia pool member.
type MemberParameters struct {
	// PoolID is the UUID of the pool to which this member belongs.
	// +kubebuilder:validation:Required
	PoolID string `json:"poolId"`

	// Address is the IP address of the member.
	// +kubebuilder:validation:Required
	Address string `json:"address"`

	// ProtocolPort is the port on which the application is hosted.
	// +kubebuilder:validation:Required
	ProtocolPort int `json:"protocolPort"`

	// Name is the human-readable name for the member.
	// +optional
	Name string `json:"name,omitempty"`

	// Weight is the relative portion of traffic this member should receive.
	// +optional
	Weight *int `json:"weight,omitempty"`

	// SubnetID is the UUID of the subnet.
	// +optional
	SubnetID string `json:"subnetId,omitempty"`

	// AdminStateUp is the administrative state of the member (true = UP).
	// +kubebuilder:default=true
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// ProjectID is the UUID of the project who owns the member.
	// +optional
	ProjectID string `json:"projectId,omitempty"`

	// Backup indicates if the member is a backup member.
	// +kubebuilder:default=false
	// +optional
	Backup *bool `json:"backup,omitempty"`

	// MonitorAddress is an alternate IP address used for health monitoring.
	// +optional
	MonitorAddress string `json:"monitorAddress,omitempty"`

	// MonitorPort is an alternate protocol port used for health monitoring.
	// +optional
	MonitorPort *int `json:"monitorPort,omitempty"`

	// Tags is a set of resource tags.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// MemberStatus defines the observed state of an OpenStack pool member.
type MemberStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider MemberProviderStatus `json:"atProvider,omitempty"`
}

// MemberProviderStatus defines the observed state of the member at the provider.
type MemberProviderStatus struct {
	// MemberID is the unique identifier of the member in OpenStack.
	MemberID string `json:"memberId,omitempty"`

	// PoolID is the pool this member belongs to.
	PoolID string `json:"poolId,omitempty"`

	// Name is the human-readable name of the member.
	Name string `json:"name,omitempty"`

	// Address is the IP address of the member.
	Address string `json:"address,omitempty"`

	// ProtocolPort is the port on which the application is hosted.
	ProtocolPort int `json:"protocolPort,omitempty"`

	// Weight is the relative weight of the member.
	Weight int `json:"weight,omitempty"`

	// SubnetID is the UUID of the subnet.
	SubnetID string `json:"subnetId,omitempty"`

	// ProjectID is the project owner.
	ProjectID string `json:"projectId,omitempty"`

	// AdminStateUp is the administrative state.
	AdminStateUp bool `json:"adminStateUp,omitempty"`

	// ProvisioningStatus is the provisioning status.
	ProvisioningStatus string `json:"provisioningStatus,omitempty"`

	// OperatingStatus is the operating status.
	OperatingStatus string `json:"operatingStatus,omitempty"`

	// Backup indicates if the member is a backup.
	Backup bool `json:"backup,omitempty"`

	// MonitorAddress is the alternate health monitor IP.
	MonitorAddress string `json:"monitorAddress,omitempty"`

	// MonitorPort is the alternate health monitor port.
	MonitorPort int `json:"monitorPort,omitempty"`

	// Tags on the member.
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
// +kubebuilder:printcolumn:name="Address",type="string",JSONPath=".spec.forProvider.address"
// +kubebuilder:printcolumn:name="Port",type="integer",JSONPath=".spec.forProvider.protocolPort"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.provisioningStatus"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// Member is a managed resource that represents an OpenStack Octavia pool member.
type Member struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MemberSpec   `json:"spec"`
	Status MemberStatus `json:"status,omitempty"`
}

// MemberSpec defines the desired state of a Member.
type MemberSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     MemberParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// MemberList contains a list of Member resources.
type MemberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Member `json:"items"`
}

// HealthMonitorParameters define the desired state of an OpenStack Octavia health monitor.
type HealthMonitorParameters struct {
	// Type is the probe type (PING, TCP, HTTP, HTTPS, TLS-HELLO, UDP-CONNECT, SCTP).
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=PING;TCP;HTTP;HTTPS;TLS-HELLO;UDP-CONNECT;SCTP
	Type string `json:"type"`

	// Delay is the time in seconds between sending probes to members.
	// +kubebuilder:validation:Required
	Delay int `json:"delay"`

	// Timeout is the maximum number of seconds to wait for a connection before timing out.
	// +kubebuilder:validation:Required
	Timeout int `json:"timeout"`

	// MaxRetries is the number of allowed connection failures before marking a member as INACTIVE.
	// +kubebuilder:validation:Required
	MaxRetries int `json:"maxRetries"`

	// MaxRetriesDown is the number of allowed connection failures before marking a member as ERROR.
	// +optional
	MaxRetriesDown int `json:"maxRetriesDown,omitempty"`

	// URLPath is the HTTP path of the request sent by the monitor. Required for HTTP(S) types.
	// +optional
	URLPath string `json:"urlPath,omitempty"`

	// HTTPMethod is the HTTP method used for requests. Defaults to GET.
	// +optional
	HTTPMethod string `json:"httpMethod,omitempty"`

	// HTTPVersion is the HTTP version (1.0 or 1.1). Defaults to 1.0.
	// +optional
	HTTPVersion string `json:"httpVersion,omitempty"`

	// ExpectedCodes is the expected HTTP codes for a passing HTTP(S) monitor.
	// +optional
	ExpectedCodes string `json:"expectedCodes,omitempty"`

	// DomainName is the HTTP host header injected into the HTTP request.
	// +optional
	DomainName string `json:"domainName,omitempty"`

	// AdminStateUp is the administrative state of the monitor (true = UP).
	// +kubebuilder:default=true
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// Name is the human-readable name for the monitor.
	// +optional
	Name string `json:"name,omitempty"`

	// ProjectID is the UUID of the project who owns the monitor.
	// +optional
	ProjectID string `json:"projectId,omitempty"`

	// Tags is a set of resource tags.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// HealthMonitorStatus defines the observed state of an OpenStack health monitor.
type HealthMonitorStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider HealthMonitorProviderStatus `json:"atProvider,omitempty"`
}

// HealthMonitorProviderStatus defines the observed state of the monitor at the provider.
type HealthMonitorProviderStatus struct {
	// MonitorID is the unique identifier of the health monitor in OpenStack.
	MonitorID string `json:"monitorId,omitempty"`

	// Name is the human-readable name of the monitor.
	Name string `json:"name,omitempty"`

	// Type is the probe type.
	Type string `json:"type,omitempty"`

	// Delay is the time between probes.
	Delay int `json:"delay,omitempty"`

	// Timeout is the probe timeout.
	Timeout int `json:"timeout,omitempty"`

	// MaxRetries is the failure threshold.
	MaxRetries int `json:"maxRetries,omitempty"`

	// MaxRetriesDown is the error threshold.
	MaxRetriesDown int `json:"maxRetriesDown,omitempty"`

	// HTTPMethod is the HTTP method.
	HTTPMethod string `json:"httpMethod,omitempty"`

	// URLPath is the HTTP path.
	URLPath string `json:"urlPath,omitempty"`

	// ExpectedCodes is the expected HTTP codes.
	ExpectedCodes string `json:"expectedCodes,omitempty"`

	// DomainName is the HTTP host header.
	DomainName string `json:"domainName,omitempty"`

	// AdminStateUp is the administrative state.
	AdminStateUp bool `json:"adminStateUp,omitempty"`

	// ProvisioningStatus is the provisioning status.
	ProvisioningStatus string `json:"provisioningStatus,omitempty"`

	// OperatingStatus is the operating status.
	OperatingStatus string `json:"operatingStatus,omitempty"`

	// ProjectID is the project owner.
	ProjectID string `json:"projectId,omitempty"`

	// Tags on the monitor.
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
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.forProvider.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.provisioningStatus"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// HealthMonitor is a managed resource that represents an OpenStack Octavia health monitor.
type HealthMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HealthMonitorSpec   `json:"spec"`
	Status HealthMonitorStatus `json:"status,omitempty"`
}

// HealthMonitorSpec defines the desired state of a HealthMonitor.
type HealthMonitorSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     HealthMonitorParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// HealthMonitorList contains a list of HealthMonitor resources.
type HealthMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HealthMonitor `json:"items"`
}
