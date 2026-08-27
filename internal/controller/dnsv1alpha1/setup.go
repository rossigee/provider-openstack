package dnsv1alpha1

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/rossigee/provider-openstack/internal/controller/recordset"
	"github.com/rossigee/provider-openstack/internal/controller/zone"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Setup registers all DNS controllers with the given manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		zone.Setup,
		recordset.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
