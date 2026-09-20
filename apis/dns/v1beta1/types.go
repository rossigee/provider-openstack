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
	"k8s.io/apimachinery/pkg/runtime"
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

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane.io/external-name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,openstack}
// Zone is a managed resource that represents an OpenStack DNS zone.
type Zone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ZoneSpec   `json:"spec"`
	Status            ZoneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// ZoneList contains a list of Zone resources.
type ZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zone `json:"items"`
}

// DeepCopyObject implements runtime.Object
func (in *Zone) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Zone)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies the receiver into out
func (in *Zone) DeepCopyInto(out *Zone) {
	*out = *in
	out.ObjectMeta = in.ObjectMeta
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopyObject implements runtime.Object
func (in *ZoneList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(ZoneList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies the receiver into out
func (in *ZoneList) DeepCopyInto(out *ZoneList) {
	*out = *in
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]Zone, len(in.Items))
		copy(out.Items, in.Items)
	}
}

// GetCondition implements resource.Managed
func (in *Zone) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return in.Status.GetCondition(ct)
}

// SetConditions implements resource.Managed
func (in *Zone) SetConditions(c ...xpv1.Condition) {
	in.Status.SetConditions(c...)
}

// GetManagementPolicies implements resource.Managed
func (in *Zone) GetManagementPolicies() xpv1.ManagementPolicies {
	return in.Spec.ManagementPolicies
}

// SetManagementPolicies implements resource.Managed
func (in *Zone) SetManagementPolicies(p xpv1.ManagementPolicies) {
	in.Spec.ManagementPolicies = p
}
