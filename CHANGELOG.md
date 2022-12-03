# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.0] - 2022-12-03

### Added
- Builds for macOS.
- Builds for ARM64.
- Publish container images to GitHub container registry.

### Changed

- Code cleanup.
- Upgrade dependencies.

## [1.3.0] - 2020-02-22

### Added
- --version flag

### Changed

- Code cleanup.
- Upgrade dependencies.

### Removed
- version command.

## [1.2.0] - 2019-09-14

### Added

- delete command to remove a single certificate from vault.
- prune command to automatically remove expired certificates from vault. 

### Changed

- Update dependencies.
- Compile using Go 1.13.

## [1.1.0] - 2018-09-23

### Added

- list command to list all certificates stored in vault.

### Changed

- Migrate dependency management from `dep` to go modules.
- Update dependencies.

## 1.0.0 - 2018-05-23

Initial release

[Unreleased]: https://github.com/emgag/keyvault-certdeploy/compare/v1.4.0...HEAD
[1.1.0]: https://github.com/emgag/keyvault-certdeploy/compare/v1.0.0...v1.1.0
[1.2.0]: https://github.com/emgag/keyvault-certdeploy/compare/v1.1.0...v1.2.0
[1.3.0]: https://github.com/emgag/keyvault-certdeploy/compare/v1.2.0...v1.3.0
[1.4.0]: https://github.com/emgag/keyvault-certdeploy/compare/v1.3.0...v1.4.0
