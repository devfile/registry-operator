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

CACHED_CSV_VERSION=${CACHED_CSV_VERSION:-'.cache/csv_version.txt'}
CACHED_CSV_CONTAINER_IMAGE_TAG=${CACHED_CSV_CONTAINER_IMAGE_TAG:-'.cache/csv_container_image_tag.txt'}
CACHED_CSV_NAME_TAG=${CACHED_CSV_NAME_TAG:-'.cache/csv_name_tag.txt'}
CACHED_BUNDLE_VERSION=${CACHED_BUNDLE_VERSION:-'.cache/bundle_version.txt'} 
CACHED_BUNDLE_CONTAINER_IMAGE_TAG=${CACHED_BUNDLE_CONTAINER_IMAGE_TAG:-'.cache/bundle_container_image_tag.txt'}
CACHED_BUNDLE_NAME_TAG=${CACHED_BUNDLE_NAME_TAG:-'.cache/bundle_name_tag.txt'}
CACHED_MANAGER_IMAGE_TAG=${CACHED_MANAGER_IMAGE_TAG:-'.cache/manager_image_tag.txt'}

ref_name=$1
failed="false"

if [ -z ${ref_name} ]
then
    echo "expects release tag (ref_name): bash check_version.sh <ref_name> [is_ci=false]"
    exit 1
fi

if [ "${CI}" != "true" ]
then
    if [ -z $(command -v yq) ] && [ -z $(command -v ${YQ_CLI}) ]
    then
        echo "This script requires the yq tool."
        exit 1
    fi

    # Make cache directory if not exists
    mkdir -p .cache/

    # Export variables
    export CACHED_CSV_VERSION
    export CACHED_CSV_CONTAINER_IMAGE_TAG
    export CACHED_CSV_NAME_TAG
    export CACHED_BUNDLE_VERSION
    export CACHED_BUNDLE_CONTAINER_IMAGE_TAG
    export CACHED_BUNDLE_NAME_TAG
    export CACHED_MANAGER_IMAGE_TAG

    # Read references to release version under project files and cache them for checks
    bash .ci/cache_version_tags.sh
    if [ $? -ne 0 ]
    then
        exit 1
    fi
fi

# error on unset variables
set -u

## Check if all references to the release version match ##

if [ "${ref_name}" != "v$(cat ./VERSION)" ]
then
    echo "Release tag does not match VERSION: release tag = ${ref_name}, VERSION = v$(cat ./VERSION)"
    failed="true"
fi

if [ "${ref_name}" != "v$(cat ${CACHED_CSV_VERSION})" ]
then
    echo "Release tag does not match csv version: release tag = ${ref_name}, csv version = v$(cat ${CACHED_CSV_VERSION})"
    failed="true"
fi

if [ "${ref_name}" != "$(cat ${CACHED_CSV_CONTAINER_IMAGE_TAG})" ]
then
    echo "Release tag does not match csv container image tag: release tag = ${ref_name}, csv container image tag = $(cat ${CACHED_CSV_CONTAINER_IMAGE_TAG})"
    failed="true"
fi

if [ "${ref_name}" != "$(cat ${CACHED_CSV_NAME_TAG})" ]
then
    echo "Release tag does not match csv name tag: release tag = ${ref_name}, csv name tag = $(cat ${CACHED_CSV_NAME_TAG})"
    failed="true"
fi

if [ "${ref_name}" != "v$(cat ${CACHED_BUNDLE_VERSION})" ]
then
    echo "Release tag does not match bundle version: release tag = ${ref_name}, bundle version = v$(cat ${CACHED_BUNDLE_VERSION})"
    failed="true"
fi

if [ "${ref_name}" != "$(cat ${CACHED_BUNDLE_CONTAINER_IMAGE_TAG})" ]
then
    echo "Release tag does not match bundle container image tag: release tag = ${ref_name}, bundle container image tag = $(cat ${CACHED_BUNDLE_CONTAINER_IMAGE_TAG})"
    failed="true"
fi

if [ "${ref_name}" != "$(cat ${CACHED_BUNDLE_NAME_TAG})" ]
then
    echo "Release tag does not match bundle name tag: release tag = ${ref_name}, bundle name tag = $(cat ${CACHED_BUNDLE_NAME_TAG})"
    failed="true"
fi

if [ "${ref_name}" != "$(cat ${CACHED_MANAGER_IMAGE_TAG})" ]
then
    echo "Release tag does not match manager image tag: release tag = ${ref_name}, manager image tag = $(cat ${CACHED_MANAGER_IMAGE_TAG})"
    failed="true"
fi

if [ ${failed} == "true" ]
then
    echo "One or more checks failed!"
    exit 1
else
    echo "All version tags match!"
fi
