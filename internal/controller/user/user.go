package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
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
	errNotUserResource   = "managed resource is not a User"
	errTrackUsage        = "cannot track ProviderConfig usage"
	errGetProviderConfig = "cannot get ProviderConfig"
	errGetCredentials    = "cannot get credentials"
	errNewClient         = "cannot create OpenStack client"

	errGetUser    = "cannot get User"
	errCreateUser = "cannot create User"
	errUpdateUser = "cannot update User"
	errDeleteUser = "cannot delete User"
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
	return &openstackUserClient{
		client:   identityClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackUserClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.UserGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.UserGroupKind)),
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
		For(&v1alpha1.User{}).
		Complete(r)
}

func (e *openstackUserClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotUserResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	user, err := users.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetUser, err)
	}

	cr.Status.AtProvider = v1alpha1.UserProviderStatus{
		UserID:           user.ID,
		Name:             user.Name,
		Description:      user.Description,
		DomainID:         user.DomainID,
		DefaultProjectID: user.DefaultProjectID,
		Enabled:          &user.Enabled,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackUserClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotUserResource)
	}

	createOpts := users.CreateOpts{
		Name:             cr.Spec.ForProvider.Name,
		Description:      cr.Spec.ForProvider.Description,
		DomainID:         cr.Spec.ForProvider.DomainID,
		DefaultProjectID: cr.Spec.ForProvider.DefaultProjectID,
		Password:         cr.Spec.ForProvider.Password,
		Enabled:          cr.Spec.ForProvider.Enabled,
	}

	user, err := users.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateUser, err)
	}

	meta.SetExternalName(cr, user.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackUserClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotUserResource)
	}

	desc := cr.Spec.ForProvider.Description
	updateOpts := users.UpdateOpts{
		Name:             cr.Spec.ForProvider.Name,
		Description:      &desc,
		DefaultProjectID: cr.Spec.ForProvider.DefaultProjectID,
		Enabled:          cr.Spec.ForProvider.Enabled,
	}

	_, err := users.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateUser, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackUserClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotUserResource)
	}

	result := users.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackUserClient) Disconnect(ctx context.Context) error {
	return nil
}
