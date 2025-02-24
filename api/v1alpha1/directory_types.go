/*
Copyright 2025.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// directoryAvailableCondition represents the current status of a directory
	DirectoryAvailableCondition = "Available"
	// directoryDegradedCondition represents the status of a directory while resources are being deleted
	DirectoryDegradedCondition = "Degraded"
)

// DirectorySpec defines the desired state of Directory.
type DirectorySpec struct {
}

// DirectoryStatus defines the observed state of Directory.
type DirectoryStatus struct {
	// Slice of conditions storing the condition of the directory
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Directory is the Schema for the directories API.
type Directory struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DirectorySpec   `json:"spec,omitempty"`
	Status DirectoryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DirectoryList contains a list of Directory.
type DirectoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Directory `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Directory{}, &DirectoryList{})
}
