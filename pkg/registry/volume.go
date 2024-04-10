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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
)

// GenerateDevfileRegistryPVC returns a PVC for providing storage on the OCI registry container
func GeneratePVC(cr *registryv1alpha1.DevfileRegistry, scheme *runtime.Scheme, labels map[string]string) *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: generateObjectMeta(PVCName(cr), cr.Namespace, labels),
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(getDevfileRegistryVolumeSize(cr)),
				},
			},
		},
	}

	// Set DevfileRegistry instance as the owner and controller
	_ = ctrl.SetControllerReference(cr, pvc, scheme)
	return pvc
}
