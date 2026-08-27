package keypair

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
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
	errNotKeyPairResource = "managed resource is not a KeyPair"
	errTrackUsage         = "cannot track ProviderConfig usage"
	errGetProviderConfig  = "cannot get ProviderConfig"
	errGetCredentials     = "cannot get credentials"
	errNewClient          = "cannot create OpenStack client"

	errGetKeyPair    = "cannot get KeyPair"
	errCreateKeyPair = "cannot create KeyPair"
	errUpdateKeyPair = "cannot update KeyPair"
	errDeleteKeyPair = "cannot delete KeyPair"
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
	return &openstackKeyPairClient{
		client:   computeClient,
		kube:     e.kube,
		recorder: e.recorder,
	}, nil
}

type openstackKeyPairClient struct {
	client   *gophercloud.ServiceClient
	kube     client.Client
	recorder event.Recorder
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.KeyPairGroupKind)
	rec := event.NewNopRecorder()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.KeyPairGroupKind)),
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
		For(&v1alpha1.KeyPair{}).
		Complete(r)
}

func (e *openstackKeyPairClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.KeyPair)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotKeyPairResource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	kp, err := keypairs.Get(ctx, e.client, meta.GetExternalName(cr), keypairs.GetOpts{}).Extract()
	if gophercloud.ResponseCodeIs(err, 404) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	if err != nil {
		return managed.ExternalObservation{}, fmt.Errorf("%s: %w", errGetKeyPair, err)
	}

	cr.Status.AtProvider = v1alpha1.KeyPairProviderStatus{
		Name:        kp.Name,
		PublicKey:   kp.PublicKey,
		Fingerprint: kp.Fingerprint,
		Type:        kp.Type,
		TenantID:    kp.UserID,
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *openstackKeyPairClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.KeyPair)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotKeyPairResource)
	}

	createOpts := keypairs.CreateOpts{
		Name:      cr.Spec.ForProvider.Name,
		PublicKey: cr.Spec.ForProvider.PublicKey,
		Type:      cr.Spec.ForProvider.Type,
	}

	kp, err := keypairs.Create(ctx, e.client, createOpts).Extract()
	if err != nil {
		return managed.ExternalCreation{}, fmt.Errorf("%s: %w", errCreateKeyPair, err)
	}

	meta.SetExternalName(cr, kp.Name)

	return managed.ExternalCreation{
		ConnectionDetails: nil,
	}, nil
}

func (e *openstackKeyPairClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}

func (e *openstackKeyPairClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.KeyPair)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotKeyPairResource)
	}

	result := keypairs.Delete(ctx, e.client, meta.GetExternalName(cr), keypairs.DeleteOpts{})
	return managed.ExternalDelete{}, result.Err
}

func (e *openstackKeyPairClient) Disconnect(ctx context.Context) error {
	return nil
}
