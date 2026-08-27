package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/rossigee/provider-openstack/apis/compute/v1alpha1"
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
	errNotServerResource = "managed resource is not a Server"
	errTrackUsage        = "cannot track ProviderConfig usage"
	errGetProviderConfig = "cannot get ProviderConfig"
	errGetCredentials    = "cannot get credentials"
	errNewClient         = "cannot create OpenStack client"

	errGetServer    = "cannot get Server"
	errCreateServer = "cannot create Server"
	errUpdateServer = "cannot update Server"
	errDeleteServer = "cannot delete Server"
)

type External struct {
	kube     client.Client
	recorder event.Recorder
}

func (e *External) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	computeClient, err := clients.NewComputeClient(ctx, e.kube, mg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errNewClient, err)
	}
	return &openstackServerClient{
		client:   computeClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackServerClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.ServerGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ServerGroupKind)),
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
		For(&v1alpha1.Server{}).
		Complete(r)
}

func (e *openstackServerClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Server)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotServerResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	srv, err := servers.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetServer, err)
	}

	flavorStr := ""
	if f, ok := srv.Flavor["original_name"].(string); ok {
		flavorStr = f
	}
	imageStr := ""
	if img, ok := srv.Image["name"].(string); ok {
		imageStr = img
	}

	cr.Status.AtProvider = v1alpha1.ServerProviderStatus{
		ServerID:  srv.ID,
		Status:    srv.Status,
		TenantID:  srv.TenantID,
		HostID:    srv.HostID,
		Flavor:    flavorStr,
		Image:     imageStr,
		Addresses: convertAddresses(srv.Addresses),
		KeyName:   srv.KeyName,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackServerClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Server)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotServerResource)
	}

	networks := make([]servers.Network, len(cr.Spec.ForProvider.Networks))
	for i, net := range cr.Spec.ForProvider.Networks {
		networks[i] = servers.Network{
			UUID:    net.UUID,
			FixedIP: net.FixedIP,
		}
	}

	createServerOpts := servers.CreateOpts{
		Name:      cr.Spec.ForProvider.Name,
		ImageRef:  cr.Spec.ForProvider.ImageRef,
		FlavorRef: cr.Spec.ForProvider.FlavorRef,
	}
	createServerOpts.Networks = networks

	createOpts := keypairs.CreateOptsExt{
		CreateOptsBuilder: &createServerOpts,
		KeyName:           cr.Spec.ForProvider.KeyName,
	}

	srv, err := servers.Create(ctx, e.client, &createOpts, nil).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateServer, err)
	}

	meta.SetExternalName(cr, srv.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackServerClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Server)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotServerResource)
	}

	updateOpts := servers.UpdateOpts{
		Name: cr.Spec.ForProvider.Name,
	}

	_, err := servers.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateServer, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackServerClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Server)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotServerResource)
	}

	result := servers.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackServerClient) Disconnect(ctx context.Context) error {
	return nil
}

func convertAddresses(addrs map[string]any) []v1alpha1.ServerAddress {
	var result []v1alpha1.ServerAddress
	for networkName, rawAddrs := range addrs {
		addrList, ok := rawAddrs.([]any)
		if !ok {
			continue
		}
		for _, rawAddr := range addrList {
			addrMap, ok := rawAddr.(map[string]any)
			if !ok {
				continue
			}
			version, _ := addrMap["version"].(float64)
			addr, _ := addrMap["addr"].(string)
			result = append(result, v1alpha1.ServerAddress{
				Network: networkName,
				Version: int(version),
				Address: addr,
			})
		}
	}
	return result
}
