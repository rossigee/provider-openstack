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
	metav1.AddToGroupVersion(s, SchemeGroupVersion)
	return nil
}

// LoadBalancer type metadata.
var (
	LoadBalancerKind             = reflect.TypeOf(LoadBalancer{}).Name()
	LoadBalancerGroupKind        = schema.GroupKind{Group: Group, Kind: LoadBalancerKind}.String()
	LoadBalancerKindAPIVersion   = LoadBalancerKind + "." + SchemeGroupVersion.String()
	LoadBalancerGroupVersionKind = SchemeGroupVersion.WithKind(LoadBalancerKind)
)

// Listener type metadata.
var (
	ListenerKind             = reflect.TypeOf(Listener{}).Name()
	ListenerGroupKind        = schema.GroupKind{Group: Group, Kind: ListenerKind}.String()
	ListenerKindAPIVersion   = ListenerKind + "." + SchemeGroupVersion.String()
	ListenerGroupVersionKind = SchemeGroupVersion.WithKind(ListenerKind)
)

// Pool type metadata.
var (
	PoolKind             = reflect.TypeOf(Pool{}).Name()
	PoolGroupKind        = schema.GroupKind{Group: Group, Kind: PoolKind}.String()
	PoolKindAPIVersion   = PoolKind + "." + SchemeGroupVersion.String()
	PoolGroupVersionKind = SchemeGroupVersion.WithKind(PoolKind)
)

// Member type metadata.
var (
	MemberKind             = reflect.TypeOf(Member{}).Name()
	MemberGroupKind        = schema.GroupKind{Group: Group, Kind: MemberKind}.String()
	MemberKindAPIVersion   = MemberKind + "." + SchemeGroupVersion.String()
	MemberGroupVersionKind = SchemeGroupVersion.WithKind(MemberKind)
)

// HealthMonitor type metadata.
var (
	HealthMonitorKind             = reflect.TypeOf(HealthMonitor{}).Name()
	HealthMonitorGroupKind        = schema.GroupKind{Group: Group, Kind: HealthMonitorKind}.String()
	HealthMonitorKindAPIVersion   = HealthMonitorKind + "." + SchemeGroupVersion.String()
	HealthMonitorGroupVersionKind = SchemeGroupVersion.WithKind(HealthMonitorKind)
)
