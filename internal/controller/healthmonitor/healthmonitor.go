package healthmonitor

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
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
	errNotHealthMonitorResource = "managed resource is not a HealthMonitor"
	errNewClient                = "cannot create OpenStack client"

	errGetHealthMonitor    = "cannot get HealthMonitor"
	errCreateHealthMonitor = "cannot create HealthMonitor"
	errUpdateHealthMonitor = "cannot update HealthMonitor"
	errDeleteHealthMonitor = "cannot delete HealthMonitor"
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
	return &openstackHealthMonitorClient{
		client:   lbClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackHealthMonitorClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.HealthMonitorGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.HealthMonitorGroupKind)),
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
		For(&v1alpha1.HealthMonitor{}).
		Complete(r)
}

func (e *openstackHealthMonitorClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.HealthMonitor)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotHealthMonitorResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	monitor, err := monitors.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetHealthMonitor, err)
	}

	cr.Status.AtProvider = v1alpha1.HealthMonitorProviderStatus{
		MonitorID:          monitor.ID,
		Name:               monitor.Name,
		Type:               monitor.Type,
		Delay:              monitor.Delay,
		Timeout:            monitor.Timeout,
		MaxRetries:         monitor.MaxRetries,
		MaxRetriesDown:     monitor.MaxRetriesDown,
		HTTPMethod:         monitor.HTTPMethod,
		URLPath:            monitor.URLPath,
		ExpectedCodes:      monitor.ExpectedCodes,
		DomainName:         monitor.DomainName,
		AdminStateUp:       monitor.AdminStateUp,
		ProvisioningStatus: monitor.ProvisioningStatus,
		OperatingStatus:    monitor.OperatingStatus,
		ProjectID:          monitor.ProjectID,
		Tags:               monitor.Tags,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackHealthMonitorClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.HealthMonitor)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotHealthMonitorResource)
	}

	createOpts := monitors.CreateOpts{
		Type:           cr.Spec.ForProvider.Type,
		Delay:          cr.Spec.ForProvider.Delay,
		Timeout:        cr.Spec.ForProvider.Timeout,
		MaxRetries:     cr.Spec.ForProvider.MaxRetries,
		MaxRetriesDown: cr.Spec.ForProvider.MaxRetriesDown,
		URLPath:        cr.Spec.ForProvider.URLPath,
		HTTPMethod:     cr.Spec.ForProvider.HTTPMethod,
		HTTPVersion:    cr.Spec.ForProvider.HTTPVersion,
		ExpectedCodes:  cr.Spec.ForProvider.ExpectedCodes,
		DomainName:     cr.Spec.ForProvider.DomainName,
		AdminStateUp:   cr.Spec.ForProvider.AdminStateUp,
		Name:           cr.Spec.ForProvider.Name,
		ProjectID:      cr.Spec.ForProvider.ProjectID,
		Tags:           cr.Spec.ForProvider.Tags,
	}

	monitor, err := monitors.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateHealthMonitor, err)
	}

	meta.SetExternalName(cr, monitor.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackHealthMonitorClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.HealthMonitor)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotHealthMonitorResource)
	}

	updateOpts := monitors.UpdateOpts{
		Delay:          cr.Spec.ForProvider.Delay,
		Timeout:        cr.Spec.ForProvider.Timeout,
		MaxRetries:     cr.Spec.ForProvider.MaxRetries,
		MaxRetriesDown: cr.Spec.ForProvider.MaxRetriesDown,
		URLPath:        cr.Spec.ForProvider.URLPath,
		HTTPMethod:     cr.Spec.ForProvider.HTTPMethod,
		ExpectedCodes:  cr.Spec.ForProvider.ExpectedCodes,
		Name:           &cr.Spec.ForProvider.Name,
		DomainName:     &cr.Spec.ForProvider.DomainName,
		AdminStateUp:   cr.Spec.ForProvider.AdminStateUp,
		Tags:           cr.Spec.ForProvider.Tags,
	}

	_, err := monitors.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateHealthMonitor, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackHealthMonitorClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.HealthMonitor)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotHealthMonitorResource)
	}

	result := monitors.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackHealthMonitorClient) Disconnect(ctx context.Context) error {
	return nil
}
