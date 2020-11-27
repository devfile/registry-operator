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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func generateObjectMeta(name string, namespace string, labels map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	}
}

// LabelsForDevfileRegistry returns the labels for selecting the resources
// belonging to the given devfileregistry CR name.
func LabelsForDevfileRegistry(name string) map[string]string {
	return map[string]string{"app": "devfileregistry", "devfileregistry_cr": name}
}
