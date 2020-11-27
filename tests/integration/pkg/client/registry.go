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
//

package client

import (
	"context"
	"time"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	"github.com/devfile/registry-operator/tests/integration/pkg/config"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/wait"
)

// GetRegistryInstance uses the Kubernetes REST API to retrieve the specified instance of the DevfileRegistry custom resource
// If there are any issues retrieving the resource or unmarshalling it, an error is returned
func (w *K8sClient) GetRegistryInstance(name string) (*registryv1alpha1.DevfileRegistry, error) {
	data, err := w.kubeClient.RESTClient().
		Get().
		AbsPath("/apis/registry.devfile.io/v1alpha1").
		Namespace(config.Namespace).
		Resource("devfileregistries").
		Name(name).
		DoRaw(context.TODO())

	if err != nil {
		return nil, err
	}

	// Unmarshall the struct
	registry := &registryv1alpha1.DevfileRegistry{}
	err = yaml.Unmarshal(data, registry)
	if err != nil {
		return nil, err
	}

	return registry, nil
}

// WaitForRegistryInstance polls up to timeout seconds for the registry's server to become active (URL set in the status)
func (w *K8sClient) WaitForRegistryInstance(name string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, func() (bool, error) {
		devfileRegistry, err := w.GetRegistryInstance(name)
		if err != nil {
			return false, err
		}
		if devfileRegistry.Status.URL != "" {
			return true, nil
		}
		return false, nil
	})
}

// WaitForURLChange polls up to timeout seconds for the registry's URL to change in the status and returns it.
// If the URL doesn't change in the specified timeout, an error is returned
func (w *K8sClient) WaitForURLChange(name string, oldURL string, timeout time.Duration) (string, error) {
	var newURL string
	err := wait.PollImmediate(time.Second, timeout, func() (bool, error) {
		devfileRegistry, err := w.GetRegistryInstance(name)
		if err != nil {
			return false, err
		}
		if devfileRegistry.Status.URL != oldURL {
			newURL = devfileRegistry.Status.URL
			return true, nil
		}
		return false, nil
	})

	return newURL, err
}
