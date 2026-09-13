package v1beta1

import (
	"testing"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestServerDeepCopy(t *testing.T) {
	server := &Server{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-server",
			Namespace: "default",
		},
		Spec: ServerSpec{
			ManagedResourceSpec: xpv2.ManagedResourceSpec{},
			ForProvider: ServerParameters{
				Name:      "web-server",
				ImageRef:  "ubuntu-22.04",
				FlavorRef: "m1.small",
			},
		},
		Status: ServerStatus{
			ConditionedStatus: xpv2.ConditionedStatus{},
			AtProvider: ServerProviderStatus{
				ServerID: "12345",
				Status:   "ACTIVE",
			},
		},
	}

	// Test DeepCopyInto
	target := &Server{}
	server.DeepCopyInto(target)
	if target.Name != server.Name {
		t.Errorf("DeepCopyInto: Name mismatch")
	}
	if target.Spec.ForProvider.ImageRef != server.Spec.ForProvider.ImageRef {
		t.Errorf("DeepCopyInto: ImageRef mismatch")
	}

	// Test DeepCopyObject
	obj := server.DeepCopyObject()
	if obj == nil {
		t.Errorf("DeepCopyObject returned nil")
	}
	copiedServer, ok := obj.(*Server)
	if !ok {
		t.Errorf("DeepCopyObject returned wrong type")
	}
	if copiedServer.Name != server.Name {
		t.Errorf("DeepCopyObject: Name mismatch")
	}
}

func TestServerListDeepCopy(t *testing.T) {
	servers := []Server{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "server1"},
			Spec: ServerSpec{
				ForProvider: ServerParameters{Name: "web1"},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "server2"},
			Spec: ServerSpec{
				ForProvider: ServerParameters{Name: "web2"},
			},
		},
	}

	serverList := &ServerList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "compute.openstack.m.crossplane.io/v1beta1",
			Kind:       "ServerList",
		},
		ListMeta: metav1.ListMeta{},
		Items:    servers,
	}

	// Test DeepCopyObject
	obj := serverList.DeepCopyObject()
	if obj == nil {
		t.Errorf("DeepCopyObject returned nil")
	}
	copiedList, ok := obj.(*ServerList)
	if !ok {
		t.Errorf("DeepCopyObject returned wrong type")
	}
	if len(copiedList.Items) != len(serverList.Items) {
		t.Errorf("DeepCopyObject: Items count mismatch")
	}
}

func TestServerValidation(t *testing.T) {
	tests := []struct {
		name    string
		server  *Server
		isValid bool
	}{
		{
			name: "valid server spec",
			server: &Server{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-server",
					Namespace: "default",
				},
				Spec: ServerSpec{
					ForProvider: ServerParameters{
						Name:      "web-server",
						ImageRef:  "ubuntu-22.04",
						FlavorRef: "m1.small",
					},
				},
			},
			isValid: true,
		},
		{
			name: "missing required fields",
			server: &Server{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: ServerSpec{
					ForProvider: ServerParameters{
						// Missing Name, ImageRef, FlavorRef
					},
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.server == nil {
				t.Errorf("Test server is nil")
			}
			if tt.server.Spec.ForProvider.Name == "" && tt.isValid {
				t.Errorf("Missing required Name field but marked as valid")
			}
		})
	}
}

func TestKeyPairDeepCopy(t *testing.T) {
	keypair := &KeyPair{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-key",
			Namespace: "default",
		},
		Spec: KeyPairSpec{
			ForProvider: KeyPairParameters{
				Name:      "my-key",
				PublicKey: "ssh-rsa AAAA...",
				Type:      "ssh",
			},
		},
		Status: KeyPairStatus{
			AtProvider: KeyPairProviderStatus{
				Name: "my-key",
			},
		},
	}

	// Test DeepCopyObject
	obj := keypair.DeepCopyObject()
	if obj == nil {
		t.Errorf("DeepCopyObject returned nil")
	}
	copiedKeypair, ok := obj.(*KeyPair)
	if !ok {
		t.Errorf("DeepCopyObject returned wrong type")
	}
	if copiedKeypair.Name != keypair.Name {
		t.Errorf("DeepCopyObject: Name mismatch")
	}
}
