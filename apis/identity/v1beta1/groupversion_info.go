package v1beta1

import (
	"reflect"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const APIGroup = "identity.openstack.m.crossplane.io"

var SchemeGroupVersion = schema.GroupVersion{Group: APIGroup, Version: "v1beta1"}
var SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
var AddToScheme = SchemeBuilder.AddToScheme

func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(SchemeGroupVersion, &Project{}, &ProjectList{}, &User{}, &UserList{})
	return nil
}

var (
	ProjectKind = reflect.TypeOf(Project{}).Name()
	ProjectGroupKind = schema.GroupKind{Group: APIGroup, Kind: ProjectKind}
	UserKind = reflect.TypeOf(User{}).Name()
	UserGroupKind = schema.GroupKind{Group: APIGroup, Kind: UserKind}
)
