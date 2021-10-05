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
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func GenerateIngress(cr *registryv1alpha1.DevfileRegistry, host string, scheme *runtime.Scheme, labels map[string]string) *networkingv1.Ingress {
	pathTypeImplementationSpecific := networkingv1.PathTypeImplementationSpecific
	ingress := &networkingv1.Ingress{
		ObjectMeta: generateObjectMeta(IngressName(cr.Name), cr.Namespace, labels),
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: "/",
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: ServiceName(cr.Name),
											Port: networkingv1.ServiceBackendPort{
												Number: int32(DevfileIndexPort),
											},
										},
									},
									//Field is required to be set based on attempt to create the ingress
									PathType: &pathTypeImplementationSpecific,
								},
							},
						},
					},
				},
			},
		},
	}

	if IsTLSEnabled(cr) && cr.Spec.TLS.SecretName != "" {
		ingress.Spec.TLS = []networkingv1.IngressTLS{
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
