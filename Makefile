#
# Copyright Red Hat
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Current Operator version
VERSION ?= `cat $(PWD)/VERSION`
# Default bundle image tag
BUNDLE_IMG ?= quay.io/devfile/registry-operator-bundle:v$(VERSION)
CERT_MANAGER_VERSION ?= v1.11.0
ENABLE_WEBHOOKS ?= true
ENABLE_WEBHOOK_HTTP2 ?= false

# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# Image URL to use all building/pushing image targets
IMG ?= quay.io/devfile/registry-operator:next
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.26
# Controller tools version number
CONTROLLER_TOOLS_VERSION ?= v0.9.2
# Kustomize version number
KUSTOMIZE_VERSION ?= v3.8.7

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Check if oc or kubectl are installed and determine which of the two to use
ifeq ($(K8S_CLI),)
ifeq (,$(shell which kubectl))
ifeq (,$(shell which oc))
$(error oc or kubectl is required to proceed)
else
K8S_CLI := oc
endif
else
K8S_CLI := kubectl
endif
endif

# operator-sdk
OPERATOR_SDK_CLI ?= operator-sdk


# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build Dependencies

## Local binary path
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Target platform
TARGET_OS ?= linux # `make manager` only
TARGET_ARCH ?= amd64

##@ Development

.PHONY: test
test: manifests generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test -timeout=1h `go list ./... | grep -v /tests/` -coverprofile cover.out -v

### test-integration: runs integration tests on the cluster set in context.
.PHONY: test-integration
test-integration:
	CGO_ENABLED=0 go test -v -c -o $(LOCALBIN)/devfileregistry-operator-integration ./tests/integration/cmd/devfileregistry_test.go
# Create junit report output file if does not exist
	mkdir -p /tmp/artifacts && touch /tmp/artifacts/junit-devfileregistry-operator.xml
	$(LOCALBIN)/devfileregistry-operator-integration -ginkgo.fail-fast --ginkgo.junit-report=/tmp/artifacts/junit-devfileregistry-operator.xml

##@ Build

.PHONY: build manager
manager: manifests generate fmt vet ## Build manager binary.
	GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} go build -o $(LOCALBIN)/manager ./main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	$(LOCALBIN)/manager

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

### check_fmt: Checks the formatting on files in repo
.PHONY: check_fmt
check_fmt:
  ifeq ($(shell command -v goimports 2> /dev/null),)
	  $(error "goimports must be installed for this rule" && exit 1)
  endif
  ifeq ($(shell command -v addlicense 2> /dev/null),)
	  $(error "error addlicense must be installed for this rule: go install github.com/google/addlicense")
  endif

	  if [[ $$(find . -not -path '*/\.*' -not -name '*zz_generated*.go' -name '*.go' -exec goimports -l {} \;) != "" ]]; then \
	    echo "Files not formatted; run 'make fmt'"; exit 1 ;\
	  fi ;\
	  if ! addlicense -check -f license_header.txt $$(find . -not -path '*/\.*' -name '*.go'); then \
	    echo "Licenses are not formatted; run 'make fmt_license'"; exit 1 ;\
	  fi \

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

### fmt_license: ensure license header is set on all files
.PHONY: fmt_license
fmt_license:
ifneq ($(shell command -v addlicense 2> /dev/null),)
	@echo 'addlicense -v -f license_header.txt **/*.go'
	@addlicense -v -f license_header.txt $$(find . -name '*.go')
else
	$(error addlicense must be installed for this rule: go install github.com/google/addlicense)
endif

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
.PHONY: docker-build
docker-build:
	docker build . -t ${IMG} --build-arg TARGETARCH=${TARGET_ARCH} \
--build-arg ENABLE_WEBHOOKS=${ENABLE_WEBHOOKS} \
--build-arg ENABLE_WEBHOOK_HTTP2=${ENABLE_WEBHOOK_HTTP2}

# Push the docker image
.PHONY: docker-push
docker-push:
	docker push ${IMG}

# PLATFORMS defines the target platforms for the manager image be build to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=quay.io/devfile/registry-operator:next). To use this option you need to:
# - able to use docker buildx . More info: https://docs.docker.com/build/buildx/
# - have enable BuildKit, More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image for your registry (i.e. if you do not inform a valid value via IMG=quay.io/<user>/registry-operator:next than the export will fail)
# To properly provided solutions that supports more than one platform you should use this option.
# **docker-buildx does not work with podman**
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: test ## Build and push docker image for the manager for cross-platform support
# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- docker buildx create --name registry-operator-builder
	docker buildx use registry-operator-builder
	- docker buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross $(shell pwd)
	- docker buildx rm registry-operator-builder
	rm Dockerfile.cross

# Build the podman image
.PHONY: podman-build
podman-build:
	podman build . -t ${IMG} --build-arg TARGETARCH=${TARGET_ARCH} \
--build-arg ENABLE_WEBHOOKS=${ENABLE_WEBHOOKS} \
--build-arg ENABLE_WEBHOOK_HTTP2=${ENABLE_WEBHOOK_HTTP2}

# Push the podman image
.PHONY: podman-push
podman-push:
	podman push ${IMG}


.PHONY: install-cert
install-cert: ## Install cert manager for webhooks
	$(K8S_CLI) apply -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml

.PHONY: uninstall-cert
uninstall-cert:
	$(K8S_CLI) delete -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml

# find or download controller-gen
# download controller-gen if necessary
.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN)
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen && $(LOCALBIN)/controller-gen --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
GOBIN=$(LOCALBIN) GOFLAGS="" go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

##@ Deployment

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | $(K8S_CLI) apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | $(K8S_CLI) delete -f -

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | $(K8S_CLI) apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | $(K8S_CLI) delete -f -

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE)
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) GOFLAGS="" go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

# Generate bundle manifests and metadata, then validate generated files.
## sh replace-alm-examples.sh is there to fix an issue with Kustomize removing List items
## More information can be found here: https://github.com/kubernetes-sigs/kustomize/issues/5042
.PHONY: bundle
bundle: manifests
	$(OPERATOR_SDK_CLI) generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK_CLI) generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	sh replace-alm-examples.sh
	$(OPERATOR_SDK_CLI) bundle validate ./bundle

# Build the bundle image.
.PHONY: bundle-build
bundle-build:
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

# Push the bundle image.
.PHONY: bundle-push
bundle-push:
	docker push $(BUNDLE_IMG)

### gosec - runs the gosec scanner for non-test files in this repo
.PHONY: gosec
gosec:
	# Run this command to install gosec, if not installed:
	# go install github.com/securego/gosec/v2/cmd/gosec@v2.14.0
	gosec -no-fail -fmt=sarif -out=gosec.sarif -exclude-dir pkg/test -exclude-dir tests ./...

### Release
# RUN: make release NEW_VERSION=x.x.x
.PHONY: 
release:
	sh make-release.sh ${NEW_VERSION}