/*
Copyright 2022 Red Hat, Inc.

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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DevfileRegistriesListReconciler reconciles a DevfileRegistriesList object
type DevfileRegistriesListReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=registry.devfile.io,resources=devfileregistrieslists,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=registry.devfile.io,resources=devfileregistrieslists/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=registry.devfile.io,resources=devfileregistrieslists/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the DevfileRegistriesList object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *DevfileRegistriesListReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("devfileregistrieslist", req.NamespacedName)
	// Fetch the DevfileRegistriesList instance
	devfileRegistriesList := registryv1alpha1.DevfileRegistriesList{}
	err := r.Get(ctx, req.NamespacedName, &devfileRegistriesList)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("DevfileRegistriesList resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get DevfileRegistriesList")
		return ctrl.Result{}, err
	}

	// check the list of registries and report on the state
	devfileRegistries := devfileRegistriesList.Spec.DevfileRegistries
	devfileRegistriesList.Status.Status = validateDevfileRegistries(devfileRegistries)
	log.Info(fmt.Sprintf("Status is being updated %s ", devfileRegistriesList.Status.Status))
	err = r.Status().Update(ctx, &devfileRegistriesList)
	if err != nil {
		log.Error(err, "Failed to update DevfileRegistriesList status")
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{RequeueAfter: time.Hour}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DevfileRegistriesListReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&registryv1alpha1.DevfileRegistriesList{}).
		Complete(r)
}
