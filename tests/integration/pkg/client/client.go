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

package client

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type K8sClient struct {
	kubeClient       *kubernetes.Clientset
	controllerClient client.Client
	cli              string
}

// NewK8sClient creates kubernetes client wrapper with helper functions and direct access to k8s go client
func NewK8sClient() (*K8sClient, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Instantiate an instance of conroller-runtime client
	controllerClient, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("failed to create client")
		os.Exit(1)
	}

	h := &K8sClient{kubeClient: kubeClient, controllerClient: controllerClient}

	h.cli, err = findCLI()
	if err != nil {
		fmt.Println("failed to find oc or kubectl cli")
		os.Exit(1)
	}
	return h, nil
}

// Kube returns the clientset for Kubernetes upstream.
func (c *K8sClient) Kube() kubernetes.Interface {
	return c.kubeClient
}

// findCLI returns the first found CLI compatible with oc/kubectl
func findCLI() (string, error) {
	for _, cli := range []string{"oc", "kubectl"} {
		cmd := exec.Command(cli, "version")
		err := cmd.Run()
		if err != nil {
			continue
		}
		return cli, nil
	}

	return "", errors.New("no oc/kubectl CLI found")
}
