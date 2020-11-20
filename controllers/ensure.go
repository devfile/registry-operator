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

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	"github.com/devfile/registry-operator/pkg/registry"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/prometheus/common/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ensureService ensures that a service for the devfile registry exists on the cluster and is up to date with the custom resource
func (r *DevfileRegistryReconciler) ensureService(ctx context.Context, cr *registryv1alpha1.DevfileRegistry, labels map[string]string) (*reconcile.Result, error) {
	// Check if the service already exists, if not create a new one
	svc := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: registry.ServiceName(cr.Name), Namespace: cr.Namespace}, svc)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		pvc := registry.GenerateService(cr, r.Scheme, labels)
		log.Info("Creating a new Service", "Service.Namespace", pvc.Namespace, "Service.Name", pvc.Name)
		err = r.Create(ctx, pvc)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", pvc.Namespace, "Service.Name", pvc.Name)
			return &ctrl.Result{}, err
		}
		return nil, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

// ensureDeployment ensures that a devfile registry deployment exists on the cluster and is up to date with the custom resource
func (r *DevfileRegistryReconciler) ensureDeployment(ctx context.Context, cr *registryv1alpha1.DevfileRegistry, labels map[string]string) (*reconcile.Result, error) {
	dep := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: registry.DeploymentName(cr.Name), Namespace: cr.Namespace}, dep)
	if err != nil && errors.IsNotFound(err) {
		// Generate a new Deployment template and create it on the cluster
		dep = registry.GenerateDeployment(cr, r.Scheme, labels)

		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &ctrl.Result{}, err
		}
		return nil, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return &ctrl.Result{}, err
	}

	err = r.updateDeployment(ctx, cr, dep)
	if err != nil {
		log.Error(err, "Failed to update Deployment")
		return &ctrl.Result{}, err
	}
	return nil, nil
}

func (r *DevfileRegistryReconciler) ensurePVC(ctx context.Context, cr *registryv1alpha1.DevfileRegistry, labels map[string]string) (*reconcile.Result, error) {
	// Check if the persistentvolumeclaim already exists, if not create a new one
	pvc := &corev1.PersistentVolumeClaim{}
	err := r.Get(ctx, types.NamespacedName{Name: registry.PVCName(cr.Name), Namespace: cr.Namespace}, pvc)
	if err != nil && errors.IsNotFound(err) {
		// Define a new PVC
		pvc := registry.GeneratePVC(cr, r.Scheme, labels)
		log.Info("Creating a new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
		err = r.Create(ctx, pvc)
		if err != nil {
			log.Error(err, "Failed to create new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
			return &ctrl.Result{}, err
		}
		return nil, nil
	} else if err != nil {
		log.Error(err, "Failed to get PersistentVolumeClaim")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func (r *DevfileRegistryReconciler) ensureRoute(ctx context.Context, cr *registryv1alpha1.DevfileRegistry, labels map[string]string) (*reconcile.Result, error) {
	route := &routev1.Route{}
	err := r.Get(ctx, types.NamespacedName{Name: registry.DevfilesRouteName(cr.Name), Namespace: cr.Namespace}, route)
	if err != nil && errors.IsNotFound(err) {
		// Define a new route exposing the devfile registry index
		route := registry.GenerateRoute(cr, r.Scheme, labels)
		log.Info("Creating a new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		err = r.Create(ctx, route)
		if err != nil {
			log.Error(err, "Failed to create new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
			return &ctrl.Result{}, err
		}
		return nil, nil
	} else if err != nil {
		log.Error(err, "Failed to get Route")
		return &ctrl.Result{}, err
	}

	err = r.updateRoute(ctx, cr, route)
	if err != nil {
		log.Error(err, "Failed to update Route")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func (r *DevfileRegistryReconciler) ensureIngress(ctx context.Context, cr *registryv1alpha1.DevfileRegistry, hostname string, labels map[string]string) (*reconcile.Result, error) {
	ingress := &v1beta1.Ingress{}
	err := r.Get(ctx, types.NamespacedName{Name: registry.IngressName(cr.Name), Namespace: cr.Namespace}, ingress)
	if err != nil && errors.IsNotFound(err) {
		// Define a new ingress exposing the devfile index and oci registry
		ingress = registry.GenerateIngress(cr, hostname, r.Scheme, labels)
		log.Info("Creating a new Ingress", "Ingress.Namespace", ingress.Namespace, "Ingress.Name", ingress.Name)
		err = r.Create(ctx, ingress)
		if err != nil {
			log.Error(err, "Failed to create new Ingress", "Ingress.Namespace", ingress.Namespace, "Ingress.Name", ingress.Name)
			return &ctrl.Result{}, err
		}
		return nil, nil
	} else if err != nil {
		log.Error(err, "Failed to get Ingress")
		return &ctrl.Result{}, err
	}

	err = r.updateIngress(ctx, cr, hostname, ingress)
	if err != nil {
		log.Error(err, "Failed to update Ingress")
		return &ctrl.Result{}, err
	}
	return nil, nil
}
