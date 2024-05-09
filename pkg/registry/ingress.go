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

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func GenerateIngress(cr *registryv1alpha1.DevfileRegistry, host string, scheme *runtime.Scheme, labels map[string]string) *networkingv1.Ingress {
	pathTypeImplementationSpecific := networkingv1.PathTypeImplementationSpecific
	ingress := &networkingv1.Ingress{
		ObjectMeta: generateObjectMeta(IngressName(cr), cr.Namespace, labels),
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
											Name: ServiceName(cr),
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
	_ = ctrl.SetControllerReference(cr, ingress, scheme)
	return ingress
}

func GetDevfileRegistryIngress(cr *registryv1alpha1.DevfileRegistry) string {
	if IsIngressSkipped(cr) {
		return localHostname
	}
	return GetHostname(cr) + "." + cr.Spec.K8s.IngressDomain
}

func GetHostname(cr *registryv1alpha1.DevfileRegistry) string {
	return fmt.Sprintf("%s-%s", getAppFullName(cr), cr.Namespace)
}
