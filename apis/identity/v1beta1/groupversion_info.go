package v1beta1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	APIGroup = "identity.openstack.m.crossplane.io"
)

var SchemeGroupVersion = schema.GroupVersion{Group: APIGroup, Version: "v1beta1"}

var SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
var AddToScheme = SchemeBuilder.AddToScheme

func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(SchemeGroupVersion, &Project{}, &ProjectList{}, &User{}, &UserList{}, &Role{}, &RoleList{})
	return nil
}

var (
	UserKind             = reflect.TypeOf(User{}).Name()
	UserGroupKind        = schema.GroupKind{Group: APIGroup, Kind: UserKind}
	UserKindAPIVersion   = UserKind + "." + SchemeGroupVersion.String()
	UserGroupVersionKind = SchemeGroupVersion.WithKind(UserKind)
)

var (
	ProjectKind             = reflect.TypeOf(Project{}).Name()
	ProjectGroupKind        = schema.GroupKind{Group: APIGroup, Kind: ProjectKind}
	ProjectKindAPIVersion   = ProjectKind + "." + SchemeGroupVersion.String()
	ProjectGroupVersionKind = SchemeGroupVersion.WithKind(ProjectKind)
)

var (
	RoleKind             = reflect.TypeOf(Role{}).Name()
	RoleGroupKind        = schema.GroupKind{Group: APIGroup, Kind: RoleKind}
	RoleKindAPIVersion   = RoleKind + "." + SchemeGroupVersion.String()
	RoleGroupVersionKind = SchemeGroupVersion.WithKind(RoleKind)
)

type UserProviderStatus struct {
	UserID            string `json:"userId,omitempty"`
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	DomainID          string `json:"domainId,omitempty"`
	DefaultProjectID  string `json:"defaultProjectId,omitempty"`
	Enabled           *bool  `json:"enabled,omitempty"`
}

// RoleProviderStatus defines the observed state of the role at the provider
type RoleProviderStatus struct {
	RoleID      string `json:"roleId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ProjectProviderStatus defines the observed state of the project at the provider
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
