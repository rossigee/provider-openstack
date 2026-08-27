package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ImageParameters struct {
	Name            string            `json:"name"`
	ContainerFormat string            `json:"containerFormat"`
	DiskFormat      string            `json:"diskFormat"`
	MinDisk         *int              `json:"minDisk,omitempty"`
	MinRAM          *int              `json:"minRam,omitempty"`
	Visibility      string            `json:"visibility"`
	Protected       *bool             `json:"protected,omitempty"`
	Tags            []string          `json:"tags,omitempty"`
	Properties      map[string]string `json:"properties,omitempty"`
}

type ImageProviderStatus struct {
	ImageID         string   `json:"imageId,omitempty"`
	Status          string   `json:"status,omitempty"`
	Size            int64    `json:"size,omitempty"`
	MinDisk         int      `json:"minDisk,omitempty"`
	MinRAM          int      `json:"minRam,omitempty"`
	Protected       *bool    `json:"protected,omitempty"`
	Visibility      string   `json:"visibility,omitempty"`
	ContainerFormat string   `json:"containerFormat,omitempty"`
	DiskFormat      string   `json:"diskFormat,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Owner           string   `json:"owner,omitempty"`
	File            string   `json:"file,omitempty"`
	Schema          string   `json:"schema,omitempty"`
	Checksum        string   `json:"checksum,omitempty"`
	VirtualSize     int64    `json:"virtualSize,omitempty"`
}

type ImageSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     ImageParameters `json:"forProvider"`
}

type ImageStatus struct {
	xpv2.ConditionedStatus `json:",inline"`
	AtProvider             ImageProviderStatus `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec"`
	Status ImageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}
