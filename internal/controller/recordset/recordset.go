package recordset

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/recordsets"
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
	errNotRecordSetResource = "managed resource is not a RecordSet"
	errGetProviderConfig    = "cannot get ProviderConfig"
	errGetCredentials       = "cannot get credentials"
	errNewClient            = "cannot create OpenStack client"
	errGetRecordSet         = "cannot get RecordSet"
	errCreateRecordSet      = "cannot create RecordSet"
	errUpdateRecordSet      = "cannot update RecordSet"
	errDeleteRecordSet      = "cannot delete RecordSet"
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
	return &openstackRecordSetClient{
		client:   dnsClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackRecordSetClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.RecordSetGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.RecordSetGroupKind)),
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
		For(&v1alpha1.RecordSet{}).
		Complete(r)
}

func (e *openstackRecordSetClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.RecordSet)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotRecordSetResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	zoneID := cr.Spec.ForProvider.ZoneID
	if zoneID == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	rs, err := recordsets.Get(ctx, e.client, zoneID, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetRecordSet, err)
	}

	cr.Status.AtProvider = v1alpha1.RecordSetProviderStatus{
		RecordSetID: rs.ID,
		Name:        rs.Name,
		ZoneID:      rs.ZoneID,
		ZoneName:    rs.ZoneName,
		Type:        rs.Type,
		TTL:         rs.TTL,
		Records:     rs.Records,
		Status:      rs.Status,
		Action:      rs.Action,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackRecordSetClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.RecordSet)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotRecordSetResource)
	}

	createOpts := recordsets.CreateOpts{
		Name:    cr.Spec.ForProvider.Name,
		Type:    cr.Spec.ForProvider.Type,
		Records: cr.Spec.ForProvider.Records,
		TTL:     cr.Spec.ForProvider.TTL,
	}

	rs, err := recordsets.Create(ctx, e.client, cr.Spec.ForProvider.ZoneID, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateRecordSet, err)
	}

	meta.SetExternalName(cr, rs.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackRecordSetClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.RecordSet)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotRecordSetResource)
	}

	desc := cr.Spec.ForProvider.Name
	updateOpts := recordsets.UpdateOpts{
		Description: &desc,
		TTL:         &cr.Spec.ForProvider.TTL,
		Records:     cr.Spec.ForProvider.Records,
	}

	_, err := recordsets.Update(ctx, e.client, cr.Spec.ForProvider.ZoneID, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateRecordSet, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackRecordSetClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.RecordSet)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotRecordSetResource)
	}

	result := recordsets.Delete(ctx, e.client, cr.Spec.ForProvider.ZoneID, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackRecordSetClient) Disconnect(ctx context.Context) error {
	return nil
}
