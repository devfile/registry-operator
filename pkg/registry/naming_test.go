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
	"testing"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenericResourceName(t *testing.T) {
	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Default App Full Name",
			cr:   &registryv1alpha1.DevfileRegistry{},
			want: DefaultAppName,
		},
		{
			name: "Case 2: Default App Full Name with overridden short App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 3: Default App Full Name with overridden long App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 4: Overridden Short App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 5: Overridden Long App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 6: Default App Full Name with short CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr",
				},
			},
			want: fmt.Sprintf("%s-%s", "dr", DefaultAppName),
		},
		{
			name: "Case 7: Default App Full Name with long CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-registry",
				},
			},
			want: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-reg",
		},
		{
			name: "Case 8: Default App Full Name with CR name contains default app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "devfile-registry-test",
				},
			},
			want: "devfile-registry-test",
		},
		{
			name: "Case 9: Default App Full Name with CR name contains overridden app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr-test",
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr-test",
		},
		{
			name: "Case 10: CR set to nil",
			cr:   nil,
			want: DefaultAppName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GenericResourceName(test.cr)
			if got != test.want {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}

func TestDeploymentName(t *testing.T) {
	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Default App Full Name",
			cr:   &registryv1alpha1.DevfileRegistry{},
			want: DefaultAppName,
		},
		{
			name: "Case 2: Default App Full Name with overridden short App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 3: Default App Full Name with overridden long App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 4: Overridden Short App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 5: Overridden Long App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 6: Default App Full Name with short CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr",
				},
			},
			want: fmt.Sprintf("%s-%s", "dr", DefaultAppName),
		},
		{
			name: "Case 7: Default App Full Name with long CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-registry",
				},
			},
			want: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-reg",
		},
		{
			name: "Case 8: Default App Full Name with CR name contains default app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "devfile-registry-test",
				},
			},
			want: "devfile-registry-test",
		},
		{
			name: "Case 9: Default App Full Name with CR name contains overridden app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr-test",
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr-test",
		},
		{
			name: "Case 10: CR set to nil",
			cr:   nil,
			want: DefaultAppName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := DeploymentName(test.cr)
			if got != test.want {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}

func TestConfigMapName(t *testing.T) {
	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Default App Full Name",
			cr:   &registryv1alpha1.DevfileRegistry{},
			want: "devfile-registry-registry-config",
		},
		{
			name: "Case 2: Default App Full Name with overridden short App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr-registry-config",
		},
		{
			name: "Case 3: Default App Full Name with overridden long App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-pr-registry-config",
		},
		{
			name: "Case 4: Overridden Short App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "dr",
				},
			},
			want: "dr-registry-config",
		},
		{
			name: "Case 5: Overridden Long App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-pr-registry-config",
		},
		{
			name: "Case 6: Default App Full Name with short CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr",
				},
			},
			want: "dr-devfile-registry-registry-config",
		},
		{
			name: "Case 7: Default App Full Name with long CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-registry",
				},
			},
			want: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0r-registry-config",
		},
		{
			name: "Case 8: Default App Full Name with CR name contains default app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "devfile-registry-test",
				},
			},
			want: "devfile-registry-test-registry-config",
		},
		{
			name: "Case 9: Default App Full Name with CR name contains overridden app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr-test",
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr-test-registry-config",
		},
		{
			name: "Case 10: CR set to nil",
			cr:   nil,
			want: "devfile-registry-registry-config",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ConfigMapName(test.cr)
			if got != test.want {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}

func TestServiceName(t *testing.T) {
	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Default App Full Name",
			cr:   &registryv1alpha1.DevfileRegistry{},
			want: DefaultAppName,
		},
		{
			name: "Case 2: Default App Full Name with overridden short App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 3: Default App Full Name with overridden long App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 4: Overridden Short App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 5: Overridden Long App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 6: Default App Full Name with short CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr",
				},
			},
			want: fmt.Sprintf("%s-%s", "dr", DefaultAppName),
		},
		{
			name: "Case 7: Default App Full Name with long CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-registry",
				},
			},
			want: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-reg",
		},
		{
			name: "Case 8: Default App Full Name with CR name contains default app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "devfile-registry-test",
				},
			},
			want: "devfile-registry-test",
		},
		{
			name: "Case 9: Default App Full Name with CR name contains overridden app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr-test",
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr-test",
		},
		{
			name: "Case 10: CR set to nil",
			cr:   nil,
			want: DefaultAppName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ServiceName(test.cr)
			if got != test.want {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}

func TestPVCName(t *testing.T) {
	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Default App Full Name",
			cr:   &registryv1alpha1.DevfileRegistry{},
			want: DefaultAppName,
		},
		{
			name: "Case 2: Default App Full Name with overridden short App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 3: Default App Full Name with overridden long App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 4: Overridden Short App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 5: Overridden Long App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 6: Default App Full Name with short CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr",
				},
			},
			want: fmt.Sprintf("%s-%s", "dr", DefaultAppName),
		},
		{
			name: "Case 7: Default App Full Name with long CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-registry",
				},
			},
			want: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-reg",
		},
		{
			name: "Case 8: Default App Full Name with CR name contains default app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "devfile-registry-test",
				},
			},
			want: "devfile-registry-test",
		},
		{
			name: "Case 9: Default App Full Name with CR name contains overridden app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr-test",
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr-test",
		},
		{
			name: "Case 10: CR set to nil",
			cr:   nil,
			want: DefaultAppName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := PVCName(test.cr)
			if got != test.want {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}

func TestIngressName(t *testing.T) {
	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Default App Full Name",
			cr:   &registryv1alpha1.DevfileRegistry{},
			want: DefaultAppName,
		},
		{
			name: "Case 2: Default App Full Name with overridden short App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 3: Default App Full Name with overridden long App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 4: Overridden Short App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 5: Overridden Long App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 6: Default App Full Name with short CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr",
				},
			},
			want: fmt.Sprintf("%s-%s", "dr", DefaultAppName),
		},
		{
			name: "Case 7: Default App Full Name with long CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-registry",
				},
			},
			want: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-reg",
		},
		{
			name: "Case 8: Default App Full Name with CR name contains default app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "devfile-registry-test",
				},
			},
			want: "devfile-registry-test",
		},
		{
			name: "Case 9: Default App Full Name with CR name contains overridden app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr-test",
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr-test",
		},
		{
			name: "Case 10: CR set to nil",
			cr:   nil,
			want: DefaultAppName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IngressName(test.cr)
			if got != test.want {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}
