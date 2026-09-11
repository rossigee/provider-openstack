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
	ImageKind             = reflect.TypeOf(Image{}).Name()
	ImageGroupKind        = schema.GroupKind{Group: APIGroup, Kind: ImageKind}
	ImageKindAPIVersion   = ImageKind + "." + SchemeGroupVersion.String()
	ImageGroupVersionKind = SchemeGroupVersion.WithKind(ImageKind)
)
