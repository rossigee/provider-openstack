package zone

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/rossigee/provider-openstack/apis/dns/v1alpha1"
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
	errNotZoneResource   = "managed resource is not a Zone"
	errGetProviderConfig = "cannot get ProviderConfig"
	errGetCredentials    = "cannot get credentials"
	errNewClient         = "cannot create OpenStack client"
	errGetZone           = "cannot get Zone"
	errCreateZone        = "cannot create Zone"
	errUpdateZone        = "cannot update Zone"
	errDeleteZone        = "cannot delete Zone"
)

type External struct {
	kube     client.Client
	recorder event.Recorder
}

func (e *External) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	dnsClient, err := clients.NewDNSClient(ctx, e.kube, mg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errNewClient, err)
	}
	return &openstackZoneClient{
		client:   dnsClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackZoneClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.ZoneGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ZoneGroupKind)),
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
		For(&v1alpha1.Zone{}).
		Complete(r)
}

func (e *openstackZoneClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Zone)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotZoneResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	zone, err := zones.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetZone, err)
	}

	cr.Status.AtProvider = v1alpha1.ZoneProviderStatus{
		ZoneID:     zone.ID,
		Name:       zone.Name,
		Type:       zone.Type,
		Email:      zone.Email,
		TTL:        zone.TTL,
		Attributes: zone.Attributes,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackZoneClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Zone)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotZoneResource)
	}

	var ttl int
	if cr.Spec.ForProvider.TTL != nil {
		ttl = *cr.Spec.ForProvider.TTL
	}

	createOpts := zones.CreateOpts{
		Name:        cr.Spec.ForProvider.Name,
		Email:       cr.Spec.ForProvider.Email,
		Description: cr.Spec.ForProvider.Description,
		Type:        cr.Spec.ForProvider.Type,
		TTL:         ttl,
		Attributes:  cr.Spec.ForProvider.Attributes,
	}

	zone, err := zones.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateZone, err)
	}

	meta.SetExternalName(cr, zone.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackZoneClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Zone)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotZoneResource)
	}

	var ttl int
	if cr.Spec.ForProvider.TTL != nil {
		ttl = *cr.Spec.ForProvider.TTL
	}

	desc := cr.Spec.ForProvider.Description
	updateOpts := zones.UpdateOpts{
		Email:       cr.Spec.ForProvider.Email,
		Description: &desc,
		TTL:         ttl,
	}

	_, err := zones.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateZone, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackZoneClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Zone)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotZoneResource)
	}

	result := zones.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackZoneClient) Disconnect(ctx context.Context) error {
	return nil
}
