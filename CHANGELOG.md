# Change Log
All notable changes to this project will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).
 
## [Unreleased]

### Added

### Changed

### Fixed

## [0.1.0]

### Added
- [devfile/api#1106](https://github.com/devfile/api/issues/1106)
  Registry operator should be in sync with operator SDK releases
  - Addition of `CHANGELOG.md`
- [devfile/api#1233](https://github.com/devfile/api/issues/1233)
  - Included `install-cert` as an required step before `install` and `deploy`
 
### Changed
- [devfile/api#1233](https://github.com/devfile/api/issues/1233)
  Update CONTRIBUTING for registry operator
  - Replaced prerequisites with a reference to the [requirements section](README.md#requirements) under the README
- [devfile/api#1211](https://github.com/devfile/api/issues/1211)
  Devfile registry operator pre-release to OperatorHub
  - Bump version to `v0.1.0`
  - Changes to CRs, CRDs, samples, and CSV to provide up to date information on the operator
    - Metadata
      - Capabilities
      - Categories
      - Container Image
      - Source Repository Link
      - Maintainer Organization
      - Samples List
      - Short Description
    - Maintainers List
    - Icon
    - Display Name
    - Long Description
    - Version
  - Update [README](README.md) with [LICENSE](LICENSE) and [CONTRIBUTING.md](CONTRIBUTING.md) links
- [devfile/api#1106](https://github.com/devfile/api/issues/1106)
  Registry operator should be in sync with operator SDK releases
  - Syncs changes up to operator sdk v1.28
    - [v1.28.0 changes](https://sdk.operatorframework.io/docs/upgrading-sdk-version/v1.28.0/)
    - [v1.25.0 changes](https://sdk.operatorframework.io/docs/upgrading-sdk-version/v1.25.0/)
    - [v1.23.0 changes](https://sdk.operatorframework.io/docs/upgrading-sdk-version/v1.23.0/)
- [devfile/api#1015](https://github.com/devfile/api/issues/1015)
  Complete documentation coverage of registry operator
  - Provided complete documentation on makefile usage
  - Provided complete documentation on deployment spec fields
- [devfile/api#881](https://github.com/devfile/api/issues/881)
  Update Ginkgo version in registry operator
 
### Fixed
- [devfile/api#1106](https://github.com/devfile/api/issues/1106)
  Registry operator should be in sync with operator SDK releases
  - Registry operator service account changed from `default` to `service-account` to fix permissions error for creating leader election leases
- [devfile/api#1211](https://github.com/devfile/api/issues/1211)
  Devfile registry operator pre-release to OperatorHub
  - Fixes for links under [PR template](.github/PULL_REQUEST_TEMPLATE.md)