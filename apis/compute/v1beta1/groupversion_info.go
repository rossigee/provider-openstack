/*
Copyright 2025 The Crossplane Authors.
Licensed under the Apache License, Version 2.0.
*/

// +groupName=compute.openstack.m.crossplane.io
// +versionName=v1beta1
package v1beta1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	APIGroup = "compute.openstack.m.crossplane.io"
)

var SchemeGroupVersion = schema.GroupVersion{Group: APIGroup, Version: "v1beta1"}

var SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
var AddToScheme = SchemeBuilder.AddToScheme

func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(SchemeGroupVersion,
		&Server{},
		&ServerList{},
		&KeyPair{},
		&KeyPairList{},
	)
	return nil
}

// Server type metadata.
var (
	ServerKind             = reflect.TypeOf(Server{}).Name()
	ServerGroupKind        = schema.GroupKind{Group: APIGroup, Kind: ServerKind}
	ServerKindAPIVersion   = ServerKind + "." + SchemeGroupVersion.String()
	ServerGroupVersionKind = SchemeGroupVersion.WithKind(ServerKind)
)

// KeyPair type metadata.
var (
	KeyPairKind             = reflect.TypeOf(KeyPair{}).Name()
	KeyPairGroupKind        = schema.GroupKind{Group: APIGroup, Kind: KeyPairKind}
	KeyPairKindAPIVersion   = KeyPairKind + "." + SchemeGroupVersion.String()
	KeyPairGroupVersionKind = SchemeGroupVersion.WithKind(KeyPairKind)
)
