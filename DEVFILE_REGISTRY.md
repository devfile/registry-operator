# Devfile Registry

The Devfile Registry custom resource allows you to create and manage your own index server and registry viewer.

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
