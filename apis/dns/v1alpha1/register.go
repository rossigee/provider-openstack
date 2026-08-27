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
		&Zone{},
		&ZoneList{},
		&RecordSet{},
		&RecordSetList{},
	)
	metav1.AddToGroupVersion(s, SchemeGroupVersion)
	return nil
}

// Zone type metadata.
var (
	ZoneKind             = reflect.TypeOf(Zone{}).Name()
	ZoneGroupKind        = schema.GroupKind{Group: Group, Kind: ZoneKind}.String()
	ZoneKindAPIVersion   = ZoneKind + "." + SchemeGroupVersion.String()
	ZoneGroupVersionKind = SchemeGroupVersion.WithKind(ZoneKind)
)

// RecordSet type metadata.
var (
	RecordSetKind             = reflect.TypeOf(RecordSet{}).Name()
	RecordSetGroupKind        = schema.GroupKind{Group: Group, Kind: RecordSetKind}.String()
	RecordSetKindAPIVersion   = RecordSetKind + "." + SchemeGroupVersion.String()
	RecordSetGroupVersionKind = SchemeGroupVersion.WithKind(RecordSetKind)
)
