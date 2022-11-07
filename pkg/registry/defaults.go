/*
Copyright 2020-2022 Red Hat, Inc.

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

	// Defaults/constants for devfile registry storages
	DefaultDevfileRegistryVolumeSize = "1Gi"
	DevfileRegistryVolumeEnabled     = true
	DevfileRegistryVolumeName        = "devfile-registry-storage"

	DevfileRegistryTLSEnabled       = true
	DevfileRegistryTelemetryEnabled = false

	// Defaults/constants for devfile registry services
	DevfileIndexPortName        = "devfile-registry-metadata"
	DevfileIndexPort            = 8080
	DevfileIndexMetricsPortName = "devfile-index-metrics"
	DevfileIndexMetricsPort     = 7071
	OCIMetricsPortName          = "oci-registry-metrics"
	OCIMetricsPort              = 5001
	RegistryViewerPort          = 3000
)

func GetRegistryViewerImage(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.RegistryViewerImage != "" {
		return cr.Spec.RegistryViewerImage
	}
	return DefaultRegistryViewerImage
}

func GetOCIRegistryImage(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.OciRegistryImage != "" {
		return cr.Spec.OciRegistryImage
	}
	return DefaultOCIRegistryImage
}

func GetDevfileIndexImage(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.DevfileIndexImage != "" {
		return cr.Spec.DevfileIndexImage
	}
	return DefaultDevfileIndexImage
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
// If it's not set, it returns true by default.
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
