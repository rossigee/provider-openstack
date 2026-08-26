/*
Copyright 2021 Upbound Inc.
*/

package providerconfig

import (
	v1event "github.com/crossplane/crossplane-runtime/pkg/event"
	v1logging "github.com/crossplane/crossplane-runtime/pkg/logging"
	v1providerconfig "github.com/crossplane/crossplane-runtime/pkg/reconciler/providerconfig"
	v1resource "github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/rossigee/provider-openstack/apis/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func Setup(mgr ctrl.Manager) error {
	name := v1providerconfig.ControllerName(v1beta1.ProviderConfigGroupKind.Kind)

	of := v1resource.ProviderConfigKinds{
		Config:    v1beta1.ProviderConfigGroupVersionKind,
		UsageList: v1beta1.ProviderConfigUsageListGroupVersionKind,
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1beta1.ProviderConfig{}).
		Watches(&v1beta1.ProviderConfigUsage{}, &v1resource.EnqueueRequestForProviderConfig{}).
		Complete(v1providerconfig.NewReconciler(mgr, of,
			v1providerconfig.WithLogger(v1logging.NewNopLogger()),
			v1providerconfig.WithRecorder(v1event.NewNopRecorder())))
}
