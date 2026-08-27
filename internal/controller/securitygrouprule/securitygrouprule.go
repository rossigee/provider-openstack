package securitygrouprule

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
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
	errNotSecurityGroupRuleResource = "managed resource is not a SecurityGroupRule"
	errTrackUsage                   = "cannot track ProviderConfig usage"
	errGetProviderConfig            = "cannot get ProviderConfig"
	errGetCredentials               = "cannot get credentials"
	errNewClient                    = "cannot create OpenStack client"

	errGetSecurityGroupRule    = "cannot get SecurityGroupRule"
	errCreateSecurityGroupRule = "cannot create SecurityGroupRule"
	errUpdateSecurityGroupRule = "cannot update SecurityGroupRule"
	errDeleteSecurityGroupRule = "cannot delete SecurityGroupRule"
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
	return &openstackSecurityGroupRuleClient{
		client:   networkClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackSecurityGroupRuleClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.SecurityGroupRuleGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.SecurityGroupRuleGroupKind)),
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
		For(&v1alpha1.SecurityGroupRule{}).
		Complete(r)
}

func (e *openstackSecurityGroupRuleClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.SecurityGroupRule)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotSecurityGroupRuleResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	rule, err := rules.Get(ctx, e.client, meta.GetExternalName(cr)).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetSecurityGroupRule, err)
	}

	portRangeMin := rule.PortRangeMin
	portRangeMax := rule.PortRangeMax
	cr.Status.AtProvider = v1alpha1.SecurityGroupRuleProviderStatus{
		RuleID:          rule.ID,
		SecurityGroupID: rule.SecGroupID,
		Direction:       rule.Direction,
		Ethertype:       rule.EtherType,
		Protocol:        rule.Protocol,
		PortRangeMin:    &portRangeMin,
		PortRangeMax:    &portRangeMax,
		RemoteIPPrefix:  rule.RemoteIPPrefix,
		RemoteGroupID:   rule.RemoteGroupID,
		TenantID:        rule.TenantID,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackSecurityGroupRuleClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.SecurityGroupRule)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotSecurityGroupRuleResource)
	}

	createOpts := rules.CreateOpts{
		SecGroupID:     cr.Spec.ForProvider.SecurityGroupID,
		Direction:      rules.RuleDirection(cr.Spec.ForProvider.Direction),
		EtherType:      rules.RuleEtherType(cr.Spec.ForProvider.Ethertype),
		Protocol:       rules.RuleProtocol(cr.Spec.ForProvider.Protocol),
		RemoteIPPrefix: cr.Spec.ForProvider.RemoteIPPrefix,
		RemoteGroupID:  cr.Spec.ForProvider.RemoteGroupID,
	}

	if cr.Spec.ForProvider.PortRangeMin != nil {
		createOpts.PortRangeMin = *cr.Spec.ForProvider.PortRangeMin
	}
	if cr.Spec.ForProvider.PortRangeMax != nil {
		createOpts.PortRangeMax = *cr.Spec.ForProvider.PortRangeMax
	}

	rule, err := rules.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateSecurityGroupRule, err)
	}

	meta.SetExternalName(cr, rule.ID)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackSecurityGroupRuleClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}

func (e *openstackSecurityGroupRuleClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.SecurityGroupRule)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotSecurityGroupRuleResource)
	}

	result := rules.Delete(ctx, e.client, meta.GetExternalName(cr))
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackSecurityGroupRuleClient) Disconnect(ctx context.Context) error {
	return nil
}
