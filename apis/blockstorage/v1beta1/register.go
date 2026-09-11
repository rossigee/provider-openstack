/*
Copyright 2024 The Crossplane Authors.
Licensed under the Apache License, Version 2.0.
*/

package v1beta1

import (
	"reflect"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	VolumeKind             = reflect.TypeOf(Volume{}).Name()
	VolumeGroupKind        = schema.GroupKind{Group: APIGroup, Kind: VolumeKind}
	VolumeKindAPIVersion   = VolumeKind + "." + SchemeGroupVersion.String()
	VolumeGroupVersionKind = SchemeGroupVersion.WithKind(VolumeKind)
)
