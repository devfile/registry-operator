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
	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
)

// GenericResourceName returns just the fully qualified app name, to be used
func GenericResourceName(cr *registryv1alpha1.DevfileRegistry) string {
	return getAppFullName(cr)
}

// DeploymentName returns the name of the deployment object associated with the DevfileRegistry CR
// Just returns the fully qualified app name right now, but extracting to a function to avoid relying on that assumption in the future
func DeploymentName(cr *registryv1alpha1.DevfileRegistry) string {
	return GenericResourceName(cr)
}

// ServiceName returns the name of the service object associated with the DevfileRegistry CR
// Just returns the fully qualified app name right now, but extracting to a function to avoid relying on that assumption in the future
func ServiceName(cr *registryv1alpha1.DevfileRegistry) string {
	return GenericResourceName(cr)
}

// ConfigMapName returns the name of the service object associated with the DevfileRegistry CR
func ConfigMapName(cr *registryv1alpha1.DevfileRegistry) string {
	const suffixLength = 15
	appFullName := getAppFullName(cr)
	configMapNameLength := (len(appFullName) + suffixLength)

	if configMapNameLength > maxTruncLength {
		return truncateNameLengthN(appFullName, len(appFullName)-suffixLength) + "-registry-config"
	}

	return appFullName + "-registry-config"
}

// PVCName returns the name of the PVC object associated with the DevfileRegistry CR
// Just returns the fully qualified app name right now, but extracting to a function to avoid relying on that assumption in the future
func PVCName(cr *registryv1alpha1.DevfileRegistry) string {
	return GenericResourceName(cr)
}

// IngressName returns the name of the Ingress object associated with the DevfileRegistry CR
// Just returns the fully qualified app name right now, but extracting to a function to avoid relying on that assumption in the future
func IngressName(cr *registryv1alpha1.DevfileRegistry) string {
	return GenericResourceName(cr)
}
