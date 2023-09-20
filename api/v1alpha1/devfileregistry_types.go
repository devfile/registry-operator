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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "make" to regenerate code after modifying this file

// DevfileRegistrySpec defines the desired state of DevfileRegistry
type DevfileRegistrySpec struct {
	// Sets the devfile index container spec to be deployed on the Devfile Registry
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	DevfileIndex DevfileRegistrySpecContainer `json:"devfileIndex,omitempty"`
	// Sets the OCI registry container spec to be deployed on the Devfile Registry
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	OciRegistry DevfileRegistrySpecContainer `json:"ociRegistry,omitempty"`
	// Sets the registry viewer container spec to be deployed on the Devfile Registry
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RegistryViewer DevfileRegistrySpecContainer `json:"registryViewer,omitempty"`

	// Sets the container image containing devfile stacks to be deployed on the Devfile Registry
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	// +deprecated
	DevfileIndexImage string `json:"devfileIndexImage,omitempty"`

	// Overrides the container image used for the OCI registry.
	// Recommended to leave blank and default to the image specified by the operator.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	// +deprecated
	OciRegistryImage string `json:"ociRegistryImage,omitempty"`
	// Overrides the container image used for the registry viewer.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	// +deprecated
	RegistryViewerImage string `json:"registryViewerImage,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Storage DevfileRegistrySpecStorage `json:"storage,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	TLS DevfileRegistrySpecTLS `json:"tls,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	K8s DevfileRegistrySpecK8sOnly `json:"k8s,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Telemetry DevfileRegistrySpecTelemetry `json:"telemetry,omitempty"`
	// Sets the registry server deployment to run under headless mode
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Headless *bool `json:"headless,omitempty"`
}

// DevfileRegistrySpecContainer defines the desired state of a container for the DevfileRegistry
type DevfileRegistrySpecContainer struct {
	// Sets the container image
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Image string `json:"image,omitempty"`
	// Sets the image pull policy for the container
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

// DevfileRegistrySpecStorage defines the desired state of the storage for the DevfileRegistry
type DevfileRegistrySpecStorage struct {
	// Instructs the operator to deploy the DevfileRegistry with persistent storage
	// Disabled by default.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Configures the size of the devfile registry's persistent volume, if enabled.
	// Defaults to 1Gi.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RegistryVolumeSize string `json:"registryVolumeSize,omitempty"`
}

// DevfileRegistrySpecTLS defines the desired state for TLS in the DevfileRegistry
type DevfileRegistrySpecTLS struct {
	// Instructs the operator to deploy the DevfileRegistry with TLS enabled.
	// Enabled by default. Disabling is only recommended for development or test.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Name of an optional, pre-existing TLS secret to use for TLS termination on ingress/route resources.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SecretName string `json:"secretName,omitempty"`
}

// DevfileRegistrySpecK8sOnly defines the desired state of the kubernetes-only fields of the DevfileRegistry
type DevfileRegistrySpecK8sOnly struct {
	// Ingress domain for a Kubernetes cluster. This MUST be explicitly specified on Kubernetes. There are no defaults
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	IngressDomain string `json:"ingressDomain,omitempty"`
}

// Telemetry defines the desired state for telemetry in the DevfileRegistry
type DevfileRegistrySpecTelemetry struct {
	// The registry name (can be any string) that is used as identifier for devfile telemetry.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RegistryName string `json:"registryName"`

	// Specify a telemetry key to allow devfile specific data to be sent to a client's own Segment analytics source.
	// If the write key is specified then telemetry will be enabled
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Key string `json:"key,omitempty"`

	// Specify a telemetry write key for the registry viewer component to allow data to be sent to a client's own Segment analytics source.
	// If the write key is specified then telemetry for the registry viewer component will be enabled
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RegistryViewerWriteKey string `json:"registryViewerWriteKey,omitempty"`
}

// DevfileRegistryStatus defines the observed state of DevfileRegistry
type DevfileRegistryStatus struct {
	// URL is the exposed URL for the Devfile Registry, and is set in the status after the registry has become available.
	// +operator-sdk:csv:customresourcedefinitions:type=status
	URL string `json:"url"`

	// Conditions shows the state devfile registries.
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +operator-sdk:csv:customresourcedefinitions:resources={{Deployment,v1,devfileregistry-deployment}}

// DevfileRegistry is a custom resource allows you to create and manage your own index server and registry viewer.
// In order to be added, the Devfile Registry must be reachable, supports the Devfile v2.0 spec and above, and is
// not using the default namespace.
// +kubebuilder:resource:path=devfileregistries,shortName=devreg;dr
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".status.url",description="The URL for the Devfile Registry"
type DevfileRegistry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevfileRegistrySpec   `json:"spec,omitempty"`
	Status DevfileRegistryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DevfileRegistryList contains a list of DevfileRegistry
type DevfileRegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevfileRegistry `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevfileRegistry{}, &DevfileRegistryList{})
}
