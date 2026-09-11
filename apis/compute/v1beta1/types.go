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

type ServerParameters struct {
	Name               string               `json:"name"`
	ImageRef           string               `json:"imageRef"`
	FlavorRef          string               `json:"flavorRef"`
	KeyName            string               `json:"keyName,omitempty"`
	Networks           []ServerNetwork      `json:"networks,omitempty"`
	SecurityGroups     []string             `json:"securityGroups,omitempty"`
	BlockDeviceMapping []BlockDeviceMapping `json:"blockDeviceMapping,omitempty"`
	UserData           string               `json:"userData,omitempty"`
	TenantID           string               `json:"tenantId,omitempty"`
	AvailabilityZone   string               `json:"availabilityZone,omitempty"`
	Description        string               `json:"description,omitempty"`
	Tags               []string             `json:"tags,omitempty"`
	ConfigDrive        *bool                `json:"configDrive,omitempty"`
	AdminStateUp       *bool                `json:"adminStateUp,omitempty"`
}

type ServerNetwork struct {
	UUID    string `json:"uuid,omitempty"`
	Name    string `json:"name,omitempty"`
	FixedIP string `json:"fixedIP,omitempty"`
	Port    string `json:"port,omitempty"`
}

type BlockDeviceMapping struct {
	UUID                string `json:"uuid,omitempty"`
	SourceType          string `json:"sourceType,omitempty"`
	DestinationType     string `json:"destinationType,omitempty"`
	VolumeSize          *int   `json:"volumeSize,omitempty"`
	VolumeType          string `json:"volumeType,omitempty"`
	DeleteOnTermination *bool  `json:"deleteOnTermination,omitempty"`
	BootIndex           *int   `json:"bootIndex,omitempty"`
	GuestFormat         string `json:"guestFormat,omitempty"`
	DeviceName          string `json:"deviceName,omitempty"`
}

type ServerObservation struct {
	ServerID   string          `json:"serverId,omitempty"`
	Status     string          `json:"status,omitempty"`
	TenantID   string          `json:"tenantId,omitempty"`
	HostID     string          `json:"hostId,omitempty"`
	Flavor     string          `json:"flavor,omitempty"`
	Image      string          `json:"image,omitempty"`
	Addresses  []ServerAddress `json:"addresses,omitempty"`
	PowerState int             `json:"powerState,omitempty"`
	VMState    string          `json:"vmState,omitempty"`
	KeyName    string          `json:"keyName,omitempty"`
	Tags       []string        `json:"tags,omitempty"`
	CreatedAt  *metav1.Time    `json:"createdAt,omitempty"`
	UpdatedAt  *metav1.Time    `json:"updatedAt,omitempty"`
}

type ServerAddress struct {
	Network string `json:"network"`
	Version int    `json:"version"`
	Address string `json:"address"`
}

type ServerSpec struct {
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider              ServerParameters `json:"forProvider"`
}

type ServerStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
	AtProvider             ServerObservation       `json:"atProvider,omitempty"`
	ConnectionDetails      ServerConnectionDetails `json:"connectionDetails,omitempty"`
}

type ServerConnectionDetails struct {
	PrivateIPv4 string `json:"privateIPv4,omitempty"`
	PublicIPv4  string `json:"publicIPv4,omitempty"`
	ServerState string `json:"serverState,omitempty"`
}

type Server struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ServerSpec   `json:"spec"`
	Status            ServerStatus `json:"status,omitempty"`
}

type ServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Server `json:"items"`
}
