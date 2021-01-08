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
	ctrl.SetControllerReference(cr, route, scheme)
	return route
}
