
# Devfile Registry Operator

<div id="header">

[![Apache2.0 License](https://img.shields.io/badge/license-Apache2.0-brightgreen.svg)](LICENSE)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/8256/badge)](https://www.bestpractices.dev/projects/8256)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/devfile/registry-operator/badge)](https://securityscorecards.dev/viewer/?uri=github.com/devfile/registry-operator)
</div>

The Devfile Registry operator manages the lifecycle of the following custom resources:
1. [Devfile Registry](DEVFILE_REGISTRY.md)
2. [Devfile Registries List](REGISTRIES_LISTS.md)
3. [Cluster Devfile Registries List](REGISTRIES_LISTS.md)

## Releases

**Minor releases** will be scheduled roughly on a _quarterly basis_ and will include non-feature breaking changes made between release cycles.

**Patch releases** are provided as needed for _critical bug fixes_ and _security patches_, it is **strongly recommended** for users running on the same minor version of the registry operator to update when these releases become available.

Releases are available on [GitHub](https://github.com/devfile/registry-operator/releases) along with bundle entries on [OperatorHub.io](https://operatorhub.io) and the OpenShift OperatorHub community catalog.

For more updates on releases, please join our [communication channels](https://devfile.io/docs/2.2.2/community#getting-involved).

## Issue Tracking

Issue tracking repo: https://github.com/devfile/api with label area/registry

## Changelog

This repository utilizes Release Notes to track and display changes, a changelog for every release can be found [here](https://github.com/devfile/registry-operator/releases).

## Requirements

Deployment cluster must meet one of the following criteria:

- OpenShift Container Platform (OCP) 4.12.x
- Kubernetes 1.25.x-1.26.x

More on the support of container orchestration systems can be 
found [here](CLUSTER_SUPPORT.md).

Deployments made by the devfile registry operator must *never* target the default namespace due to incompatibility with 
security setups.

### Development

- Go 1.19.x
- Docker / Podman
- Operator SDK 1.28.x

See [Upgrade SDK Version](https://sdk.operatorframework.io/docs/upgrading-sdk-version/) for a guide on updating the Operator SDK. Ensure the Operator SDK version and the version of Kubernetes APIs match each other when updating by checking [CLUSTER_SUPPORT.md](CLUSTER_SUPPORT.md).

## Running the controller in a cluster

Install cert-manager to provision self-signed certificates for the validating webhooks which are used specifically for the `ClusterDevfileRegistriesList` and `DevfileRegistriesList` CRs.  Cert manager needs to be installed in order for the controller manager to start.

```bash
make install-cert
```

The controller can be deployed to a cluster provided you are logged in with cluster-admin credentials:

```bash
make install && make deploy
```

The operator will be installed under the `registry-operator-system` namespace. However, devfile registries can be deployed in any namespace.


## Development

The repository contains a Makefile; building and deploying can be configured via the environment variables

|variable|purpose|default value|
|---|---|---|
| `IMG` | Image used for controller (run makefile, if `IMG` is updated) | `quay.io/devfile/registry-operator:next` |
| `BUNDLE_IMG` | Image used for bundle OLM package | `quay.io/devfile/registry-operator-bundle:<latest_version>` |
| `CERT_MANAGER_VERSION` | Version of `cert-manager` installed using `make install-cert` | `v1.11.0` |
| `ENABLE_WEBHOOKS` | If `false`, disables operator webhooks | `true` |
| `ENABLE_WEBHOOK_HTTP2` | Overrides webhook HTTP server deployment to use http/2 if set to `true`, **not recommended** | `false` |
| `BUNDLE_CHANNELS` | Sets the list channel(s) include bundle build under | `alpha` |
| `BUNDLE_DEFAULT_CHANNEL` | Sets the default channel to use when installing the bundle | |
| `ENVTEST_K8S_VERSION` | Version of k8s to use for the test environment | `1.26` (current) |
| `CONTROLLER_TOOLS_VERSION` | Version of the controller tools | `v0.9.2` |
| `KUSTOMIZE_VERSION` | Version of kustomize | `v3.8.7` | 
| `GOBIN` | Path to install Go binaries to | `${GOPATH}/bin` |
| `K8S_CLI` | Path to CLI tool to use with the target cluster environment, `kubectl` or `oc` | Either `oc` or `kubectl` if installed in that order |
| `OPERATOR_SDK_CLI` | CLI path to `operator-sdk` tool | `operator-sdk` |
| `SHELL` | Active shell to use with make | `/usr/bin/env bash -o pipefail` |
| `LOCALBIN` | Path to place project binaries | `./bin` |
| `KUSTOMIZE` | Path to target `kustomize` binary | `${LOCALBIN}/kustomize` |
| `CONTROLLER_GEN` | Path to target `controller-gen` binary | `${LOCALBIN}/controller-gen` |
| `ENVTEST` | Path to target `setup-envtest` binary | `${LOCALBIN}/setup-envtest` |
| `TARGET_ARCH` | Target architecture for operator manager builds, possible values: `amd64`, `arm64`, `s390x`, `ppc64le` | `amd64` |
| `TARGET_OS` | Target operating system for operator manager build, **only for `make manager`** | `linux` |
| `PLATFORMS` | Target architecture(s) for `make docker-buildx` | All supported: `linux/arm64,linux/amd64,linux/s390x,linux/ppc64le` |
| `KUSTOMIZE_INSTALL_SCRIPT` | URL of kustomize installation script, see [kustomize installation instructions](https://kubectl.docs.kubernetes.io/installation/kustomize/binaries/) | `https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh` |

Some of the rules supported by the makefile:

|rule|purpose|
|---|---|
| controller-gen | install the controller-gen tool, used by other commands |
| kustomize | install the kustomize tool, used by other commands |
| docker-build | build registry operator container image using docker |
| docker-push | push registry operator container image using docker |
| docker-buildx | build & push registry operator docker image for all supported architectures |
| podman-build | build registry operator container image using podman |
| podman-push | push registry operator container image using podman |
| deploy | deploy operator to cluster |
| undeploy | undeploy operator from cluster |
| install | create the devfile registry CRDs on the cluster |
| uninstall | remove the devfile registry operator and CRDs from the cluster |
| install-cert | for validating webhooks, install cert manager on the cluster |
| uninstall-cert | for validating webhooks, remove cert manager from the cluster |
| manifests | Generate manifests e.g. CRD, RBAC etc. |
| generate | Generate the API type definitions. Must be run after modifying the DevfileRegistry type. |
| bundle | Generate bundle manifests and metadata, then validate generated files. |
| test-integration | Run the integration tests for the operator. |
| test | Run the unit tests for the operator. |
| fmt | Check code formatting |
| fmt_license | Ensure license header is set on all files |
| vet | Check suspicious constructs into code |
| gosec | Check for security problems in non-test source files |

To see all rules supported by the makefile, run `make help`

## Testing

To run integration tests for the operator, run `make test-integration`. 

One of the `oc` or `kubectl` executables must be accessible. If both are present in your path, `oc` will be used, except if you
define the environment variable `K8S_CLI` with the command you prefer to use.

By default, the tests will use the default image for the operator, `quay.io/devfile/registry-operator:next`.

You can use `make <engine>-build` to build your own image, `make <engine>-push` to publish it - Replace `<engine>` with `podman` or `docker`. Then, to use your own image, run:

```
IMG=<your-operator-image> make test-integration
```

### Run operator locally
It's possible to run an instance of the operator locally while communicating with a cluster. 

1. You may need to install the `controller-gen` tool before, used when building the binary:

```bash
make controller-gen
```

2. Build the binary

```bash
make manager
```

3. Deploy the CRDs

```bash
make install
```

4. Run the controller

```bash
make run ENABLE_WEBHOOKS=false
```

## Contributing

Please see our [contributing.md](./CONTRIBUTING.md).
