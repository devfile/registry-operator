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

package v1alpha1

import (
	"context"

	. "github.com/devfile/registry-operator/pkg/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("DevfileRegistry validation webhook", func() {
	Context("Create DevfileRegistry CR with valid values", func() {
		It("Should create a new CR in a non-default namespace called 'main'", func() {
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, getDevfileRegistryCR("devfileregistry", devfileRegistriesNamespace))).Should(Succeed())
		})
	})

	Context("Create DevfileRegistry CR in forbidden default namespace", func() {
		It("Should fail to create a new CR in the default namespace", func() {
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, getDevfileRegistryCR("default-devfileregistry", "default"))).ShouldNot(Succeed())
		})
	})
})

func getDevfileRegistryCR(name string, namespace string) *DevfileRegistry {
	return &DevfileRegistry{
		TypeMeta: metav1.TypeMeta{
			APIVersion: ApiVersion,
			Kind:       "DevfileRegistry",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}
