# Devfile Registry

The Devfile Registry custom resource allows you to create and manage your own index server and registry viewer.

## Deploying a Devfile Registry

Once the Devfile Registry operator has been deployed to a cluster, you can deploy a Devfile Registry by creating custom resources. The following samples below showcase how the registry can be deployed on to an OpenShift or Kubernetes cluster.

In addition to the examples below, the `samples/` folder in this repo provides some example devfile registry yaml files for convenience.


### OpenShift

Installing the devfile registry via the devfile registry operator on an OpenShift cluster can be done with this command:

```bash
$ cat <<EOF | oc apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistry
metadata:
  name: devfile-registry
spec:
  devfileIndex:
    image: quay.io/devfile/devfile-index:next
  telemetry:
    registryName: test
EOF
```


### Kubernetes

Installing the devfile registry on a Kubernetes cluster is similar, but requires setting the `k8s.ingressDomain` field.

```bash
$ export INGRESS_DOMAIN=<my-ingress-domain>

$ cat <<EOF | kubectl apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistry
metadata:
  name: devfile-registry
spec:
  devfileIndex:
    image: quay.io/devfile/devfile-index:next
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
    registryViewerWriteKey: <your-segment-write-key>
```

## Using Specific Container Images

By default, the operator deploys a Devfile Registry using a Pod containing three containers:

- `devfile-registry` container runs the API serving the Devfile stacks. The container image used by default
is `quay.io/devfile/devfile-index:next` and you can change it with the field `spec.devfileIndex.image`.
- `oci-registry` container runs a standard OCI registry, serving the stacks in the OCI format.
The container image used by default is `quay.io/devfile/oci-registry:next` and you can change it with the field `spec.ociRegistry.image`.
- `registry-viewer` container runs a web frontend, to help the user brwose the Devfile stacks.
The container image used by default is `quay.io/devfile/registry-viewer:next` and you can change it with the field `spec.registryViewer.image`.

```bash
$ cat <<EOF | oc apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistry
metadata:
  name: my-devfile-registry
spec:
  devfileIndex:
    image: quay.io/my-devfile/devfile-index:next
  ociRegistry:
    image: quay.io/my-devfile/oci-registry:next
  registryViewer:
    image: quay.io/my-devfile/registry-viewer:next
EOF
```

### Defining the ImagePullPolicy for pulling containers

By default, the containers will be pulled depending on the policy set on the Kubernetes or OpenShift cluster the registry is deployed on.

You can set a specific pull policy for each container, by setting the fields `spec.devfileIndex.imagePullPolicy`, `spec.ociRegistry.imagePullPolicy` and `spec.registryViewer.imagePullPolicy` to a value `IfNotPresent`, `Always` or `Never`. See the [Kubernetes documentation](https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy) for more information.


```bash
$ cat <<EOF | oc apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistry
metadata:
  name: my-devfile-registry
spec:
  devfileIndex:
    image: quay.io/my-devfile/devfile-index:next
    imagePullPolicy: Always
  ociRegistry:
    image: quay.io/my-devfile/oci-registry:next
    imagePullPolicy: Always
  registryViewer:
    image: quay.io/my-devfile/registry-viewer:next
    imagePullPolicy: Always
EOF
```

## Disabling the web frontend 

You can ask the operator to deploy the Devfile Registry without the `registry-viewer` container, by setting the field `spec.headless` to `true`.

```bash
$ cat <<EOF | oc apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistry
metadata:
  name: headless-devfile-registry
spec:
  devfileIndex:
    image: quay.io/devfile/devfile-index:next
  telemetry:
    registryName: test
  headless: true
EOF
```

## Configuring TLS for Ingress/Route resource

The operator creates a Route resource (on OpenShift) or an Ingress resources (on Kubernetes)
to give access to the Web frontend.

By default, the Ingress or Route is secured by TLS.

You can ask the operator to disable the use of TLS by setting the field `spec.tls.enabled` to false.

```bash
$ cat <<EOF | kubectl apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistry
metadata:
  name: devfile-registry
spec:
  devfileIndex:
    image: quay.io/devfile/devfile-index:next
  telemetry:
    registryName: test
  tls:
    enabled: false
EOF
```

You can ask the operator to configure the TLS with a specific certificate, by specifying a secret
containing the certificate and the associated private key using the field `spec.tls.secretName`:

```bash
$ kubectl create secret tls my-tls-secret --key=certs/ingress-tls.key --cert=certs/ingress-tls.crt

$ cat <<EOF | kubectl apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistry
metadata:
  name: devfile-registry
spec:
  devfileIndex:
    image: quay.io/devfile/devfile-index:next
  telemetry:
    registryName: test
  tls:
    enabled: true
    secretName: my-tls-secret
EOF
```

## Configuring the Ingress Domain

On Kubernetes, the operator needs to know the Domain associated with the cluster, to create an Ingress
with this specific domain. You need to indicate the Ingress domain with the field `spec.k8s.ingressDomain`.


```bash
$ export INGRESS_DOMAIN=<my-ingress-domain>

$ cat <<EOF | kubectl apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistry
metadata:
  name: devfile-registry
spec:
  devfileIndex:
    image: quay.io/devfile/devfile-index:next
  telemetry:
    registryName: test
  k8s:
    ingressDomain: $INGRESS_DOMAIN
EOF
```
