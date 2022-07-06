#!/bin/bash
#
# Copyright (c) 2020 Red Hat, Inc.
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#

#!/usr/bin/env bash
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
make test-integration
