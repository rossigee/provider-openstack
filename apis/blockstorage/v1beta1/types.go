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

type VolumeParameters struct {
	Size             int               `json:"size"`
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	AvailabilityZone string            `json:"availabilityZone,omitempty"`
	VolumeType       string            `json:"volumeType,omitempty"`
	ImageID          string            `json:"imageId,omitempty"`
	SnapshotID       string            `json:"snapshotId,omitempty"`
	SourceVolID      string            `json:"sourceVolId,omitempty"`
	SourceReplica    string            `json:"sourceReplica,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	Bootable         *bool             `json:"bootable,omitempty"`
	Multiattach      *bool             `json:"multiattach,omitempty"`
}

type VolumeObservation struct {
	ID               string       `json:"id,omitempty"`
	Status           string       `json:"status,omitempty"`
	Size             int          `json:"size,omitempty"`
	AvailabilityZone string       `json:"availabilityZone,omitempty"`
	VolumeType       string       `json:"volumeType,omitempty"`
	CreatedAt        *metav1.Time `json:"createdAt,omitempty"`
	AttachedDevices  []string     `json:"attachedDevices,omitempty"`
}

type VolumeSpec struct {
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider              VolumeParameters `json:"forProvider"`
}

type VolumeStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
	AtProvider             VolumeObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,openstack}
// +genclient
// +groupName=blockstorage.openstack.m.crossplane.io

type Volume struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              VolumeSpec   `json:"spec"`
	Status            VolumeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type VolumeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Volume `json:"items"`
}
