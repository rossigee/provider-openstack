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

type LoadBalancerParameters struct {
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	VipSubnetID      string   `json:"vipSubnetId,omitempty"`
	VipNetworkID     string   `json:"vipNetworkId,omitempty"`
	VipPortID        string   `json:"vipPortId,omitempty"`
	VipAddress       string   `json:"vipAddress,omitempty"`
	VipQosPolicyID   string   `json:"vipQosPolicyId,omitempty"`
	AdminStateUp     *bool    `json:"adminStateUp,omitempty"`
	ProjectID        string   `json:"projectId,omitempty"`
	FlavorID         string   `json:"flavorId,omitempty"`
	AvailabilityZone string   `json:"availabilityZone,omitempty"`
	Provider         string   `json:"provider,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

type LoadBalancerObservation struct {
	LoadBalancerID     string       `json:"loadBalancerId,omitempty"`
	Name               string       `json:"name,omitempty"`
	Description        string       `json:"description,omitempty"`
	ProvisioningStatus string       `json:"provisioningStatus,omitempty"`
	OperatingStatus    string       `json:"operatingStatus,omitempty"`
	VipAddress         string       `json:"vipAddress,omitempty"`
	VipPortID          string       `json:"vipPortId,omitempty"`
	VipSubnetID        string       `json:"vipSubnetId,omitempty"`
	VipNetworkID       string       `json:"vipNetworkId,omitempty"`
	ProjectID          string       `json:"projectId,omitempty"`
	AdminStateUp       bool         `json:"adminStateUp,omitempty"`
	FlavorID           string       `json:"flavorId,omitempty"`
	Provider           string       `json:"provider,omitempty"`
	Tags               []string     `json:"tags,omitempty"`
	CreatedAt          *metav1.Time `json:"createdAt,omitempty"`
	UpdatedAt          *metav1.Time `json:"updatedAt,omitempty"`
}

type LoadBalancerSpec struct {
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider              LoadBalancerParameters `json:"forProvider"`
}

type LoadBalancerStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
	AtProvider             LoadBalancerObservation `json:"atProvider,omitempty"`
}

type LoadBalancer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              LoadBalancerSpec   `json:"spec"`
	Status            LoadBalancerStatus `json:"status,omitempty"`
}

type LoadBalancerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoadBalancer `json:"items"`
}
