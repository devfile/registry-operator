/*
Copyright 2022-2023 Red Hat, Inc.

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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	typeAvailableClusterDevfileRegistriesList = "Available"
	typeDegradedClusterDevfileRegistriesList  = "Degraded"
)

// ClusterDevfileRegistriesListReconciler reconciles a ClusterDevfileRegistriesList object
type ClusterDevfileRegistriesListReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=registry.devfile.io,resources=clusterdevfileregistrieslists,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=registry.devfile.io,resources=clusterdevfileregistrieslists/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=registry.devfile.io,resources=clusterdevfileregistrieslists/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the ClusterDevfileRegistriesList object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ClusterDevfileRegistriesListReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("clusterdevfileregistrieslist", req.NamespacedName)

	// Fetch the DevfileRegistriesList instance
	clusterDevfileRegistriesList := &registryv1alpha1.ClusterDevfileRegistriesList{}
	err := r.Get(ctx, req.NamespacedName, clusterDevfileRegistriesList)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("ClusterDevfileRegistriesList resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ClusterDevfileRegistriesList")
		return ctrl.Result{}, err
	}

	if clusterDevfileRegistriesList.Status.Conditions == nil || len(clusterDevfileRegistriesList.Status.Conditions) == 0 {
		meta.SetStatusCondition(&clusterDevfileRegistriesList.Status.Conditions, metav1.Condition{
			Type:    typeAvailableClusterDevfileRegistriesList,
			Status:  metav1.ConditionUnknown,
			Reason:  "Reconciling",
			Message: "Starting reconciliation",
		})
		if err = r.Status().Update(ctx, clusterDevfileRegistriesList); err != nil {
			log.Error(err, "Failed to update ClusterDevfileRegistriesList status")
			return ctrl.Result{}, err
		}

		// re-fetch the Custom Resource after update the status
		// so that we have the latest state of the resource on the cluster and we will avoid
		// raise the issue "the object has been modified, please apply
		// your changes to the latest version and try again" which would re-trigger the reconciliation
		if err := r.Get(ctx, req.NamespacedName, clusterDevfileRegistriesList); err != nil {
			log.Error(err, "Failed to re-fetch ClusterDevfileRegistriesList")
			return ctrl.Result{}, err
		}
	}

	clusterDevfileRegistries := clusterDevfileRegistriesList.Spec.DevfileRegistries
	validateMessage := validateDevfileRegistries(clusterDevfileRegistries)
	newCondition := metav1.Condition{
		Type:    typeAvailableClusterDevfileRegistriesList,
		Reason:  "Reconciling",
		Message: validateMessage,
	}
	if validateMessage != allRegistriesReachable {
		newCondition.Status = metav1.ConditionFalse
	} else {
		newCondition.Status = metav1.ConditionTrue
	}
	log.Info(fmt.Sprintf("Status is being updated %s ", newCondition.Message))
	meta.SetStatusCondition(&clusterDevfileRegistriesList.Status.Conditions, newCondition)
	if err = r.Status().Update(ctx, clusterDevfileRegistriesList); err != nil {
		log.Error(err, "Failed to update ClusterDevfileRegistriesList status")
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{RequeueAfter: time.Hour}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterDevfileRegistriesListReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&registryv1alpha1.ClusterDevfileRegistriesList{}).
		Complete(r)
}
