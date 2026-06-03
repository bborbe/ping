# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v0.1.8

- Bump Go to 1.26.4
- Update ginkgo/v2 v2.28.3 → v2.29.0
- Update gomega v1.40.0 → v1.41.0
- Exclude cloud.google.com/go v0.26.0
- Add .maintainer.yaml with autoRelease/autoApprove

## v0.1.7

- bump go 1.26.2 → 1.26.3
- bump github.com/bborbe/errors v1.5.12 → v1.5.13

## v0.1.6

- chore: Migrate to tools.env + Makefile @version pattern; remove tools.go and obsolete replace block. go.mod reduced from 504 to 26 lines

## v0.1.5

- chore: Add golangci-lint v2 with standard linter config, modernize Makefile with lint/gosec/osv-scanner/trivy targets, update tools.go with current tool dependencies

## v0.1.4

- chore: Verify project health — all tests pass, linting clean, precommit succeeds

## v0.1.3

- go mod update

## v0.1.2

- go mod update

## v0.1.1

- go mod update

## v0.1.0

- Initial release
