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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
)

// DevfileRegistryReconciler reconciles a DevfileRegistry object
type DevfileRegistryReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=registry.devfile.io,resources=devfileregistries,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=registry.devfile.io,resources=devfileregistries/status,verbs=get;update;patch

func (r *DevfileRegistryReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("devfileregistry", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *DevfileRegistryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&registryv1alpha1.DevfileRegistry{}).
		Complete(r)
}
