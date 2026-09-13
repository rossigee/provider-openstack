package identityv1alpha1

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/rossigee/provider-openstack/internal/controller/project"
	"github.com/rossigee/provider-openstack/internal/controller/role"
	"github.com/rossigee/provider-openstack/internal/controller/user"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Setup registers all identity controllers with the given manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		project.Setup,
		user.Setup,
		role.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
