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
	"fmt"

	corev1 "k8s.io/api/core/v1"
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
	// Image to use for slapd container
	// +kubebuilder:validation:Optional
	Image string `json:"image,omitempty"`
	// Configuration for slapd daemon
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	SlapdConfig *SlapdConfigSpec `json:"slapd,omitempty"`
	// Service to create for directory
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Service *DirectoryServiceSpec `json:"service,omitempty"`
}

type SlapdConfigSpec struct {
	// Schemas to include in olcSchemaConfig
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={"core","cosine","nis"}
	Schemas []Schema `json:"schemas,omitempty"`
	// Overlays to include in olcSchemaConfig
	// +kubebuilder:validation:Optional
	Overlays []Overlay `json:"overlays,omitempty"`
	// Global frontend database configuration. Refers to olcDatabase=frontend,cn=config
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	FrontendDatabase *FrontendDatabaseConfig `json:"frontendDatabase,omitempty"`
	// Global config database configuration. Refers to olcDatabase=config,cn=config
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	ConfigDatabase *ConfigDatabaseConfig `json:"configDatabase,omitempty"`
}

// +kubebuilder:validation:Enum:=collective;cobra;core;cosine;dsee;duaconf;dyngroup;inetorgperson;java;misc;msuser;namedobject;nis;openldap;pmi
type Schema string

// +kbebuilder:validation:Enum:=accesslog;auditlog;autoca;collect;constraint;dds;deref;dyngroup;;dynlist;homedir;memberof;nestgroup;otp;pcache;ppolicy;refint;remoteauth;retcode;rwm;seqmod;sssvlv;syncprov;translucent;unique;valsort
type Overlay string

type FrontendDatabaseConfig struct {
	// Access controls for frontend database
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={"to * by * read"}
	Access []string `json:"access,omitempty"`
}

type ConfigDatabaseConfig struct {
	// Access controls for config database
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={"to * by * none"}
	Access []string `json:"access,omitempty"`
}

type DirectoryServiceSpec struct {
	// Type of service to create. Defaults to ClusterIP
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum:=ClusterIP;LoadBalancer;NodePort
	// +kubebuilder:default:=ClusterIP
	Type corev1.ServiceType `json:"type,omitempty"`
}

// DirectoryStatus defines the observed state of Directory.
type DirectoryStatus struct {
	// Slice of conditions storing the condition of the directory
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Available",type="string",JSONPath=`.status.conditions[?(@.type=="Available")].status`
// +kubebuilder:printcolumn:name="Age", type="date",JSONPath=`.metadata.creationTimestamp`
// Directory is the Schema for the directories API.
type Directory struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DirectorySpec   `json:"spec,omitempty"`
	Status DirectoryStatus `json:"status,omitempty"`
}

func (directory *Directory) SecretName() string {
	return fmt.Sprintf("%s-slapd-config", directory.Name)
}

func (directory *Directory) ServiceName() string {
	return directory.Name
}

func (directory *Directory) StatefulSetName() string {
	return fmt.Sprintf("%s-slapd", directory.Name)
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
