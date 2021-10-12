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

package registry

import (
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

	configMapData["registry-config.yml"] = registryConfig

	viewerConfig := `{
  "Community": {
    "url": "http://localhost:8080"
  }
}`

	configMapData["devfile-registry-hosts.json"] = viewerConfig

	cm := &corev1.ConfigMap{
		ObjectMeta: generateObjectMeta(ConfigMapName(cr.Name), cr.Namespace, labels),
		Data:       configMapData,
	}

	// Set DevfileRegistry instance as the owner and controller
	ctrl.SetControllerReference(cr, cm, scheme)
	return cm
}
