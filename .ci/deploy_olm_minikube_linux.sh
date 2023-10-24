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

# Re-make bundle
export CHANNELS=${REGISTRY_OPERATOR_CHANNELS}
make bundle

# Build the registry operator bundle image
export BUNDLE_IMG=${REGISTRY_OPERATOR_BUNDLE}
make bundle-build
if [ $? -ne 0 ]; then
    echo "Error building registry operator bundle image"
    exit 1;
fi

# Install OLM
operator-sdk olm install

# OLM install registry operator
operator-sdk run bundle ${BUNDLE_IMG}

# Wait for the registry operator to become ready
kubectl wait deploy/registry-operator-controller-manager --for=condition=Available --timeout=600s
if [ $? -ne 0 ]; then
    echo "manager container logs:"
    kubectl logs -l app=devfileregistry-operator --container manager
    echo "kube-rbac-proxy container logs:"
    kubectl logs -l app=devfileregistry-operator --container kube-rbac-proxy

    # Return the description of every pod
    kubectl describe pods
    exit 1
fi

