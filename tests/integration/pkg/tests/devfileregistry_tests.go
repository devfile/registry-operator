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

package tests

import (
	"fmt"
	"time"

	"github.com/devfile/registry-operator/pkg/util"
	"github.com/devfile/registry-operator/tests/integration/pkg/client"
	"github.com/devfile/registry-operator/tests/integration/pkg/config"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Integration/e2e test logic based on https://github.com/devfile/devworkspace-operator/tree/master/test/e2e

var K8sClient *client.K8sClient

var _ = ginkgo.Describe("[Create Devfile Registry resource]", func() {
	ginkgo.It("Should deploy a devfile registry on to the cluster", func() {
		crName := "devfileregistry"
		label := "devfileregistry_cr=" + crName

		// Deploy the devfileregistry resource for this test case and wait for the pod to be running
		err := K8sClient.OcApplyResource("tests/integration/examples/create/devfileregistry.yaml")
		if err != nil {
			ginkgo.Fail("Failed to create devfileregistry instance: " + err.Error())
			return
		}

		waitForPodsToFullyStartInUIMode(label)

		// Wait for the registry instance to become ready
		err = K8sClient.WaitForRegistryInstance(crName, 5*time.Minute)
		Expect(err).NotTo(HaveOccurred())

		// Retrieve the registry URL and verify the server is up and running
		registry, err := K8sClient.GetRegistryInstance(crName)
		Expect(err).NotTo(HaveOccurred())
		err = util.WaitForServer(registry.Status.URL, 5*time.Minute, false)
		Expect(err).NotTo(HaveOccurred())

		podList, err := K8sClient.ListPods(config.Namespace, label)
		Expect(err).NotTo(HaveOccurred())
		registryPod := podList.Items[0]

		// Verify that the index metrics endpoint is running
		indexMetricsURL := "http://localhost:7071/metrics"
		output, err := K8sClient.CurlEndpointInContainer(registryPod.Name, "devfile-registry", indexMetricsURL)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("promhttp_metric_handler_requests_total"))

		// Verify that the oci metrics endpoint is running
		ociMetricsURL := "http://localhost:5001/metrics"
		output, err = K8sClient.CurlEndpointInContainer(registryPod.Name, "devfile-registry", ociMetricsURL)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("registry_storage_cache_total"))

		//verify the registry viewer is running
		registryViewerURL := "http://localhost:8080/viewer"
		output, err = K8sClient.CurlEndpointInContainer(registryPod.Name, "devfile-registry", registryViewerURL)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(ContainSubstring("registry viewer is not available in headless mode"))

		validateEnvVariables(label, registry.Name, registry.Status.URL, "", "")

	})

	var _ = ginkgo.AfterEach(func() {
		K8sClient.OcDeleteResource("tests/integration/examples/create/devfileregistry.yaml")
	})
})

var _ = ginkgo.Describe("[Create Devfile Registry resource with TLS enabled]", func() {
	ginkgo.It("Should deploy a devfile registry on to the cluster with HTTPS", func() {
		crName := "devfileregistry-tls"
		label := "devfileregistry_cr=" + crName

		// Deploy the devfileregistry resource for this test case and wait for the pod to be running
		err := K8sClient.OcApplyResource("tests/integration/examples/create/devfileregistry-tls.yaml")
		if err != nil {
			ginkgo.Fail("Failed to create devfileregistry instance: " + err.Error())
			return
		}

		waitForPodsToFullyStartInUIMode(label)

		// Wait for the registry instance to become ready
		err = K8sClient.WaitForRegistryInstance(crName, 5*time.Minute)
		Expect(err).NotTo(HaveOccurred())

		// Retrieve the registry URL and verify that its protocol is https
		registry, err := K8sClient.GetRegistryInstance(crName)
		Expect(err).NotTo(HaveOccurred())
		Expect(registry.Status.URL).To(ContainSubstring("https://"))

		// Verify that the server is accessible.
		err = util.WaitForServer(registry.Status.URL, 5*time.Minute, false)
		Expect(err).NotTo(HaveOccurred())

		validateEnvVariables(label, registry.Name, registry.Status.URL, "", "")
	})

	var _ = ginkgo.AfterEach(func() {
		K8sClient.OcDeleteResource("tests/integration/examples/create/devfileregistry-tls.yaml")
	})
})

var _ = ginkgo.Describe("[Create Devfile Registry resource with headless enabled]", func() {
	ginkgo.It("Should deploy a headless devfile registry on to the cluster", func() {
		crName := "devfileregistry-headless"
		label := "devfileregistry_cr=" + crName

		// Deploy the devfileregistry resource for this test case and wait for the pod to be running
		err := K8sClient.OcApplyResource("tests/integration/examples/create/devfileregistry-headless.yaml")
		if err != nil {
			ginkgo.Fail("Failed to create devfileregistry instance: " + err.Error())
			return
		}
		deploy, err := K8sClient.WaitForPodRunningByLabel(label)
		if !deploy {
			fmt.Println("Devfile Registry didn't start properly")
		}
		Expect(err).NotTo(HaveOccurred())

		// Wait for the registry instance to become ready
		err = K8sClient.WaitForRegistryInstance(crName, 5*time.Minute)
		Expect(err).NotTo(HaveOccurred())

		// Retrieve the registry URL and verify the server is up and running
		registry, err := K8sClient.GetRegistryInstance(crName)
		Expect(err).NotTo(HaveOccurred())
		err = util.WaitForServer(registry.Status.URL, 5*time.Minute, false)
		Expect(err).NotTo(HaveOccurred())

		// Registry Viewer should not be accessible
		podList, err := K8sClient.ListPods(config.Namespace, label)
		Expect(err).NotTo(HaveOccurred())
		registryPod := podList.Items[0]
		registryViewerURL := "http://localhost:8080/viewer"
		output, err := K8sClient.CurlEndpointInContainer(registryPod.Name, "devfile-registry", registryViewerURL)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("registry viewer is not available in headless mode"))
	})

	var _ = ginkgo.AfterEach(func() {
		K8sClient.OcDeleteResource("tests/integration/examples/create/devfileregistry-headless.yaml")
	})
})

var _ = ginkgo.Describe("[Update Devfile Registry resource]", func() {
	ginkgo.It("Should deploy a devfile registry on to the cluster and properly update it", func() {
		crName := "devfileregistry-update"
		label := "devfileregistry_cr=" + crName

		// Deploy the devfileregistry resource for this test case and wait for the pod to be running
		err := K8sClient.OcApplyResource("tests/integration/examples/update/devfileregistry-old.yaml")
		if err != nil {
			ginkgo.Fail("Failed to create devfileregistry instance: " + err.Error())
			return
		}

		waitForPodsToFullyStartInUIMode(label)

		// Wait for the registry instance to become ready
		err = K8sClient.WaitForRegistryInstance(crName, 5*time.Minute)
		Expect(err).NotTo(HaveOccurred())

		// Retrieve the registry URL and verify the server is up and running
		registry, err := K8sClient.GetRegistryInstance(crName)
		Expect(err).NotTo(HaveOccurred())
		err = util.WaitForServer(registry.Status.URL, 5*time.Minute, false)
		Expect(err).NotTo(HaveOccurred())

		// Update the devfileregistry resource for this test case
		fmt.Printf("Applying update...")
		err = K8sClient.OcApplyResource("tests/integration/examples/update/devfileregistry-new.yaml")
		if err != nil {
			ginkgo.Fail("Failed to create devfileregistry instance: " + err.Error())
			return
		}

		//Wait for termination and then restart after update
		fmt.Println("Wait for Devfile Registry pods to terminate")
		//wait for pod to terminate and run again.
		failed, err := K8sClient.WaitForPodFailedByLabel(label)
		if !failed {
			fmt.Println("Devfile Registry did not terminate")
		}
		Expect(err).NotTo(HaveOccurred())

		//wait for pod to run again
		fmt.Println("Wait for Devfile Registry pods to start again")
		deploy, err := K8sClient.WaitForPodRunningByLabel(label)
		if !deploy {
			fmt.Println("Devfile Registry didn't start properly")
		}
		Expect(err).NotTo(HaveOccurred())

		// Retrieve the registry URL and verify that its protocol is https
		url, err := K8sClient.WaitForURLChange(crName, registry.Status.URL, 5*time.Minute)
		Expect(err).NotTo(HaveOccurred())
		Expect(url).To(ContainSubstring("https://"))

		// Verify that the server is accessible.
		err = util.WaitForServer(url, 5*time.Minute, false)
		Expect(err).NotTo(HaveOccurred())

		validateEnvVariables(label, registry.Name, url, "viewer-key", "registry-key")
	})

	var _ = ginkgo.AfterEach(func() {
		K8sClient.OcDeleteResource("tests/integration/examples/update/devfileregistry-new.yaml")
	})
})

// waitForPodsToFullyStartInUIMode.  The registry viewer relies on the resolved DevfileRegistry URL which will cause an update when writing to the environment variable.
// We need to wait for pod restarts in order to fully test the values
func waitForPodsToFullyStartInUIMode(label string) {
	deploy, err := K8sClient.WaitForPodRunningByLabel(label)
	if !deploy {
		fmt.Println("Devfile Registry didn't start properly")
	}
	Expect(err).NotTo(HaveOccurred())

	fmt.Println("Wait for Devfile Registry pods to terminate")
	//wait for pod to terminate and run again.
	failed, err := K8sClient.WaitForPodFailedByLabel(label)
	if !failed {
		fmt.Println("Devfile Registry did not terminate")
	}
	Expect(err).NotTo(HaveOccurred())

	//wait for pod to run again
	fmt.Println("Wait for Devfile Registry pods to start again")
	deploy, err = K8sClient.WaitForPodRunningByLabel(label)
	if !deploy {
		fmt.Println("Devfile Registry didn't start properly")
	}
	Expect(err).NotTo(HaveOccurred())
}

// validateEnvVariables validates the set environment variables of the devfile registry containers
func validateEnvVariables(label, registryName, registryURL, viewerWriteKey, telemetryKey string) {
	//Determine if the viewer pod contains the resolved DevfileRegistryURL
	newPodList, err := K8sClient.ListPods(config.Namespace, label)
	Expect(err).NotTo(HaveOccurred())
	registryPod := newPodList.Items[0]
	containers := registryPod.Spec.Containers
	if len(containers) == 3 {
		//verify the telemetry key is set in the devfile registry
		registryEnvVars := containers[0].Env
		for _, regEnvs := range registryEnvVars {
			if regEnvs.Name == "TELEMETRY_KEY" {
				Expect(regEnvs.Value).Should(Equal(telemetryKey))
			}
		}

		//registry viewer is the last container
		viewerEnvVars := containers[2].Env
		for _, env := range viewerEnvVars {
			if env.Name == "NEXT_PUBLIC_ANALYTICS_WRITE_KEY" {
				Expect(env.Value).Should(Equal(viewerWriteKey))
			}

			if env.Name == "DEVFILE_REGISTRIES" {
				Expect(env.Value).Should(ContainSubstring(registryName))
				Expect(env.Value).Should(ContainSubstring(registryURL))
			}
		}
	} else {
		ginkgo.Fail("There should be 3 containers, got %d ", len(containers))
	}
}
