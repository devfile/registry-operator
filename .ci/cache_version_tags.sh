#!/bin/sh

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

VERSION_PATTERN="v[0-9]+\.[0-9]+\.[0-9]+(-rc\.[0-9]+)?"
CONFIG_CSV='config/manifests/bases/registry-operator.clusterserviceversion.yaml'
CONFIG_MANAGER_KUSTOMIZE='config/manager/kustomization.yaml'
BUNDLE_CSV='bundle/manifests/registry-operator.clusterserviceversion.yaml'
YQ_CLI=${YQ_CLI:-yq}

# error on unset variables
set -u

${YQ_CLI} '.spec.version' ${CONFIG_CSV} > ${CACHED_CSV_VERSION} && \
${YQ_CLI} '(.metadata.annotations.containerImage | split(":") | .[1])' ${CONFIG_CSV} > ${CACHED_CSV_CONTAINER_IMAGE_TAG} && \
${YQ_CLI} "(.metadata.name | capture(\"(?P<tag>${VERSION_PATTERN})\") | .tag)" ${CONFIG_CSV} > ${CACHED_CSV_NAME_TAG} && \
${YQ_CLI} '.spec.version' ${BUNDLE_CSV} > ${CACHED_BUNDLE_VERSION} && \
${YQ_CLI} '(.metadata.annotations.containerImage | split(":") | .[1])' ${BUNDLE_CSV} > ${CACHED_BUNDLE_CONTAINER_IMAGE_TAG} && \
${YQ_CLI} "(.metadata.name | capture(\"(?P<tag>${VERSION_PATTERN})\") | .tag)" ${BUNDLE_CSV} > ${CACHED_BUNDLE_NAME_TAG} && \
${YQ_CLI} '.images[0].newTag' ${CONFIG_MANAGER_KUSTOMIZE} > ${CACHED_MANAGER_IMAGE_TAG}
