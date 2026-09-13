package clients

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetScope(t *testing.T) {
	tests := []struct {
		name     string
		creds    map[string]string
		wantNil  bool
		wantSys  bool
		wantProj string
		wantDom  string
	}{
		{
			name:    "nil creds",
			creds:   nil,
			wantNil: true,
		},
		{
			name:    "empty creds",
			creds:   map[string]string{},
			wantNil: true,
		},
		{
			name:    "system_scope true",
			creds:   map[string]string{"system_scope": "true"},
			wantSys: true,
		},
		{
			name:     "system_scope false does not set system",
			creds:    map[string]string{"system_scope": "false", "tenant_name": "myproject"},
			wantNil:  false,
			wantProj: "myproject",
		},
		{
			name:     "project_domain_id",
			creds:    map[string]string{"project_domain_id": "dom123", "tenant_id": "proj456"},
			wantProj: "proj456",
			wantDom:  "dom123",
		},
		{
			name:     "project_domain_name",
			creds:    map[string]string{"project_domain_name": "mydomain", "tenant_name": "myproject"},
			wantProj: "myproject",
			wantDom:  "mydomain",
		},
		{
			name:     "tenant_name only",
			creds:    map[string]string{"tenant_name": "myproject"},
			wantProj: "myproject",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := getScope(tt.creds)
			if tt.wantNil {
				require.Nil(t, scope)
				return
			}
			require.NotNil(t, scope)
			if tt.wantSys {
				require.True(t, scope.System)
				return
			}
			if tt.wantDom != "" {
				if tt.creds["project_domain_id"] != "" {
					require.Equal(t, tt.creds["project_domain_id"], scope.DomainID)
					require.Equal(t, tt.creds["tenant_id"], scope.ProjectID)
				} else {
					require.Equal(t, tt.creds["project_domain_name"], scope.DomainName)
					require.Equal(t, tt.creds["tenant_name"], scope.ProjectName)
				}
			} else {
				require.Equal(t, tt.wantProj, scope.ProjectName)
			}
		})
	}
}

func TestExternalName(t *testing.T) {
	t.Skip("Skipped after v1beta1 migration - v1alpha1 types removed")
}

func TestExternalNameOrDefault(t *testing.T) {
	t.Skip("Skipped after v1beta1 migration - v1alpha1 types removed")
}
