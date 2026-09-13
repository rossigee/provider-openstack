// Package controller contains the provider controllers.
package controller

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/rossigee/provider-openstack/internal/controller/blockstoragev1beta1"
	"github.com/rossigee/provider-openstack/internal/controller/computev1beta1"
	"github.com/rossigee/provider-openstack/internal/controller/dnsv1beta1"
	"github.com/rossigee/provider-openstack/internal/controller/identityv1beta1"
	"github.com/rossigee/provider-openstack/internal/controller/imagev1beta1"
	"github.com/rossigee/provider-openstack/internal/controller/loadbalancingv1beta1"
	"github.com/rossigee/provider-openstack/internal/controller/networkingv1beta1"
	"github.com/rossigee/provider-openstack/internal/controller/providerconfig"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Setup registers all controllers with the given manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	if err := providerconfig.Setup(mgr, o); err != nil {
		return err
	}
	if err := networkingv1beta1.Setup(mgr, o); err != nil {
		return err
	}
	if err := blockstoragev1beta1.Setup(mgr, o); err != nil {
		return err
	}
	if err := computev1beta1.Setup(mgr, o); err != nil {
		return err
	}
	if err := imagev1beta1.Setup(mgr, o); err != nil {
		return err
	}
	if err := identityv1beta1.Setup(mgr, o); err != nil {
		return err
	}
	if err := loadbalancingv1beta1.Setup(mgr, o); err != nil {
		return err
	}
	return dnsv1beta1.Setup(mgr, o)
}
