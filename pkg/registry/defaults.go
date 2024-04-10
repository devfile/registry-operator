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

package registry

import (
	"fmt"
	"strings"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

	// Default memory limits
	DefaultDevfileIndexMemoryLimit   = "256Mi"
	DefaultRegistryViewerMemoryLimit = "256Mi"
	DefaultOCIRegistryMemoryLimit    = "256Mi"

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

	// Default kubernetes-only fields
	DefaultK8sIngressClass = "nginx"

	// Override defaults (should be empty)
	DefaultHostnameOverride = ""
	DefaultNameOverride     = ""
	DefaultFullnameOverride = ""

	// App name default
	DefaultAppName = "devfile-registry"
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

// GetRegistryViewerMemoryLimit returns the memory limit for the registry viewer container.
// In case of invalid quantity given, it returns the default value.
// Default: resource.Quantity{s: "256Mi"}
func GetRegistryViewerMemoryLimit(cr *registryv1alpha1.DevfileRegistry) resource.Quantity {
	return getDevfileRegistrySpecContainer(cr.Spec.RegistryViewer.MemoryLimit, DefaultRegistryViewerMemoryLimit)
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

// GetOCIRegistryMemoryLimit returns the memory limit for the OCI registry container.
// In case of invalid quantity given, it returns the default value.
// Default: resource.Quantity{s: "256Mi"}
func GetOCIRegistryMemoryLimit(cr *registryv1alpha1.DevfileRegistry) resource.Quantity {
	return getDevfileRegistrySpecContainer(cr.Spec.OciRegistry.MemoryLimit, DefaultOCIRegistryMemoryLimit)
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

// GetDevfileIndexMemoryLimit returns the memory limit for the devfile index container.
// In case of invalid quantity given, it returns the default value.
// Default: resource.Quantity{s: "256Mi"}
func GetDevfileIndexMemoryLimit(cr *registryv1alpha1.DevfileRegistry) resource.Quantity {
	return getDevfileRegistrySpecContainer(cr.Spec.DevfileIndex.MemoryLimit, DefaultDevfileIndexMemoryLimit)
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

// GetK8sIngressClass returns ingress class used for the k8s ingress class field.
// Default: "nginx"
func GetK8sIngressClass(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.K8s.IngressClass != "" {
		return cr.Spec.K8s.IngressClass
	}
	return DefaultK8sIngressClass
}

// GetHostnameOverride returns hostname override used to override the hostname and domain of a devfile registry
// Default: ""
func GetHostnameOverride(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.HostnameOverride != "" {
		return cr.Spec.HostnameOverride
	}

	return DefaultHostnameOverride
}

// GetNameOverride returns name override used to override the app name of a devfile registry
// Default: ""
func GetNameOverride(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.NameOverride != "" {
		return cr.Spec.NameOverride
	}

	return DefaultNameOverride
}

// GetFullnameOverride returns full name override used to override the fully qualified app name of a devfile registry
// Default: ""
func GetFullnameOverride(cr *registryv1alpha1.DevfileRegistry) string {
	if cr.Spec.FullnameOverride != "" {
		return cr.Spec.FullnameOverride
	}

	return DefaultFullnameOverride
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

func getDevfileRegistrySpecContainer(quantity string, defaultValue string) resource.Quantity {
	if quantity != "" {
		resourceQuantity, err := resource.ParseQuantity(quantity)
		if err == nil {
			return resourceQuantity
		}
	}
	return resource.MustParse(defaultValue)
}

// getAppName returns app name of a devfile registry
// truncated to 63 characters max, if `DevfileRegistry.NameOverride`
// is set it will return the override name truncated to 63 characters max
func getAppName(cr *registryv1alpha1.DevfileRegistry) string {
	if cr != nil {
		nameOverride := GetNameOverride(cr)

		if nameOverride == DefaultNameOverride {
			return truncateName(DefaultAppName)
		}

		return truncateName(nameOverride)
	}

	return truncateName(DefaultAppName)
}

// getAppFullName returns fully qualified app name of a devfile registry
// truncated to 63 characters max, if `DevfileRegistry.FullnameOverride`
// is set it will return the override name truncated to 63 characters max
func getAppFullName(cr *registryv1alpha1.DevfileRegistry) string {
	if cr != nil {
		fullNameOverride := GetFullnameOverride(cr)

		if fullNameOverride == DefaultFullnameOverride {
			appName := getAppName(cr)
			if cr.Name == "" {
				return truncateName(appName)
			} else if strings.Contains(appName, cr.Name) {
				return truncateName(cr.Name)
			} else {
				return truncateName(fmt.Sprintf("%s-%s", cr.Name, appName))
			}
		}

		return truncateName(fullNameOverride)
	}

	return truncateName(DefaultAppName)
}
