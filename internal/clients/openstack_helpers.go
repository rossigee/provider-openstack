package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/rossigee/provider-openstack/apis/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
)

const (
	errNoProviderConfig     = "providerConfigRef is not set"
	errGetProviderConfig    = "cannot get ProviderConfig"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal credentials"
)

// GetProviderConfigReference extracts the ProviderConfig reference from a
// managed resource that embeds xpv2.ClusterManagedResourceSpec.
// This is needed because resource.Managed in crossplane-runtime v2 does not
// expose GetProviderConfigReference().
type providerConfigRefGetter interface {
	GetProviderConfigReference() *xpv2.Reference
}

// NewOpenStackClient creates a new gophercloud authenticated client from the
// ProviderConfig referenced by the given managed resource.
func NewOpenStackClient(ctx context.Context, kube client.Client, mg resource.Managed) (*gophercloud.ProviderClient, error) {
	pcg, ok := mg.(providerConfigRefGetter)
	if !ok {
		return nil, fmt.Errorf("managed resource does not implement GetProviderConfigReference")
	}
	configRef := pcg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}
	pc := &v1beta1.ProviderConfig{}
	if err := kube.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
		return nil, fmt.Errorf("%s: %w", errGetProviderConfig, err)
	}

	data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, kube, pc.Spec.Credentials.CommonCredentialSelectors)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errExtractCredentials, err)
	}
	creds := map[string]string{}
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalCredentials, err)
	}

	ao := gophercloud.AuthOptions{
		IdentityEndpoint:            creds["auth_url"],
		Username:                    creds["user_name"],
		UserID:                      creds["user_id"],
		Password:                    creds["password"],
		TenantID:                    creds["tenant_id"],
		TenantName:                  creds["tenant_name"],
		DomainID:                    creds["domain_id"],
		DomainName:                  creds["domain_name"],
		ApplicationCredentialID:     creds["application_credential_id"],
		ApplicationCredentialName:   creds["application_credential_name"],
		ApplicationCredentialSecret: creds["application_credential_secret"],
		TokenID:                     creds["token"],
		Scope:                       getScope(creds),
	}

	oc, err := openstack.AuthenticatedClient(ctx, ao)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with OpenStack: %w", err)
	}

	return oc, nil
}

func getScope(creds map[string]string) *gophercloud.AuthScope {
	if creds["system_scope"] == "true" {
		return &gophercloud.AuthScope{System: true}
	}
	if id := creds["project_domain_id"]; id != "" {
		return &gophercloud.AuthScope{ProjectID: creds["tenant_id"], DomainID: id}
	}
	if name := creds["project_domain_name"]; name != "" {
		return &gophercloud.AuthScope{ProjectName: creds["tenant_name"], DomainName: name}
	}
	if creds["tenant_name"] != "" {
		return &gophercloud.AuthScope{ProjectName: creds["tenant_name"]}
	}
	return nil
}

// NewNetworkClient creates a new Neutron v2 network client.
func NewNetworkClient(ctx context.Context, kube client.Client, mg resource.Managed) (*gophercloud.ServiceClient, error) {
	pc, err := NewOpenStackClient(ctx, kube, mg)
	if err != nil {
		return nil, err
	}
	return openstack.NewNetworkV2(pc, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

// NewLoadBalancerClient creates a new Octavia v2 load balancer client.
func NewLoadBalancerClient(ctx context.Context, kube client.Client, mg resource.Managed) (*gophercloud.ServiceClient, error) {
	pc, err := NewOpenStackClient(ctx, kube, mg)
	if err != nil {
		return nil, err
	}
	return openstack.NewLoadBalancerV2(pc, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

// NewComputeClient creates a new Nova v2 compute client.
func NewComputeClient(ctx context.Context, kube client.Client, mg resource.Managed) (*gophercloud.ServiceClient, error) {
	pc, err := NewOpenStackClient(ctx, kube, mg)
	if err != nil {
		return nil, err
	}
	return openstack.NewComputeV2(pc, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

// NewBlockStorageClient creates a new Cinder v3 block storage client.
func NewBlockStorageClient(ctx context.Context, kube client.Client, mg resource.Managed) (*gophercloud.ServiceClient, error) {
	pc, err := NewOpenStackClient(ctx, kube, mg)
	if err != nil {
		return nil, err
	}
	return openstack.NewBlockStorageV3(pc, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

// NewImageClient creates a new Glance v2 image client.
func NewImageClient(ctx context.Context, kube client.Client, mg resource.Managed) (*gophercloud.ServiceClient, error) {
	pc, err := NewOpenStackClient(ctx, kube, mg)
	if err != nil {
		return nil, err
	}
	return openstack.NewImageV2(pc, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

// NewIdentityClient creates a new Keystone v3 identity client.
func NewIdentityClient(ctx context.Context, kube client.Client, mg resource.Managed) (*gophercloud.ServiceClient, error) {
	pc, err := NewOpenStackClient(ctx, kube, mg)
	if err != nil {
		return nil, err
	}
	return openstack.NewIdentityV3(pc, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

// NewDNSClient creates a new Designate v2 DNS client.
func NewDNSClient(ctx context.Context, kube client.Client, mg resource.Managed) (*gophercloud.ServiceClient, error) {
	pc, err := NewOpenStackClient(ctx, kube, mg)
	if err != nil {
		return nil, err
	}
	return openstack.NewDNSV2(pc, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

// ExternalName returns the external name annotation value.
func ExternalName(mg resource.Managed) string {
	if en := mg.GetAnnotations()["crossplane.io/external-name"]; en != "" {
		return en
	}
	return mg.GetName()
}

// ExternalNameOrDefault returns the external name or the provided default.
func ExternalNameOrDefault(mg resource.Managed, def string) string {
	if en := mg.GetAnnotations()["crossplane.io/external-name"]; en != "" {
		return en
	}
	if def != "" {
		return def
	}
	return mg.GetName()
}
