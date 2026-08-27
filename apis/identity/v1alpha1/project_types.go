package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProjectParameters struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	DomainID    string   `json:"domainId,omitempty"`
	IsDomain    *bool    `json:"isDomain,omitempty"`
	ParentID    string   `json:"parentId,omitempty"`
	Enabled     *bool    `json:"enabled,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type ProjectProviderStatus struct {
	ProjectID   string   `json:"projectId,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	DomainID    string   `json:"domainId,omitempty"`
	Enabled     *bool    `json:"enabled,omitempty"`
	ParentID    string   `json:"parentId,omitempty"`
	IsDomain    *bool    `json:"isDomain,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type ProjectSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     ProjectParameters `json:"forProvider"`
}

type ProjectStatus struct {
	xpv2.ConditionedStatus `json:",inline"`
	AtProvider             ProjectProviderStatus `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectSpec   `json:"spec"`
	Status ProjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}
