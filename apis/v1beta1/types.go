package v1beta1

import (
	xpv1 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// A ProviderConfigSpec defines the desired state of a ProviderConfig.
type ProviderConfigSpec struct {
	Credentials ProviderCredentials `json:"credentials"`
}

// ProviderCredentials required to authenticate.
type ProviderCredentials struct {
	// +kubebuilder:validation:Enum=None;Secret;InjectedIdentity;Environment;Filesystem
	Source xpv1.CredentialsSource `json:"source"`

	xpv1.CommonCredentialSelectors `json:",inline"`
}

// A ProviderConfigStatus reflects the observed state of a ProviderConfig.
type ProviderConfigStatus struct {
	xpv1.ProviderConfigStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,provider,openstack}
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig.
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}

// +kubebuilder:object:root=true

// A ProviderConfigUsage indicates that a resource is using a ProviderConfig.
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="CONFIG-NAME",type="string",JSONPath=".providerConfigRef.name"
// +kubebuilder:printcolumn:name="RESOURCE-KIND",type="string",JSONPath=".resourceRef.kind"
// +kubebuilder:printcolumn:name="RESOURCE-NAME",type="string",JSONPath=".resourceRef.name"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,provider,openstack}
type ProviderConfigUsage struct {
	metav1.TypeMeta               `json:",inline"`
	metav1.ObjectMeta             `json:"metadata,omitempty"`
	xpv1.TypedProviderConfigUsage `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderConfigUsageList contains a list of ProviderConfigUsage
type ProviderConfigUsageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfigUsage `json:"items"`
}

// DeepCopyObject implements runtime.Object
func (in *ProviderConfig) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(ProviderConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies the receiver into out
func (in *ProviderConfig) DeepCopyInto(out *ProviderConfig) {
	*out = *in
	out.ObjectMeta = in.ObjectMeta
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopyObject implements runtime.Object
func (in *ProviderConfigList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(ProviderConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies the receiver into out
func (in *ProviderConfigList) DeepCopyInto(out *ProviderConfigList) {
	*out = *in
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]ProviderConfig, len(in.Items))
		copy(out.Items, in.Items)
	}
}

// DeepCopyObject implements runtime.Object
func (in *ProviderConfigUsage) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(ProviderConfigUsage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies the receiver into out
func (in *ProviderConfigUsage) DeepCopyInto(out *ProviderConfigUsage) {
	*out = *in
	out.ObjectMeta = in.ObjectMeta
	out.TypedProviderConfigUsage = in.TypedProviderConfigUsage
}

// DeepCopyObject implements runtime.Object
func (in *ProviderConfigUsageList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(ProviderConfigUsageList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies the receiver into out
func (in *ProviderConfigUsageList) DeepCopyInto(out *ProviderConfigUsageList) {
	*out = *in
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]ProviderConfigUsage, len(in.Items))
		copy(out.Items, in.Items)
	}
}
