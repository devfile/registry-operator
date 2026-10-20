# Supported Container Orchestration Systems

The devfile registry operator supports container orchestration 
systems based on which Operator SDK version the registry operator
is using as well as the API overlap between OpenShift and 
Kubernetes. 

Operator currently targets [Kubernetes 1.31 API](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.31.md) and is tested on [OpenShift 4.18](https://docs.redhat.com/en/documentation/openshift_container_platform/4.18/html/release_notes/ocp-4-18-release-notes).

## Operator SDK

The Registry Operator uses [Operator SDK v1.39.1](https://sdk.operatorframework.io/docs/upgrading-sdk-version/v1.39.1/) and Kubernetes 1.31 APIs, completing the planned SDK alignment for OpenShift 4.18 ([devfile/api#1626](https://github.com/devfile/api/issues/1626)).

To update the Operator SDK, refer to [Upgrade SDK Version](https://sdk.operatorframework.io/docs/upgrading-sdk-version/) 
and change the following
1. Update the following Kubernetes API dependencies under 
`go.mod` to the version supported by the chosen Operator SDK 
version 
    - `k8s.io/api`
    - `k8s.io/apimachinery`
    - `k8s.io/client-go`
2. Update `ENVTEST_K8S_VERSION` under `Makefile` to the 
Kubernetes version supported by the chosen Operator SDK version
3. Make any changes necessary to the source based on the updates 
to Kubernetes API dependencies
4. Update the Operator SDK CLI to the version chosen
5. Run `make bundle` to update fields to include the changes
for updating Operator SDK and Kubernetes APIs
