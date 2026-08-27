package router

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
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
	errNotRouterResource = "managed resource is not a Router"
	errTrackUsage        = "cannot track ProviderConfig usage"
	errGetProviderConfig = "cannot get ProviderConfig"
	errGetCredentials    = "cannot get credentials"
	errNewClient         = "cannot create OpenStack client"

	errGetRouter    = "cannot get Router"
	errCreateRouter = "cannot create Router"
	errUpdateRouter = "cannot update Router"
	errDeleteRouter = "cannot delete Router"
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
	return &openstackRouterClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackRouterClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.RouterGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.RouterGroupKind)),
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
		For(&v1alpha1.Router{}).
		Complete(r)
}

func (e *openstackRouterClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Router)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotRouterResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	router, err := routers.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetRouter, err)
	}

	cr.Status.AtProvider = v1alpha1.RouterProviderStatus{
		RouterID:         router.ID,
		Status:           router.Status,
		TenantID:         router.TenantID,
		AdminStateUp:     router.AdminStateUp,
		Distributed:      router.Distributed,
		GatewayNetworkID: router.GatewayInfo.NetworkID,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackRouterClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Router)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotRouterResource)
	}

	createOpts := routers.CreateOpts{
		Name:         cr.Spec.ForProvider.Name,
		TenantID:     cr.Spec.ForProvider.TenantID,
		AdminStateUp: cr.Spec.ForProvider.AdminStateUp,
	}

	router, err := routers.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateRouter, err)
	}

	meta.SetExternalName(cr, router.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackRouterClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Router)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotRouterResource)
	}

	updateOpts := routers.UpdateOpts{
		Name:         cr.Spec.ForProvider.Name,
		AdminStateUp: cr.Spec.ForProvider.AdminStateUp,
	}

	_, err := routers.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateRouter, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackRouterClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Router)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotRouterResource)
	}

	result := routers.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackRouterClient) Disconnect(ctx context.Context) error {
	return nil
}
