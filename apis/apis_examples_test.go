package apis

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestExampleYAMLsAreValid(t *testing.T) {
	examplesDir := "../examples"
	scheme := runtime.NewScheme()
	AddToScheme(scheme)

	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Skipf("examples directory not found: %v", err)
	}

	exampleCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			subdir := filepath.Join(examplesDir, entry.Name())
			subentries, err := os.ReadDir(subdir)
			if err != nil {
				t.Logf("failed to read subdir %s: %v", subdir, err)
				continue
			}

			for _, subfile := range subentries {
				if !subfile.IsDir() && strings.HasSuffix(subfile.Name(), ".yaml") {
					path := filepath.Join(subdir, subfile.Name())
					exampleCount++

					data, err := os.ReadFile(path)
					if err != nil {
						t.Errorf("failed to read %s: %v", path, err)
						continue
					}

					// Parse YAML as unstructured object
					obj := &unstructured.Unstructured{}
					err = yaml.Unmarshal(data, obj)
					if err != nil {
						t.Errorf("failed to parse YAML %s: %v", path, err)
						continue
					}

					// Validate required fields
					if obj.GetAPIVersion() == "" {
						t.Errorf("%s: missing apiVersion", path)
					}
					if obj.GetKind() == "" {
						t.Errorf("%s: missing kind", path)
					}
					if obj.GetName() == "" {
						t.Errorf("%s: missing metadata.name", path)
					}

					// Validate API group is v1beta1
					apiVersion := obj.GetAPIVersion()
					if !strings.Contains(apiVersion, "v1beta1") {
						t.Errorf("%s: invalid API version %s (expected *v1beta1)", path, apiVersion)
					}
				}
			}
		} else if strings.HasSuffix(entry.Name(), ".yaml") {
			// Top-level yaml file
			path := filepath.Join(examplesDir, entry.Name())
			exampleCount++

			data, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("failed to read %s: %v", path, err)
				continue
			}

			obj := &unstructured.Unstructured{}
			err = yaml.Unmarshal(data, obj)
			if err != nil {
				t.Errorf("failed to parse YAML %s: %v", path, err)
				continue
			}

			if obj.GetAPIVersion() == "" {
				t.Errorf("%s: missing apiVersion", path)
			}
		}
	}

	if exampleCount == 0 {
		t.Skip("no example YAML files found")
	}

	t.Logf("validated %d example YAML files", exampleCount)
}

func TestExampleYAMLsHaveCorrectAPIGroups(t *testing.T) {
	// Map of resource kinds to expected API groups
	expectedGroups := map[string]string{
		"Server":               "compute.openstack.m.crossplane.io",
		"KeyPair":              "compute.openstack.m.crossplane.io",
		"Network":              "networking.openstack.m.crossplane.io",
		"Subnet":               "networking.openstack.m.crossplane.io",
		"Router":               "networking.openstack.m.crossplane.io",
		"SecurityGroup":        "networking.openstack.m.crossplane.io",
		"FloatingIP":           "networking.openstack.m.crossplane.io",
		"Port":                 "networking.openstack.m.crossplane.io",
		"Volume":               "blockstorage.openstack.m.crossplane.io",
		"VolumeType":           "blockstorage.openstack.m.crossplane.io",
		"VolumeSnapshot":       "blockstorage.openstack.m.crossplane.io",
		"Image":                "image.openstack.m.crossplane.io",
		"User":                 "identity.openstack.m.crossplane.io",
		"Project":              "identity.openstack.m.crossplane.io",
		"Role":                 "identity.openstack.m.crossplane.io",
		"LoadBalancer":         "loadbalancing.openstack.m.crossplane.io",
		"Listener":             "loadbalancing.openstack.m.crossplane.io",
		"Pool":                 "loadbalancing.openstack.m.crossplane.io",
		"Member":               "loadbalancing.openstack.m.crossplane.io",
		"HealthMonitor":        "loadbalancing.openstack.m.crossplane.io",
		"Zone":                 "dns.openstack.m.crossplane.io",
		"ProviderConfig":       "openstack.m.crossplane.io",
	}

	examplesDir := "../examples"
	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Skipf("examples directory not found: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subdir := filepath.Join(examplesDir, entry.Name())
			subentries, err := os.ReadDir(subdir)
			if err != nil {
				continue
			}

			for _, subfile := range subentries {
				if !subfile.IsDir() && strings.HasSuffix(subfile.Name(), ".yaml") {
					path := filepath.Join(subdir, subfile.Name())
					data, err := os.ReadFile(path)
					if err != nil {
						continue
					}

					obj := &unstructured.Unstructured{}
					yaml.Unmarshal(data, obj)

					kind := obj.GetKind()
					apiVersion := obj.GetAPIVersion()

					if expectedGroup, found := expectedGroups[kind]; found {
						if !strings.Contains(apiVersion, expectedGroup) {
							t.Errorf("%s (kind: %s): expected API group %s, got %s",
								path, kind, expectedGroup, apiVersion)
						}
						if !strings.Contains(apiVersion, "v1beta1") {
							t.Errorf("%s (kind: %s): expected v1beta1, got %s",
								path, kind, apiVersion)
						}
					}
				}
			}
		}
	}
}
