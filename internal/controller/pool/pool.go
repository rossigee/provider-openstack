package pool

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"github.com/rossigee/provider-openstack/apis/loadbalancing/v1alpha1"
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
	errNotPoolResource = "managed resource is not a Pool"
	errNewClient       = "cannot create OpenStack client"

	errGetPool    = "cannot get Pool"
	errCreatePool = "cannot create Pool"
	errUpdatePool = "cannot update Pool"
	errDeletePool = "cannot delete Pool"
)

type External struct {
	kube     client.Client
	recorder event.Recorder
}

func (e *External) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	lbClient, err := clients.NewLoadBalancerClient(ctx, e.kube, mg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errNewClient, err)
	}
	return &openstackPoolClient{
		client:   lbClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackPoolClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.PoolGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.PoolGroupKind)),
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
		For(&v1alpha1.Pool{}).
		Complete(r)
}

func (e *openstackPoolClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Pool)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotPoolResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	pool, err := pools.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetPool, err)
	}

	cr.Status.AtProvider = v1alpha1.PoolProviderStatus{
		PoolID:             pool.ID,
		Name:               pool.Name,
		Description:        pool.Description,
		LBMethod:           pool.LBMethod,
		Protocol:           pool.Protocol,
		AdminStateUp:       pool.AdminStateUp,
		ProvisioningStatus: pool.ProvisioningStatus,
		OperatingStatus:    pool.OperatingStatus,
		ProjectID:          pool.ProjectID,
		MonitorID:          pool.MonitorID,
		Tags:               pool.Tags,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackPoolClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Pool)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotPoolResource)
	}

	createOpts := pools.CreateOpts{
		LBMethod:       pools.LBMethod(cr.Spec.ForProvider.LBMethod),
		Protocol:       pools.Protocol(cr.Spec.ForProvider.Protocol),
		LoadbalancerID: cr.Spec.ForProvider.LoadbalancerID,
		ListenerID:     cr.Spec.ForProvider.ListenerID,
		Name:           cr.Spec.ForProvider.Name,
		Description:    cr.Spec.ForProvider.Description,
		AdminStateUp:   cr.Spec.ForProvider.AdminStateUp,
		ProjectID:      cr.Spec.ForProvider.ProjectID,
		Tags:           cr.Spec.ForProvider.Tags,
	}

	if cr.Spec.ForProvider.Persistence != nil {
		createOpts.Persistence = &pools.SessionPersistence{
			Type:       cr.Spec.ForProvider.Persistence.Type,
			CookieName: cr.Spec.ForProvider.Persistence.CookieName,
		}
	}

	pool, err := pools.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreatePool, err)
	}

	meta.SetExternalName(cr, pool.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackPoolClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Pool)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotPoolResource)
	}

	updateOpts := pools.UpdateOpts{
		Name:         &cr.Spec.ForProvider.Name,
		Description:  &cr.Spec.ForProvider.Description,
		LBMethod:     pools.LBMethod(cr.Spec.ForProvider.LBMethod),
		AdminStateUp: cr.Spec.ForProvider.AdminStateUp,
		Tags:         &cr.Spec.ForProvider.Tags,
	}

	if cr.Spec.ForProvider.Persistence != nil {
		updateOpts.Persistence = &pools.SessionPersistence{
			Type:       cr.Spec.ForProvider.Persistence.Type,
			CookieName: cr.Spec.ForProvider.Persistence.CookieName,
		}
	}

	_, err := pools.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdatePool, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackPoolClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Pool)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotPoolResource)
	}

	result := pools.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackPoolClient) Disconnect(ctx context.Context) error {
	return nil
}
