package routerinterface

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
	errNotRouterInterfaceResource = "managed resource is not a RouterInterface"
	errTrackUsage                 = "cannot track ProviderConfig usage"
	errGetProviderConfig          = "cannot get ProviderConfig"
	errGetCredentials             = "cannot get credentials"
	errNewClient                  = "cannot create OpenStack client"

	errGetRouterInterface    = "cannot get RouterInterface"
	errCreateRouterInterface = "cannot create RouterInterface"
	errUpdateRouterInterface = "cannot update RouterInterface"
	errDeleteRouterInterface = "cannot delete RouterInterface"
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
	return &openstackRouterInterfaceClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackRouterInterfaceClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.RouterInterfaceGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.RouterInterfaceGroupKind)),
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
		For(&v1alpha1.RouterInterface{}).
		Complete(r)
}

func (e *openstackRouterInterfaceClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.RouterInterface)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotRouterInterfaceResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	routerID := cr.Spec.ForProvider.RouterID
	if routerID == "" {
		routerID = meta.GetExternalName(cr)
	}

	router, err := routers.Get(ctx, e.client, routerID).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetRouterInterface, err)
	}

	cr.Status.AtProvider = v1alpha1.RouterInterfaceProviderStatus{
		RouterID: router.ID,
		SubnetID: cr.Spec.ForProvider.SubnetID,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackRouterInterfaceClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.RouterInterface)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotRouterInterfaceResource)
	}

	routerID := cr.Spec.ForProvider.RouterID
	if routerID == "" {
		routerID = meta.GetExternalName(cr)
	}

	iface, err := routers.AddInterface(ctx, e.client, routerID, routers.AddInterfaceOpts{
		SubnetID: cr.Spec.ForProvider.SubnetID,
	}).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateRouterInterface, err)
	}

	meta.SetExternalName(cr, iface.PortID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackRouterInterfaceClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}

func (e *openstackRouterInterfaceClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.RouterInterface)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotRouterInterfaceResource)
	}

	routerID := cr.Spec.ForProvider.RouterID
	if routerID == "" {
		routerID = meta.GetExternalName(cr)
	}

	result := routers.RemoveInterface(ctx, e.client, routerID, routers.RemoveInterfaceOpts{
		SubnetID: cr.Spec.ForProvider.SubnetID,
	})
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackRouterInterfaceClient) Disconnect(ctx context.Context) error {
	return nil
}
