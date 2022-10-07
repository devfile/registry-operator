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
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
)

// GenerateRoute returns a route exposing the devfile registry index
func GenerateRoute(cr *registryv1alpha1.DevfileRegistry, scheme *runtime.Scheme, labels map[string]string) *routev1.Route {
	weight := int32(100)

	route := &routev1.Route{
		ObjectMeta: generateObjectMeta(IngressName(cr.Name), cr.Namespace, labels),
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind:   "Service",
				Name:   ServiceName(cr.Name),
				Weight: &weight,
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString(DevfileIndexPortName),
			},
			Path: "/",
		},
	}

	if IsTLSEnabled(cr) {
		route.Spec.TLS = &routev1.TLSConfig{Termination: routev1.TLSTerminationEdge}
	}

	// Set DevfileRegistry instance as the owner and controller
	_ = ctrl.SetControllerReference(cr, route, scheme)
	return route
}
