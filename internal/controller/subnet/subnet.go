package subnet

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
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
	errNotSubnetResource = "managed resource is not a Subnet"
	errTrackUsage        = "cannot track ProviderConfig usage"
	errGetProviderConfig = "cannot get ProviderConfig"
	errGetCredentials    = "cannot get credentials"
	errNewClient         = "cannot create OpenStack client"

	errGetSubnet    = "cannot get Subnet"
	errCreateSubnet = "cannot create Subnet"
	errUpdateSubnet = "cannot update Subnet"
	errDeleteSubnet = "cannot delete Subnet"
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
	return &openstackSubnetClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackSubnetClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.SubnetGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.SubnetGroupKind)),
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
		For(&v1alpha1.Subnet{}).
		Complete(r)
}

func (e *openstackSubnetClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Subnet)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotSubnetResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	subnet, err := subnets.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetSubnet, err)
	}

	cr.Status.AtProvider = v1alpha1.SubnetProviderStatus{
		SubnetID:        subnet.ID,
		NetworkID:       subnet.NetworkID,
		CIDR:            subnet.CIDR,
		Gateway:         subnet.GatewayIP,
		DNSNameservers:  subnet.DNSNameservers,
		IPVersion:       subnet.IPVersion,
		EnableDHCP:      subnet.EnableDHCP,
		TenantID:        subnet.TenantID,
		AllocationPools: convertAllocationPools(subnet.AllocationPools),
		HostRoutes:      convertHostRoutes(subnet.HostRoutes),
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackSubnetClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Subnet)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotSubnetResource)
	}

	createOpts := subnets.CreateOpts{
		Name:            cr.Spec.ForProvider.Name,
		Description:     cr.Spec.ForProvider.Description,
		CIDR:            cr.Spec.ForProvider.CIDR,
		IPVersion:       gophercloud.IPVersion(cr.Spec.ForProvider.IPVersion),
		NetworkID:       cr.Spec.ForProvider.NetworkID,
		TenantID:        cr.Spec.ForProvider.TenantID,
		EnableDHCP:      cr.Spec.ForProvider.EnableDHCP,
		DNSNameservers:  cr.Spec.ForProvider.DNSNameservers,
		IPv6AddressMode: cr.Spec.ForProvider.IPv6AddressMode,
		IPv6RAMode:      cr.Spec.ForProvider.IPv6RAMode,
		SubnetPoolID:    cr.Spec.ForProvider.SubnetPoolID,
	}

	if cr.Spec.ForProvider.Gateway != "" {
		gw := cr.Spec.ForProvider.Gateway
		createOpts.GatewayIP = &gw
	}

	subnet, err := subnets.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateSubnet, err)
	}

	meta.SetExternalName(cr, subnet.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackSubnetClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Subnet)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotSubnetResource)
	}

	updateOpts := subnets.UpdateOpts{
		Name:           &cr.Spec.ForProvider.Name,
		Description:    &cr.Spec.ForProvider.Description,
		DNSNameservers: &cr.Spec.ForProvider.DNSNameservers,
		EnableDHCP:     cr.Spec.ForProvider.EnableDHCP,
	}

	_, err := subnets.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateSubnet, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackSubnetClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Subnet)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotSubnetResource)
	}

	result := subnets.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackSubnetClient) Disconnect(ctx context.Context) error {
	return nil
}

func convertAllocationPools(pools []subnets.AllocationPool) []v1alpha1.AllocationPool {
	if pools == nil {
		return nil
	}
	result := make([]v1alpha1.AllocationPool, len(pools))
	for i, p := range pools {
		result[i] = v1alpha1.AllocationPool{
			Start: p.Start,
			End:   p.End,
		}
	}
	return result
}

func convertHostRoutes(routes []subnets.HostRoute) []v1alpha1.HostRoute {
	if routes == nil {
		return nil
	}
	result := make([]v1alpha1.HostRoute, len(routes))
	for i, r := range routes {
		result[i] = v1alpha1.HostRoute{
			Destination: r.DestinationCIDR,
			Nexthop:     r.NextHop,
		}
	}
	return result
}
