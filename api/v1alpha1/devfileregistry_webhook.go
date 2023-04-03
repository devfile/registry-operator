/*
Copyright 2023 Red Hat, Inc.

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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var (
	devfileregistrylog = logf.Log.WithName("devfileregistry-resource")
)

func (r *DevfileRegistry) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-registry-devfile-io-v1alpha1-devfileregistry,mutating=true,failurePolicy=fail,sideEffects=None,groups=registry.devfile.io,resources=devfileregistries,verbs=create;update,versions=v1alpha1,name=mdevfileregistry.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &DevfileRegistry{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *DevfileRegistry) Default() {
	devfileregistrylog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-registry-devfile-io-v1alpha1-devfileregistry,mutating=false,failurePolicy=fail,sideEffects=None,groups=registry.devfile.io,resources=devfileregistries,verbs=create;update,versions=v1alpha1,name=vdevfileregistry.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &DevfileRegistry{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *DevfileRegistry) ValidateCreate() error {
	devfileregistrylog.Info("validate create", "name", r.Name)

	return IsNamespaceValid(r.Namespace)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *DevfileRegistry) ValidateUpdate(old runtime.Object) error {
	devfileregistrylog.Info("validate update", "name", r.Name)

	// Validate if namespace is valid
	if err := IsNamespaceValid(r.Namespace); err != nil {
		return err
	}

	//re-validate the entire list to ensure existing URL has not gone stale
	if r.Spec.TLS.Enabled != nil {
		return IsRegistryValid(*r.Spec.TLS.Enabled, r.Status.URL)
	} else {
		return IsRegistryValid(true, r.Status.URL)
	}
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *DevfileRegistry) ValidateDelete() error {
	devfileregistrylog.Info("validate delete", "name", r.Name)
	return nil
}
