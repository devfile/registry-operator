
# Devfile Registry Operator

The Devfile Registry operator manages the lifecycle of the following custom resources:
1. [Devfile Registry](DEVFILE_REGISTRY.md)
2. [Devfile Registries List](REGISTRIES_LISTS.md)
3. [Cluster Devfile Registries List](REGISTRIES_LISTS.md)

## Issue Tracking

Issue tracking repo: https://github.com/devfile/api with label area/registry

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

Some of the rules supported by the makefile:

|rule|purpose|
|---|---|
| docker-build | build registry operator docker image |
| docker-push | push registry operator docker image |
| deploy | deploy operator to cluster |
| undeploy | undeploy operator from cluster |
| install | create the devfile registry CRDs on the cluster |
| uninstall | remove the devfile registry operator and CRDs from the cluster |
| install-cert | for validating webhooks, install cert manager on the cluster |
| uninstall-cert | for validating webhooks, remove cert manager from the cluster |
| manifests | Generate manifests e.g. CRD, RBAC etc. |
| generate | Generate the API type definitions. Must be run after modifying the DevfileRegistry type. |
| bundle | Generate bundle manifests and metadata, then validate generated files. |
| test_integration | Run the integration tests for the operator. |
| test | Run the unit tests for the operator. |

To see all rules supported by the makefile, run `make help`

## Testing

To run integration tests for the operator, run `make test-integration`. 

By default, the tests will use the default image for the operator, `quay.io/devfile/registry-operator:next`. To use your own image, run:

```
export IMG=<your-operator-image>
make test-integration
```

### Run operator locally
It's possible to run an instance of the operator locally while communicating with a cluster. 

1. Build the binary

```bash
make manager
```

2. Run the controller

```bash
export NAMESPACE=devfileregistry-operator
make run ENABLE_WEBHOOKS=false
```