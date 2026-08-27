package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/rossigee/provider-openstack/apis/identity/v1alpha1"
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
	errNotProjectResource = "managed resource is not a Project"
	errTrackUsage         = "cannot track ProviderConfig usage"
	errGetProviderConfig  = "cannot get ProviderConfig"
	errGetCredentials     = "cannot get credentials"
	errNewClient          = "cannot create OpenStack client"

	errGetProject    = "cannot get Project"
	errCreateProject = "cannot create Project"
	errUpdateProject = "cannot update Project"
	errDeleteProject = "cannot delete Project"
)

type External struct {
	kube     client.Client
	recorder event.Recorder
}

func (e *External) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	identityClient, err := clients.NewIdentityClient(ctx, e.kube, mg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errNewClient, err)
	}
	return &openstackProjectClient{
		client:   identityClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackProjectClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.ProjectGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ProjectGroupKind)),
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
		For(&v1alpha1.Project{}).
		Complete(r)
}

func (e *openstackProjectClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Project)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotProjectResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	project, err := projects.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetProject, err)
	}

	cr.Status.AtProvider = v1alpha1.ProjectProviderStatus{
		ProjectID:   project.ID,
		Name:        project.Name,
		Description: project.Description,
		DomainID:    project.DomainID,
		Enabled:     &project.Enabled,
		ParentID:    project.ParentID,
		IsDomain:    &project.IsDomain,
		Tags:        project.Tags,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackProjectClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Project)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotProjectResource)
	}

	createOpts := projects.CreateOpts{
		Name:        cr.Spec.ForProvider.Name,
		Description: cr.Spec.ForProvider.Description,
		DomainID:    cr.Spec.ForProvider.DomainID,
		IsDomain:    cr.Spec.ForProvider.IsDomain,
		ParentID:    cr.Spec.ForProvider.ParentID,
		Enabled:     cr.Spec.ForProvider.Enabled,
		Tags:        cr.Spec.ForProvider.Tags,
	}

	project, err := projects.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateProject, err)
	}

	meta.SetExternalName(cr, project.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackProjectClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Project)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotProjectResource)
	}

	desc := cr.Spec.ForProvider.Description
	updateOpts := projects.UpdateOpts{
		Name:        cr.Spec.ForProvider.Name,
		Description: &desc,
		Enabled:     cr.Spec.ForProvider.Enabled,
		Tags:        &cr.Spec.ForProvider.Tags,
	}

	_, err := projects.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateProject, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackProjectClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Project)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotProjectResource)
	}

	result := projects.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackProjectClient) Disconnect(ctx context.Context) error {
	return nil
}
