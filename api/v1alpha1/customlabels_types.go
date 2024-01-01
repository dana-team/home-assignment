/*
Copyright 2023.

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

// CustomLabelsSpec defines the desired state of CustomLabels
type CustomLabelsSpec struct {
	// map of custom labels to add to namespace
	CustomLabels map[string]string `json:"customLabels,omitempty"`
}

// CustomLabelsStatus defines the observed state of CustomLabels
type CustomLabelsStatus struct {
	// TRUE if labels have been added, else FALSE
	Applied bool   `json:"applied"`
	Message string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Applied",type=boolean,JSONPath=`.status.applied`

// CustomLabels is the Schema for the customlabels API
type CustomLabels struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomLabelsSpec   `json:"spec,omitempty"`
	Status CustomLabelsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CustomLabelsList contains a list of CustomLabels
type CustomLabelsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomLabels `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomLabels{}, &CustomLabelsList{})
}
