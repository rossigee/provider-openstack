package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/rossigee/provider-openstack/apis/image/v1alpha1"
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
	errNotImageResource  = "managed resource is not an Image"
	errTrackUsage        = "cannot track ProviderConfig usage"
	errGetProviderConfig = "cannot get ProviderConfig"
	errGetCredentials    = "cannot get credentials"
	errNewClient         = "cannot create OpenStack client"

	errGetImage    = "cannot get Image"
	errCreateImage = "cannot create Image"
	errUpdateImage = "cannot update Image"
	errDeleteImage = "cannot delete Image"
)

type External struct {
	kube     client.Client
	recorder event.Recorder
}

func (e *External) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	imageClient, err := clients.NewImageClient(ctx, e.kube, mg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errNewClient, err)
	}
	return &openstackImageClient{
		client:   imageClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackImageClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.ImageGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ImageGroupKind)),
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
		For(&v1alpha1.Image{}).
		Complete(r)
}

func (e *openstackImageClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Image)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotImageResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	image, err := images.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetImage, err)
	}

	protected := image.Protected
	cr.Status.AtProvider = v1alpha1.ImageProviderStatus{
		ImageID:         image.ID,
		Status:          string(image.Status),
		Size:            image.SizeBytes,
		MinDisk:         image.MinDiskGigabytes,
		MinRAM:          image.MinRAMMegabytes,
		Protected:       &protected,
		Visibility:      string(image.Visibility),
		ContainerFormat: image.ContainerFormat,
		DiskFormat:      image.DiskFormat,
		Tags:            image.Tags,
		Owner:           image.Owner,
		File:            image.File,
		Schema:          image.Schema,
		Checksum:        image.Checksum,
		VirtualSize:     image.VirtualSize,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackImageClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Image)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotImageResource)
	}

	minDisk := 0
	if cr.Spec.ForProvider.MinDisk != nil {
		minDisk = *cr.Spec.ForProvider.MinDisk
	}
	minRAM := 0
	if cr.Spec.ForProvider.MinRAM != nil {
		minRAM = *cr.Spec.ForProvider.MinRAM
	}

	vis := images.ImageVisibility(cr.Spec.ForProvider.Visibility)

	createOpts := images.CreateOpts{
		Name:            cr.Spec.ForProvider.Name,
		ContainerFormat: cr.Spec.ForProvider.ContainerFormat,
		DiskFormat:      cr.Spec.ForProvider.DiskFormat,
		MinDisk:         minDisk,
		MinRAM:          minRAM,
		Visibility:      &vis,
		Tags:            cr.Spec.ForProvider.Tags,
		Protected:       cr.Spec.ForProvider.Protected,
		Properties:      cr.Spec.ForProvider.Properties,
	}

	result, err := images.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateImage, err)
	}

	meta.SetExternalName(cr, result.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackImageClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Image)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotImageResource)
	}

	updateOpts := images.UpdateOpts{
		images.ReplaceImageName{NewName: cr.Spec.ForProvider.Name},
	}

	_, err := images.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateImage, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackImageClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Image)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotImageResource)
	}

	result := images.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackImageClient) Disconnect(ctx context.Context) error {
	return nil
}
