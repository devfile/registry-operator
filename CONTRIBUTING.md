# Contributing

Thank you for your interest in contributing to the Devfile Registry Operator! We welcome your additions to this project.

## Code of Conduct

Before contributing to this repository for the first time, please review our project's [Code of Conduct](https://github.com/devfile/api/blob/main/CODE_OF_CONDUCT.md)

## How to Contribute:

### Issues

If you spot a problem with the devfile registry, [search if an issue already exists](https://github.com/devfile/api/issues). If a related issue doesn't exist, you can open a new issue using a relevant [issue form](https://github.com/devfile/api/issues/new/choose).

You can tag Devfile Registry related issues with the `/area registry` text in your issue.

### Development

#### First Time Setup
1. Install prerequisites:
   - Go 1.13 or higher
   - Docker or Podman
   - Operator-SDK 1.11.0 or higher (including `controller-gen` 0.6.0 or higher)

2. Fork and clone this repository.

3. `cd registry-operator`

4. Open the folder in the IDE of your choice (VS Code with Go extension, or GoLand is recommended)

#### Build and Run the Operator
1. Log in to an OpenShfit or Kubernetes cluster

2. Run `export IMG=<operator-image>` where `<operator-image>` is the image repository where you would like to push the image to (e.g. `quay.io/user/registry-operator:latest`).

3. Run `make docker-build` to build the devfile registry operator.

4. Run `make docker-push` to push the devfile registry operator image.

4. Run `make install` to install the CRDs

5. Run `make deploy` to deploy the operator.

### Testing your Changes

All changes delivered to the Devfile Registry operator are expected to be sufficiently tested. This may include validating that existing tests pass, updating tests, or adding new tests.

#### Unit Tests

The unit tests for this repository are located under the `pkg/` folder and are denoted with the `_test.go` suffix. The tests can be run by running `make test`.

#### Integration Tests

The integration tests for this repository are located under the `tests/integration` folder and contain tests that validate the Operator's functionality when running on an OpenShift cluster.

To run these tests, run the following commands :

```bash
export IMG=<your-built-operator-image>
make test_integration
```

# Contact us

If you have questions, please visit us on `#devfile` on the [Kubernetes Slack](https://slack.k8s.io).
