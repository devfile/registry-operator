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
	"reflect"
	"testing"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_truncateName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Case 1: Short name",
			input: "devfile-registry-test",
			want:  "devfile-registry-test",
		},
		{
			name:  "Case 2: Long name",
			input: "devfile-registry-testregistry-devfile-io-k8s-prow-test-environment-afdfs2345j2234j2k42ljl234",
			want:  "devfile-registry-testregistry-devfile-io-k8s-prow-test-environm",
		},
		{
			name:  "Case 3: Short name with leftover suffix",
			input: "devfile-registry-test-",
			want:  "devfile-registry-test",
		},
		{
			name:  "Case 4: Long name with leftover suffix",
			input: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
			want:  "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := truncateName(test.input)
			if got != test.want {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}

func Test_truncateNameLengthN(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		inputLength int
		want        string
	}{
		{
			name:        "Case 1: Short name",
			inputName:   "devfile-registry-test",
			want:        "devfile-registry",
			inputLength: 17,
		},
		{
			name:        "Case 2: Long name",
			inputName:   "devfile-registry-testregistry-devfile-io-k8s-prow-test-environment-afdfs2345j2234j2k42ljl234",
			want:        "devfile-registry-testregistry-devfile-io-k8s-prow",
			inputLength: 49,
		},
		{
			name:        "Case 3: Short name with leftover suffix",
			inputName:   "devfile-registry-test-",
			want:        "devfile-registry-test",
			inputLength: 30,
		},
		{
			name:        "Case 4: Long name with leftover suffix",
			inputName:   "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
			want:        "devfile-registry-testregistry",
			inputLength: 30,
		},
		{
			name:        "Case 5: Negative truncation length",
			inputName:   "devfile-registry-test",
			want:        "",
			inputLength: -17,
		},
		{
			name:        "Case 6: Truncation length zero",
			inputName:   "devfile-registry-test",
			want:        "",
			inputLength: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := truncateNameLengthN(test.inputName, test.inputLength)
			if got != test.want {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}

func TestLabelsForDevfileRegistry(t *testing.T) {
	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want map[string]string
	}{
		{
			name: "Case 1: Labels with set CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "devfile-registry-test",
				},
			},
			want: map[string]string{
				"app":                DefaultAppName,
				"devfileregistry_cr": "devfile-registry-test",
			},
		},
		{
			name: "Case 2: Labels with set CR name and app name override",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "devfile-registry-test",
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: map[string]string{
				"app":                "dr",
				"devfileregistry_cr": "devfile-registry-test",
			},
		},
		{
			name: "Case 3: Labels with set app name override",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: map[string]string{
				"app":                "dr",
				"devfileregistry_cr": "",
			},
		},
		{
			name: "Case 4: Labels with empty CR",
			cr:   &registryv1alpha1.DevfileRegistry{},
			want: map[string]string{
				"app":                DefaultAppName,
				"devfileregistry_cr": "",
			},
		},
		{
			name: "Case 5: Labels with nil passed as CR",
			cr:   nil,
			want: map[string]string{
				"app": DefaultAppName,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := LabelsForDevfileRegistry(test.cr)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}
