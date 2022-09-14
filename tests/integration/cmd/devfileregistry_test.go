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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/devfile/registry-operator/tests/integration/pkg/config"
	"github.com/devfile/registry-operator/tests/integration/pkg/deploy"
	"github.com/devfile/registry-operator/tests/integration/pkg/tests"

	"github.com/devfile/registry-operator/tests/integration/pkg/client"
	_ "github.com/devfile/registry-operator/tests/integration/pkg/tests"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
)

// Integration/e2e test logic based on https://github.com/devfile/devworkspace-operator/tree/master/test/e2e

// Create Constant file
const (
	testResultsDirectory = "/tmp/artifacts"
	jUnitOutputFilename  = "junit-devfileregistry-operator.xml"
)

// SynchronizedBeforeSuite blocks is executed before run all test suites
var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	fmt.Println("Starting to setup objects before run ginkgo suite")
	namespace := os.Getenv("TEST_NAMESPACE")
	if namespace != "" {
		config.Namespace = namespace
	} else {
		config.Namespace = "registry-operator-system"
	}

	k8sClient, err := client.NewK8sClient()
	if err != nil {
		fmt.Println("Failed to create k8s client")
		panic(err)
	}

	operator := deploy.NewDeployment(k8sClient)

	err = operator.CreateNamespace()
	if err != nil {
		panic(err)
	}

	if err := operator.InstallCustomResourceDefinitions(); err != nil {
		fmt.Println("Failed to install custom resources definitions on cluster")
		panic(err)
	}

	if err := operator.DeployDevfileRegistryOperator(); err != nil {
		fmt.Println("Failed to deploy DevfileRegistry operator")
		panic(err)
	}

	tests.K8sClient, err = client.NewK8sClient()
	if err != nil {
		fmt.Println("Failed to create k8s client: " + err.Error())
		panic(err)
	}

	return nil
}, func(data []byte) {})

func TestDevfileRegistryController(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)

	fmt.Println("Creating ginkgo reporter for Test Harness: Junit and Debug Detail reporter")
	var r []ginkgo.Reporter
	r = append(r, reporters.NewJUnitReporter(filepath.Join(testResultsDirectory, jUnitOutputFilename)))

	fmt.Println("Running Devfile Registry integration tests...")
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "Devfile Registry Operator Tests", r)
}
