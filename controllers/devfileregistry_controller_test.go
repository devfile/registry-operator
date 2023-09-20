/*
Copyright Red Hat

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

package controllers

import (
	"github.com/devfile/registry-operator/api/v1alpha1"
	. "github.com/devfile/registry-operator/pkg/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
)

// The following test is only intended to test creation and updates to the DevfileRegistry CR, while ignoring the fact that routes will fail because they will
// not resolve to an actual URL

var _ = Describe("DevfileRegistry controller test", func() {
	Context("When Creating and updating a DevfileRegistry CR with valid values", func() {
		It("Should create and update a DevfileRegistry CR Name with Status URL not found", func() {
			Expect(k8sClient.Create(ctx, getDevfileRegistryCR(devfileRegistryName, devfileRegistriesNamespace,
				image))).Should(Succeed())
			drlLookupKey := types.NamespacedName{Name: devfileRegistryName, Namespace: devfileRegistriesNamespace}
			dr := &v1alpha1.DevfileRegistry{}
			Eventually(func() error {
				return k8sClient.Get(ctx, drlLookupKey, dr)
			}, Timeout, Interval).Should(Succeed())

			Expect(dr.Name).Should(Equal(devfileRegistryName))
			Expect(dr.Status.URL).Should(Equal("")) // an empty URL is an indicator that a URL was not resolved

			// update values
			dr.Spec.Telemetry.RegistryName = devfileRegistryName
			dr.Spec.Telemetry.Key = "abcdefghijk"
			dr.Spec.DevfileIndex.Image = "quay.io/xyz/devfile-index:next"

			Expect(k8sClient.Update(ctx, dr)).Should(Succeed())

			//try to update the name - it's immutable, so it should fail
			dr.Name = devfileRegistryName + "1"
			Expect(k8sClient.Update(ctx, dr)).ShouldNot(Succeed())

			//delete all crs
			deleteCRList(drlLookupKey, DevfileRegistryType)
		})
	})

})
