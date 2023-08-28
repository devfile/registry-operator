# Change Log
All notable changes to this project will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).
 
## [Unreleased]
 
### Added
- [devfile/api#1106](https://github.com/devfile/api/issues/1106)
  Registry operator should be in sync with operator SDK releases
  - Addition of `CHANGELOG.md`
 
### Changed

- [devfile/api#1106](https://github.com/devfile/api/issues/1106)
  Registry operator should be in sync with operator SDK releases
  - Syncs changes up to operator sdk v1.28
    - [v1.28.0 changes](https://sdk.operatorframework.io/docs/upgrading-sdk-version/v1.28.0/)
    - [v1.25.0 changes](https://sdk.operatorframework.io/docs/upgrading-sdk-version/v1.25.0/)
    - [v1.23.0 changes](https://sdk.operatorframework.io/docs/upgrading-sdk-version/v1.23.0/)
- [devfile/api#881](https://github.com/devfile/api/issues/881)
  Update Ginkgo version in registry operator
 
### Fixed
- [devfile/api#1106](https://github.com/devfile/api/issues/1106)
  Registry operator should be in sync with operator SDK releases
  - Registry operator service account changed from `default` to `service-account` to fix permissions error for creating leader election leases
