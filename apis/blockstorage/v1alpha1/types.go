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

// VolumeParameters define the desired state of an OpenStack Cinder volume.
type VolumeParameters struct {
	// Size is the size of the volume, in GB.
	// +kubebuilder:validation:Required
	Size int `json:"size"`

	// Name is the human-readable name for the volume.
	// +optional
	Name string `json:"name,omitempty"`

	// Description of the volume.
	// +optional
	Description string `json:"description,omitempty"`

	// AvailabilityZone is the availability zone for the volume.
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// VolumeType is the associated volume type.
	// +optional
	VolumeType string `json:"volumeType,omitempty"`

	// ImageID is the ID of the image from which to create a bootable volume.
	// +optional
	ImageID string `json:"imageId,omitempty"`

	// SnapshotID is the ID of an existing volume snapshot to create the volume from.
	// +optional
	SnapshotID string `json:"snapshotId,omitempty"`

	// SourceVolID is the ID of an existing volume to create the volume from.
	// +optional
	SourceVolID string `json:"sourceVolId,omitempty"`

	// Metadata is a set of key/value pairs to associate with the volume.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}

// VolumeStatus defines the observed state of an OpenStack volume.
type VolumeStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider VolumeProviderStatus `json:"atProvider,omitempty"`

	ConnectionDetails []byte `json:"connectionDetails,omitempty"`
}

// VolumeProviderStatus defines the observed state of the volume at the provider.
type VolumeProviderStatus struct {
	// VolumeID is the unique identifier of the volume.
	VolumeID string `json:"volumeId,omitempty"`

	// Status is the current status of the volume (available, creating, error, etc.).
	Status string `json:"status,omitempty"`

	// Size is the size of the volume in GB.
	Size int `json:"size,omitempty"`

	// AvailabilityZone is the availability zone of the volume.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// VolumeType is the volume type.
	VolumeType string `json:"volumeType,omitempty"`

	// Description of the volume.
	Description string `json:"description,omitempty"`

	// Bootable indicates whether this is a bootable volume.
	Bootable string `json:"bootable,omitempty"`

	// Encrypted denotes if the volume is encrypted.
	Encrypted bool `json:"encrypted,omitempty"`

	// SnapshotID is the snapshot the volume was created from.
	SnapshotID string `json:"snapshotId,omitempty"`

	// SourceVolID is the source volume the volume was created from.
	SourceVolID string `json:"sourceVolId,omitempty"`

	// TenantID is the project that owns the volume.
	TenantID string `json:"tenantId,omitempty"`

	// Host is the identifier of the host holding the volume.
	Host string `json:"host,omitempty"`

	// Metadata is the user-defined metadata of the volume.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Size",type="integer",JSONPath=".status.atProvider.size"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// Volume is a managed resource that represents an OpenStack Cinder volume.
type Volume struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VolumeSpec   `json:"spec"`
	Status VolumeStatus `json:"status,omitempty"`
}

// VolumeSpec defines the desired state of a Volume.
type VolumeSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     VolumeParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// VolumeList contains a list of Volume resources.
type VolumeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Volume `json:"items"`
}

// VolumeTypeParameters define the desired state of an OpenStack Cinder volume type.
type VolumeTypeParameters struct {
	// Name is the human-readable name for the volume type.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description of the volume type.
	// +optional
	Description string `json:"description,omitempty"`

	// IsPublic indicates whether the volume type is publicly visible.
	// +kubebuilder:default=true
	// +optional
	IsPublic *bool `json:"isPublic,omitempty"`

	// ExtraSpecs is a set of extra-spec key/value pairs defined by the user.
	// +optional
	ExtraSpecs map[string]string `json:"extraSpecs,omitempty"`
}

// VolumeTypeStatus defines the observed state of an OpenStack volume type.
type VolumeTypeStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider VolumeTypeProviderStatus `json:"atProvider,omitempty"`
}

// VolumeTypeProviderStatus defines the observed state of the volume type at the provider.
type VolumeTypeProviderStatus struct {
	// VolumeTypeID is the unique identifier of the volume type.
	VolumeTypeID string `json:"volumeTypeId,omitempty"`

	// Name is the name of the volume type.
	Name string `json:"name,omitempty"`

	// Description of the volume type.
	Description string `json:"description,omitempty"`

	// IsPublic indicates whether the volume type is publicly visible.
	IsPublic bool `json:"isPublic,omitempty"`

	// ExtraSpecs is the set of extra-spec key/value pairs.
	ExtraSpecs map[string]string `json:"extraSpecs,omitempty"`

	// QosSpecID is the QoS spec ID.
	QosSpecID string `json:"qosSpecId,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// VolumeType is a managed resource that represents an OpenStack Cinder volume type.
type VolumeType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VolumeTypeSpec   `json:"spec"`
	Status VolumeTypeStatus `json:"status,omitempty"`
}

// VolumeTypeSpec defines the desired state of a VolumeType.
type VolumeTypeSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     VolumeTypeParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// VolumeTypeList contains a list of VolumeType resources.
type VolumeTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VolumeType `json:"items"`
}

// VolumeSnapshotParameters define the desired state of an OpenStack Cinder snapshot.
type VolumeSnapshotParameters struct {
	// VolumeID is the ID of the volume to snapshot.
	// +kubebuilder:validation:Required
	VolumeID string `json:"volumeId"`

	// Name is the human-readable name for the snapshot.
	// +optional
	Name string `json:"name,omitempty"`

	// Description of the snapshot.
	// +optional
	Description string `json:"description,omitempty"`

	// Force indicates whether to snapshot a volume even if it is attached.
	// +kubebuilder:default=false
	// +optional
	Force bool `json:"force,omitempty"`

	// Metadata is a set of key/value pairs to associate with the snapshot.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}

// VolumeSnapshotStatus defines the observed state of an OpenStack snapshot.
type VolumeSnapshotStatus struct {
	xpv2.ConditionedStatus `json:",inline"`

	AtProvider VolumeSnapshotProviderStatus `json:"atProvider,omitempty"`
}

// VolumeSnapshotProviderStatus defines the observed state of the snapshot at the provider.
type VolumeSnapshotProviderStatus struct {
	// SnapshotID is the unique identifier of the snapshot.
	SnapshotID string `json:"snapshotId,omitempty"`

	// Name is the name of the snapshot.
	Name string `json:"name,omitempty"`

	// Description of the snapshot.
	Description string `json:"description,omitempty"`

	// VolumeID is the volume the snapshot was created from.
	VolumeID string `json:"volumeId,omitempty"`

	// Status is the current status (available, creating, error, etc.).
	Status string `json:"status,omitempty"`

	// Size is the size of the snapshot in GB.
	Size int `json:"size,omitempty"`

	// Progress of the snapshot creation.
	Progress string `json:"progress,omitempty"`

	// TenantID is the project that owns the snapshot.
	TenantID string `json:"tenantId,omitempty"`

	// Metadata is the user-defined metadata of the snapshot.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
// VolumeSnapshot is a managed resource that represents an OpenStack Cinder snapshot.
type VolumeSnapshot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VolumeSnapshotSpec   `json:"spec"`
	Status VolumeSnapshotStatus `json:"status,omitempty"`
}

// VolumeSnapshotSpec defines the desired state of a VolumeSnapshot.
type VolumeSnapshotSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     VolumeSnapshotParameters `json:"forProvider,omitempty"`
}

// +kubebuilder:object:root=true
// VolumeSnapshotList contains a list of VolumeSnapshot resources.
type VolumeSnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VolumeSnapshot `json:"items"`
}
