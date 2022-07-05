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

package controllers

import (
	"fmt"
	"github.com/devfile/registry-operator/api/v1alpha1"
	"strings"
)

const (
	registryUnreachable = "Devfile %s Registry cannot be reached"
	//default status
	allRegistriesReachable = "All devfile registries are active and reachable"
	emptyStatus            = "CR list does not contain any entries"
)

// validateDevfileRegistries validates the URLs in the CR to determine if they are still reachable.
func validateDevfileRegistries(devfileRegistries []v1alpha1.DevfileRegistryService) (status string) {

	var updatedStatus []string

	if len(devfileRegistries) == 0 {
		updatedStatus = append(updatedStatus, emptyStatus)
	} else {
		for i := range devfileRegistries {
			registry := devfileRegistries[i]
			url := registry.URL
			err := v1alpha1.IsRegistryValid(registry.SkipTLSVerify, url)
			if err != nil {
				updatedStatus = append(updatedStatus, fmt.Sprintf(registryUnreachable, url))
			}
		}

		if len(updatedStatus) == 0 {
			updatedStatus = append(updatedStatus, allRegistriesReachable)
		}
	}

	return strings.Join(updatedStatus, "")

}
