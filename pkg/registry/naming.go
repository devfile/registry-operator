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

// genericResourceName returns just the name of the custom resource, to be used
func GenericResourceName(devfileRegistryName string) string {
	return devfileRegistryName
}

// DeploymentName returns the name of the deployment object associated with the DevfileRegistry CR
// Just returns the CR name right now, but extracting to a function to avoid relying on that assumption in the future
func DeploymentName(devfileRegistryName string) string {
	return GenericResourceName(devfileRegistryName)
}

// ServiceName returns the name of the service object associated with the DevfileRegistry CR
// Just returns the CR name right now, but extracting to a function to avoid relying on that assumption in the future
func ServiceName(devfileRegistryName string) string {
	return GenericResourceName(devfileRegistryName)
}

// ConfigMapName returns the name of the service object associated with the DevfileRegistry CR
func ConfigMapName(devfileRegistryName string) string {
	return devfileRegistryName + "-registry-config"
}

// PVCName returns the name of the PVC object associated with the DevfileRegistry CR
// Just returns the CR name right now, but extracting to a function to avoid relying on that assumption in the future
func PVCName(devfileRegistryName string) string {
	return GenericResourceName(devfileRegistryName)
}

// IngressName returns the name of the Ingress object associated with the DevfileRegistry CR
// Just returns the CR name right now, but extracting to a function to avoid relying on that assumption in the future
func IngressName(devfileRegistryName string) string {
	return GenericResourceName(devfileRegistryName)
}
