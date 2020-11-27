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
	"fmt"
	"os/exec"
	"strings"

	"github.com/devfile/registry-operator/tests/integration/pkg/config"
)

// KubectlApplyResource applies resources on the cluster, corresponding to the specified file(s)
func (w *K8sClient) KubectlApplyResource(filePath string) (err error) {
	cmd := exec.Command("kubectl", "apply", "--namespace", config.Namespace, "-f", filePath)
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil && !strings.Contains(output, "AlreadyExists") {
		fmt.Println(err)
	}
	return err
}

// KubectlDeleteResource deletes the resources from the cluster that the specified file(s) correspond to
func (w *K8sClient) KubectlDeleteResource(filePath string) (err error) {
	cmd := exec.Command("kubectl", "delete", "--namespace", config.Namespace, "-f", filePath)
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil && !strings.Contains(output, "AlreadyExists") {
		fmt.Println(err)
	}
	return err
}
