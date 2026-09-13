package blockstoragev1alpha1

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/rossigee/provider-openstack/internal/controller/volume"
	"github.com/rossigee/provider-openstack/internal/controller/volumesnapshot"
	"github.com/rossigee/provider-openstack/internal/controller/volumetype"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Setup registers all block storage controllers with the given manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		volume.Setup,
		volumetype.Setup,
		volumesnapshot.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
