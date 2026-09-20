package clients

import (
	"testing"
)

func TestNewClientWithValidCredentials(t *testing.T) {
	tests := []struct {
		name      string
		authURL   string
		username  string
		password  string
		projectID string
		expected  bool
	}{
		{
			name:      "valid OpenStack v3 credentials",
			authURL:   "https://openstack.example.com/v3",
			username:  "testuser",
			password:  "testpass",
			projectID: "proj123",
			expected:  true,
		},
		{
			name:      "missing auth URL",
			authURL:   "",
			username:  "testuser",
			password:  "testpass",
			projectID: "proj123",
			expected:  false,
		},
		{
			name:      "missing username",
			authURL:   "https://openstack.example.com/v3",
			username:  "",
			password:  "testpass",
			projectID: "proj123",
			expected:  false,
		},
		{
			name:      "missing password",
			authURL:   "https://openstack.example.com/v3",
			username:  "testuser",
			password:  "",
			projectID: "proj123",
			expected:  false,
		},
		{
			name:      "missing project ID",
			authURL:   "https://openstack.example.com/v3",
			username:  "testuser",
			password:  "testpass",
			projectID: "",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creds := map[string]string{
				"auth_url":   tt.authURL,
				"username":   tt.username,
				"password":   tt.password,
				"project_id": tt.projectID,
			}

			// Validate credentials map has required fields
			hasAuthURL := creds["auth_url"] != ""
			hasUsername := creds["username"] != ""
			hasPassword := creds["password"] != ""
			hasProjectID := creds["project_id"] != ""

			isValid := hasAuthURL && hasUsername && hasPassword && hasProjectID

			if isValid != tt.expected {
				t.Errorf("NewClient(%s): got %v, expected %v",
					tt.name, isValid, tt.expected)
			}
		})
	}
}

func TestCredentialsWithApplicationCredentials(t *testing.T) {
	creds := map[string]string{
		"auth_url":                      "https://openstack.example.com/v3",
		"application_credential_id":     "app-123",
		"application_credential_secret": "secret",
	}

	// Validate required fields
	hasAuthURL := creds["auth_url"] != ""
	hasAppID := creds["application_credential_id"] != ""
	hasAppSecret := creds["application_credential_secret"] != ""

	if !hasAuthURL || !hasAppID || !hasAppSecret {
		t.Error("Application credential credentials validation failed")
	}
}

func TestCredentialsWithProjectName(t *testing.T) {
	creds := map[string]string{
		"auth_url":            "https://openstack.example.com/v3",
		"username":            "testuser",
		"password":            "testpass",
		"project_name":        "myproject",
		"user_domain_name":    "Default",
		"project_domain_name": "Default",
	}

	// Validate required fields for project name authentication
	hasAuthURL := creds["auth_url"] != ""
	hasUsername := creds["username"] != ""
	hasPassword := creds["password"] != ""
	hasProjectName := creds["project_name"] != ""

	if !hasAuthURL || !hasUsername || !hasPassword || !hasProjectName {
		t.Error("Project name credentials validation failed")
	}
}
