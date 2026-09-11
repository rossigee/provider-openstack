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
	ServerKind             = reflect.TypeOf(Server{}).Name()
	ServerGroupKind        = schema.GroupKind{Group: APIGroup, Kind: ServerKind}
	ServerKindAPIVersion   = ServerKind + "." + SchemeGroupVersion.String()
	ServerGroupVersionKind = SchemeGroupVersion.WithKind(ServerKind)
)
