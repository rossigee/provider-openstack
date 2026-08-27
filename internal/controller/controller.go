// Package controller contains the provider controllers.
package controller

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/rossigee/provider-openstack/internal/controller/blockstoragev1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/computev1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/dnsv1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/identityv1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/imagev1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/loadbalancingv1alpha1"
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
	if err := blockstoragev1alpha1.Setup(mgr, o); err != nil {
		return err
	}
	if err := computev1alpha1.Setup(mgr, o); err != nil {
		return err
	}
	if err := imagev1alpha1.Setup(mgr, o); err != nil {
		return err
	}
	if err := identityv1alpha1.Setup(mgr, o); err != nil {
		return err
	}
	if err := loadbalancingv1alpha1.Setup(mgr, o); err != nil {
		return err
	}
	return dnsv1alpha1.Setup(mgr, o)
}
