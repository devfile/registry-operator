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

// OCIConfigMapName returns the name of the service object associated with the DevfileRegistry CR
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
