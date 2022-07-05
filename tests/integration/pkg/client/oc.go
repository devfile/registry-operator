/*
Copyright 2020-2022 Red Hat, Inc.

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

package client

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/devfile/registry-operator/tests/integration/pkg/config"
)

// OcApplyResource applies resources on the cluster, corresponding to the specified file(s)
func (w *K8sClient) OcApplyResource(filePath string) (err error) {
	cmd := exec.Command("oc", "apply", "--namespace", config.Namespace, "-f", filePath)
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil && !strings.Contains(output, "AlreadyExists") {
		fmt.Println(err)
	}
	return err
}

// OcDeleteResource deletes the resources from the cluster that the specified file(s) correspond to
func (w *K8sClient) OcDeleteResource(filePath string) (err error) {
	cmd := exec.Command("oc", "delete", "--namespace", config.Namespace, "-f", filePath)
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil && !strings.Contains(output, "AlreadyExists") {
		fmt.Println(err)
	}
	return err
}

// CurlEndpointInContainer execs into the given container in the pod and uses curl to hit the specified endpoint
func (w *K8sClient) CurlEndpointInContainer(pod string, container string, endpoint string) (string, error) {
	cmd := exec.Command("oc", "exec", pod, "--namespace", config.Namespace, "-c", container, "--", "curl", endpoint)
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	return output, err
}
