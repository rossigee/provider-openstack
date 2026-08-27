package rbacpolicy

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/rbacpolicies"
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
	errNotRBACPolicyResource = "managed resource is not a RBACPolicy"
	errNewClient             = "cannot create OpenStack client"

	errGetRBACPolicy    = "cannot get RBACPolicy"
	errCreateRBACPolicy = "cannot create RBACPolicy"
	errUpdateRBACPolicy = "cannot update RBACPolicy"
	errDeleteRBACPolicy = "cannot delete RBACPolicy"
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
	return &openstackRBACPolicyClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackRBACPolicyClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.RBACPolicyGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.RBACPolicyGroupKind)),
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
		For(&v1alpha1.RBACPolicy{}).
		Complete(r)
}

func (e *openstackRBACPolicyClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.RBACPolicy)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotRBACPolicyResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	p, err := rbacpolicies.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetRBACPolicy, err)
	}

	cr.Status.AtProvider = v1alpha1.RBACPolicyProviderStatus{
		RBACPolicyID: p.ID,
		Action:       string(p.Action),
		ObjectType:   p.ObjectType,
		ObjectID:     p.ObjectID,
		TargetTenant: p.TargetTenant,
		TenantID:     p.TenantID,
		ProjectID:    p.ProjectID,
	}

	if !strings.EqualFold(cr.Spec.ForProvider.TargetTenant, p.TargetTenant) {
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

func (e *openstackRBACPolicyClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.RBACPolicy)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotRBACPolicyResource)
	}

	createOpts := rbacpolicies.CreateOpts{
		Action:       rbacpolicies.PolicyAction(cr.Spec.ForProvider.Action),
		ObjectType:   cr.Spec.ForProvider.ObjectType,
		TargetTenant: cr.Spec.ForProvider.TargetTenant,
		ObjectID:     cr.Spec.ForProvider.ObjectID,
	}

	p, err := rbacpolicies.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateRBACPolicy, err)
	}

	meta.SetExternalName(cr, p.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackRBACPolicyClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.RBACPolicy)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotRBACPolicyResource)
	}

	updateOpts := rbacpolicies.UpdateOpts{
		TargetTenant: cr.Spec.ForProvider.TargetTenant,
	}

	result := rbacpolicies.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts)
	if result.Err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateRBACPolicy, result.Err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackRBACPolicyClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.RBACPolicy)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotRBACPolicyResource)
	}

	result := rbacpolicies.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackRBACPolicyClient) Disconnect(ctx context.Context) error {
	return nil
}
