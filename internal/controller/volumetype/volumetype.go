package volumetype

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
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
	errNotVolumeTypeResource = "managed resource is not a VolumeType"
	errNewClient             = "cannot create OpenStack client"

	errGetVolumeType    = "cannot get VolumeType"
	errCreateVolumeType = "cannot create VolumeType"
	errUpdateVolumeType = "cannot update VolumeType"
	errDeleteVolumeType = "cannot delete VolumeType"
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
	return &openstackVolumeTypeClient{
		client:   blockStorageClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackVolumeTypeClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.VolumeTypeGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.VolumeTypeGroupKind)),
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
		For(&v1alpha1.VolumeType{}).
		Complete(r)
}

func (e *openstackVolumeTypeClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.VolumeType)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotVolumeTypeResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	vt, err := volumetypes.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetVolumeType, err)
	}

	cr.Status.AtProvider = v1alpha1.VolumeTypeProviderStatus{
		VolumeTypeID: vt.ID,
		Name:         vt.Name,
		Description:  vt.Description,
		IsPublic:     vt.IsPublic,
		ExtraSpecs:   vt.ExtraSpecs,
		QosSpecID:    vt.QosSpecID,
	}

	if !strings.EqualFold(cr.Spec.ForProvider.Name, vt.Name) {
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

func (e *openstackVolumeTypeClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.VolumeType)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotVolumeTypeResource)
	}

	createOpts := volumetypes.CreateOpts{
		Name:        cr.Spec.ForProvider.Name,
		Description: cr.Spec.ForProvider.Description,
		IsPublic:    cr.Spec.ForProvider.IsPublic,
		ExtraSpecs:  cr.Spec.ForProvider.ExtraSpecs,
	}

	vt, err := volumetypes.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateVolumeType, err)
	}

	meta.SetExternalName(cr, vt.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackVolumeTypeClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.VolumeType)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotVolumeTypeResource)
	}

	name := cr.Spec.ForProvider.Name
	desc := cr.Spec.ForProvider.Description
	updateOpts := volumetypes.UpdateOpts{
		Name:        &name,
		Description: &desc,
		IsPublic:    cr.Spec.ForProvider.IsPublic,
	}

	_, err := volumetypes.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateVolumeType, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackVolumeTypeClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.VolumeType)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotVolumeTypeResource)
	}

	result := volumetypes.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackVolumeTypeClient) Disconnect(ctx context.Context) error {
	return nil
}
