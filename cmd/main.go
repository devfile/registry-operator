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

package main

import (
	"crypto/tls"
	"flag"
	"os"

	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	"github.com/devfile/registry-operator/internal/controller"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(routev1.AddToScheme(scheme))
	utilruntime.Must(registryv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", "0", "The metrics address: :8443 for HTTPS, :8080 for HTTP, or 0 to disable.")
	flag.BoolVar(&secureMetrics, "metrics-secure", true, "Serve metrics over HTTPS with authentication and authorization.")
	flag.BoolVar(&enableHTTP2, "enable-http2", false, "Enable HTTP/2 on the metrics and webhook servers.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	setupLog.Info(("logger set up"))
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// Configure TLS before starting either server. Preserve the existing webhook override.
	disableHTTP2 := func(c *tls.Config) { c.NextProtos = []string{"http/1.1"} }
	var tlsOpts []func(*tls.Config)
	if !enableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}
	webhookTLSOpts := tlsOpts
	if value, set := os.LookupEnv("ENABLE_WEBHOOK_HTTP2"); set {
		webhookTLSOpts = nil
		if value != "true" {
			webhookTLSOpts = append(webhookTLSOpts, disableHTTP2)
		}
	}
	metricsOptions := newMetricsOptions(metricsAddr, secureMetrics, tlsOpts)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		LeaderElection:         enableLeaderElection,
		Scheme:                 scheme,
		HealthProbeBindAddress: probeAddr,
		LeaderElectionID:       "1984829e.devfile.io",
		WebhookServer: webhook.NewServer(webhook.Options{
			Port:    9443,
			TLSOpts: webhookTLSOpts,
		}),
		Metrics: metricsOptions,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controller.DevfileRegistryReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("DevfileRegistry"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DevfileRegistry")
		os.Exit(1)
	}
	if err = (&controller.DevfileRegistriesListReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("DevfileRegistriesList"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DevfileRegistriesList")
		os.Exit(1)
	}
	if err = (&controller.ClusterDevfileRegistriesListReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("ClusterDevfileRegistriesList"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ClusterDevfileRegistriesList")
		os.Exit(1)
	}

	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		setupLog.Info("setting up webhooks")
		if err = (&registryv1alpha1.DevfileRegistry{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "DevfileRegistry")
			os.Exit(1)
		}
		if err = (&registryv1alpha1.DevfileRegistriesList{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "DevfileRegistriesList")
			os.Exit(1)
		}
		if err = (&registryv1alpha1.ClusterDevfileRegistriesList{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "ClusterDevfileRegistriesList")
			os.Exit(1)
		}
	}

	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// newMetricsOptions keeps HTTPS metrics protected by Kubernetes authentication and authorization.
func newMetricsOptions(address string, secure bool, tlsOpts []func(*tls.Config)) metricsserver.Options {
	options := metricsserver.Options{BindAddress: address, SecureServing: secure, TLSOpts: tlsOpts}
	if secure {
		options.FilterProvider = filters.WithAuthenticationAndAuthorization
	}
	return options
}
