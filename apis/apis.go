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

// Package apis contains Kubernetes API groups for the OpenStack provider.
package apis

import (
	"k8s.io/apimachinery/pkg/runtime"

	blockstoragev1alpha1 "github.com/rossigee/provider-openstack/apis/blockstorage/v1alpha1"
	computev1alpha1 "github.com/rossigee/provider-openstack/apis/compute/v1alpha1"
	dnsv1alpha1 "github.com/rossigee/provider-openstack/apis/dns/v1alpha1"
	identityv1alpha1 "github.com/rossigee/provider-openstack/apis/identity/v1alpha1"
	imagev1alpha1 "github.com/rossigee/provider-openstack/apis/image/v1alpha1"
	loadbalancingv1alpha1 "github.com/rossigee/provider-openstack/apis/loadbalancing/v1alpha1"
	networkingv1alpha1 "github.com/rossigee/provider-openstack/apis/networking/v1alpha1"
	v1beta1 "github.com/rossigee/provider-openstack/apis/v1beta1"
)

// AddToSchemes may be used to register all exported APIs with a Scheme.
var AddToSchemes runtime.SchemeBuilder

func init() {
	AddToSchemes = append(AddToSchemes,
		v1beta1.AddToScheme,
		networkingv1alpha1.AddToScheme,
		blockstoragev1alpha1.AddToScheme,
		loadbalancingv1alpha1.AddToScheme,
		computev1alpha1.AddToScheme,
		imagev1alpha1.AddToScheme,
		identityv1alpha1.AddToScheme,
		dnsv1alpha1.AddToScheme,
	)
}

// AddToScheme registers all exported APIs with a Scheme.
var AddToScheme = AddToSchemes.AddToScheme
