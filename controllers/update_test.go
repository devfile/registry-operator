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
	"testing"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestUpdateDeploymentForHeadlessChange(t *testing.T) {
	r := &DevfileRegistryReconciler{}

	tests := []struct {
		name    string
		cr      *registryv1alpha1.DevfileRegistry
		dep     *appsv1.Deployment
		want    bool
		wantErr bool
	}{
		{
			name: "Headless true - REGISTRY_HEADLESS set correctly, viewer not present",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Headless: func(b bool) *bool { return &b }(true),
				},
			},
			dep: func() *appsv1.Deployment {
				dep := &appsv1.Deployment{}
				dep.Spec.Template.Spec.Containers = []corev1.Container{
					{
						Name: "devfile-registry",
						Env:  []corev1.EnvVar{{Name: "REGISTRY_HEADLESS", Value: "true"}},
					},
				}
				return dep
			}(),
			want:    false, // No changes needed, already correct
			wantErr: false,
		},
		{
			name: "Headless true - REGISTRY_HEADLESS incorrect, viewer present",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Headless: func(b bool) *bool { return &b }(true),
				},
			},
			dep: func() *appsv1.Deployment {
				dep := &appsv1.Deployment{}
				dep.Spec.Template.Spec.Containers = []corev1.Container{
					{
						Name: "devfile-registry",
						Env:  []corev1.EnvVar{{Name: "REGISTRY_HEADLESS", Value: "false"}},
					},
					{
						Name: viewerContainerName,
					},
				}
				return dep
			}(),
			want:    true, // Changes required: update ENV and remove viewer
			wantErr: false,
		},
		{
			name: "Headless false - REGISTRY_HEADLESS set correctly, viewer present",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Headless: func(b bool) *bool { return &b }(false),
				},
			},
			dep: func() *appsv1.Deployment {
				dep := &appsv1.Deployment{}
				dep.Spec.Template.Spec.Containers = []corev1.Container{
					{
						Name: "devfile-registry",
						Env:  []corev1.EnvVar{{Name: "REGISTRY_HEADLESS", Value: "false"}},
					},
					{
						Name: viewerContainerName,
					},
				}
				return dep
			}(),
			want:    false, // No changes needed, already correct
			wantErr: false,
		},
		{
			name: "Headless false - REGISTRY_HEADLESS incorrect, viewer missing",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Headless: func(b bool) *bool { return &b }(false),
				},
			},
			dep: func() *appsv1.Deployment {
				dep := &appsv1.Deployment{}
				dep.Spec.Template.Spec.Containers = []corev1.Container{
					{
						Name: "devfile-registry",
						Env:  []corev1.EnvVar{{Name: "REGISTRY_HEADLESS", Value: "true"}},
					},
				}
				return dep
			}(),
			want:    true, // Changes required: update ENV and add viewer
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the original cr to check if it's modified
			crCopy := tt.cr.DeepCopy()

			// Call the method and check for errors
			value, err := r.updateDeploymentForHeadlessChange(crCopy, tt.dep)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateDeploymentForHeadlessChange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare the return value with the expected value
			if value != tt.want {
				t.Errorf("updateDeploymentForHeadlessChange() = %v, want %v", value, tt.want)
			}
		})
	}
}
