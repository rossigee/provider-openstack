package listener

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
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
	errNotListenerResource = "managed resource is not a Listener"
	errNewClient           = "cannot create OpenStack client"

	errGetListener    = "cannot get Listener"
	errCreateListener = "cannot create Listener"
	errUpdateListener = "cannot update Listener"
	errDeleteListener = "cannot delete Listener"
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
	return &openstackListenerClient{
		client:   lbClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackListenerClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.ListenerGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ListenerGroupKind)),
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
		For(&v1alpha1.Listener{}).
		Complete(r)
}

func (e *openstackListenerClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Listener)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotListenerResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	listener, err := listeners.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetListener, err)
	}

	cr.Status.AtProvider = v1alpha1.ListenerProviderStatus{
		ListenerID:         listener.ID,
		Name:               listener.Name,
		Description:        listener.Description,
		LoadbalancerID:     findLoadBalancerID(listener.Loadbalancers),
		Protocol:           listener.Protocol,
		ProtocolPort:       listener.ProtocolPort,
		DefaultPoolID:      listener.DefaultPoolID,
		ConnLimit:          listener.ConnLimit,
		AdminStateUp:       listener.AdminStateUp,
		ProvisioningStatus: listener.ProvisioningStatus,
		OperatingStatus:    listener.OperatingStatus,
		ProjectID:          listener.ProjectID,
		Tags:               listener.Tags,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func findLoadBalancerID(lbs []listeners.LoadBalancerID) string {
	if len(lbs) > 0 {
		return lbs[0].ID
	}
	return ""
}

func (e *openstackListenerClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Listener)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotListenerResource)
	}

	createOpts := listeners.CreateOpts{
		LoadbalancerID:         cr.Spec.ForProvider.LoadbalancerID,
		Protocol:               listeners.Protocol(cr.Spec.ForProvider.Protocol),
		ProtocolPort:           cr.Spec.ForProvider.ProtocolPort,
		Name:                   cr.Spec.ForProvider.Name,
		Description:            cr.Spec.ForProvider.Description,
		DefaultPoolID:          cr.Spec.ForProvider.DefaultPoolID,
		ConnLimit:              cr.Spec.ForProvider.ConnLimit,
		DefaultTlsContainerRef: cr.Spec.ForProvider.DefaultTlsContainerRef,
		SniContainerRefs:       cr.Spec.ForProvider.SniContainerRefs,
		AdminStateUp:           cr.Spec.ForProvider.AdminStateUp,
		ProjectID:              cr.Spec.ForProvider.ProjectID,
		TimeoutClientData:      cr.Spec.ForProvider.TimeoutClientData,
		TimeoutMemberData:      cr.Spec.ForProvider.TimeoutMemberData,
		TimeoutMemberConnect:   cr.Spec.ForProvider.TimeoutMemberConnect,
		TimeoutTCPInspect:      cr.Spec.ForProvider.TimeoutTCPInspect,
		InsertHeaders:          cr.Spec.ForProvider.InsertHeaders,
		AllowedCIDRs:           cr.Spec.ForProvider.AllowedCIDRs,
		Tags:                   cr.Spec.ForProvider.Tags,
	}

	listener, err := listeners.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateListener, err)
	}

	meta.SetExternalName(cr, listener.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackListenerClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Listener)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotListenerResource)
	}

	updateOpts := listeners.UpdateOpts{
		Name:                   &cr.Spec.ForProvider.Name,
		Description:            &cr.Spec.ForProvider.Description,
		DefaultPoolID:          &cr.Spec.ForProvider.DefaultPoolID,
		ConnLimit:              cr.Spec.ForProvider.ConnLimit,
		DefaultTlsContainerRef: &cr.Spec.ForProvider.DefaultTlsContainerRef,
		SniContainerRefs:       &cr.Spec.ForProvider.SniContainerRefs,
		AdminStateUp:           cr.Spec.ForProvider.AdminStateUp,
		TimeoutClientData:      cr.Spec.ForProvider.TimeoutClientData,
		TimeoutMemberData:      cr.Spec.ForProvider.TimeoutMemberData,
		TimeoutMemberConnect:   cr.Spec.ForProvider.TimeoutMemberConnect,
		TimeoutTCPInspect:      cr.Spec.ForProvider.TimeoutTCPInspect,
		InsertHeaders:          &cr.Spec.ForProvider.InsertHeaders,
		AllowedCIDRs:           &cr.Spec.ForProvider.AllowedCIDRs,
		Tags:                   &cr.Spec.ForProvider.Tags,
	}

	_, err := listeners.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateListener, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackListenerClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Listener)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotListenerResource)
	}

	result := listeners.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackListenerClient) Disconnect(ctx context.Context) error {
	return nil
}
