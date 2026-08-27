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
		&Volume{},
		&VolumeList{},
		&VolumeType{},
		&VolumeTypeList{},
		&VolumeSnapshot{},
		&VolumeSnapshotList{},
	)
	metav1.AddToGroupVersion(s, SchemeGroupVersion)
	return nil
}

// Volume type metadata.
var (
	VolumeKind             = reflect.TypeOf(Volume{}).Name()
	VolumeGroupKind        = schema.GroupKind{Group: Group, Kind: VolumeKind}.String()
	VolumeKindAPIVersion   = VolumeKind + "." + SchemeGroupVersion.String()
	VolumeGroupVersionKind = SchemeGroupVersion.WithKind(VolumeKind)
)

// VolumeType type metadata.
var (
	VolumeTypeKind             = reflect.TypeOf(VolumeType{}).Name()
	VolumeTypeGroupKind        = schema.GroupKind{Group: Group, Kind: VolumeTypeKind}.String()
	VolumeTypeKindAPIVersion   = VolumeTypeKind + "." + SchemeGroupVersion.String()
	VolumeTypeGroupVersionKind = SchemeGroupVersion.WithKind(VolumeTypeKind)
)

// VolumeSnapshot type metadata.
var (
	VolumeSnapshotKind             = reflect.TypeOf(VolumeSnapshot{}).Name()
	VolumeSnapshotGroupKind        = schema.GroupKind{Group: Group, Kind: VolumeSnapshotKind}.String()
	VolumeSnapshotKindAPIVersion   = VolumeSnapshotKind + "." + SchemeGroupVersion.String()
	VolumeSnapshotGroupVersionKind = SchemeGroupVersion.WithKind(VolumeSnapshotKind)
)
