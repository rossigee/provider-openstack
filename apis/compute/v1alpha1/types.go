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

// Server parameters define the desired state of an OpenStack Nova server.
type ServerParameters struct {
	// Name is the human-readable name for the server.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// ImageRef is the ID or name of the image to use for the server.
	// +kubebuilder:validation:Required
	ImageRef string `json:"imageRef"`

	// FlavorRef is the ID or name of the flavor (compute sizing).
	// +kubebuilder:validation:Required
	FlavorRef string `json:"flavorRef"`

	// KeyName is the name of the key pair for SSH access.
	// +optional
	KeyName string `json:"keyName,omitempty"`

	// Networks defines the networks to attach to the server.
	// +optional
	Networks []ServerNetwork `json:"networks,omitempty"`

	// SecurityGroups is a list of security group names or IDs to apply.
	// +optional
	SecurityGroups []string `json:"securityGroups,omitempty"`

	// BlockDeviceMapping defines block devices to attach to the server.
	// +optional
	BlockDeviceMapping []BlockDeviceMapping `json:"blockDeviceMapping,omitempty"`

	// UserData is base64-encoded user data for cloud-init.
	// +optional
	UserData string `json:"userData,omitempty"`

	// TenantID is the project owner of the server.
	// +optional
	TenantID string `json:"tenantId,omitempty"`

	// AvailabilityZone is the AZ in which to create the server.
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// Description of the server.
	// +optional
	Description string `json:"description,omitempty"`

	// Tags is a set of arbitrary key/value pairs for server tagging.
	// +optional
	Tags []string `json:"tags,omitempty"`

	// ConfigDrive enables config drive for the server.
	// +kubebuilder:default=false
	// +optional
	ConfigDrive *bool `json:"configDrive,omitempty"`

	// AdminStateUp controls the administrative state of the server.
	// +kubebuilder:default=true
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`
}

// ServerNetwork defines a network attachment for a server.
type ServerNetwork struct {
	// UUID is the ID of the network to attach.
	// +optional
	UUID string `json:"uuid,omitempty"`

	// Name is the name of the network (resolved to UUID if UUID is not specified).
	// +optional
	Name string `json:"name,omitempty"`

	// FixedIP is a specific fixed IP to use from the network's subnet.
	// +optional
	FixedIP string `json:"fixedIP,omitempty"`

	// Port is the ID of an existing port to attach to the server.
	// +optional
	Port string `json:"port,omitempty"`
}

// BlockDeviceMapping defines a block device to attach to a server.
type BlockDeviceMapping struct {
	// UUID is the ID of the source image, volume, or snapshot.
	// +optional
	UUID string `json:"uuid,omitempty"`

	// SourceType is the source type (image, volume, snapshot, blank).
	// +kubebuilder:default=image
	// +kubebuilder:validation:Enum=image;volume;snapshot;blank
	// +optional
	SourceType string `json:"sourceType,omitempty"`

	// DestinationType is the destination type (local, volume).
	// +kubebuilder:default=local
	// +kubebuilder:validation:Enum=local;volume
	// +optional
	DestinationType string `json:"destinationType,omitempty"`

	// VolumeSize is the size in GB when creating a new volume.
	// +optional
	VolumeSize *int `json:"volumeSize,omitempty"`

	// VolumeType is the type of volume to create.
	// +optional
	VolumeType string `json:"volumeType,omitempty"`

	// DeleteOnTermination deletes the block device when the server is deleted.
	// +kubebuilder:default=true
	// +optional
	DeleteOnTermination *bool `json:"deleteOnTermination,omitempty"`

	// BootIndex is the boot order (-1 for disabled).
	// +kubebuilder:default=0
	// +optional
	BootIndex *int `json:"bootIndex,omitempty"`

	// GuestFormat is the format of the device (e.g., "swap", "ephemeral").
	// +optional
	GuestFormat string `json:"guestFormat,omitempty"`

	// DeviceName is the device name (e.g., "vda", "vdb").
	// +optional
	DeviceName string `json:"deviceName,omitempty"`
}

// ServerStatus defines the observed state of an OpenStack server.
type ServerStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider ServerProviderStatus `json:"atProvider,omitempty"`

	// ConnectionDetails stores connection details for the server.
	ConnectionDetails ServerConnectionDetails `json:"connectionDetails,omitempty"`
}

// ServerConnectionDetails stores connection details for a server.
type ServerConnectionDetails struct {
	// PrivateIPv4 is the private IPv4 address of the server.
	PrivateIPv4 string `json:"privateIPv4,omitempty"`

	// PublicIPv4 is the public (floating) IPv4 address if one is associated.
	PublicIPv4 string `json:"publicIPv4,omitempty"`

	// ServerState is the current state of the server (ACTIVE, SHUTOFF, etc.).
	ServerState string `json:"serverState,omitempty"`
}

// ServerProviderStatus defines the observed state of the server at the provider.
type ServerProviderStatus struct {
	// ServerID is the unique identifier of the server in OpenStack.
	ServerID string `json:"serverId,omitempty"`

	// Status is the current status of the server.
	Status string `json:"status,omitempty"`

	// TenantID is the project owner.
	TenantID string `json:"tenantId,omitempty"`

	// HostID is the host where the server runs.
	HostID string `json:"hostId,omitempty"`

	// Flavor is the flavor used for the server.
	Flavor string `json:"flavor,omitempty"`

	// Image is the image used for the server.
	Image string `json:"image,omitempty"`

	// Addresses is the list of IP addresses assigned to the server.
	Addresses []ServerAddress `json:"addresses,omitempty"`

	// PowerState is the current power state.
	PowerState int `json:"powerState,omitempty"`

	// VMState is the VM state (active, stopped, paused, etc.).
	VMState string `json:"vmState,omitempty"`

	// KeyName is the key pair name.
	KeyName string `json:"keyName,omitempty"`

	// Tags on the server.
	Tags []string `json:"tags,omitempty"`

	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`

	// UpdatedAt is the last update timestamp.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// ServerAddress describes an IP address on a server.
type ServerAddress struct {
	// Network is the network name.
	Network string `json:"network"`

	// Version is the IP version (4 or 6).
	Version int `json:"version"`

	// Address is the IP address.
	Address string `json:"address"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Private IPv4",type="string",JSONPath=".status.connectionDetails.privateIPv4"
// +kubebuilder:printcolumn:name="Flavor",type="string",JSONPath=".status.atProvider.flavor"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// Server is a managed resource that represents an OpenStack Nova server.
type Server struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerSpec   `json:"spec"`
	Status ServerStatus `json:"status,omitempty"`
}

// ServerSpec defines the desired state of a Server.
type ServerSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     ServerParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// ServerList contains a list of Server resources.
type ServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Server `json:"items"`
}

// KeyPair parameters define the desired state of an OpenStack Nova key pair.
type KeyPairParameters struct {
	// Name is the name of the key pair.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// PublicKey is the public key material. If not specified, OpenStack
	// generates a new key pair.
	// +optional
	PublicKey string `json:"publicKey,omitempty"`

	// Type is the key pair type (ssh or x509).
	// +kubebuilder:default=ssh
	// +kubebuilder:validation:Enum=ssh;x509
	// +optional
	Type string `json:"type,omitempty"`

	// TenantID is the project owner.
	// +optional
	TenantID string `json:"tenantId,omitempty"`
}

// KeyPairStatus defines the observed state of an OpenStack key pair.
type KeyPairStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider KeyPairProviderStatus `json:"atProvider,omitempty"`

	// ConnectionDetails stores the private key for generated key pairs.
	ConnectionDetails KeyPairConnectionDetails `json:"connectionDetails,omitempty"`
}

// KeyPairConnectionDetails stores connection details for a key pair.
type KeyPairConnectionDetails struct {
	// PrivateKey is the private key material (only available when generated by OpenStack).
	PrivateKey string `json:"privateKey,omitempty"`

	// PublicKey is the public key material.
	PublicKey string `json:"publicKey,omitempty"`

	// Fingerprint is the key fingerprint.
	Fingerprint string `json:"fingerprint,omitempty"`
}

// KeyPairProviderStatus defines the observed state of the key pair at the provider.
type KeyPairProviderStatus struct {
	// Name is the name of the key pair.
	Name string `json:"name,omitempty"`

	// PublicKey is the public key.
	PublicKey string `json:"publicKey,omitempty"`

	// Fingerprint is the key fingerprint.
	Fingerprint string `json:"fingerprint,omitempty"`

	// Type is the key pair type.
	Type string `json:"type,omitempty"`

	// TenantID is the project owner.
	TenantID string `json:"tenantId,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Fingerprint",type="string",JSONPath=".status.atProvider.fingerprint"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// KeyPair is a managed resource that represents an OpenStack Nova key pair.
type KeyPair struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeyPairSpec   `json:"spec"`
	Status KeyPairStatus `json:"status,omitempty"`
}

// KeyPairSpec defines the desired state of a KeyPair.
type KeyPairSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     KeyPairParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// KeyPairList contains a list of KeyPair resources.
type KeyPairList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeyPair `json:"items"`
}
