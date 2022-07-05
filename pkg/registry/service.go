/*
Copyright 2020-2022 Red Hat, Inc.

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

package registry

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
)

// GenerateDevfileRegistryService returns a devfileregistry Service object
func GenerateService(cr *registryv1alpha1.DevfileRegistry, scheme *runtime.Scheme, labels map[string]string) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: generateObjectMeta(ServiceName(cr.Name), cr.Namespace, labels),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: DevfileIndexPortName,
					Port: DevfileIndexPort,
				},
				{
					Name: DevfileIndexMetricsPortName,
					Port: DevfileIndexMetricsPort,
				},
				{
					Name: OCIMetricsPortName,
					Port: OCIMetricsPort,
				},
			},
			Selector: labels,
		},
	}

	// Set DevfileRegistry instance as the owner and controller
	ctrl.SetControllerReference(cr, svc, scheme)
	return svc
}
