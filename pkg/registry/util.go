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
	"strings"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// truncateName truncates given name to default 63 characters,
// trims "-" from suffix if last character
func truncateName(name string) string {
	return truncateNameLengthN(name, maxTruncLength)
}

// truncateNameLengthN truncates given name to given n characters,
// trims "-" from suffix if last character
func truncateNameLengthN(name string, n int) string {
	if n < 0 {
		n = 0
	}

	if len(name) > n {
		return strings.TrimSuffix(name[:n], "-")
	}

	return strings.TrimSuffix(name, "-")
}

func generateObjectMeta(name string, namespace string, labels map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	}
}

// LabelsForDevfileRegistry returns the labels for selecting the resources
// belonging to the given devfileregistry CR name.
func LabelsForDevfileRegistry(cr *registryv1alpha1.DevfileRegistry) map[string]string {
	if cr != nil {
		return map[string]string{"app": getAppName(cr), "devfileregistry_cr": cr.Name}
	}

	return map[string]string{"app": DefaultAppName}
}
