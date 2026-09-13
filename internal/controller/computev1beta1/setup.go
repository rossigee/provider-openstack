package computev1alpha1

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/rossigee/provider-openstack/internal/controller/keypair"
	"github.com/rossigee/provider-openstack/internal/controller/server"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Setup registers all compute controllers with the given manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		server.Setup,
		keypair.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
