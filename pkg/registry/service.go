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
			},
			Selector: labels,
		},
	}

	// Set DevfileRegistry instance as the owner and controller
	ctrl.SetControllerReference(cr, svc, scheme)
	return svc
}
