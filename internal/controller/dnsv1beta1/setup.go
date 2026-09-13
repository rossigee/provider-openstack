package dnsv1beta1

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/rossigee/provider-openstack/internal/controller/zone"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Setup registers all DNS controllers with the given manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		zone.Setup,
		// TODO: RecordSet type needs to be defined in v1beta1/types.go
		// recordset.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
