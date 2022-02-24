# Devfile Registry Operator

Devfile Registry operator repository that contains the operator for the DevfileRegistry Custom Resource. 

## Issue Tracking

Issue tracking repo: https://github.com/devfile/api with label area/registry

## Running the controller in a cluster

The controller can be deployed to a cluster provided you are logged in with cluster-admin credentials:

```bash
make install && make deploy
```

The operator will be installed under the `registry-operator-system` namespace. However, devfile registries can be deployed in any namespace.

## Deploying a Devfile Registry

Once the Devfile Registry operator has been deployed to a cluster, it's straightforward to deploy a Devfile Registry. The following samples below showcase how the registry can be deployed on to an OpenShift or Kubernetes cluster. 

In addition to the examples below, the `samples/` folder in this repo provides some example devfile registry yaml files for convenience.


### OpenShift

Installing the devfile registry via the devfile registry operator on an OpenShift cluster can be done in one easy command:

```bash
$ cat <<EOF | oc apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistry
metadata:
  name: devfile-registry
spec:
  devfileIndexImage: quay.io/devfile/devfile-index:next
  telemetry:
    registryName: test
EOF
```


### Kubernetes

Installing the devfile registry on a Kubernetes cluster is similar, but requires setting the `k8s.ingressDomain` field first.

```bash
$ export INGRESS_DOMAIN=<my-ingress-domain>

$ cat <<EOF | kubectl apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistry
metadata:
  name: devfile-registry
spec:
  devfileIndexImage: quay.io/devfile/devfile-index:next
  tls:
    enabled: false
  k8s:
    ingressDomain: $INGRESS_DOMAIN
  telemetry:
    registryName: test
EOF
```

## Telemetry
If you want to send telemetry information to your own Segment instance, specify the write key in the telemetry object

```bash
  telemetry:
    registryName: test
    key: <your-segment-write-key>
```


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
| install | create the devfile registry CRDs on the cluster |
| uninstall | remove the devfile registry operator and CRDs from the cluster |
| manifests | Generate manifests e.g. CRD, RBAC etc. |
| generate | Generate the API type definitions. Must be run after modifying the DevfileRegistry type. |
| bundle | Generate bundle manifests and metadata, then validate generated files. |
| test_integration | Run the integration tests for the operator. |

To see all rules supported by the makefile, run `make help`

## Testing

To run integration tests for the operator, run `make test_integration`. 

By default, the tests will use the default image for the operator, `quay.io/devfile/registry-operator:next`. To use your own image, run:

```
export IMG=<your-operator-image>
make test_integration
```

### Run operator locally
It's possible to run an instance of the operator locally while communicating with a cluster. 

```bash
export NAMESPACE=devfileregistry-operator
make run ENABLE_WEBHOOKS=false
```