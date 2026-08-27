// Package controller contains the provider controllers.
package controller

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/rossigee/provider-openstack/internal/controller/computev1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/networkingv1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/providerconfig"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Setup registers all controllers with the given manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	if err := providerconfig.Setup(mgr, o); err != nil {
		return err
	}
	if err := networkingv1alpha1.Setup(mgr, o); err != nil {
		return err
	}
	return computev1alpha1.Setup(mgr, o)
}
