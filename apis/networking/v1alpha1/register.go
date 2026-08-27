/*
Copyright 2025 The Crossplane Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// +kubebuilder:object:generate=true
// +groupName=openstack.crossplane.io
// +versionName=v1alpha1

package v1alpha1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Group   = "openstack.crossplane.io"
	Version = "v1alpha1"
)

var (
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
)

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
		&SecurityGroup{},
		&SecurityGroupList{},
		&SecurityGroupRule{},
		&SecurityGroupRuleList{},
		&FloatingIP{},
		&FloatingIPList{},
	)
	metav1.AddToGroupVersion(s, SchemeGroupVersion)
	return nil
}

// Network type metadata.
var (
	NetworkKind             = reflect.TypeOf(Network{}).Name()
	NetworkGroupKind        = schema.GroupKind{Group: Group, Kind: NetworkKind}.String()
	NetworkKindAPIVersion   = NetworkKind + "." + SchemeGroupVersion.String()
	NetworkGroupVersionKind = SchemeGroupVersion.WithKind(NetworkKind)
)

// Subnet type metadata.
var (
	SubnetKind             = reflect.TypeOf(Subnet{}).Name()
	SubnetGroupKind        = schema.GroupKind{Group: Group, Kind: SubnetKind}.String()
	SubnetKindAPIVersion   = SubnetKind + "." + SchemeGroupVersion.String()
	SubnetGroupVersionKind = SchemeGroupVersion.WithKind(SubnetKind)
)

// Router type metadata.
var (
	RouterKind             = reflect.TypeOf(Router{}).Name()
	RouterGroupKind        = schema.GroupKind{Group: Group, Kind: RouterKind}.String()
	RouterKindAPIVersion   = RouterKind + "." + SchemeGroupVersion.String()
	RouterGroupVersionKind = SchemeGroupVersion.WithKind(RouterKind)
)

// RouterInterface type metadata.
var (
	RouterInterfaceKind             = reflect.TypeOf(RouterInterface{}).Name()
	RouterInterfaceGroupKind        = schema.GroupKind{Group: Group, Kind: RouterInterfaceKind}.String()
	RouterInterfaceKindAPIVersion   = RouterInterfaceKind + "." + SchemeGroupVersion.String()
	RouterInterfaceGroupVersionKind = SchemeGroupVersion.WithKind(RouterInterfaceKind)
)

// SecurityGroup type metadata.
var (
	SecurityGroupKind             = reflect.TypeOf(SecurityGroup{}).Name()
	SecurityGroupGroupKind        = schema.GroupKind{Group: Group, Kind: SecurityGroupKind}.String()
	SecurityGroupKindAPIVersion   = SecurityGroupKind + "." + SchemeGroupVersion.String()
	SecurityGroupGroupVersionKind = SchemeGroupVersion.WithKind(SecurityGroupKind)
)

// SecurityGroupRule type metadata.
var (
	SecurityGroupRuleKind             = reflect.TypeOf(SecurityGroupRule{}).Name()
	SecurityGroupRuleGroupKind        = schema.GroupKind{Group: Group, Kind: SecurityGroupRuleKind}.String()
	SecurityGroupRuleKindAPIVersion   = SecurityGroupRuleKind + "." + SchemeGroupVersion.String()
	SecurityGroupRuleGroupVersionKind = SchemeGroupVersion.WithKind(SecurityGroupRuleKind)
)

// FloatingIP type metadata.
var (
	FloatingIPKind             = reflect.TypeOf(FloatingIP{}).Name()
	FloatingIPGroupKind        = schema.GroupKind{Group: Group, Kind: FloatingIPKind}.String()
	FloatingIPKindAPIVersion   = FloatingIPKind + "." + SchemeGroupVersion.String()
	FloatingIPGroupVersionKind = SchemeGroupVersion.WithKind(FloatingIPKind)
)
