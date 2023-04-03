/*
Copyright 2020-2023 Red Hat, Inc.

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

package registry

import (
	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

const (
	// Default image:tags
	DefaultDevfileIndexImage   = "quay.io/devfile/devfile-index:next"
	DefaultRegistryViewerImage = "quay.io/devfile/registry-viewer:next"
	DefaultOCIRegistryImage    = "quay.io/devfile/oci-registry:next"

	// Default image pull policies
	DefaultDevfileIndexImagePullPolicy   = corev1.PullAlways
	DefaultRegistryViewerImagePullPolicy = corev1.PullAlways
	DefaultOCIRegistryImagePullPolicy    = corev1.PullAlways

	// Defaults/constants for devfile registry storages
	DefaultDevfileRegistryVolumeSize = "1Gi"
	DevfileRegistryVolumeEnabled     = false
	DevfileRegistryVolumeName        = "devfile-registry-storage"

	DevfileRegistryTLSEnabled       = true
	DevfileRegistryTelemetryEnabled = false

	DefaultDevfileRegistryHeadlessEnabled = false

	// Defaults/constants for devfile registry services
	DevfileIndexPortName        = "devfile-registry-metadata"
	DevfileIndexPort            = 8080
	DevfileIndexMetricsPortName = "devfile-index-metrics"
	DevfileIndexMetricsPort     = 7071
	OCIMetricsPortName          = "oci-registry-metrics"
	OCIMetricsPort              = 5001
	OCIServerPort               = 5000
	RegistryViewerPort          = 3000
)

// GetRegistryViewerImage returns the container image for the registry viewer to be deployed on the Devfile Registry.
// Default: "quay.io/devfile/registry-viewer:next"
func GetRegistryViewerImage(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.RegistryViewer.Image != "" {
		return cr.Spec.RegistryViewer.Image
	} else if cr.Spec.RegistryViewerImage != "" {
		return cr.Spec.RegistryViewerImage
	}
	return DefaultRegistryViewerImage
}

// GetRegistryViewerImagePullPolicy returns the image pull policy for the registry viewer container.
// Default: "Always"
func GetRegistryViewerImagePullPolicy(cr *registryv1alpha1.DevfileRegistry) corev1.PullPolicy {
	if cr.Spec.RegistryViewer.ImagePullPolicy != "" {
		return cr.Spec.RegistryViewer.ImagePullPolicy
	}
	return DefaultRegistryViewerImagePullPolicy
}

// GetOCIRegistryImage returns the container image for the OCI registry to be deployed on the Devfile Registry.
// Default: "quay.io/devfile/oci-registry:next"
func GetOCIRegistryImage(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.OciRegistry.Image != "" {
		return cr.Spec.OciRegistry.Image
	} else if cr.Spec.OciRegistryImage != "" {
		return cr.Spec.OciRegistryImage
	}
	return DefaultOCIRegistryImage
}

// GetOCIRegistryImagePullPolicy returns the image pull policy for the OCI registry container.
// Default: "Always"
func GetOCIRegistryImagePullPolicy(cr *registryv1alpha1.DevfileRegistry) corev1.PullPolicy {
	if cr.Spec.OciRegistry.ImagePullPolicy != "" {
		return cr.Spec.OciRegistry.ImagePullPolicy
	}
	return DefaultOCIRegistryImagePullPolicy
}

// GetDevfileIndexImage returns the container image for the devfile index server to be deployed on the Devfile Registry.
// Default: "quay.io/devfile/devfile-index:next"
func GetDevfileIndexImage(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.DevfileIndex.Image != "" {
		return cr.Spec.DevfileIndex.Image
	} else if cr.Spec.DevfileIndexImage != "" {
		return cr.Spec.DevfileIndexImage
	}
	return DefaultDevfileIndexImage
}

// GetDevfileIndexImagePullPolicy returns the image pull policy for the devfile index container.
// Default: "Always"
func GetDevfileIndexImagePullPolicy(cr *registryv1alpha1.DevfileRegistry) corev1.PullPolicy {
	if cr.Spec.DevfileIndex.ImagePullPolicy != "" {
		return cr.Spec.DevfileIndex.ImagePullPolicy
	}
	return DefaultDevfileIndexImagePullPolicy
}

func getDevfileRegistryVolumeSize(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.Storage.RegistryVolumeSize != "" {
		return cr.Spec.Storage.RegistryVolumeSize
	}
	return DefaultDevfileRegistryVolumeSize
}

func GetDevfileRegistryVolumeSource(cr *registryv1alpha1.DevfileRegistry) corev1.VolumeSource {
	if IsStorageEnabled(cr) {
		return corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: PVCName(cr.Name),
			},
		}
	}
	// If persistence is not enabled, return an empty dir volume source
	return corev1.VolumeSource{}
}

// IsStorageEnabled returns true if storage.enabled is set in the DevfileRegistry CR
// If it's not set, it returns false by default.
func IsStorageEnabled(cr *registryv1alpha1.DevfileRegistry) bool {
	if cr.Spec.Storage.Enabled != nil {
		return *cr.Spec.Storage.Enabled
	}
	return DevfileRegistryVolumeEnabled
}

// IsTLSEnabled returns true if tls.enabled is set in the DevfileRegistry CR
// If it's not set, it returns true by default.
func IsTLSEnabled(cr *registryv1alpha1.DevfileRegistry) bool {
	if cr.Spec.TLS.Enabled != nil {
		return *cr.Spec.TLS.Enabled
	}
	return DevfileRegistryTLSEnabled
}

// IsTelemetryEnabled returns true if telemetry.key is set in the DevfileRegistry CR
// If it's not set, it returns false by default
func IsTelemetryEnabled(cr *registryv1alpha1.DevfileRegistry) bool {
	if len(cr.Spec.Telemetry.Key) > 0 {
		return true
	}
	return DevfileRegistryTelemetryEnabled
}

// IsHeadlessEnabled returns value (true/false) set under spec attribute headless
// If it's not set, it returns false by default
func IsHeadlessEnabled(cr *registryv1alpha1.DevfileRegistry) bool {
	if cr.Spec.Headless != nil {
		return *cr.Spec.Headless
	}
	return DefaultDevfileRegistryHeadlessEnabled
}
