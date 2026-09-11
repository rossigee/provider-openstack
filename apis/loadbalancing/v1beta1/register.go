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
	LoadBalancerKind             = reflect.TypeOf(LoadBalancer{}).Name()
	LoadBalancerGroupKind        = schema.GroupKind{Group: APIGroup, Kind: LoadBalancerKind}
	LoadBalancerKindAPIVersion   = LoadBalancerKind + "." + SchemeGroupVersion.String()
	LoadBalancerGroupVersionKind = SchemeGroupVersion.WithKind(LoadBalancerKind)
)
