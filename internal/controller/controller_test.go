package controller

import (
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/rossigee/provider-openstack/internal/controller/blockstoragev1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/computev1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/dnsv1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/identityv1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/imagev1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/loadbalancingv1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/networkingv1alpha1"
	"github.com/rossigee/provider-openstack/internal/controller/providerconfig"
)

// Verify all Setup functions satisfy the expected signature at compile time.
var (
	_ func(ctrl.Manager, controller.Options) error = providerconfig.Setup
	_ func(ctrl.Manager, controller.Options) error = networkingv1alpha1.Setup
	_ func(ctrl.Manager, controller.Options) error = blockstoragev1alpha1.Setup
	_ func(ctrl.Manager, controller.Options) error = computev1alpha1.Setup
	_ func(ctrl.Manager, controller.Options) error = imagev1alpha1.Setup
	_ func(ctrl.Manager, controller.Options) error = identityv1alpha1.Setup
	_ func(ctrl.Manager, controller.Options) error = loadbalancingv1alpha1.Setup
	_ func(ctrl.Manager, controller.Options) error = dnsv1alpha1.Setup
	_ func(ctrl.Manager, controller.Options) error = Setup
)

// TestSetupSignature is a no-op that confirms the compile-time checks above are
// evaluated (Go only evaluates package-level vars when the package is compiled).
func TestSetupSignature(t *testing.T) {
	t.Log("all Setup function signatures verified at compile time")
}
