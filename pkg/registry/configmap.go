/*
Copyright 2020-2023 Red Hat, Inc.

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

package registry

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
)

// GenerateRegistryConfigMap returns a configmap that is used to configure the devfile registry
func GenerateRegistryConfigMap(cr *registryv1alpha1.DevfileRegistry, scheme *runtime.Scheme, labels map[string]string) *corev1.ConfigMap {
	configMapData := make(map[string]string, 0)

	registryConfig := `
version: 0.1
log:
  fields:
    service: registry
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
  debug:
    addr: :5001
    prometheus:
      enabled: true
      path: /metrics`

	viewerEnvfile := fmt.Sprintf(`
ANALYTICS_WRITE_KEY=%s
DEVFILE_REGISTRIES=[{"name":"Community","url":"http://localhost:8080","fqdn":"%s"}]`,
		cr.Spec.Telemetry.RegistryViewerWriteKey, cr.Status.URL)

	configMapData["registry-config.yml"] = registryConfig
	configMapData[".env.registry-viewer"] = viewerEnvfile

	cm := &corev1.ConfigMap{
		ObjectMeta: generateObjectMeta(ConfigMapName(cr.Name), cr.Namespace, labels),
		Data:       configMapData,
	}

	// Set DevfileRegistry instance as the owner and controller
	_ = ctrl.SetControllerReference(cr, cm, scheme)
	return cm
}
