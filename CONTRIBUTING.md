# Contributing

Thank you for your interest in contributing to the Devfile Registry Operator! We welcome your additions to this project.

## Code of Conduct

Before contributing to this repository for the first time, please review our project's [Code of Conduct](https://github.com/devfile/api/blob/main/CODE_OF_CONDUCT.md)

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](DCO) file for details.

In order to show your agreement with the DCO you should include at the end of the commit message,
the following line:
```console
Signed-off-by: Firstname Lastname <email@email.com>
```

Once you set your user.name and user.email in your git config, you can sign your commit automatically with `git commit -s`.

## How to Contribute

### Issues

If you spot a problem with the **devfile registry operator**, [search if an issue already exists](https://github.com/devfile/api/issues?q=is%3Aissue+is%3Aopen+label%3Aarea%2Fregistry). 

If a related issue doesn't exist, you can open a new issue using a relevant [issue form](https://github.com/devfile/api/issues/new/choose). You can tag Devfile Registry Operator related issues with the `/area registry` text in your issue.

### Development

#### First Time Setup
1. Install prerequisites: see [Requirements section in README](README.md#requirements).

2. Fork and clone this repository.

3. Open the folder in the IDE of your choice (VS Code with Go extension, or GoLand is recommended)

#### Build and Run the Operator
The Makefile currently supports both Docker and Podman. To run the proper command replace `<engine>` with either `podman` or `docker` depending on your container engine.
1. Log in to an OpenShift or Kubernetes cluster

2. Run `export IMG=<operator-image>` where `<operator-image>` is the image repository to where you would like to push the image (e.g. `quay.io/user/registry-operator:latest`).

3. Run `make <engine>-build` to build the devfile registry operator.

4. Run `make <engine>-push` to push the devfile registry operator image.

5. (Optional, **docker only**) Run `make docker-buildx` to build and push the devfile registry operator multi-architecture image.

6. Run `make install-cert` to install the cert-manager. (Allow time for these services to spin up before moving on to step 7 & 8).

7. Run `make install` to install the CRDs.

8. Run `make deploy` to deploy the operator.

##### Enabling HTTP/2 on the Webhook Server

By default, http/2 on the webhook server is disabled due to [CVE-2023-44487](https://github.com/advisories/GHSA-qppj-fm5r-hxr3).

If you want to enable http/2 for the webhook server, build with `ENABLE_WEBHOOK_HTTP2=true make <engine>-build` or with 
`ENABLE_WEBHOOK_HTTP2=true make run` if running locally.

##### Using other platforms

If you need to target another platform for container builds, such as Apple silicon, you can use `TARGET_ARCH=<architecture> make <engine>-build`.

For example, to target container build to `arm64` run the following:

```sh
TARGET_ARCH=arm64 make <engine>-build
```

**Note:** Container builds only use `linux` as the operating system as local cluster runtime environments, such as `minikube` environments, run under Linux virtual machines for other operating systems. For example, _Apple silicon_ would just use the `arm64` container build.

For local builds, you can also set the target operating system:

**Apple silicon**

```sh
export TARGET_OS=darwin
export TARGET_ARCH=arm64
make manager
```

**Linux ARM**

```sh
export TARGET_ARCH=arm64
make manager
```

**Windows**

```sh
export TARGET_OS=windows
export TARGET_ARCH=amd64
make manager
```

By default, `amd64` is used for the target architecture and `linux` is used for the target operating system.

### Testing your Changes

All changes delivered to the Devfile Registry Operator are expected to be sufficiently tested. This may include validating that existing tests pass, updating tests, or adding new tests.

#### Unit Tests

The unit tests for this repository are located under the `pkg/` folder and are denoted with the `_test.go` suffix. The tests can be run by running `make test`.

#### Integration Tests

The integration tests for this repository are located under the `tests/integration` folder and contain tests that validate the Operator's functionality when running on an OpenShift cluster.

To run these tests, run the following commands:

```bash
export IMG=<your-built-operator-image>
make test-integration
```

### Submitting Pull Request

**Note:** All commits must be signed off with the footer:
```
Signed-off-by: First Lastname <email@email.com>
```

You can easily add this footer to your commits by adding `-s` when running `git commit`. When you think the code is ready for review, create a pull request and link the issue associated with it.

Owners of the repository will watch out for new PRs and provide reviews to them.

For each change in the PR, GitHub Actions and OpenShift CI will by default run checks against your changes (linting, unit testing, and integration tests).

If comments have been given in a review, they have to be addressed before merging.

After addressing review comments, don't forget to add a comment in the PR with the reviewer mentioned afterward, so they get notified by Github to provide a re-review.

# Contact us

If you have any questions, please visit the `#devfile` channel under the [Kubernetes Slack](https://slack.k8s.io) workspace.
