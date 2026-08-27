package apis

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestSchemeRegistration(t *testing.T) {
	s := runtime.NewScheme()
	require.NoError(t, AddToScheme(s))

	kinds := []struct {
		apiVersion string
		kind       string
	}{
		// v1beta1 - ProviderConfig
		{"openstack.crossplane.io/v1beta1", "ProviderConfig"},
		{"openstack.crossplane.io/v1beta1", "ProviderConfigUsage"},

		// networking v1alpha1
		{"openstack.crossplane.io/v1alpha1", "Network"},
		{"openstack.crossplane.io/v1alpha1", "Subnet"},
		{"openstack.crossplane.io/v1alpha1", "Router"},
		{"openstack.crossplane.io/v1alpha1", "RouterInterface"},
		{"openstack.crossplane.io/v1alpha1", "SecurityGroup"},
		{"openstack.crossplane.io/v1alpha1", "SecurityGroupRule"},
		{"openstack.crossplane.io/v1alpha1", "FloatingIP"},
		{"openstack.crossplane.io/v1alpha1", "Port"},
		{"openstack.crossplane.io/v1alpha1", "SubnetPool"},
		{"openstack.crossplane.io/v1alpha1", "Trunk"},
		{"openstack.crossplane.io/v1alpha1", "RBACPolicy"},

		// blockstorage v1alpha1
		{"openstack.crossplane.io/v1alpha1", "Volume"},
		{"openstack.crossplane.io/v1alpha1", "VolumeType"},
		{"openstack.crossplane.io/v1alpha1", "VolumeSnapshot"},

		// compute v1alpha1
		{"openstack.crossplane.io/v1alpha1", "Server"},
		{"openstack.crossplane.io/v1alpha1", "KeyPair"},

		// loadbalancing v1alpha1
		{"openstack.crossplane.io/v1alpha1", "LoadBalancer"},
		{"openstack.crossplane.io/v1alpha1", "Listener"},
		{"openstack.crossplane.io/v1alpha1", "Pool"},
		{"openstack.crossplane.io/v1alpha1", "Member"},
		{"openstack.crossplane.io/v1alpha1", "HealthMonitor"},

		// image v1alpha1
		{"openstack.crossplane.io/v1alpha1", "Image"},

		// identity v1alpha1
		{"openstack.crossplane.io/v1alpha1", "Project"},
		{"openstack.crossplane.io/v1alpha1", "User"},
		{"openstack.crossplane.io/v1alpha1", "Role"},

		// dns v1alpha1
		{"openstack.crossplane.io/v1alpha1", "Zone"},
		{"openstack.crossplane.io/v1alpha1", "RecordSet"},
	}

	for _, k := range kinds {
		t.Run(k.kind, func(t *testing.T) {
			gv, err := schema.ParseGroupVersion(k.apiVersion)
			require.NoError(t, err)

			obj, err := s.New(gv.WithKind(k.kind))
			require.NoError(t, err)
			require.NotNil(t, obj)
		})
	}
}

func TestSchemeRoundTrip(t *testing.T) {
	s := runtime.NewScheme()
	require.NoError(t, AddToScheme(s))

	// Verify we can create and list all resource types
	gv := schema.GroupVersion{Group: "openstack.crossplane.io", Version: "v1alpha1"}

	listKinds := []string{
		"NetworkList", "SubnetList", "RouterList", "RouterInterfaceList",
		"SecurityGroupList", "SecurityGroupRuleList", "FloatingIPList",
		"PortList", "SubnetPoolList", "TrunkList", "RBACPolicyList",
		"VolumeList", "VolumeTypeList", "VolumeSnapshotList",
		"ServerList", "KeyPairList",
		"LoadBalancerList", "ListenerList", "PoolList", "MemberList", "HealthMonitorList",
		"ImageList",
		"ProjectList", "UserList", "RoleList",
		"ZoneList", "RecordSetList",
	}

	for _, kind := range listKinds {
		t.Run(kind, func(t *testing.T) {
			obj, err := s.New(gv.WithKind(kind))
			require.NoError(t, err)
			require.NotNil(t, obj)
		})
	}

	// Verify ProviderConfig list
	gvb := schema.GroupVersion{Group: "openstack.crossplane.io", Version: "v1beta1"}
	for _, kind := range []string{"ProviderConfigList", "ProviderConfigUsageList"} {
		t.Run(kind, func(t *testing.T) {
			obj, err := s.New(gvb.WithKind(kind))
			require.NoError(t, err)
			require.NotNil(t, obj)
		})
	}
}
