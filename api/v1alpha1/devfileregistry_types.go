//
// Copyright (c) 2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DevfileRegistrySpec defines the desired state of DevfileRegistry
type DevfileRegistrySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Sets the container image containing devfile stacks to be deployed on the Devfile Registry
	DevfileIndexImage string `json:"devfileIndexImage,omitempty"`

	// Overrides the container image used for the OCI registry.
	// Recommended to leave blank and default to the image specified by the operator.
	// +optional
	OciRegistryImage string                     `json:"ociRegistryImage,omitempty"`
	Storage          DevfileRegistrySpecStorage `json:"storage,omitempty"`
	TLS              DevfileRegistrySpecTLS     `json:"tls,omitempty"`
	K8s              DevfileRegistrySpecK8sOnly `json:"k8s,omitempty"`
}

// DevfileRegistrySpecStorage defines the desired state of the storage for the DevfileRegistry
type DevfileRegistrySpecStorage struct {
	// Instructs the operator to deploy the DevfileRegistry with persistent storage
	// Enabled by default. Disabling is only recommended for development or test.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Configures the size of the devfile registry's persistent volume, if enabled.
	// Defaults to 1Gi.
	// +optional
	RegistryVolumeSize string `json:"ociRegistryImage,omitempty"`
}

// DevfileRegistrySpecTLS defines the desired state for TLS in the DevfileRegistry
type DevfileRegistrySpecTLS struct {
	// Instructs the operator to deploy the DevfileRegistry with TLS enabled.
	// Enabled by default. Disabling is only recommended for development or test.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Name of an optional, pre-existing TLS secret to use for TLS termination on ingress/route resources.
	// +optional
	SecretName string `json:"ociRegistryImage,omitempty"`
}

// DevfileRegistrySpecK8sOnly defines the desired state of the kubernetes-only fields of the DevfileRegistry
type DevfileRegistrySpecK8sOnly struct {
	// Ingress domain for a Kubernetes cluster. This MUST be explicitly specified on Kubernetes. There are no defaults
	IngressDomain string `json:"ingressDomain,omitempty"`
}

// DevfileRegistryStatus defines the observed state of DevfileRegistry
type DevfileRegistryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	URL string `json:"url"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DevfileRegistry is the Schema for the devfileregistries API
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
