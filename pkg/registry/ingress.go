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
	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func GenerateIngress(cr *registryv1alpha1.DevfileRegistry, host string, scheme *runtime.Scheme, labels map[string]string) *v1beta1.Ingress {
	ingress := &v1beta1.Ingress{
		ObjectMeta: generateObjectMeta(IngressName(cr.Name), cr.Namespace, labels),
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: host,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: v1beta1.IngressBackend{
										ServiceName: ServiceName(cr.Name),
										ServicePort: intstr.FromInt(int(DevfileIndexPort)),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if IsTLSEnabled(cr) && cr.Spec.TLS.SecretName != "" {
		ingress.Spec.TLS = []v1beta1.IngressTLS{
			{
				Hosts:      []string{host},
				SecretName: cr.Spec.TLS.SecretName,
			},
		}
	}

	// Set DevfileRegistry instance as the owner and controller
	ctrl.SetControllerReference(cr, ingress, scheme)
	return ingress
}

func GetDevfileRegistryIngress(cr *registryv1alpha1.DevfileRegistry) string {
	return cr.Name + "." + cr.Spec.K8s.IngressDomain
}
