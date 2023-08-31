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

package deploy

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/devfile/registry-operator/tests/integration/pkg/config"
)

const OperatorNamespace = "registry-operator-system"

// CreateNamespace ensures that the namespace that the tests will run in already exiss
func (w *Deployment) CreateNamespace() error {
	_, err := w.kubeClient.Kube().CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.Namespace,
		},
	}, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

// DeployDevfileRegistryOperator deploys the DevfileRegistry operator
func (w *Deployment) DeployDevfileRegistryOperator() error {
	label := "app.kubernetes.io/name=devfileregistry-operator"
	cmd := exec.Command("make", "deploy")
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil && !strings.Contains(string(output), "AlreadyExists") {
		fmt.Println(err.Error())
		return err
	}

	deploy, err := w.kubeClient.WaitForPodRunningByLabelWithNamespace(label, OperatorNamespace)
	fmt.Println("Devfile Registry pod to be ready")
	if !deploy || err != nil {
		fmt.Println("Devfile Registry not deployed")
		return err
	}
	return nil
}

func (w *Deployment) InstallCustomResourceDefinitions() error {
	devfileRegistryCRD := exec.Command("make", "install")
	output, err := devfileRegistryCRD.CombinedOutput()
	fmt.Println(string(output))
	if err != nil && !strings.Contains(string(output), "AlreadyExists") {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
