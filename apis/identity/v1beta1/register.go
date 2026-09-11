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
	UserKind             = reflect.TypeOf(User{}).Name()
	UserGroupKind        = schema.GroupKind{Group: APIGroup, Kind: UserKind}
	UserKindAPIVersion   = UserKind + "." + SchemeGroupVersion.String()
	UserGroupVersionKind = SchemeGroupVersion.WithKind(UserKind)
)
