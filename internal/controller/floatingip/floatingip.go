package floatingip

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/rossigee/provider-openstack/apis/networking/v1alpha1"
	"github.com/rossigee/provider-openstack/internal/clients"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/event"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
)

const (
	errNotFloatingIPResource = "managed resource is not a FloatingIP"
	errTrackUsage            = "cannot track ProviderConfig usage"
	errGetProviderConfig     = "cannot get ProviderConfig"
	errGetCredentials        = "cannot get credentials"
	errNewClient             = "cannot create OpenStack client"

	errGetFloatingIP    = "cannot get FloatingIP"
	errCreateFloatingIP = "cannot create FloatingIP"
	errUpdateFloatingIP = "cannot update FloatingIP"
	errDeleteFloatingIP = "cannot delete FloatingIP"
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
	return &openstackFloatingIPClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackFloatingIPClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.FloatingIPGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.FloatingIPGroupKind)),
		managed.WithExternalConnector(&External{
			kube:     mgr.GetClient(),
			recorder: rec,
		}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(rec),
		managed.WithPollInterval(o.PollInterval),
	)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.FloatingIP{}).
		Complete(r)
}

func (e *openstackFloatingIPClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.FloatingIP)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotFloatingIPResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	fip, err := floatingips.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetFloatingIP, err)
	}

	cr.Status.AtProvider = v1alpha1.FloatingIPProviderStatus{
		FloatingIPID: fip.ID,
		FloatingIP:   fip.FloatingIP,
		FixedIP:      fip.FixedIP,
		PortID:       fip.PortID,
		Status:       fip.Status,
		TenantID:     fip.TenantID,
		RouterID:     fip.RouterID,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackFloatingIPClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.FloatingIP)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotFloatingIPResource)
	}

	createOpts := floatingips.CreateOpts{
		FloatingNetworkID: cr.Spec.ForProvider.FloatingNetworkID,
		FloatingIP:        cr.Spec.ForProvider.FloatingIP,
		PortID:            cr.Spec.ForProvider.PortID,
		TenantID:          cr.Spec.ForProvider.TenantID,
		Description:       cr.Spec.ForProvider.Description,
	}

	fip, err := floatingips.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateFloatingIP, err)
	}

	meta.SetExternalName(cr, fip.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackFloatingIPClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.FloatingIP)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotFloatingIPResource)
	}

	portID := cr.Spec.ForProvider.PortID
	updateOpts := floatingips.UpdateOpts{
		PortID:      &portID,
		Description: &cr.Spec.ForProvider.Description,
	}

	_, err := floatingips.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateFloatingIP, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackFloatingIPClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.FloatingIP)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotFloatingIPResource)
	}

	result := floatingips.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackFloatingIPClient) Disconnect(ctx context.Context) error {
	return nil
}
