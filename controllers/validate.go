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

package controllers

import (
	"fmt"
	"strings"

	"github.com/devfile/registry-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// validateDevfileRegistriesAndUpdateCondition runs validateDevfileRegistries and updates a status condition based on the result
func validateDevfileRegistriesAndUpdateCondition(devfileRegistries []v1alpha1.DevfileRegistryService, condition metav1.Condition, updateConditionFn func(metav1.Condition)) {
	validateMessage := validateDevfileRegistries(devfileRegistries)

	condition.Message = validateMessage
	if validateMessage != allRegistriesReachable {
		condition.Status = metav1.ConditionFalse
	} else {
		condition.Status = metav1.ConditionTrue
	}

	updateConditionFn(condition)
}
