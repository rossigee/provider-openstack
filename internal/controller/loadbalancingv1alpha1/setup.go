package loadbalancingv1alpha1

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/rossigee/provider-openstack/internal/controller/healthmonitor"
	"github.com/rossigee/provider-openstack/internal/controller/listener"
	"github.com/rossigee/provider-openstack/internal/controller/loadbalancer"
	"github.com/rossigee/provider-openstack/internal/controller/member"
	"github.com/rossigee/provider-openstack/internal/controller/pool"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Setup registers all loadbalancing controllers with the given manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		loadbalancer.Setup,
		listener.Setup,
		pool.Setup,
		member.Setup,
		healthmonitor.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
