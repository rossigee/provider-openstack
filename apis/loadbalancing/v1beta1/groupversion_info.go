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
		&Listener{},
		&ListenerList{},
		&Pool{},
		&PoolList{},
		&Member{},
		&MemberList{},
		&HealthMonitor{},
		&HealthMonitorList{},
	)
	return nil
}

var (
	LoadBalancerKind             = reflect.TypeOf(LoadBalancer{}).Name()
	LoadBalancerGroupKind        = schema.GroupKind{Group: APIGroup, Kind: LoadBalancerKind}
	LoadBalancerKindAPIVersion   = LoadBalancerKind + "." + SchemeGroupVersion.String()
	LoadBalancerGroupVersionKind = SchemeGroupVersion.WithKind(LoadBalancerKind)
)

var (
	ListenerKind             = reflect.TypeOf(Listener{}).Name()
	ListenerGroupKind        = schema.GroupKind{Group: APIGroup, Kind: ListenerKind}
	ListenerKindAPIVersion   = ListenerKind + "." + SchemeGroupVersion.String()
	ListenerGroupVersionKind = SchemeGroupVersion.WithKind(ListenerKind)
)

var (
	PoolKind             = reflect.TypeOf(Pool{}).Name()
	PoolGroupKind        = schema.GroupKind{Group: APIGroup, Kind: PoolKind}
	PoolKindAPIVersion   = PoolKind + "." + SchemeGroupVersion.String()
	PoolGroupVersionKind = SchemeGroupVersion.WithKind(PoolKind)
)

var (
	MemberKind             = reflect.TypeOf(Member{}).Name()
	MemberGroupKind        = schema.GroupKind{Group: APIGroup, Kind: MemberKind}
	MemberKindAPIVersion   = MemberKind + "." + SchemeGroupVersion.String()
	MemberGroupVersionKind = SchemeGroupVersion.WithKind(MemberKind)
)

var (
	HealthMonitorKind             = reflect.TypeOf(HealthMonitor{}).Name()
	HealthMonitorGroupKind        = schema.GroupKind{Group: APIGroup, Kind: HealthMonitorKind}
	HealthMonitorKindAPIVersion   = HealthMonitorKind + "." + SchemeGroupVersion.String()
	HealthMonitorGroupVersionKind = SchemeGroupVersion.WithKind(HealthMonitorKind)
)
