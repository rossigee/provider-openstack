/*
Copyright 2025 The Crossplane Authors.
Licensed under the Apache License, Version 2.0.
*/

package v1beta1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	APIGroup = "networking.openstack.m.crossplane.io"
)

var SchemeGroupVersion = schema.GroupVersion{Group: APIGroup, Version: "v1beta1"}

var SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
var AddToScheme = SchemeBuilder.AddToScheme

func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(SchemeGroupVersion,
		&Network{},
		&NetworkList{},
		&Port{},
		&PortList{},
	)
	return nil
}

// Network type metadata.
var (
	NetworkKind             = reflect.TypeOf(Network{}).Name()
	NetworkGroupKind        = schema.GroupKind{Group: APIGroup, Kind: NetworkKind}
	NetworkKindAPIVersion   = NetworkKind + "." + SchemeGroupVersion.String()
	NetworkGroupVersionKind = SchemeGroupVersion.WithKind(NetworkKind)
)

// Port type metadata.
var (
	PortKind             = reflect.TypeOf(Port{}).Name()
	PortGroupKind        = schema.GroupKind{Group: APIGroup, Kind: PortKind}
	PortKindAPIVersion   = PortKind + "." + SchemeGroupVersion.String()
	PortGroupVersionKind = SchemeGroupVersion.WithKind(PortKind)
)
