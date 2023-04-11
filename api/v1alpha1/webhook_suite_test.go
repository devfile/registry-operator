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

package v1alpha1

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"path/filepath"
	"testing"
	"time"

	. "github.com/devfile/registry-operator/pkg/test"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	//+kubebuilder:scaffold:imports
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var ctx context.Context
var cancel context.CancelFunc

const (
	devfileRegistriesListName  = "main-namespace-list"
	devfileRegistriesNamespace = "main"
	devfileStagingRegistryName = "StagingRegistry"
	devfileStagingRegistryURL  = "https://registry.stage.devfile.io"
	localRegistryName          = "localRegistry"
)

var (
	devfileRegistriesNs = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: devfileRegistriesNamespace,
		},
	}
	testNs = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webhook Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: false,
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{filepath.Join("..", "..", "config", "webhook")},
		},
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	scheme := runtime.NewScheme()
	err = AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	err = admissionv1beta1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	err = corev1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
	//create devfileregistries namespace
	Expect(k8sClient.Create(ctx, devfileRegistriesNs)).Should(Succeed())
	//create test namespace
	Expect(k8sClient.Create(ctx, testNs)).Should(Succeed())

	// start webhook server using Manager
	webhookInstallOptions := &testEnv.WebhookInstallOptions
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme,
		Host:               webhookInstallOptions.LocalServingHost,
		Port:               webhookInstallOptions.LocalServingPort,
		CertDir:            webhookInstallOptions.LocalServingCertDir,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	Expect(err).NotTo(HaveOccurred())

	err = (&DevfileRegistry{}).SetupWebhookWithManager(mgr)
	Expect(err).NotTo(HaveOccurred())

	err = (&DevfileRegistriesList{}).SetupWebhookWithManager(mgr)
	Expect(err).NotTo(HaveOccurred())

	err = (&ClusterDevfileRegistriesList{}).SetupWebhookWithManager(mgr)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:webhook

	go func() {
		err = mgr.Start(ctx)
		if err != nil {
			Expect(err).NotTo(HaveOccurred())
		}
	}()

	// wait for the webhook server to get ready
	dialer := &net.Dialer{Timeout: time.Second}
	addrPort := fmt.Sprintf("%s:%d", webhookInstallOptions.LocalServingHost, webhookInstallOptions.LocalServingPort)
	Eventually(func() error {
		conn, err := tls.DialWithDialer(dialer, "tcp", addrPort, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}
		conn.Close()
		return nil
	}).Should(Succeed())

}, 60)

var _ = AfterSuite(func() {
	// delete the devfileregistries namespace
	Expect(k8sClient.Delete(ctx, devfileRegistriesNs)).Should(Succeed())
	// delete the test namespace
	Expect(k8sClient.Delete(ctx, testNs)).Should(Succeed())
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

// deleteCRList removes the cluster or namespace CR list from the cluster
func deleteCRList(drlLookupKey types.NamespacedName, f ListType) {

	cl := &ClusterDevfileRegistriesList{}
	nl := &DevfileRegistriesList{}
	// Delete
	Eventually(func() error {
		if f == NamespaceListType {
			k8sClient.Get(context.Background(), drlLookupKey, nl)
			return k8sClient.Delete(context.Background(), nl)
		} else {
			k8sClient.Get(context.Background(), drlLookupKey, cl)
			return k8sClient.Delete(context.Background(), cl)
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

// deleteFromDevfileRegistriesListCR validates that a  DevfileRegistryService object can be deleted from the list as an update to the CR
func deleteFromDevfileRegistriesService(lookupKey types.NamespacedName, rName string, lType ListType) error {
	ctx := context.Background()
	nl := &DevfileRegistriesList{}
	cl := &ClusterDevfileRegistriesList{}
	var registriesList []DevfileRegistryService
	var err error

	if lType == NamespaceListType {
		err = k8sClient.Get(ctx, lookupKey, nl)
		if err != nil {
			return err
		}
		registriesList = nl.Spec.DevfileRegistries
	} else {
		err = k8sClient.Get(ctx, lookupKey, cl)
		if err != nil {
			return err
		}
		registriesList = cl.Spec.DevfileRegistries
	}

	//update list in existing CR
	var newList []DevfileRegistryService
	newList = make([]DevfileRegistryService, 0, len(registriesList))

	for i := range registriesList {
		if registriesList[i].Name != rName {
			newList = append(newList, registriesList[i])
		}
	}

	if lType == NamespaceListType {
		nl.Spec.DevfileRegistries = newList
		err = k8sClient.Update(ctx, nl)
	} else {
		cl.Spec.DevfileRegistries = newList
		err = k8sClient.Update(ctx, cl)
	}

	return err
}

// appendToDevfileRegistriesService validates that a new DevfileRegistryService object can be added to update an existing CR
func appendToDevfileRegistriesService(lookupKey types.NamespacedName, rName string, rUrl string, lType ListType) error {
	ctx := context.Background()
	cl := &ClusterDevfileRegistriesList{}
	nl := &DevfileRegistriesList{}

	var err error
	if lType == NamespaceListType {
		err = k8sClient.Get(ctx, lookupKey, nl)
		if err != nil {
			return err
		}
		//update list in existing CR
		registriesList := nl.Spec.DevfileRegistries
		registriesList = append(registriesList, DevfileRegistryService{Name: rName, URL: rUrl})
		nl.Spec.DevfileRegistries = registriesList
		err = k8sClient.Update(ctx, nl)
	} else {
		err = k8sClient.Get(ctx, lookupKey, cl)
		if err != nil {
			return err
		}
		//update list in existing CR
		registriesList := cl.Spec.DevfileRegistries
		registriesList = append(registriesList, DevfileRegistryService{Name: rName, URL: rUrl})
		cl.Spec.DevfileRegistries = registriesList
		err = k8sClient.Update(ctx, cl)
	}
	return err
}
