package subnetpool

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/subnetpools"
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
	errNotSubnetPoolResource = "managed resource is not a SubnetPool"
	errNewClient             = "cannot create OpenStack client"

	errGetSubnetPool    = "cannot get SubnetPool"
	errCreateSubnetPool = "cannot create SubnetPool"
	errUpdateSubnetPool = "cannot update SubnetPool"
	errDeleteSubnetPool = "cannot delete SubnetPool"
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
	return &openstackSubnetPoolClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackSubnetPoolClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.SubnetPoolGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.SubnetPoolGroupKind)),
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
		For(&v1alpha1.SubnetPool{}).
		Complete(r)
}

func (e *openstackSubnetPoolClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.SubnetPool)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotSubnetPoolResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	sp, err := subnetpools.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetSubnetPool, err)
	}

	cr.Status.AtProvider = v1alpha1.SubnetPoolProviderStatus{
		SubnetPoolID:     sp.ID,
		Name:             sp.Name,
		Prefixes:         sp.Prefixes,
		DefaultQuota:     sp.DefaultQuota,
		TenantID:         sp.TenantID,
		DefaultPrefixLen: sp.DefaultPrefixLen,
		MinPrefixLen:     sp.MinPrefixLen,
		MaxPrefixLen:     sp.MaxPrefixLen,
		AddressScopeID:   sp.AddressScopeID,
		IPVersion:        sp.IPversion,
		Shared:           sp.Shared,
		Description:      sp.Description,
		IsDefault:        sp.IsDefault,
		RevisionNumber:   sp.RevisionNumber,
	}

	if !strings.EqualFold(cr.Spec.ForProvider.Name, sp.Name) {
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

func (e *openstackSubnetPoolClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.SubnetPool)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotSubnetPoolResource)
	}

	createOpts := subnetpools.CreateOpts{
		Name:             cr.Spec.ForProvider.Name,
		Prefixes:         cr.Spec.ForProvider.Prefixes,
		DefaultQuota:     cr.Spec.ForProvider.DefaultQuota,
		TenantID:         cr.Spec.ForProvider.TenantID,
		DefaultPrefixLen: cr.Spec.ForProvider.DefaultPrefixLen,
		MinPrefixLen:     cr.Spec.ForProvider.MinPrefixLen,
		MaxPrefixLen:     cr.Spec.ForProvider.MaxPrefixLen,
		AddressScopeID:   cr.Spec.ForProvider.AddressScopeID,
		Shared:           boolValue(cr.Spec.ForProvider.Shared),
		Description:      cr.Spec.ForProvider.Description,
		IsDefault:        boolValue(cr.Spec.ForProvider.IsDefault),
	}

	sp, err := subnetpools.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateSubnetPool, err)
	}

	meta.SetExternalName(cr, sp.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackSubnetPoolClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.SubnetPool)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotSubnetPoolResource)
	}

	updateOpts := subnetpools.UpdateOpts{
		Name:             cr.Spec.ForProvider.Name,
		Prefixes:         cr.Spec.ForProvider.Prefixes,
		DefaultQuota:     ptrInt(cr.Spec.ForProvider.DefaultQuota),
		TenantID:         cr.Spec.ForProvider.TenantID,
		DefaultPrefixLen: cr.Spec.ForProvider.DefaultPrefixLen,
		MinPrefixLen:     cr.Spec.ForProvider.MinPrefixLen,
		MaxPrefixLen:     cr.Spec.ForProvider.MaxPrefixLen,
		AddressScopeID:   ptrString(cr.Spec.ForProvider.AddressScopeID),
		Description:      ptrString(cr.Spec.ForProvider.Description),
		IsDefault:        cr.Spec.ForProvider.IsDefault,
	}

	_, err := subnetpools.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateSubnetPool, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackSubnetPoolClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.SubnetPool)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotSubnetPoolResource)
	}

	result := subnetpools.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackSubnetPoolClient) Disconnect(ctx context.Context) error {
	return nil
}

func boolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func ptrInt(i int) *int {
	return &i
}

func ptrString(s string) *string {
	return &s
}
