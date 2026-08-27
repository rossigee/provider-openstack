package network

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
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
	errNotNetworkResource = "managed resource is not a Network"
	errTrackUsage         = "cannot track ProviderConfig usage"
	errGetProviderConfig  = "cannot get ProviderConfig"
	errGetCredentials     = "cannot get credentials"
	errNewClient          = "cannot create OpenStack client"

	errGetNetwork    = "cannot get Network"
	errCreateNetwork = "cannot create Network"
	errUpdateNetwork = "cannot update Network"
	errDeleteNetwork = "cannot delete Network"
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
	return &openstackNetworkClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackNetworkClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.NetworkGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.NetworkGroupKind)),
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
		For(&v1alpha1.Network{}).
		Complete(r)
}

func (e *openstackNetworkClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Network)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotNetworkResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	network, err := networks.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetNetwork, err)
	}

	cr.Status.AtProvider = v1alpha1.NetworkProviderStatus{
		NetworkID:    network.ID,
		Status:       network.Status,
		Subnets:      network.Subnets,
		TenantID:     network.TenantID,
		AdminStateUp: network.AdminStateUp,
		Shared:       network.Shared,
		Tags:         network.Tags,
	}

	if !strings.EqualFold(cr.Spec.ForProvider.Name, network.Name) {
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

func (e *openstackNetworkClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Network)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotNetworkResource)
	}

	createOpts := networks.CreateOpts{
		Name:         cr.Spec.ForProvider.Name,
		Shared:       cr.Spec.ForProvider.Shared,
		AdminStateUp: cr.Spec.ForProvider.AdminStateUp,
		TenantID:     cr.Spec.ForProvider.TenantID,
	}

	network, err := networks.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateNetwork, err)
	}

	meta.SetExternalName(cr, network.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackNetworkClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Network)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotNetworkResource)
	}

	updateOpts := networks.UpdateOpts{
		Name:         &cr.Spec.ForProvider.Name,
		Shared:       cr.Spec.ForProvider.Shared,
		AdminStateUp: cr.Spec.ForProvider.AdminStateUp,
	}

	_, err := networks.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateNetwork, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackNetworkClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Network)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotNetworkResource)
	}

	result := networks.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackNetworkClient) Disconnect(ctx context.Context) error {
	return nil
}
