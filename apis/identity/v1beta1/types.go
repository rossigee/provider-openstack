/*
Copyright 2025 The Crossplane Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	xpv1 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type UserParameters struct {
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	DomainID         string            `json:"domainId,omitempty"`
	DefaultProjectID string            `json:"defaultProjectId,omitempty"`
	Email            string            `json:"email,omitempty"`
	Password         string            `json:"password,omitempty"`
	Enabled          *bool             `json:"enabled,omitempty"`
	Options          map[string]string `json:"options,omitempty"`
}

type UserObservation struct {
	UserID            string       `json:"userId,omitempty"`
	Name              string       `json:"name,omitempty"`
	Description       string       `json:"description,omitempty"`
	DomainID          string       `json:"domainId,omitempty"`
	DefaultProjectID  string       `json:"defaultProjectId,omitempty"`
	Enabled           *bool        `json:"enabled,omitempty"`
	PasswordExpiresAt *metav1.Time `json:"passwordExpiresAt,omitempty"`
}

type UserSpec struct {
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider              UserParameters `json:"forProvider"`
}

type UserStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
	AtProvider             UserObservation `json:"atProvider,omitempty"`
}

type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              UserSpec   `json:"spec"`
	Status            UserStatus `json:"status,omitempty"`
}

type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}

type ProjectParameters struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	DomainID    string   `json:"domainId,omitempty"`
	IsDomain    *bool    `json:"isDomain,omitempty"`
	ParentID    string   `json:"parentId,omitempty"`
	Enabled     *bool    `json:"enabled,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type ProjectObservation struct {
	ProjectID   string   `json:"projectId,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	DomainID    string   `json:"domainId,omitempty"`
	Enabled     *bool    `json:"enabled,omitempty"`
	ParentID    string   `json:"parentId,omitempty"`
	IsDomain    *bool    `json:"isDomain,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type ProjectSpec struct {
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider              ProjectParameters `json:"forProvider"`
}

type ProjectStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
	AtProvider             ProjectObservation `json:"atProvider,omitempty"`
}

type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProjectSpec   `json:"spec"`
	Status            ProjectStatus `json:"status,omitempty"`
}

type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}
