package v1beta1

import (
	xpv1 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider              ImageParameters `json:"forProvider"`
}

type ImageStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
	AtProvider             ImageProviderStatus `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,openstack}
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

// DeepCopyObject implements runtime.Object
func (in *Image) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Image)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies the receiver into out
func (in *Image) DeepCopyInto(out *Image) {
	*out = *in
	out.ObjectMeta = in.ObjectMeta
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopyObject implements runtime.Object
func (in *ImageList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(ImageList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies the receiver into out
func (in *ImageList) DeepCopyInto(out *ImageList) {
	*out = *in
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]Image, len(in.Items))
		copy(out.Items, in.Items)
	}
}

// GetCondition implements resource.Managed
func (in *Image) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return in.Status.GetCondition(ct)
}

// SetConditions implements resource.Managed
func (in *Image) SetConditions(c ...xpv1.Condition) {
	in.Status.SetConditions(c...)
}

// GetManagementPolicies implements resource.Managed
func (in *Image) GetManagementPolicies() xpv1.ManagementPolicies {
	return in.Spec.ManagementPolicies
}

// SetManagementPolicies implements resource.Managed
func (in *Image) SetManagementPolicies(p xpv1.ManagementPolicies) {
	in.Spec.ManagementPolicies = p
}
