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

# exit immediately when a command fails
set -e
# only exit with zero if all commands of the pipeline exit successfully
set -o pipefail
# error on unset variables
set -u
# print each command before executing it
set -x

# Make sure we're running the integration tests with the image built by OpenShift CI
export IMG=${REGISTRY_OPERATOR}

# For some reason go on PROW force usage vendor folder
# This workaround is here until we don't figure out cause
go mod tidy
go mod vendor

# Make sure kustomize and controller-gen are installed before running the tests
# ToDo: Remove later, should not be required.
make kustomize
make controller-gen
# need to have cert-manager installed to run tests
make install-cert
# wait one minute for cert manager to get set up
sleep 60
# need to deploy the registry operator to run tests
# ToDo: Remove later after the addition of readiness check, integration tests can deploy the operator however tests fail if 
# pod is not ready in time.
make install && make deploy
# wait 15 seconds for registry operator to get set up
sleep 15

# run integration test suite
make test-integration
