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
	APIGroup = "blockstorage.openstack.m.crossplane.io"
)

var SchemeGroupVersion = schema.GroupVersion{Group: APIGroup, Version: "v1beta1"}

var SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
var AddToScheme = SchemeBuilder.AddToScheme

func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(SchemeGroupVersion,
		&Volume{},
		&VolumeList{},
		&VolumeSnapshot{},
		&VolumeSnapshotList{},
		&VolumeType{},
		&VolumeTypeList{},
	)
	return nil
}

// Volume type metadata.
var (
	VolumeKind             = reflect.TypeOf(Volume{}).Name()
	VolumeGroupKind        = schema.GroupKind{Group: APIGroup, Kind: VolumeKind}
	VolumeKindAPIVersion   = VolumeKind + "." + SchemeGroupVersion.String()
	VolumeGroupVersionKind = SchemeGroupVersion.WithKind(VolumeKind)
)

// VolumeSnapshot type metadata.
var (
	VolumeSnapshotKind             = reflect.TypeOf(VolumeSnapshot{}).Name()
	VolumeSnapshotGroupKind        = schema.GroupKind{Group: APIGroup, Kind: VolumeSnapshotKind}
	VolumeSnapshotKindAPIVersion   = VolumeSnapshotKind + "." + SchemeGroupVersion.String()
	VolumeSnapshotGroupVersionKind = SchemeGroupVersion.WithKind(VolumeSnapshotKind)
)

// VolumeType type metadata.
var (
	VolumeTypeKind             = reflect.TypeOf(VolumeType{}).Name()
	VolumeTypeGroupKind        = schema.GroupKind{Group: APIGroup, Kind: VolumeTypeKind}
	VolumeTypeKindAPIVersion   = VolumeTypeKind + "." + SchemeGroupVersion.String()
	VolumeTypeGroupVersionKind = SchemeGroupVersion.WithKind(VolumeTypeKind)
)
