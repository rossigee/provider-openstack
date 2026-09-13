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
		&Subnet{},
		&SubnetList{},
		&Router{},
		&RouterList{},
		&RouterInterface{},
		&RouterInterfaceList{},
		&Port{},
		&PortList{},
		&SecurityGroup{},
		&SecurityGroupList{},
		&SecurityGroupRule{},
		&SecurityGroupRuleList{},
		&FloatingIP{},
		&FloatingIPList{},
		&SubnetPool{},
		&SubnetPoolList{},
		&Trunk{},
		&TrunkList{},
		&RBACPolicy{},
		&RBACPolicyList{},
	)
	return nil
}

var (
	NetworkKind             = reflect.TypeOf(Network{}).Name()
	NetworkGroupKind        = schema.GroupKind{Group: APIGroup, Kind: NetworkKind}
	NetworkKindAPIVersion   = NetworkKind + "." + SchemeGroupVersion.String()
	NetworkGroupVersionKind = SchemeGroupVersion.WithKind(NetworkKind)
)

var (
	SubnetKind             = reflect.TypeOf(Subnet{}).Name()
	SubnetGroupKind        = schema.GroupKind{Group: APIGroup, Kind: SubnetKind}
	SubnetKindAPIVersion   = SubnetKind + "." + SchemeGroupVersion.String()
	SubnetGroupVersionKind = SchemeGroupVersion.WithKind(SubnetKind)
)

var (
	RouterKind             = reflect.TypeOf(Router{}).Name()
	RouterGroupKind        = schema.GroupKind{Group: APIGroup, Kind: RouterKind}
	RouterKindAPIVersion   = RouterKind + "." + SchemeGroupVersion.String()
	RouterGroupVersionKind = SchemeGroupVersion.WithKind(RouterKind)
)

var (
	RouterInterfaceKind             = reflect.TypeOf(RouterInterface{}).Name()
	RouterInterfaceGroupKind        = schema.GroupKind{Group: APIGroup, Kind: RouterInterfaceKind}
	RouterInterfaceKindAPIVersion   = RouterInterfaceKind + "." + SchemeGroupVersion.String()
	RouterInterfaceGroupVersionKind = SchemeGroupVersion.WithKind(RouterInterfaceKind)
)

var (
	PortKind             = reflect.TypeOf(Port{}).Name()
	PortGroupKind        = schema.GroupKind{Group: APIGroup, Kind: PortKind}
	PortKindAPIVersion   = PortKind + "." + SchemeGroupVersion.String()
	PortGroupVersionKind = SchemeGroupVersion.WithKind(PortKind)
)

var (
	SecurityGroupKind             = reflect.TypeOf(SecurityGroup{}).Name()
	SecurityGroupGroupKind        = schema.GroupKind{Group: APIGroup, Kind: SecurityGroupKind}
	SecurityGroupKindAPIVersion   = SecurityGroupKind + "." + SchemeGroupVersion.String()
	SecurityGroupGroupVersionKind = SchemeGroupVersion.WithKind(SecurityGroupKind)
)

var (
	SecurityGroupRuleKind             = reflect.TypeOf(SecurityGroupRule{}).Name()
	SecurityGroupRuleGroupKind        = schema.GroupKind{Group: APIGroup, Kind: SecurityGroupRuleKind}
	SecurityGroupRuleKindAPIVersion   = SecurityGroupRuleKind + "." + SchemeGroupVersion.String()
	SecurityGroupRuleGroupVersionKind = SchemeGroupVersion.WithKind(SecurityGroupRuleKind)
)

var (
	FloatingIPKind             = reflect.TypeOf(FloatingIP{}).Name()
	FloatingIPGroupKind        = schema.GroupKind{Group: APIGroup, Kind: FloatingIPKind}
	FloatingIPKindAPIVersion   = FloatingIPKind + "." + SchemeGroupVersion.String()
	FloatingIPGroupVersionKind = SchemeGroupVersion.WithKind(FloatingIPKind)
)

var (
	SubnetPoolKind             = reflect.TypeOf(SubnetPool{}).Name()
	SubnetPoolGroupKind        = schema.GroupKind{Group: APIGroup, Kind: SubnetPoolKind}
	SubnetPoolKindAPIVersion   = SubnetPoolKind + "." + SchemeGroupVersion.String()
	SubnetPoolGroupVersionKind = SchemeGroupVersion.WithKind(SubnetPoolKind)
)

var (
	TrunkKind             = reflect.TypeOf(Trunk{}).Name()
	TrunkGroupKind        = schema.GroupKind{Group: APIGroup, Kind: TrunkKind}
	TrunkKindAPIVersion   = TrunkKind + "." + SchemeGroupVersion.String()
	TrunkGroupVersionKind = SchemeGroupVersion.WithKind(TrunkKind)
)

var (
	RBACPolicyKind             = reflect.TypeOf(RBACPolicy{}).Name()
	RBACPolicyGroupKind        = schema.GroupKind{Group: APIGroup, Kind: RBACPolicyKind}
	RBACPolicyKindAPIVersion   = RBACPolicyKind + "." + SchemeGroupVersion.String()
	RBACPolicyGroupVersionKind = SchemeGroupVersion.WithKind(RBACPolicyKind)
)
