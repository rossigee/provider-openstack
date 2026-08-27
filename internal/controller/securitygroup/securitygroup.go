package securitygroup

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
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
	errNotSecurityGroupResource = "managed resource is not a SecurityGroup"
	errTrackUsage               = "cannot track ProviderConfig usage"
	errGetProviderConfig        = "cannot get ProviderConfig"
	errGetCredentials           = "cannot get credentials"
	errNewClient                = "cannot create OpenStack client"

	errGetSecurityGroup    = "cannot get SecurityGroup"
	errCreateSecurityGroup = "cannot create SecurityGroup"
	errUpdateSecurityGroup = "cannot update SecurityGroup"
	errDeleteSecurityGroup = "cannot delete SecurityGroup"
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
	return &openstackSecurityGroupClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackSecurityGroupClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.SecurityGroupGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.SecurityGroupGroupKind)),
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
		For(&v1alpha1.SecurityGroup{}).
		Complete(r)
}

func (e *openstackSecurityGroupClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.SecurityGroup)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotSecurityGroupResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	group, err := groups.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetSecurityGroup, err)
	}

	cr.Status.AtProvider = v1alpha1.SecurityGroupProviderStatus{
		SecurityGroupID: group.ID,
		Name:            group.Name,
		Description:     group.Description,
		TenantID:        group.TenantID,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackSecurityGroupClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.SecurityGroup)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotSecurityGroupResource)
	}

	createOpts := groups.CreateOpts{
		Name:        cr.Spec.ForProvider.Name,
		Description: cr.Spec.ForProvider.Description,
		TenantID:    cr.Spec.ForProvider.TenantID,
	}

	group, err := groups.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateSecurityGroup, err)
	}

	meta.SetExternalName(cr, group.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackSecurityGroupClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.SecurityGroup)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotSecurityGroupResource)
	}

	desc := cr.Spec.ForProvider.Description
	updateOpts := groups.UpdateOpts{
		Name:        cr.Spec.ForProvider.Name,
		Description: &desc,
	}

	_, err := groups.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateSecurityGroup, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackSecurityGroupClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.SecurityGroup)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotSecurityGroupResource)
	}

	result := groups.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackSecurityGroupClient) Disconnect(ctx context.Context) error {
	return nil
}
