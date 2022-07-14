module github.com/devfile/registry-operator

go 1.13

require (
	github.com/devfile/registry-support/index/generator v0.0.0-20220222194908-7a90a4214f3e
	github.com/devfile/registry-support/registry-library v0.0.0-20220505191145-973af4b38f31
	github.com/go-logr/logr v0.4.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/openshift/api v0.0.0-20200930075302-db52bc4ef99f
	github.com/prometheus/common v0.26.0
	github.com/stretchr/testify v1.7.2
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.5
	k8s.io/apiextensions-apiserver v0.22.5 // indirect
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	sigs.k8s.io/controller-runtime v0.10.3
)
