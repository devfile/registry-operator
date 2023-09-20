//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DevfileRegistriesListSpec defines the desired state of DevfileRegistriesList
type DevfileRegistriesListSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// DevfileRegistries is a list of devfile registry services
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	DevfileRegistries []DevfileRegistryService `json:"devfileRegistries"`
}

// DevfileRegistryService represents the properties used to identify a devfile registry service.
type DevfileRegistryService struct {
	// Name is the unique Name of the devfile registry.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`
	// URL is the unique URL of the devfile registry.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	URL string `json:"url"`
	// SkipTLSVerify defaults to false.  Set to true in a non-production environment to bypass certificate checking
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SkipTLSVerify bool `json:"skipTLSVerify"`
}

// DevfileRegistriesListStatus defines the observed state of DevfileRegistriesList
type DevfileRegistriesListStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions shows the state of this CR's devfile registry list.  If registries are no longer reachable, they will be listed here
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The status for the Devfile Registries List"
// +operator-sdk:csv:customresourcedefinitions:resources={{Deployment,v1,devfileregistrieslist-deployment}}

// DevfileRegistriesList is a custom resource where namespace users can add a list of Devfile Registries to allow devfiles to be visible
// at the namespace level.  In order to be added to the list, the Devfile Registries must be reachable, supports the Devfile v2.0 spec
// and above, and is not using the default namespace.
type DevfileRegistriesList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevfileRegistriesListSpec   `json:"spec,omitempty"`
	Status DevfileRegistriesListStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DevfileRegistriesListList contains a list of DevfileRegistriesList
type DevfileRegistriesListList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevfileRegistriesList `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevfileRegistriesList{}, &DevfileRegistriesListList{})
}
