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
	"testing"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetDevfileRegistryIngress(t *testing.T) {

	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Correct Conjunction",
			cr: registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-name",
					Namespace: "test-namespace",
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					K8s: registryv1alpha1.DevfileRegistrySpecK8sOnly{
						IngressDomain: "my-domain",
					},
				}},
			want: "test-name-devfile-registry-test-namespace.my-domain",
		},
		{
			name: "Case 2: Unset Ingress",
			cr: registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-name",
					Namespace: "test-namespace",
				},
			},
			want: localHostname,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ingress := GetDevfileRegistryIngress(&tt.cr)
			if ingress != tt.want {
				t.Errorf("expected: %v got: %v", tt.want, ingress)
			}
		})
	}

}

func TestGetHostname(t *testing.T) {

	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Correct Hostname",
			cr: registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-name",
					Namespace: "test-namespace",
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					K8s: registryv1alpha1.DevfileRegistrySpecK8sOnly{
						IngressDomain: "my-domain",
					},
				}},
			want: "test-name-devfile-registry-test-namespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hostname := GetHostname(&tt.cr)
			if hostname != tt.want {
				t.Errorf("expected: %v got: %v", tt.want, hostname)
			}
		})
	}

}
