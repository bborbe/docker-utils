# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## Unreleased

- Migrate to tools.env + Makefile @version pattern; remove tools.go (note: indirect tool deps persist via unmigrated bborbe/argument v1)

## v1.7.8

- Update Go version from 1.24.2 to 1.26.0
- Update direct dependencies (glog, addlicense, goimports-reviser, counterfeiter, gomega)
- Update indirect dependencies (golang.org/x/*, google.golang.org/protobuf)
- Add CI workflow for automated testing on pull requests and pushes
- Add dependency exclusions for problematic k8s.io and go-logr versions

## v1.7.7

- go mod update

## v1.7.6

- go mod update

## v1.7.5

- go mod update

## v1.7.4

- go mod update

## v1.7.3

- go mod update

## v1.7.2

- go mod update

## v1.7.1

- go mod upgrade

## v1.7.0

- add dockerhub-cleaner

## v1.6.0

- make username and password optional
- remove subpackages
- use go modules

## v1.5.0

- Rename docker-remote-size to docker-remote-size-repositories
- Add docker-remote-size-tags command

## v1.4.0

- Add docker-remote-size command

## v1.3.1

- Improve logging and error handling

## v1.3.0

- Add docker-remote-sha-for-tag command

## v1.2.1

- Add Jenkinsfile
- Use deps instead glide

## v1.2.0

- Add flag to read Docker credentials from ~/.docker/config.json

## v1.1.0

- Add delete tag command

## v1.0.2

- Allow read password from file
