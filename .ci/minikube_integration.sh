#!/bin/bash

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

# NOTE: This script assumes that minikube is installed and running, and using the docker driver on Linux
# Due to networking issues with the docker driver and ingress on macOS/Windows, this script must be run on Linux

# Share docker env with Minikube
eval $(minikube docker-env)

# error on unset variables
set -u
# print each command before executing it
set -x

# Build the registry operator image
export IMG=${REGISTRY_OPERATOR}
make docker-build
if [ $? -ne 0 ]; then
    echo "Error building registry operator image"
    exit 1;
fi

# Install cert-manager
make install-cert

# Wait for the cert-manager to become ready
kubectl wait deploy/cert-manager --namespace cert-manager --for=condition=Available --timeout=600s
kubectl wait deploy/cert-manager-cainjector --namespace cert-manager --for=condition=Available --timeout=600s
kubectl wait deploy/cert-manager-webhook --namespace cert-manager --for=condition=Available --timeout=600s
if [ $? -ne 0 ]; then
    echo "cert-manager-controller container logs:"
    kubectl logs -l app=cert-manager --namespace cert-manager --container cert-manager-controller
    echo "cert-manager-cainjector container logs:"
    kubectl logs -l app=cainjector --namespace cert-manager --container cert-manager-cainjector
    echo "cert-manager-webhook container logs:"
    kubectl logs -l app=webhook --namespace cert-manager --container cert-manager-webhook

    # Return the description of every pod
    kubectl describe pods --namespace cert-manager
    exit 1
fi

# Install CRDs & deploy registry operator
make install && make deploy

# Wait for the registry operator to become ready
kubectl wait deploy/registry-operator-controller-manager --namespace registry-operator-system --for=condition=Available --timeout=600s
if [ $? -ne 0 ]; then
    echo "manager container logs:"
    kubectl logs -l app=devfileregistry-operator --namespace registry-operator-system --container manager
    echo "kube-rbac-proxy container logs:"
    kubectl logs -l app=devfileregistry-operator --namespace registry-operator-system --container kube-rbac-proxy

    # Return the description of every pod
    kubectl describe pods --namespace registry-operator-system
    exit 1
fi

# run integration test suite
make test-integration
