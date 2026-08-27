package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RecordSetParameters struct {
	Name    string   `json:"name"`
	ZoneID  string   `json:"zoneId"`
	Type    string   `json:"type"`
	TTL     int      `json:"ttl,omitempty"`
	Records []string `json:"records"`
}

type RecordSetProviderStatus struct {
	RecordSetID string   `json:"recordSetId,omitempty"`
	Name        string   `json:"name,omitempty"`
	ZoneID      string   `json:"zoneId,omitempty"`
	ZoneName    string   `json:"zoneName,omitempty"`
	Type        string   `json:"type,omitempty"`
	TTL         int      `json:"ttl,omitempty"`
	Records     []string `json:"records,omitempty"`
	Status      string   `json:"status,omitempty"`
	Action      string   `json:"action,omitempty"`
}

type RecordSetSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     RecordSetParameters `json:"forProvider"`
}

type RecordSetStatus struct {
	xpv2.ConditionedStatus `json:",inline"`
	AtProvider             RecordSetProviderStatus `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
type RecordSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordSetSpec   `json:"spec"`
	Status RecordSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type RecordSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecordSet `json:"items"`
}
