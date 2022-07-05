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

package v1alpha1

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var clusterdevfileregistrieslistlog = logf.Log.WithName("clusterdevfileregistrieslist-resource")

func (r *ClusterDevfileRegistriesList) SetupWebhookWithManager(mgr ctrl.Manager) error {
	kubeClient = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-registry-devfile-io-v1alpha1-clusterdevfileregistrieslist,mutating=true,failurePolicy=fail,sideEffects=None,groups=registry.devfile.io,resources=clusterdevfileregistrieslists,verbs=create;update,versions=v1alpha1,name=mclusterdevfileregistrieslist.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ClusterDevfileRegistriesList{}

const multiCRError = "A ClusterDevfileRegistriesList instance already exists. Only one instance can exist in a cluster"

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ClusterDevfileRegistriesList) Default() {
	clusterdevfileregistrieslistlog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-registry-devfile-io-v1alpha1-clusterdevfileregistrieslist,mutating=false,failurePolicy=fail,sideEffects=None,groups=registry.devfile.io,resources=clusterdevfileregistrieslists,verbs=create;update,versions=v1alpha1,name=vclusterdevfileregistrieslist.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ClusterDevfileRegistriesList{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterDevfileRegistriesList) ValidateCreate() error {
	clusterdevfileregistrieslistlog.Info("validate create", "name", r.Name)
	//limit CR creation to one per cluster
	clusterDevfileRegistriesList := &ClusterDevfileRegistriesListList{}
	listOpts := []client.ListOption{
		client.InNamespace(corev1.NamespaceAll),
	}

	if err := kubeClient.List(context.TODO(), clusterDevfileRegistriesList, listOpts...); err != nil {
		return fmt.Errorf("Error listing clusterDevfileRegistriesList custom resources: %v", err)
	}

	if len(clusterDevfileRegistriesList.Items) == 1 {
		return fmt.Errorf(multiCRError)
	}

	return validateURLs(r.Spec.DevfileRegistries)

}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterDevfileRegistriesList) ValidateUpdate(old runtime.Object) error {
	clusterdevfileregistrieslistlog.Info("validate update", "name", r.Name)
	return validateURLs(r.Spec.DevfileRegistries)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterDevfileRegistriesList) ValidateDelete() error {
	clusterdevfileregistrieslistlog.Info("validate delete", "name", r.Name)
	return nil
}
