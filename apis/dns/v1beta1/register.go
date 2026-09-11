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
	ZoneKind             = reflect.TypeOf(Zone{}).Name()
	ZoneGroupKind        = schema.GroupKind{Group: APIGroup, Kind: ZoneKind}
	ZoneKindAPIVersion   = ZoneKind + "." + SchemeGroupVersion.String()
	ZoneGroupVersionKind = SchemeGroupVersion.WithKind(ZoneKind)
)
