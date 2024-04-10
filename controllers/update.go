//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"
	"fmt"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	"github.com/devfile/registry-operator/pkg/registry"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// updateDeployment ensures that a devfile registry deployment exists on the cluster and is up to date with the custom resource
func (r *DevfileRegistryReconciler) updateDeployment(ctx context.Context, cr *registryv1alpha1.DevfileRegistry, dep *appsv1.Deployment) error {
	// Check to see if the existing devfile registry deployment needs to be updated
	needsUpdating := false

	indexImage := registry.GetDevfileIndexImage(cr)
	indexImageContainer := dep.Spec.Template.Spec.Containers[0]
	if indexImageContainer.Image != indexImage {
		indexImageContainer.Image = indexImage
		needsUpdating = true
	} else {
		//check Telemetry config to see updates are needed
		registryName := cr.Spec.Telemetry.RegistryName
		registryKey := cr.Spec.Telemetry.Key
		if indexImageContainer.Env[0].Value != registryName {
			indexImageContainer.Env[0].Value = registryName
			needsUpdating = true
		}

		if indexImageContainer.Env[1].Value != registryKey {
			indexImageContainer.Env[1].Value = registryKey
			needsUpdating = true
		}

		if indexImagePullPolicy := registry.GetDevfileIndexImagePullPolicy(cr); indexImageContainer.ImagePullPolicy != indexImagePullPolicy {
			indexImageContainer.ImagePullPolicy = indexImagePullPolicy
			needsUpdating = true
		}
	}

	ociImage := registry.GetOCIRegistryImage(cr)
	ociImageContainer := dep.Spec.Template.Spec.Containers[1]
	if ociImageContainer.Image != ociImage {
		ociImageContainer.Image = ociImage
		needsUpdating = true
	} else {
		if ociImagePullPolicy := registry.GetOCIRegistryImagePullPolicy(cr); ociImageContainer.ImagePullPolicy != ociImagePullPolicy {
			ociImageContainer.ImagePullPolicy = ociImagePullPolicy
			needsUpdating = true
		}
	}

	if registry.IsStorageEnabled(cr) {
		if dep.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim == nil {
			dep.Spec.Template.Spec.Volumes[0].VolumeSource = registry.GetDevfileRegistryVolumeSource(cr)
			needsUpdating = true
		}
	} else {
		if dep.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim != nil {
			dep.Spec.Template.Spec.Volumes[0].VolumeSource = registry.GetDevfileRegistryVolumeSource(cr)
			needsUpdating = true
		}
	}

	if len(dep.Spec.Template.Spec.Containers) > 2 {
		viewerImage := registry.GetRegistryViewerImage(cr)
		viewerImageContainer := dep.Spec.Template.Spec.Containers[2]

		//determine if the NEXT_PUBLIC_ANALYTICS_WRITE_KEY env needs updating
		viewerKey := cr.Spec.Telemetry.RegistryViewerWriteKey
		if viewerImageContainer.Env[0].Value != viewerKey {
			r.Log.Info("Updating NEXT_PUBLIC_ANALYTICS_WRITE_KEY ", "value", viewerKey)
			viewerImageContainer.Env[0].Value = viewerKey
			needsUpdating = true
		}

		//determine if the DEVFILE_REGISTRIES env needs updating.  This will only occur on initial deployment since object name is unique
		newDRValue := fmt.Sprintf(`[{"name": "%s","url": "http://localhost:8080","fqdn": "%s"}]`, cr.ObjectMeta.Name, cr.Status.URL)
		if viewerImageContainer.Env[1].Value != newDRValue {
			r.Log.Info("Updating DEVFILE_REGISTRIES ", "value", newDRValue)
			viewerImageContainer.Env[1].Value = newDRValue
			needsUpdating = true
		}

		if viewerImageContainer.Image != viewerImage {
			viewerImageContainer.Image = viewerImage
			needsUpdating = true
		} else {
			if viewerImagePullPolicy := registry.GetRegistryViewerImagePullPolicy(cr); viewerImageContainer.ImagePullPolicy != viewerImagePullPolicy {
				viewerImageContainer.ImagePullPolicy = viewerImagePullPolicy
				needsUpdating = true
			}
		}
	}

	if needsUpdating {
		r.Log.Info("Updating the DevfileRegistry deployment")
		return r.Update(ctx, dep)
	}
	return nil
}

// updateRoute checks to see if any of the fields in an existing devfile index route needs updating
func (r *DevfileRegistryReconciler) updateRoute(ctx context.Context, cr *registryv1alpha1.DevfileRegistry, route *routev1.Route) error {
	needsUpdating := false

	// Check to see if TLS fields were updated
	if registry.IsTLSEnabled(cr) {
		if route.Spec.TLS == nil {
			route.Spec.TLS = &routev1.TLSConfig{Termination: routev1.TLSTerminationEdge}
			needsUpdating = true
		}
	} else {
		if route.Spec.TLS != nil {
			route.Spec.TLS = nil
			needsUpdating = true
		}
	}

	if needsUpdating {
		return r.Update(ctx, route)
	}
	return nil
}

// updateIngress checks to see if any of the fields in an existing ingress resouorce need to be updated
func (r *DevfileRegistryReconciler) updateIngress(ctx context.Context, cr *registryv1alpha1.DevfileRegistry, hostname string, ingress *networkingv1.Ingress) error {
	needsUpdating := false
	// Check to see if TLS fields were updated
	if registry.IsTLSEnabled(cr) {
		if len(ingress.Spec.TLS) == 0 {
			// TLS was toggled on, so enable it in the ingress spec
			ingress.Spec.TLS = []networkingv1.IngressTLS{
				{
					Hosts:      []string{hostname},
					SecretName: cr.Spec.TLS.SecretName,
				},
			}
			needsUpdating = true
		}
		if ingress.Spec.TLS[0].SecretName != cr.Spec.TLS.SecretName {
			// TLS secret name was updated, so update it in the ingress spec
			ingress.Spec.TLS[0].SecretName = cr.Spec.TLS.SecretName
			needsUpdating = true
		}
	} else {
		if len(ingress.Spec.TLS) > 0 {
			// TLS was disabled, so disable it in the ingress spec
			ingress.Spec.TLS = []networkingv1.IngressTLS{}
			needsUpdating = true
		}
	}

	// Check to see if the ingress domain was updated
	if ingress.Spec.Rules[0].Host != hostname {
		ingress.Spec.Rules[0].Host = hostname

		// If TLS is enabled, need to update the hostname there too
		if registry.IsTLSEnabled(cr) {
			ingress.Spec.TLS[0].Hosts = []string{hostname}
		}
		needsUpdating = true
	}

	if needsUpdating {
		return r.Update(ctx, ingress, &client.UpdateOptions{})
	}

	return nil
}

// deletePVCIfNeeded deletes the PVC for the devfile registry if one exists and if persistent storage was disabled
func (r *DevfileRegistryReconciler) deleteOldPVCIfNeeded(ctx context.Context, cr *registryv1alpha1.DevfileRegistry) error {
	// Check to see if a PVC exists, if so, need to clean it up because storage was disabled
	if !registry.IsStorageEnabled(cr) {
		pvc := &corev1.PersistentVolumeClaim{}
		err := r.Get(ctx, types.NamespacedName{Name: registry.PVCName(cr), Namespace: cr.Namespace}, pvc)
		if err != nil {
			if errors.IsNotFound(err) {
				// PVC not found, so there's no old PVC to delete. Just return nil, nothing to do.
				return nil
			} else {
				// Some other error occurred when listing PVCs, so log and return an error
				r.Log.Error(err, "Error listing PersistentVolumeClaims")
				return err
			}
		} else {
			// PVC found despite storage being disable, so delete it
			r.Log.Info("Old PersistentVolumeClaim " + pvc.Name + " found. Deleting it as storage has been disabled.")
			err = r.Delete(ctx, pvc)
			if err != nil {
				r.Log.Error(err, "Error deleting PersistentVolumeClaim", pvc.Name)
				return err
			}
		}
	}
	return nil
}
