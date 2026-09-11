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

type ImageParameters struct {
	Name            string            `json:"name"`
	ContainerFormat string            `json:"containerFormat"`
	DiskFormat      string            `json:"diskFormat"`
	MinDisk         *int              `json:"minDisk,omitempty"`
	MinRAM          *int              `json:"minRam,omitempty"`
	Visibility      string            `json:"visibility"`
	Protected       *bool             `json:"protected,omitempty"`
	Tags            []string          `json:"tags,omitempty"`
	Properties      map[string]string `json:"properties,omitempty"`
}

type ImageObservation struct {
	ImageID         string   `json:"imageId,omitempty"`
	Status          string   `json:"status,omitempty"`
	Size            int64    `json:"size,omitempty"`
	MinDisk         int      `json:"minDisk,omitempty"`
	MinRAM          int      `json:"minRam,omitempty"`
	Protected       *bool    `json:"protected,omitempty"`
	Visibility      string   `json:"visibility,omitempty"`
	ContainerFormat string   `json:"containerFormat,omitempty"`
	DiskFormat      string   `json:"diskFormat,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Owner           string   `json:"owner,omitempty"`
	File            string   `json:"file,omitempty"`
	Schema          string   `json:"schema,omitempty"`
	Checksum        string   `json:"checksum,omitempty"`
	VirtualSize     int64    `json:"virtualSize,omitempty"`
}

type ImageSpec struct {
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider              ImageParameters `json:"forProvider"`
}

type ImageStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
	AtProvider             ImageObservation `json:"atProvider,omitempty"`
}

type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ImageSpec   `json:"spec"`
	Status            ImageStatus `json:"status,omitempty"`
}

type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}
