# Contributing

Thank you for your interest in contributing to the Devfile Registry Operator! We welcome your additions to this project.

## Code of Conduct

Before contributing to this repository for the first time, please review our project's [Code of Conduct](https://github.com/devfile/api/blob/main/CODE_OF_CONDUCT.md)

## How to Contribute

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

3. Open the folder in the IDE of your choice (VS Code with Go extension, or GoLand is recommended)

#### Build and Run the Operator
1. Log in to an OpenShfit or Kubernetes cluster

2. Run `export IMG=<operator-image>` where `<operator-image>` is the image repository to where you would like to push the image (e.g. `quay.io/user/registry-operator:latest`).

3. Run `make docker-build` to build the devfile registry operator.

4. Run `make docker-push` to push the devfile registry operator image.

5. Run `make install` to install the CRDs

6. Run `make deploy` to deploy the operator.

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


### Submitting Pull Request

**Note:** All commits must be signed off with the footer:
```
Signed-off-by: First Lastname <email@email.com>
```

You can easily add this footer to your commits by adding `-s` when running `git commit`.When you think the code is ready for review, create a pull request and link the issue associated with it.

Owners of the repository will watch out for and review new PRs. 

By default for each change in the PR, GitHub Actions and OpenShift CI will run checks against your changes (linting, unit testing, and integration tests)

If comments have been given in a review, they have to be addressed before merging.

After addressing review comments, don’t forget to add a comment in the PR afterward, so everyone gets notified by Github and knows to re-review.


# Contact us

If you have questions, please visit us on `#devfile` on the [Kubernetes Slack](https://slack.k8s.io).
