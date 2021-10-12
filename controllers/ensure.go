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

package controllers

import (
	"context"
	"reflect"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	"github.com/devfile/registry-operator/pkg/registry"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/prometheus/common/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *DevfileRegistryReconciler) ensure(ctx context.Context, cr *registryv1alpha1.DevfileRegistry, resource runtime.Object, labels map[string]string, ingressDomain string) (*reconcile.Result, error) {
	resourceType := reflect.TypeOf(resource).Elem().Name()
	resourceName := getResourceName(resource, cr.Name)

	// Check to see if the requested resource exists on the cluster. If it doesn't exist, create it and return.
	err := r.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: cr.Namespace}, resource)
	if err != nil && errors.IsNotFound(err) {
		generatedResource := r.generateResourceObject(cr, resource, labels, ingressDomain)
		log.Info("Creating a new "+resourceType, resourceType+".Namespace", cr.Namespace+".Name", resourceName)
		err = r.Create(ctx, generatedResource)
		if err != nil {
			log.Error(err, "Failed to create new "+resourceType, resourceType+".Namespace", cr.Namespace, "Service.Name", cr.Namespace+".Name", resourceName)
			return &ctrl.Result{}, err
		}
		return nil, nil
	} else if err != nil {
		log.Error(err, "Failed to get "+resourceType)
		return &ctrl.Result{}, err
	}

	// Update the given resource, if needed
	// At this moment, only registry deployments, routes and ingresses need to be updated.
	switch resource.(type) {
	case *appsv1.Deployment:
		dep, _ := resource.(*appsv1.Deployment)
		err = r.updateDeployment(ctx, cr, dep)
	case *routev1.Route:
		route, _ := resource.(*routev1.Route)
		err = r.updateRoute(ctx, cr, route)
	case *networkingv1.Ingress:
		ingress, _ := resource.(*networkingv1.Ingress)
		err = r.updateIngress(ctx, cr, ingressDomain, ingress)
	}
	if err != nil {
		log.Error(err, "Failed to update "+resourceType)
		return &ctrl.Result{}, err
	}
	return nil, nil
}

func getResourceName(resource runtime.Object, crName string) string {
	switch resource.(type) {
	case *appsv1.Deployment:
		return registry.DeploymentName(crName)
	case *corev1.ConfigMap:
		return registry.ConfigMapName(crName)
	case *corev1.PersistentVolumeClaim:
		return registry.PVCName(crName)
	case *corev1.Service:
		return registry.ServiceName(crName)
	case *routev1.Route, *networkingv1.Ingress:
		return registry.IngressName(crName)
	}
	return registry.GenericResourceName(crName)
}

func (r *DevfileRegistryReconciler) generateResourceObject(cr *registryv1alpha1.DevfileRegistry, resource runtime.Object, labels map[string]string, ingressDomain string) runtime.Object {
	switch resource.(type) {
	case *appsv1.Deployment:
		return registry.GenerateDeployment(cr, r.Scheme, labels)
	case *corev1.ConfigMap:
		return registry.GenerateRegistryConfigMap(cr, r.Scheme, labels)
	case *corev1.PersistentVolumeClaim:
		return registry.GeneratePVC(cr, r.Scheme, labels)
	case *corev1.Service:
		return registry.GenerateService(cr, r.Scheme, labels)
	case *routev1.Route:
		return registry.GenerateRoute(cr, r.Scheme, labels)
	case *networkingv1.Ingress:
		return registry.GenerateIngress(cr, ingressDomain, r.Scheme, labels)
	}
	return nil
}
