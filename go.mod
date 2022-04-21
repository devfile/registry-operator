module github.com/devfile/registry-operator

go 1.13

require (
	github.com/go-logr/logr v0.4.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/openshift/api v0.0.0-20200205133042-34f0ec8dab87
	github.com/prometheus/common v0.26.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.5
	k8s.io/apiextensions-apiserver v0.22.5 // indirect
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	sigs.k8s.io/controller-runtime v0.10.3
)
