/*
Copyright 2021 The Clusternet Authors.

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

// Important: Run "make generated" to regenerate code after modifying this file

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope="Namespaced",shortName=loc;local,categories=clusternet
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// Localization represents the override rule for a group of resources.
type Localization struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec LocalizationSpec `json:"spec"`
}

// LocalizationSpec defines the desired state of Localization
type LocalizationSpec struct {
	// OverrideType specifies the override type for this Localization.
	//
	// +optional
	OverrideType OverrideType `json:"overrideType,omitempty"`

	// Overrides holds all the OverrideRule.
	//
	// +optional
	Overrides []OverrideRule `json:"overrides,omitempty"`

	// Priority is an integer defining the relative importance of this Localization compared to others. Lower
	// numbers are considered lower priority.
	//
	// +optional
	// +kubebuilder:default=1000
	Priority int32 `json:"priority,omitempty"`

	// Feed holds references to the objects the Localization applies to.
	//
	// +optional
	Feed `json:"feed,omitempty"`
}

type OverrideType string

const (
	HelmOverride    OverrideType = "Helm"
	DefaultOverride OverrideType = "Default"
)

// OverrideRule holds information that describes a override rule.
type OverrideRule struct {
	// RuleName indicate the OverrideRule name.
	//
	// +optional
	RuleName string `json:"ruleName,omitempty"`

	// OverrideValue represents override values.
	//
	// +optional
	OverrideValue map[string]string `json:"value,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalizationList contains a list of Localization
type LocalizationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Localization `json:"items"`
}

// Important: Run "make generated" to regenerate code after modifying this file

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope="Cluster",shortName=cloc;clocal,categories=clusternet
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// ClusterLocalization represents the cluster-scoped override rule for a group of resources.
type ClusterLocalization struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ClusterLocalizationSpec `json:"spec"`
}

// ClusterLocalizationSpec defines the desired state of ClusterLocalization
type ClusterLocalizationSpec struct {
	// OverrideType specifies the override type for this ClusterLocalization.
	//
	// +optional
	OverrideType OverrideType `json:"overrideType,omitempty"`

	// Overrides holds all the OverrideRule.
	// +optional
	Overrides []OverrideRule `json:"overrides,omitempty"`

	// Priority is an integer defining the relative importance of this ClusterLocalization compared to others. Lower
	// numbers are considered lower priority.
	//
	// +optional
	// +kubebuilder:default=1000
	Priority int32 `json:"priority,omitempty"`

	// Feed holds references to the objects the ClusterLocalization applies to.
	//
	// +optional
	Feed `json:"feed,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterLocalizationList contains a list of ClusterLocalization
type ClusterLocalizationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterLocalization `json:"items"`
}
