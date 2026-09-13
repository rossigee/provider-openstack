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
	APIGroup = "loadbalancing.openstack.m.crossplane.io"
)

var SchemeGroupVersion = schema.GroupVersion{Group: APIGroup, Version: "v1beta1"}

var SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
var AddToScheme = SchemeBuilder.AddToScheme

func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(SchemeGroupVersion,
		&LoadBalancer{},
		&LoadBalancerList{},
	)
	return nil
}

// LoadBalancer type metadata.
var (
	LoadBalancerKind             = reflect.TypeOf(LoadBalancer{}).Name()
	LoadBalancerGroupKind        = schema.GroupKind{Group: APIGroup, Kind: LoadBalancerKind}
	LoadBalancerKindAPIVersion   = LoadBalancerKind + "." + SchemeGroupVersion.String()
	LoadBalancerGroupVersionKind = SchemeGroupVersion.WithKind(LoadBalancerKind)
)
