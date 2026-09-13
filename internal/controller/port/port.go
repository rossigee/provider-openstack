package port

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	networkingv1beta1 "github.com/rossigee/provider-openstack/apis/networking/v1beta1"
	"github.com/rossigee/provider-openstack/internal/clients"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/event"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/rossigee/provider-openstack/internal/features"
)

const (
	errNotPortResource = "managed resource is not a Port"
	errNewClient       = "cannot create OpenStack client"

	errGetPort    = "cannot get Port"
	errCreatePort = "cannot create Port"
	errUpdatePort = "cannot update Port"
	errDeletePort = "cannot delete Port"
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
	return &openstackPortClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackPortClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(networkingv1beta1.PortGroupKind)
	rec := event.NewNopRecorder()

	opts := []managed.ReconcilerOption{
		managed.WithExternalConnector(&External{
			kube:     mgr.GetClient(),
			recorder: rec,
		}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(rec),
		managed.WithPollInterval(o.PollInterval),
	}
	if o.Features.Enabled(features.EnableAlphaManagementPolicies) {
		opts = append(opts, managed.WithManagementPolicies())
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(networkingv1beta1.SchemeGroupVersion.WithKind(networkingv1beta1.PortGroupKind)),
		opts...)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&networkingv1beta1.Port{}).
		Complete(r)
}

func (e *openstackPortClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*networkingv1beta1.Port)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotPortResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	port, err := ports.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetPort, err)
	}

	cr.Status.AtProvider = networkingv1beta1.PortProviderStatus{
		PortID:              port.ID,
		NetworkID:           port.NetworkID,
		Name:                port.Name,
		Description:         port.Description,
		AdminStateUp:        port.AdminStateUp,
		Status:              port.Status,
		MACAddress:          port.MACAddress,
		TenantID:            port.TenantID,
		DeviceOwner:         port.DeviceOwner,
		DeviceID:            port.DeviceID,
		SecurityGroups:      port.SecurityGroups,
		AllowedAddressPairs: toAddressPairStatus(port.AllowedAddressPairs),
	}

	for _, fip := range port.FixedIPs {
		cr.Status.AtProvider.FixedIPs = append(cr.Status.AtProvider.FixedIPs, networkingv1beta1.FixedIP{
			SubnetID:  fip.SubnetID,
			IPAddress: fip.IPAddress,
		})
	}

	if !strings.EqualFold(cr.Spec.ForProvider.Name, port.Name) {
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

func (e *openstackPortClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*networkingv1beta1.Port)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotPortResource)
	}

	createOpts := ports.CreateOpts{
		NetworkID:           cr.Spec.ForProvider.NetworkID,
		Name:                cr.Spec.ForProvider.Name,
		Description:         cr.Spec.ForProvider.Description,
		AdminStateUp:        cr.Spec.ForProvider.AdminStateUp,
		MACAddress:          cr.Spec.ForProvider.MACAddress,
		DeviceID:            cr.Spec.ForProvider.DeviceID,
		DeviceOwner:         cr.Spec.ForProvider.DeviceOwner,
		TenantID:            cr.Spec.ForProvider.TenantID,
		SecurityGroups:      cr.Spec.ForProvider.SecurityGroups,
		AllowedAddressPairs: toAddressPairs(cr.Spec.ForProvider.AllowedAddressPairs),
	}

	if cr.Spec.ForProvider.FixedIPs != nil {
		fips := make([]ports.IP, 0, len(cr.Spec.ForProvider.FixedIPs))
		for _, fip := range cr.Spec.ForProvider.FixedIPs {
			fips = append(fips, ports.IP{
				SubnetID:  fip.SubnetID,
				IPAddress: fip.IPAddress,
			})
		}
		createOpts.FixedIPs = fips
	}

	port, err := ports.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreatePort, err)
	}

	meta.SetExternalName(cr, port.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackPortClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*networkingv1beta1.Port)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotPortResource)
	}

	updateOpts := ports.UpdateOpts{
		Name:                &cr.Spec.ForProvider.Name,
		Description:         &cr.Spec.ForProvider.Description,
		AdminStateUp:        cr.Spec.ForProvider.AdminStateUp,
		DeviceID:            &cr.Spec.ForProvider.DeviceID,
		DeviceOwner:         &cr.Spec.ForProvider.DeviceOwner,
		SecurityGroups:      cr.Spec.ForProvider.SecurityGroups,
		AllowedAddressPairs: ptrAddressPairs(toAddressPairs(cr.Spec.ForProvider.AllowedAddressPairs)),
	}

	_, err := ports.Update(ctx, e.client, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdatePort, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackPortClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*networkingv1beta1.Port)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotPortResource)
	}

	result := ports.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackPortClient) Disconnect(ctx context.Context) error {
	return nil
}

func toAddressPairs(in []networkingv1beta1.AddressPair) []ports.AddressPair {
	if in == nil {
		return nil
	}
	out := make([]ports.AddressPair, 0, len(in))
	for _, a := range in {
		out = append(out, ports.AddressPair{
			IPAddress:  a.IPAddress,
			MACAddress: a.MACAddress,
		})
	}
	return out
}

func toAddressPairStatus(in []ports.AddressPair) []networkingv1beta1.AddressPair {
	if in == nil {
		return nil
	}
	out := make([]networkingv1beta1.AddressPair, 0, len(in))
	for _, a := range in {
		out = append(out, networkingv1beta1.AddressPair{
			IPAddress:  a.IPAddress,
			MACAddress: a.MACAddress,
		})
	}
	return out
}

func ptrAddressPairs(in []ports.AddressPair) *[]ports.AddressPair {
	return &in
}
