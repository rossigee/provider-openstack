package networkingv1alpha1

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/rossigee/provider-openstack/internal/controller/floatingip"
	"github.com/rossigee/provider-openstack/internal/controller/network"
	"github.com/rossigee/provider-openstack/internal/controller/router"
	"github.com/rossigee/provider-openstack/internal/controller/routerinterface"
	"github.com/rossigee/provider-openstack/internal/controller/securitygroup"
	"github.com/rossigee/provider-openstack/internal/controller/securitygrouprule"
	"github.com/rossigee/provider-openstack/internal/controller/subnet"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Setup registers all networking controllers with the given manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		network.Setup,
		subnet.Setup,
		router.Setup,
		routerinterface.Setup,
		securitygroup.Setup,
		securitygrouprule.Setup,
		floatingip.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
