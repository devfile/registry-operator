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

package controllers

import (
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"path/filepath"
	"testing"

	. "github.com/devfile/registry-operator/api/v1alpha1"
	. "github.com/devfile/registry-operator/pkg/test"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"context"

	"github.com/devfile/registry-operator/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var testEnv *envtest.Environment
var ctx context.Context
var cancel context.CancelFunc
var k8sClient client.Client

const (
	clusterdevfileRegistriesListName = "cluster-default-namespace-list"
	devfileRegistriesListName        = "default-namespace-list"
	devfileRegistriesNamespace       = "default"
	devfileStagingRegistryName       = "StagingRegistry"
	devfileStagingRegistryURL        = "https://registry.stage.devfile.io"
	localRegistryName                = "LocalRegistry"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = v1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})

	err = (&DevfileRegistriesListReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("DevfileRegistriesList"),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&ClusterDevfileRegistriesListReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("ClusterDevfileRegistriesList"),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()

})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

// getClusterDevfileRegistriesListCR returns a minimally populated DevfileRegistriesList object for testing
func getClusterDevfileRegistriesListCR(name string, namespace string, registryName string, registryURL string) *ClusterDevfileRegistriesList {
	return &ClusterDevfileRegistriesList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: ApiVersion,
			Kind:       string(ClusterListType),
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

// getDevfileRegistriesListCR returns a minimally populated DevfileRegistriesList object for testing
func getDevfileRegistriesListCR(name string, namespace string, registryName string, registryURL string) *DevfileRegistriesList {

	return &DevfileRegistriesList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: ApiVersion,
			Kind:       string(NamespaceListType),
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

// deleteCRList removes the cluster or namespace CR list from the cluster
func deleteCRList(drlLookupKey types.NamespacedName, f ListType) {

	cl := &ClusterDevfileRegistriesList{}
	nl := &DevfileRegistriesList{}

	// Delete
	Eventually(func() error {
		if f == ClusterListType {
			k8sClient.Get(context.Background(), drlLookupKey, cl)
			return k8sClient.Delete(context.Background(), cl)
		} else {
			k8sClient.Get(context.Background(), drlLookupKey, nl)
			return k8sClient.Delete(context.Background(), nl)
		}
	}, Timeout, Interval).Should(Succeed())

	// Wait for delete to finish
	Eventually(func() error {
		if f == ClusterListType {
			return k8sClient.Get(context.Background(), drlLookupKey, cl)
		} else {
			return k8sClient.Get(context.Background(), drlLookupKey, nl)
		}
	}, Timeout, Interval).ShouldNot(Succeed())

}

// validateStatus validates the controller status
func validateStatus(lookupKey types.NamespacedName, lType ListType, expectedStatus string) {
	cl := v1alpha1.ClusterDevfileRegistriesList{}
	nl := v1alpha1.DevfileRegistriesList{}
	var status string

	Eventually(func() string {
		if lType == NamespaceListType {
			k8sClient.Get(ctx, lookupKey, &nl)
			if condition := meta.FindStatusCondition(nl.Status.Conditions, typeValidateDevfileRegistries); condition != nil {
				status = condition.Message
			} else {
				status = ""
			}
		} else {
			k8sClient.Get(ctx, lookupKey, &cl)
			if condition := meta.FindStatusCondition(cl.Status.Conditions, typeValidateDevfileRegistries); condition != nil {
				status = condition.Message
			} else {
				status = ""
			}
		}
		return status
	}, Timeout, Interval).Should(ContainSubstring(expectedStatus))
}
