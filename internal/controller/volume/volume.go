package volume

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/rossigee/provider-openstack/apis/blockstorage/v1alpha1"
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
	errNotVolumeResource = "managed resource is not a Volume"
	errNewClient         = "cannot create OpenStack client"

	errGetVolume    = "cannot get Volume"
	errCreateVolume = "cannot create Volume"
	errUpdateVolume = "cannot update Volume"
	errDeleteVolume = "cannot delete Volume"
)

type External struct {
	kube     client.Client
	recorder event.Recorder
}

func (e *External) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	blockStorageClient, err := clients.NewBlockStorageClient(ctx, e.kube, mg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errNewClient, err)
	}
	return &openstackVolumeClient{
		client:   blockStorageClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackVolumeClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.VolumeGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.VolumeGroupKind)),
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
		For(&v1alpha1.Volume{}).
		Complete(r)
}

func (e *openstackVolumeClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Volume)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotVolumeResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	vol, err := volumes.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetVolume, err)
	}

	cr.Status.AtProvider = v1alpha1.VolumeProviderStatus{
		VolumeID:         vol.ID,
		Status:           vol.Status,
		Size:             vol.Size,
		AvailabilityZone: vol.AvailabilityZone,
		VolumeType:       vol.VolumeType,
		Description:      vol.Description,
		Bootable:         vol.Bootable,
		Encrypted:        vol.Encrypted,
		SnapshotID:       vol.SnapshotID,
		SourceVolID:      vol.SourceVolID,
		TenantID:         vol.TenantID,
		Host:             vol.Host,
		Metadata:         vol.Metadata,
	}

	if !strings.EqualFold(cr.Spec.ForProvider.Name, vol.Name) {
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

func (e *openstackVolumeClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Volume)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotVolumeResource)
	}

	createOpts := volumes.CreateOpts{
		Size:             cr.Spec.ForProvider.Size,
		Name:             cr.Spec.ForProvider.Name,
		Description:      cr.Spec.ForProvider.Description,
		AvailabilityZone: cr.Spec.ForProvider.AvailabilityZone,
		VolumeType:       cr.Spec.ForProvider.VolumeType,
		ImageID:          cr.Spec.ForProvider.ImageID,
		SnapshotID:       cr.Spec.ForProvider.SnapshotID,
		SourceVolID:      cr.Spec.ForProvider.SourceVolID,
		Metadata:         cr.Spec.ForProvider.Metadata,
	}

	vol, err := volumes.Create(ctx, e.client, createOpts, nil).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateVolume, err)
	}

	meta.SetExternalName(cr, vol.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackVolumeClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Volume)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotVolumeResource)
	}

	name := cr.Spec.ForProvider.Name
	desc := cr.Spec.ForProvider.Description
	updateOpts := volumes.UpdateOpts{
		Name:        &name,
		Description: &desc,
		Metadata:    cr.Spec.ForProvider.Metadata,
	}

	_, err := volumes.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateVolume, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackVolumeClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Volume)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotVolumeResource)
	}

	result := volumes.Delete(ctx, e.client, meta.GetExternalName(cr), nil)
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackVolumeClient) Disconnect(ctx context.Context) error {
	return nil
}
