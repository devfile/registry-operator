#!/bin/bash
#
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

set -eu

usage ()
{   echo "Usage: ./make-release.sh <schema-version>"
    exit
}

if [[ $# -lt 1 ]]; then usage; fi
SCHEMA_VERSION=$1
FIRST_DIGIT="${SCHEMA_VERSION%%.*}"
RELEASE_BRANCH="release-v${FIRST_DIGIT}"

if ! command -v hub > /dev/null; then
  echo "[ERROR] The hub CLI needs to be installed. See https://github.com/github/hub/releases"
  exit
fi
if [[ -z "${GITHUB_TOKEN}" ]]; then
  echo "[ERROR] The GITHUB_TOKEN environment variable must be set."
  exit
fi
if ! [[ "$SCHEMA_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
	echo >&2 "$SCHEMA_VERSION isn't a valid semver tag for the schema. Aborting..."
	exit 1
fi


## Ensures local branch matches the remote
resetChanges() {
  echo "[INFO] Reset changes in $1 branch"
  git reset --hard
  git checkout $1
  git fetch origin --prune
  git pull origin $1
}

## Branch containing releases and tags in main upstream repo will be named 'release-vx' where 'vx' is the major release
## All minor and patch releases will be contained within a major release branch
## On the local side (the release engineers), the branches will be their full versioning name e.g. x.x.x
## Local branch will create a PR to its respective major release branch (if exists) or create a new one
checkoutToReleaseBranch() {
  echo "[INFO] Checking out to $SCHEMA_VERSION branch."
  if git ls-remote -q --heads | grep -q $SCHEMA_VERSION ; then
    echo "[INFO] $SCHEMA_VERSION exists."
    #resetChanges $SCHEMA_VERSION --- commented out so it doesnt delete my work when testing
  else
    echo "[INFO] $SCHEMA_VERSION does not exist. Will create a new one from main."
    #resetChanges main --- commented out so it doesnt delete my work when testing
    git push origin main:$SCHEMA_VERSION
  fi
  git checkout -B $SCHEMA_VERSION
}


updateVersionNumbers() {
    SHORT_UNAME=$(uname -s)

    ## Updating version.md based off of operating system
    if [ "$(uname)" == "Darwin" ]; then
      sed -i '' "s/^.*$/$SCHEMA_VERSION/" VERSION
    elif [ "${SHORT_UNAME:0:5}" == "Linux" ]; then
      sed -i "s/^.*$/$SCHEMA_VERSION/" VERSION
    fi

    ## Remaining version number updates to yaml files
    yq eval ".metadata.annotations.containerImage = \"quay.io/devfile/registry-operator:v$SCHEMA_VERSION\"" --inplace ./config/manifests/bases/registry-operator.clusterserviceversion.yaml
    yq eval ".metadata.name = \"registry-operator.v$SCHEMA_VERSION\"" --inplace ./config/manifests/bases/registry-operator.clusterserviceversion.yaml
    yq eval ".spec.version = \"$SCHEMA_VERSION\"" --inplace ./config/manifests/bases/registry-operator.clusterserviceversion.yaml
}

exportEnvironmentVariables() {
    CHANNEL=$(yq eval '.annotations."operators.operatorframework.io.bundle.channels.v1"' ./bundle/metadata/annotations.yaml)
    export IMG=quay.io/devfile/registry-operator:v$SCHEMA_VERSION
    export CHANNELS=$CHANNEL
    
}

commitChanges() {
  echo "[INFO] Pushing changes to $SCHEMA_VERSION branch"
  git add -A
  git commit -s -m "$1"
  git push origin $SCHEMA_VERSION
}

# Creates a new branch in the registry-operator repo for a new major release
# with the name release-vX
## This func will be used when we have a new major release and there is no branch in the upstream repo
createNewReleaseBranch(){
  git checkout -b "${RELEASE_BRANCH}" main
  git push origin "${RELEASE_BRANCH}"
  #hub sync
}

createPullRequest(){
  echo "[INFO] Creating a PR"
  hub pull-request --base jdubrick:${RELEASE_BRANCH} --head ${SCHEMA_VERSION} -m "$1"
}
 
main(){
  # checkoutToReleaseBranch
  # updateVersionNumbers
  # exportEnvironmentVariables
  # make bundle
  # commitChanges "chore(release): release version ${SCHEMA_VERSION}"
  #createPullRequest "v${SCHEMA_VERSION} Release"
  # Check if the branch exists in the remote repository
  if git ls-remote --exit-code --heads origin "$RELEASE_BRANCH" >/dev/null 2>&1; then
      echo "Branch $RELEASE_BRANCH exists in the remote repository."
  else
      echo "Branch $RELEASE_BRANCH does not exist in the remote repository."
      createNewReleaseBranch
  fi
}

main