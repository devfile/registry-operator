# Registries Lists

The Cluster/Devfile Registries Lists allow admins to specify multiple devfile registries to expose devfiles from various sources for use within the cluster or namespace.  In order to be added to the list, the devfile registries must be reachable and support the Devfile v2.0 spec and above.

>**ClusterDevfileRegistriesList** is a custom resource where cluster admins can add a list of devfile registries to allow devfiles to be visible at the cluster level.  


> **DevfileRegistriesList** is a custom resource where cluster admins can add a list of devfile registries to allow devfiles to be visible at the namespace level.  Registries in this list will take precedence over the ones in the ClusterDevfileRegistriesList if there is a conflict.

## Deploying a Cluster or Devfile Registries List

Note the following limitations when deploying a new list type:
* Only one ClusterDevfileRegistriesList can be installed per cluster.  If there is an existing list, you will encounter a validation error.
* Only one DevfileRegistriesList can be installed per namespace.  If there is an existing list, you will encounter a validation error.


### Openshift or Kubernetes

The following examples are also available under the `samples/` folder in `yaml` format for your convenience

#### Installing a ClusterDevfileRegistriesList

```bash
$ cat <<EOF | kubectl apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: ClusterDevfileRegistriesList
metadata:
  name: cluster-list
spec:
  devfileRegistries:
    - name: devfile-staging
      url: 'https://registry.stage.devfile.io'
EOF
```

#### Installing a DevfileRegistriesList


```bash
$ cat <<EOF | kubectl apply -f -
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistriesList
metadata:
  name: namespace-list
spec:
  devfileRegistries:
    - name: devfile-staging
      url: 'https://registry.stage.devfile.io'
EOF
```



## Updating the Cluster or Devfile Registries List

To update the list of devfile registries in a CR, it's worthwhile to note that [strategic patch merge is not supported on custom resources](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/).  This limitation means we can't append to the list if we want to add a new entry for example.  As an alternative, we can use [jq](https://stedolan.github.io/jq/) to output the contents of the CR as json, modify the json, and then re-apply the config to update.  This would in effect, result in a replacement of the deployment config.


### Update the properties of an existing index

Using the example from above, we can change the name of the existing registry for the following CRs:

#####  ClusterDevfileRegistriesList  ##### 

```bash
$ kubectl get ClusterDevfileRegistriesList cluster-list -o json -n default | jq '.spec.devfileRegistries[0] = {"name": "my-staging-registry", "url":"https://registry.stage.devfile.io"}' | kubectl apply -f -
```

Similarly, for the DevfileRegistriesList run the same command but replace the CR type and name

##### DevfileRegistriesList ####
```bash
$ kubectl get DevfileRegistriesList namespace-list -o json -n default | jq '.spec.devfileRegistries[0] = {"name": "my-staging-registry", "url": "https://registry.stage.devfile.io"}' | kubectl apply -f -
```

### Append a new registry

```bash
$ kubectl get DevfileRegistriesList namespace-list -o json | jq '.spec.devfileRegistries += [{"name": "devfile-registry-community",  "url": "https://registry.devfile.io/"}]' | kubectl apply -f -
```

### Delete an existing registry

```bash
$ kubectl get DevfileRegistriesList namespace-list -o json | jq '.spec.devfileRegistries -= [{"name": "devfile-registry-community", "skipTLSVerify": false, "url":"https://registry.devfile.io/"}]' | kubectl apply -f -
```


## Consuming the CRs

### Tooling providers 

Tooling providers can query the cluster for these CR types by using the [controller-runtime client ](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client#Client) 

#### Example: 
Get a list of registries from ClusterDevfileRegistriesList

```go
...

//create a lookup key with the name and namespace of the CR
lookupKey := types.NamespacedName{Name: "cluster-name", Namespace: "default"}
cdrl := &v1alpha1.ClusterDevfileRegistriesList{}  
err := k8sClient.Get(ctx, lookupKey, cdrl)
if err ==  nil {
    registriesList := drl.Spec.DevfileRegistries
    for i:= range registriesList {
    	// If considering live URLs, check for availability. 
    	// Can use https://pkg.go.dev/github.com/devfile/registry-operator/api/v1alpha1#IsRegistryValid to verify 
    	regErr := isRegistryValid(registriesList[i].SkipTLSVerify, registries[i].URL)
        if regErr == nil {
        	// add to tooling catalog
        	...
        }
    }
}

```

### Cluster admins

Cluster admins can query their local cluster for the existence of a CR list by running the following commands:

#### Example1: 
Check for existence of ClusterDevfileRegistriesList

```bash
$ kubectl get ClusterDevfileRegistriesList
```

Expected response

``` 
NAME           STATUS
cluster-list   All devfile registries are active and reachable
```


#### Example2:
Check the contents of the ClusterDevfileRegistriesList

```bash
$ kubectl describe ClusterDevfileRegistriesList cluster-list
```

Expected response

``` 
Name:         cluster-list
Namespace:    
Labels:       <none>
Annotations:  <none>
API Version:  registry.devfile.io/v1alpha1
Kind:         ClusterDevfileRegistriesList
Metadata:
  Creation Timestamp:  2022-06-23T18:01:14Z
  Generation:          1
  Managed Fields:
    API Version:  registry.devfile.io/v1alpha1
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .:
          f:kubectl.kubernetes.io/last-applied-configuration:
      f:spec:
        .:
        f:devfileRegistries:
    Manager:      kubectl-client-side-apply
    Operation:    Update
    Time:         2022-06-23T18:01:14Z
    API Version:  registry.devfile.io/v1alpha1
    Fields Type:  FieldsV1
    fieldsV1:
      f:status:
        .:
        f:status:
    Manager:         manager
    Operation:       Update
    Subresource:     status
    Time:            2022-06-23T18:01:15Z
  Resource Version:  45143
  UID:               1c24fbe2-c154-41bd-bf80-f38a30980314
Spec:
  Devfile Registries:
    Name:             devfile-staging
    Skip TLS Verify:  false
    URL:              https://registry.stage.devfile.io
Status:
  Status:  All devfile registries are active and reachable
Events:    <none>

```