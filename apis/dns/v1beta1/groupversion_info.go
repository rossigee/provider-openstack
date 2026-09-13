package v1beta1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	APIGroup = "dns.openstack.m.crossplane.io"
)

var SchemeGroupVersion = schema.GroupVersion{Group: APIGroup, Version: "v1beta1"}

var SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
var AddToScheme = SchemeBuilder.AddToScheme

func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(SchemeGroupVersion, &Zone{}, &ZoneList{})
	return nil
}

var (
	ZoneKind             = reflect.TypeOf(Zone{}).Name()
	ZoneGroupKind        = schema.GroupKind{Group: APIGroup, Kind: ZoneKind}
	ZoneKindAPIVersion   = ZoneKind + "." + SchemeGroupVersion.String()
	ZoneGroupVersionKind = SchemeGroupVersion.WithKind(ZoneKind)
)

// ZoneProviderStatus defines the observed state of the zone at the provider
type ZoneProviderStatus struct {
	ZoneID     string            `json:"zoneId,omitempty"`
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Email      string            `json:"email,omitempty"`
	TTL        int               `json:"ttl,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Status     string            `json:"status,omitempty"`
}
