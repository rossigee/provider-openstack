package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type UserParameters struct {
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	DomainID         string            `json:"domainId,omitempty"`
	DefaultProjectID string            `json:"defaultProjectId,omitempty"`
	Email            string            `json:"email,omitempty"`
	Password         string            `json:"password,omitempty"`
	Enabled          *bool             `json:"enabled,omitempty"`
	Options          map[string]string `json:"options,omitempty"`
}

type UserProviderStatus struct {
	UserID            string      `json:"userId,omitempty"`
	Name              string      `json:"name,omitempty"`
	Description       string      `json:"description,omitempty"`
	DomainID          string      `json:"domainId,omitempty"`
	DefaultProjectID  string      `json:"defaultProjectId,omitempty"`
	Enabled           *bool       `json:"enabled,omitempty"`
	PasswordExpiresAt metav1.Time `json:"passwordExpiresAt,omitempty"`
}

type UserSpec struct {
	xpv2.ClusterManagedResourceSpec `json:",inline"`
	ForProvider                     UserParameters `json:"forProvider"`
}

type UserStatus struct {
	xpv2.ConditionedStatus `json:",inline"`
	AtProvider             UserProviderStatus `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,openstack}
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserSpec   `json:"spec"`
	Status UserStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}
