//
//
// Copyright 2022-2023 Red Hat, Inc.
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

package v1alpha1

import (
	"context"
	"fmt"

	. "github.com/devfile/registry-operator/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("DevfileRegistriesList validation webhook", func() {

	Context("Create DevfileRegistriesList CR with valid values", func() {
		It("Should create a new CR in a non-default namespace called 'main'", func() {
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, getDevfileRegistriesListCR(devfileRegistriesListName, devfileRegistriesNamespace,
				devfileStagingRegistryName, devfileStagingRegistryURL))).Should(Succeed())
		})
	})

	Context("Create DevfileRegistriesList CR in forbidden default namespace", func() {
		It("Should fail to create a new CR in the default namespace", func() {
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, getDevfileRegistriesListCR("default-namespace-list", "default",
				devfileStagingRegistryName, devfileStagingRegistryURL))).ShouldNot(Succeed())
		})
	})

	Context("Update DevfileRegistriesList CR with an invalid URL", func() {
		It("Should fail to update and issue an invalid registry URL error message", func() {
			drlLookupKey := types.NamespacedName{Name: devfileRegistriesListName, Namespace: devfileRegistriesNamespace}
			err := appendToDevfileRegistriesService(drlLookupKey, "registryName", "registryURL", NamespaceListType)
			Expect(err.Error()).Should(ContainSubstring(fmt.Sprintf(InvalidRegistry, "registryURL")))
		})
	})

	Context("Update DevfileRegistriesList CR with a registry with v2index ", func() {
		It("Should succeed with entry added to the CR", func() {
			//start mock index server
			testServer := GetNewUnstartedTestServer()
			Expect(testServer).ShouldNot(BeNil())
			testServer.Start()
			defer testServer.Close()
			drlLookupKey := types.NamespacedName{Name: devfileRegistriesListName, Namespace: devfileRegistriesNamespace}
			Expect(appendToDevfileRegistriesService(drlLookupKey, "localRegistry", testServer.URL, NamespaceListType)).Should(Succeed())
		})
	})

	Context("Update DevfileRegistriesList CR with a deleted entry ", func() {
		It("Should succeed with entry added to the CR", func() {
			drlLookupKey := types.NamespacedName{Name: devfileRegistriesListName, Namespace: devfileRegistriesNamespace}
			Expect(deleteFromDevfileRegistriesService(drlLookupKey, localRegistryName, NamespaceListType)).Should(Succeed())
		})
	})

	Context("Create a second DevfileRegistry CR with valid values in same namespace", func() {
		It("Should fail to create a new CR and return an error message", func() {
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, getDevfileRegistriesListCR(devfileRegistriesListName+"2", devfileRegistriesNamespace,
				devfileStagingRegistryName, devfileStagingRegistryURL))).ShouldNot(Succeed())
		})
	})

	Context("Create a second DevfileRegistry CR with valid fields in a different namespace", func() {
		It("Should fail to create a new CR and return an error message", func() {
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, getDevfileRegistriesListCR(devfileRegistriesListName, testNs.Name,
				devfileStagingRegistryName, devfileStagingRegistryURL))).Should(Succeed())
			//delete all crs
			drlLookupKey := types.NamespacedName{Name: devfileRegistriesListName, Namespace: devfileRegistriesNamespace}
			deleteCRList(drlLookupKey, NamespaceListType)
			drlLookupKey = types.NamespacedName{Name: devfileRegistriesListName, Namespace: testNs.Name}
			deleteCRList(drlLookupKey, NamespaceListType)
		})
	})
})

// getDevfileRegistriesListCR returns a minimally populated DevfileRegistriesList object for testing
func getDevfileRegistriesListCR(name string, namespace string, registryName string, registryURL string) *DevfileRegistriesList {

	return &DevfileRegistriesList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: ApiVersion,
			Kind:       "DevfileRegistriesList",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: DevfileRegistriesListSpec{
			DevfileRegistries: []DevfileRegistryService{
				{
					Name: registryName,
					URL:  registryURL,
				},
			},
		},
	}

}
