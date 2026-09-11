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
	NetworkKind             = reflect.TypeOf(Network{}).Name()
	NetworkGroupKind        = schema.GroupKind{Group: APIGroup, Kind: NetworkKind}
	NetworkKindAPIVersion   = NetworkKind + "." + SchemeGroupVersion.String()
	NetworkGroupVersionKind = SchemeGroupVersion.WithKind(NetworkKind)
)
