package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ZoneParameters struct {
	Name        string            `json:"name"`
	Email       string            `json:"email,omitempty"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"type,omitempty"`
	TTL         *int              `json:"ttl,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

type ZoneProviderStatus struct {
	ZoneID     string            `json:"zoneId,omitempty"`
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Email      string            `json:"email,omitempty"`
	TTL        int               `json:"ttl,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type ZoneSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     ZoneParameters `json:"forProvider"`
}

type ZoneStatus struct {
	xpv2.ConditionedStatus `json:",inline"`
	AtProvider             ZoneProviderStatus `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
type Zone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZoneSpec   `json:"spec"`
	Status ZoneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zone `json:"items"`
}
