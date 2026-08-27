package loadbalancer

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
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
	errNotLoadBalancerResource = "managed resource is not a LoadBalancer"
	errTrackUsage              = "cannot track ProviderConfig usage"
	errGetProviderConfig       = "cannot get ProviderConfig"
	errGetCredentials          = "cannot get credentials"
	errNewClient               = "cannot create OpenStack client"

	errGetLoadBalancer    = "cannot get LoadBalancer"
	errCreateLoadBalancer = "cannot create LoadBalancer"
	errUpdateLoadBalancer = "cannot update LoadBalancer"
	errDeleteLoadBalancer = "cannot delete LoadBalancer"
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
	return &openstackLoadBalancerClient{
		client:   lbClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackLoadBalancerClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.LoadBalancerGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.LoadBalancerGroupKind)),
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
		For(&v1alpha1.LoadBalancer{}).
		Complete(r)
}

func (e *openstackLoadBalancerClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.LoadBalancer)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotLoadBalancerResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	lb, err := loadbalancers.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetLoadBalancer, err)
	}

	cr.Status.AtProvider = v1alpha1.LoadBalancerProviderStatus{
		LoadBalancerID:     lb.ID,
		Name:               lb.Name,
		Description:        lb.Description,
		ProvisioningStatus: lb.ProvisioningStatus,
		OperatingStatus:    lb.OperatingStatus,
		VipAddress:         lb.VipAddress,
		VipPortID:          lb.VipPortID,
		VipSubnetID:        lb.VipSubnetID,
		VipNetworkID:       lb.VipNetworkID,
		ProjectID:          lb.ProjectID,
		AdminStateUp:       lb.AdminStateUp,
		FlavorID:           lb.FlavorID,
		Provider:           lb.Provider,
		Tags:               lb.Tags,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackLoadBalancerClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.LoadBalancer)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotLoadBalancerResource)
	}

	createOpts := loadbalancers.CreateOpts{
		Name:             cr.Spec.ForProvider.Name,
		Description:      cr.Spec.ForProvider.Description,
		VipSubnetID:      cr.Spec.ForProvider.VipSubnetID,
		VipNetworkID:     cr.Spec.ForProvider.VipNetworkID,
		VipPortID:        cr.Spec.ForProvider.VipPortID,
		VipAddress:       cr.Spec.ForProvider.VipAddress,
		VipQosPolicyID:   cr.Spec.ForProvider.VipQosPolicyID,
		AdminStateUp:     cr.Spec.ForProvider.AdminStateUp,
		ProjectID:        cr.Spec.ForProvider.ProjectID,
		FlavorID:         cr.Spec.ForProvider.FlavorID,
		AvailabilityZone: cr.Spec.ForProvider.AvailabilityZone,
		Provider:         cr.Spec.ForProvider.Provider,
		Tags:             cr.Spec.ForProvider.Tags,
	}

	lb, err := loadbalancers.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateLoadBalancer, err)
	}

	meta.SetExternalName(cr, lb.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackLoadBalancerClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.LoadBalancer)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotLoadBalancerResource)
	}

	updateOpts := loadbalancers.UpdateOpts{
		Name:           &cr.Spec.ForProvider.Name,
		Description:    &cr.Spec.ForProvider.Description,
		AdminStateUp:   cr.Spec.ForProvider.AdminStateUp,
		VipQosPolicyID: &cr.Spec.ForProvider.VipQosPolicyID,
		Tags:           &cr.Spec.ForProvider.Tags,
	}

	_, err := loadbalancers.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateLoadBalancer, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackLoadBalancerClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.LoadBalancer)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotLoadBalancerResource)
	}

	result := loadbalancers.Delete(ctx, e.client, meta.GetExternalName(cr), loadbalancers.DeleteOpts{Cascade: true})
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackLoadBalancerClient) Disconnect(ctx context.Context) error {
	return nil
}
