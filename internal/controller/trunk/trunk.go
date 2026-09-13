package trunk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/trunks"
	networkingv1beta1 "github.com/rossigee/provider-openstack/apis/networking/v1beta1"
	"github.com/rossigee/provider-openstack/internal/clients"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/event"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/rossigee/provider-openstack/internal/features"
)

const (
	errNotTrunkResource = "managed resource is not a Trunk"
	errNewClient        = "cannot create OpenStack client"

	errGetTrunk    = "cannot get Trunk"
	errCreateTrunk = "cannot create Trunk"
	errUpdateTrunk = "cannot update Trunk"
	errDeleteTrunk = "cannot delete Trunk"
)

type External struct {
	kube     client.Client
	recorder event.Recorder
}

func (e *External) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	networkClient, err := clients.NewNetworkClient(ctx, e.kube, mg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errNewClient, err)
	}
	return &openstackTrunkClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackTrunkClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(networkingv1beta1.TrunkGroupKind)
	rec := event.NewNopRecorder()

	opts := []managed.ReconcilerOption{
		managed.WithExternalConnector(&External{
			kube:     mgr.GetClient(),
			recorder: rec,
		}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(rec),
		managed.WithPollInterval(o.PollInterval),
	}
	if o.Features.Enabled(features.EnableAlphaManagementPolicies) {
		opts = append(opts, managed.WithManagementPolicies())
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(networkingv1beta1.SchemeGroupVersion.WithKind(networkingv1beta1.TrunkGroupKind)),
		opts...)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&networkingv1beta1.Trunk{}).
		Complete(r)
}

func (e *openstackTrunkClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*networkingv1beta1.Trunk)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotTrunkResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	trunk, err := trunks.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetTrunk, err)
	}

	cr.Status.AtProvider = networkingv1beta1.TrunkProviderStatus{
		TrunkID:        trunk.ID,
		Name:           trunk.Name,
		Description:    trunk.Description,
		AdminStateUp:   trunk.AdminStateUp,
		Status:         trunk.Status,
		PortID:         trunk.PortID,
		TenantID:       trunk.TenantID,
		RevisionNumber: trunk.RevisionNumber,
		Subports:       toSubports(trunk.Subports),
	}

	if !strings.EqualFold(cr.Spec.ForProvider.Name, trunk.Name) {
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: false,
		}, nil
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackTrunkClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*networkingv1beta1.Trunk)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotTrunkResource)
	}

	createOpts := trunks.CreateOpts{
		PortID:       cr.Spec.ForProvider.PortID,
		Name:         cr.Spec.ForProvider.Name,
		Description:  cr.Spec.ForProvider.Description,
		AdminStateUp: cr.Spec.ForProvider.AdminStateUp,
		TenantID:     cr.Spec.ForProvider.TenantID,
		Subports:     toTrunkSubports(cr.Spec.ForProvider.Subports),
	}

	trunk, err := trunks.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateTrunk, err)
	}

	meta.SetExternalName(cr, trunk.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackTrunkClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*networkingv1beta1.Trunk)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotTrunkResource)
	}

	updateOpts := trunks.UpdateOpts{
		Name:         ptrString(cr.Spec.ForProvider.Name),
		Description:  ptrString(cr.Spec.ForProvider.Description),
		AdminStateUp: cr.Spec.ForProvider.AdminStateUp,
	}

	_, err := trunks.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateTrunk, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackTrunkClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*networkingv1beta1.Trunk)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotTrunkResource)
	}

	result := trunks.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackTrunkClient) Disconnect(ctx context.Context) error {
	return nil
}

func toTrunkSubports(in []networkingv1beta1.Subport) []trunks.Subport {
	if in == nil {
		return nil
	}
	out := make([]trunks.Subport, 0, len(in))
	for _, s := range in {
		out = append(out, trunks.Subport{
			SegmentationID:   s.SegmentationID,
			SegmentationType: s.SegmentationType,
			PortID:           s.PortID,
		})
	}
	return out
}

func toSubports(in []trunks.Subport) []networkingv1beta1.Subport {
	if in == nil {
		return nil
	}
	out := make([]networkingv1beta1.Subport, 0, len(in))
	for _, s := range in {
		out = append(out, networkingv1beta1.Subport{
			SegmentationID:   s.SegmentationID,
			SegmentationType: s.SegmentationType,
			PortID:           s.PortID,
		})
	}
	return out
}

func ptrString(s string) *string {
	return &s
}
