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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var (
	devfileregistrieslistlog = logf.Log.WithName("devfileregistrieslist-resource")
	kubeClient               client.Client
)

func (r *DevfileRegistriesList) SetupWebhookWithManager(mgr ctrl.Manager) error {
	kubeClient = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-registry-devfile-io-v1alpha1-devfileregistrieslist,mutating=true,failurePolicy=fail,sideEffects=None,groups=registry.devfile.io,resources=devfileregistrieslists,verbs=create;update,versions=v1alpha1,name=mdevfileregistrieslist.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &DevfileRegistriesList{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *DevfileRegistriesList) Default() {
	devfileregistrieslistlog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-registry-devfile-io-v1alpha1-devfileregistrieslist,mutating=false,failurePolicy=fail,sideEffects=None,groups=registry.devfile.io,resources=devfileregistrieslists,verbs=create;update,versions=v1alpha1,name=vdevfileregistrieslist.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &DevfileRegistriesList{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *DevfileRegistriesList) ValidateCreate() error {
	devfileregistrieslistlog.Info("validate create", "name", r.Name)

	//limit CR creation to one per namespace
	devfileRegistriesList := &DevfileRegistriesListList{}
	listOpts := []client.ListOption{
		client.InNamespace(r.GetNamespace()),
	}

	if err := kubeClient.List(context.TODO(), devfileRegistriesList, listOpts...); err != nil {
		return fmt.Errorf("Error listing devfileRegistriesList custom resources: %v", err)
	}

	if len(devfileRegistriesList.Items) == 1 {
		return fmt.Errorf("A DevfileRegistriesList instance already exists. Only one instance can exist on a namespace")
	}

	return validateURLs(r.Spec.DevfileRegistries)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *DevfileRegistriesList) ValidateUpdate(old runtime.Object) error {
	devfileregistrieslistlog.Info("validate update", "name", r.Name)
	//re-validate the entire list to ensure existing URLs have not gone stale
	return validateURLs(r.Spec.DevfileRegistries)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *DevfileRegistriesList) ValidateDelete() error {
	devfileregistrieslistlog.Info("validate delete", "name", r.Name)
	return nil
}
