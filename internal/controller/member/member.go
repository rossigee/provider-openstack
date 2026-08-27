package member

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
	errNotMemberResource = "managed resource is not a Member"
	errNewClient         = "cannot create OpenStack client"

	errGetMember    = "cannot get Member"
	errCreateMember = "cannot create Member"
	errUpdateMember = "cannot update Member"
	errDeleteMember = "cannot delete Member"
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
	return &openstackMemberClient{
		client:   lbClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackMemberClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.MemberGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.MemberGroupKind)),
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
		For(&v1alpha1.Member{}).
		Complete(r)
}

func (e *openstackMemberClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Member)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotMemberResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	member, err := pools.GetMember(ctx, e.client, cr.Spec.ForProvider.PoolID, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetMember, err)
	}

	cr.Status.AtProvider = v1alpha1.MemberProviderStatus{
		MemberID:           member.ID,
		PoolID:             member.PoolID,
		Name:               member.Name,
		Address:            member.Address,
		ProtocolPort:       member.ProtocolPort,
		Weight:             member.Weight,
		SubnetID:           member.SubnetID,
		ProjectID:          member.ProjectID,
		AdminStateUp:       member.AdminStateUp,
		ProvisioningStatus: member.ProvisioningStatus,
		OperatingStatus:    member.OperatingStatus,
		Backup:             member.Backup,
		MonitorAddress:     member.MonitorAddress,
		MonitorPort:        member.MonitorPort,
		Tags:               member.Tags,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackMemberClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Member)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotMemberResource)
	}

	createOpts := pools.CreateMemberOpts{
		Address:        cr.Spec.ForProvider.Address,
		ProtocolPort:   cr.Spec.ForProvider.ProtocolPort,
		Name:           cr.Spec.ForProvider.Name,
		Weight:         cr.Spec.ForProvider.Weight,
		SubnetID:       cr.Spec.ForProvider.SubnetID,
		AdminStateUp:   cr.Spec.ForProvider.AdminStateUp,
		ProjectID:      cr.Spec.ForProvider.ProjectID,
		Backup:         cr.Spec.ForProvider.Backup,
		MonitorAddress: cr.Spec.ForProvider.MonitorAddress,
		MonitorPort:    cr.Spec.ForProvider.MonitorPort,
		Tags:           cr.Spec.ForProvider.Tags,
	}

	member, err := pools.CreateMember(ctx, e.client, cr.Spec.ForProvider.PoolID, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateMember, err)
	}

	meta.SetExternalName(cr, member.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackMemberClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Member)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotMemberResource)
	}

	updateOpts := pools.UpdateMemberOpts{
		Name:           &cr.Spec.ForProvider.Name,
		Weight:         cr.Spec.ForProvider.Weight,
		AdminStateUp:   cr.Spec.ForProvider.AdminStateUp,
		Backup:         cr.Spec.ForProvider.Backup,
		MonitorAddress: &cr.Spec.ForProvider.MonitorAddress,
		MonitorPort:    cr.Spec.ForProvider.MonitorPort,
		Tags:           cr.Spec.ForProvider.Tags,
	}

	_, err := pools.UpdateMember(ctx, e.client, cr.Spec.ForProvider.PoolID, meta.GetExternalName(cr), updateOpts).Extract()
	if err != nil {
		return managed.ExternalUpdate{}, fmt.Errorf("%s: %w", errUpdateMember, err)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *openstackMemberClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Member)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotMemberResource)
	}

	result := pools.DeleteMember(ctx, e.client, cr.Spec.ForProvider.PoolID, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackMemberClient) Disconnect(ctx context.Context) error {
	return nil
}
