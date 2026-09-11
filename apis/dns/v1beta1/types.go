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

type ZoneParameters struct {
	Name        string            `json:"name"`
	Email       string            `json:"email,omitempty"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"type,omitempty"`
	TTL         *int              `json:"ttl,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

type ZoneObservation struct {
	ZoneID     string            `json:"zoneId,omitempty"`
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Email      string            `json:"email,omitempty"`
	TTL        int               `json:"ttl,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Status     string            `json:"status,omitempty"`
	CreatedAt  *metav1.Time      `json:"createdAt,omitempty"`
}

type ZoneSpec struct {
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider              ZoneParameters `json:"forProvider"`
}

type ZoneStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
	AtProvider             ZoneObservation `json:"atProvider,omitempty"`
}

type Zone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ZoneSpec   `json:"spec"`
	Status            ZoneStatus `json:"status,omitempty"`
}

type ZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zone `json:"items"`
}
