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

package tests

import (
	"fmt"
	"time"

	"github.com/devfile/registry-operator/pkg/util"
	"github.com/devfile/registry-operator/tests/integration/pkg/client"
	"github.com/devfile/registry-operator/tests/integration/pkg/config"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
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
		deploy, err := K8sClient.WaitForPodRunningByLabel(label)
		if !deploy {
			fmt.Println("Devfile Registry didn't start properly")
		}
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait for the registry instance to become ready
		err = K8sClient.WaitForRegistryInstance(crName, 5*time.Minute)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Retrieve the registry URL and verify the server is up and running
		registry, err := K8sClient.GetRegistryInstance(crName)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = util.WaitForServer(registry.Status.URL, 5*time.Minute, false)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Verify that the index metrics endpoint is running
		podList, err := K8sClient.ListPods(config.Namespace, label)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		registryPod := podList.Items[0]

		indexMetricsURL := "http://localhost:7071/metrics"
		output, err := K8sClient.CurlEndpointInContainer(registryPod.Name, "devfile-registry", indexMetricsURL)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(output).To(gomega.ContainSubstring("promhttp_metric_handler_requests_total"))

		// Verify that the oci metrics endpoint is running
		ociMetricsURL := "http://localhost:5001/metrics"
		output, err = K8sClient.CurlEndpointInContainer(registryPod.Name, "devfile-registry", ociMetricsURL)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(output).To(gomega.ContainSubstring("registry_storage_cache_total"))
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
		deploy, err := K8sClient.WaitForPodRunningByLabel(label)
		if !deploy {
			fmt.Println("Devfile Registry didn't start properly")
		}
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait for the registry instance to become ready
		err = K8sClient.WaitForRegistryInstance(crName, 5*time.Minute)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Retrieve the registry URL and verify that its protocol is https
		registry, err := K8sClient.GetRegistryInstance(crName)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(registry.Status.URL).To(gomega.ContainSubstring("https://"))

		// Verify that the server is accessible.
		err = util.WaitForServer(registry.Status.URL, 5*time.Minute, false)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
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
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait for the registry instance to become ready
		err = K8sClient.WaitForRegistryInstance(crName, 5*time.Minute)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Retrieve the registry URL and verify the server is up and running
		registry, err := K8sClient.GetRegistryInstance(crName)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = util.WaitForServer(registry.Status.URL, 5*time.Minute, false)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Registry Viewer should not be accessible
		podList, err := K8sClient.ListPods(config.Namespace, label)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		registryPod := podList.Items[0]
		registryViewerURL := "http://localhost:8080/viewer"
		output, err := K8sClient.CurlEndpointInContainer(registryPod.Name, "devfile-registry", registryViewerURL)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(output).To(gomega.ContainSubstring("registry viewer is not available in headless mode"))
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
		deploy, err := K8sClient.WaitForPodRunningByLabel(label)
		if !deploy {
			fmt.Println("Devfile Registry didn't start properly")
		}
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait for the registry instance to become ready
		err = K8sClient.WaitForRegistryInstance(crName, 5*time.Minute)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Retrieve the registry URL and verify the server is up and running
		registry, err := K8sClient.GetRegistryInstance(crName)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = util.WaitForServer(registry.Status.URL, 5*time.Minute, false)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Update the devfileregistry resource for this test case
		err = K8sClient.OcApplyResource("tests/integration/examples/update/devfileregistry-new.yaml")
		if err != nil {
			ginkgo.Fail("Failed to create devfileregistry instance: " + err.Error())
			return
		}

		// Retrieve the registry URL and verify that its protocol is https
		url, err := K8sClient.WaitForURLChange(crName, registry.Status.URL, 5*time.Minute)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(url).To(gomega.ContainSubstring("https://"))

		// Verify that the server is accessible.
		err = util.WaitForServer(url, 5*time.Minute, false)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	var _ = ginkgo.AfterEach(func() {
		K8sClient.OcDeleteResource("tests/integration/examples/update/devfileregistry-new.yaml")
	})
})
