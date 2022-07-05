/*
Copyright 2022 Red Hat, Inc.

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

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The status for the Cluster Devfile Registries List"

// ClusterDevfileRegistriesList is a custom resource where cluster admins can add a list of Devfile Registries to allow devfiles to be visible
// at the cluster level.  In order to be added to the list, the Devfile Registries must be reachable and support the Devfile v2.0 spec and above.
type ClusterDevfileRegistriesList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevfileRegistriesListSpec   `json:"spec,omitempty"`
	Status DevfileRegistriesListStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterDevfileRegistriesListList contains a list of ClusterDevfileRegistriesList
type ClusterDevfileRegistriesListList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterDevfileRegistriesList `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterDevfileRegistriesList{}, &ClusterDevfileRegistriesListList{})
}
