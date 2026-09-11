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

type NetworkParameters struct {
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	AdminStateUp        *bool    `json:"adminStateUp,omitempty"`
	TenantID            string   `json:"tenantId,omitempty"`
	Shared              *bool    `json:"shared,omitempty"`
	DNSDomain           string   `json:"dnsDomain,omitempty"`
	PortSecurityEnabled *bool    `json:"portSecurityEnabled,omitempty"`
	Tags                []string `json:"tags,omitempty"`
}

type NetworkObservation struct {
	NetworkID               string       `json:"networkId,omitempty"`
	Status                  string       `json:"status,omitempty"`
	Subnets                 []string     `json:"subnets,omitempty"`
	TenantID                string       `json:"tenantId,omitempty"`
	AdminStateUp            bool         `json:"adminStateUp,omitempty"`
	Shared                  bool         `json:"shared,omitempty"`
	ProviderNetworkType     string       `json:"providerNetworkType,omitempty"`
	ProviderPhysicalNetwork string       `json:"providerPhysicalNetwork,omitempty"`
	ProviderSegmentationID  int          `json:"providerSegmentationId,omitempty"`
	PortSecurityEnabled     bool         `json:"portSecurityEnabled,omitempty"`
	DNSDomain               string       `json:"dnsDomain,omitempty"`
	RevisionNumber          int          `json:"revisionNumber,omitempty"`
	Tags                    []string     `json:"tags,omitempty"`
	CreatedAt               *metav1.Time `json:"createdAt,omitempty"`
	UpdatedAt               *metav1.Time `json:"updatedAt,omitempty"`
}

type NetworkSpec struct {
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider              NetworkParameters `json:"forProvider"`
}

type NetworkStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
	AtProvider             NetworkObservation `json:"atProvider,omitempty"`
	ConnectionDetails      []byte             `json:"connectionDetails,omitempty"`
}

type Network struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NetworkSpec   `json:"spec"`
	Status            NetworkStatus `json:"status,omitempty"`
}

type NetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Network `json:"items"`
}

type PortParameters struct {
	NetworkID           string        `json:"networkId"`
	Name                string        `json:"name,omitempty"`
	Description         string        `json:"description,omitempty"`
	AdminStateUp        *bool         `json:"adminStateUp,omitempty"`
	MACAddress          string        `json:"macAddress,omitempty"`
	FixedIPs            []FixedIP     `json:"fixedIPs,omitempty"`
	DeviceID            string        `json:"deviceId,omitempty"`
	DeviceOwner         string        `json:"deviceOwner,omitempty"`
	TenantID            string        `json:"tenantId,omitempty"`
	SecurityGroups      *[]string     `json:"securityGroups,omitempty"`
	AllowedAddressPairs []AddressPair `json:"allowedAddressPairs,omitempty"`
}

type FixedIP struct {
	SubnetID  string `json:"subnetId,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
}

type AddressPair struct {
	IPAddress  string `json:"ipAddress,omitempty"`
	MACAddress string `json:"macAddress,omitempty"`
}

type PortObservation struct {
	PortID              string        `json:"portId,omitempty"`
	NetworkID           string        `json:"networkId,omitempty"`
	Name                string        `json:"name,omitempty"`
	Description         string        `json:"description,omitempty"`
	AdminStateUp        bool          `json:"adminStateUp,omitempty"`
	Status              string        `json:"status,omitempty"`
	MACAddress          string        `json:"macAddress,omitempty"`
	FixedIPs            []FixedIP     `json:"fixedIPs,omitempty"`
	TenantID            string        `json:"tenantId,omitempty"`
	DeviceOwner         string        `json:"deviceOwner,omitempty"`
	DeviceID            string        `json:"deviceId,omitempty"`
	SecurityGroups      []string      `json:"securityGroups,omitempty"`
	AllowedAddressPairs []AddressPair `json:"allowedAddressPairs,omitempty"`
}

type PortSpec struct {
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider              PortParameters `json:"forProvider"`
}

type PortStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
	AtProvider             PortObservation `json:"atProvider,omitempty"`
}

type Port struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PortSpec   `json:"spec"`
	Status            PortStatus `json:"status,omitempty"`
}

type PortList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Port `json:"items"`
}
