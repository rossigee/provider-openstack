package volumesnapshot

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/snapshots"
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
	errNotVolumeSnapshotResource = "managed resource is not a VolumeSnapshot"
	errNewClient                 = "cannot create OpenStack client"

	errGetVolumeSnapshot    = "cannot get VolumeSnapshot"
	errCreateVolumeSnapshot = "cannot create VolumeSnapshot"
	errUpdateVolumeSnapshot = "cannot update VolumeSnapshot"
	errDeleteVolumeSnapshot = "cannot delete VolumeSnapshot"
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
	return &openstackVolumeSnapshotClient{
		client:   blockStorageClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackVolumeSnapshotClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.VolumeSnapshotGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.VolumeSnapshotGroupKind)),
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
		For(&v1alpha1.VolumeSnapshot{}).
		Complete(r)
}

func (e *openstackVolumeSnapshotClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.VolumeSnapshot)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotVolumeSnapshotResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	snap, err := snapshots.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetVolumeSnapshot, err)
	}

	cr.Status.AtProvider = v1alpha1.VolumeSnapshotProviderStatus{
		SnapshotID:  snap.ID,
		Name:        snap.Name,
		Description: snap.Description,
		VolumeID:    snap.VolumeID,
		Status:      snap.Status,
		Size:        snap.Size,
		Progress:    snap.Progress,
		TenantID:    snap.ProjectID,
		Metadata:    snap.Metadata,
	}

	if !strings.EqualFold(cr.Spec.ForProvider.Name, snap.Name) {
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

func (e *openstackVolumeSnapshotClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.VolumeSnapshot)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotVolumeSnapshotResource)
	}

	createOpts := snapshots.CreateOpts{
		VolumeID:    cr.Spec.ForProvider.VolumeID,
		Force:       cr.Spec.ForProvider.Force,
		Name:        cr.Spec.ForProvider.Name,
		Description: cr.Spec.ForProvider.Description,
		Metadata:    cr.Spec.ForProvider.Metadata,
	}

	snap, err := snapshots.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateVolumeSnapshot, err)
	}

	meta.SetExternalName(cr, snap.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackVolumeSnapshotClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.VolumeSnapshot)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotVolumeSnapshotResource)
	}

	name := cr.Spec.ForProvider.Name
	desc := cr.Spec.ForProvider.Description
	updateOpts := snapshots.UpdateOpts{
		Name:        &name,
		Description: &desc,
	}

	_, err := snapshots.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateVolumeSnapshot, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackVolumeSnapshotClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.VolumeSnapshot)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotVolumeSnapshotResource)
	}

	result := snapshots.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackVolumeSnapshotClient) Disconnect(ctx context.Context) error {
	return nil
}
