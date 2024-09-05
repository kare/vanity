# kkn.fi/vanity
[![CI](https://github.com/kare/vanity/workflows/CI/badge.svg)](https://github.com/kare/vanity/actions?query=workflow%3ACI)
[![Go Reference](https://pkg.go.dev/badge/kkn.fi/vanity.svg)](https://pkg.go.dev/kkn.fi/vanity)
[![GoReportCard](https://goreportcard.com/badge/kkn.fi/vanity)](https://goreportcard.com/report/kkn.fi/vanity)
    
## Concepts
- VCS is Version Control System (such as `git` or `hg`)
	- Repo root is the root path the source code repository (such as "https://github.com/kare")
- Domain is the internet address where the Go Vanity server is hosted (such as
  9fans.net or kkn.fi). Domain is deduced from HTTP request or can be set as a
  parameter.
- Path is the path component of the Go package (such as `/cmd/tcpproxy` in `kkn.fi/cmd/tcpproxy`)

## Specification
- [Go 1.4 Custom Import Path Checking](https://docs.google.com/document/d/1jVFkZTcYbNLaTxXD9OcGfn7vYv5hWtPx9--lTx1gPMs/edit)

## Features
- Zero dependencies.
- Redirects Go tool to VCS.
- Redirects browsers to [pkg.go.dev](https://pkg.go.dev) module server by default. Module Server URL is configurable.
- Automatic configuration of cmd packages:
	- All packages are redirected without sub-packages to VCS root.
	- Packages whose path is prefixed with `/cmd/` redirect automatically to VCS root by stripping the `/cmd` prefix from the package path.
	- Examples:
		- Redirect request `kkn.fi/cmd/tcpproxy` to `github.com/kare/tcpproxy`
		- Redirect request `kkn.fi/project/sub/package` to `github.com/kare/project`

## Vanity configurable options
Vanity package supports configurable [Option](https://pkg.go.dev/kkn.fi/vanity#Option)s via the [constructor](https://pkg.go.dev/kkn.fi/vanity#NewHandlerWithOptions). Use Option types to configure vanity handler features. Basic Options are documented below:
- Set [Version Control](https://pkg.go.dev/kkn.fi/vanity/#VCS) System type.
- Configurable [Version Control System HTTP URL](https://pkg.go.dev/kkn.fi/vanity/#VCSURL)
- [Module server URL](https://pkg.go.dev/kkn.fi/vanity/#ModuleServerURL) options are:
	- https://pkg.go.dev/
	- https://github.com/YOUR_USERNAME/
- Vanity server domain name defaults to request hostname, but it can also be configured.
- [Configurable](https://pkg.go.dev/kkn.fi/vanity/#Log) [Logger](https://pkg.go.dev/kkn.fi/vanity/#Logger) which is
  compatible with the standard [log.Logger](https://pkg.go.dev/log#Logger). Default output goes to standard error.
- Configurable [static content directory](https://pkg.go.dev/vanity/#StaticDir) for images, CSS, and etc.
- Configurable [IndexPageHandler](https://pkg.go.dev/vanity/#IndexPageHandler). Defaults to index.html file in the static content directory root.
- Configurable [robots.txt](https://pkg.go.dev/vanity/#RobotsTxt) [file](https://www.robotstxt.org).

## Installation
```
go get kkn.fi/vanity
```

## Development
New features or bug fixes must include comprehensive unit tests.

### Bugs
- Search [GitHub Issues](https://github.com/kare/kkn.fi-srv/issues) for existing bugs or open a new issue to report a new bug.
- To fix an existing bug from GitHub Issues open a new [GitHub Pull Request](https://github.com/kare/kkn.fi-srv/pulls).

### Building new features
1. Open a new [GitHub Issue](https://github.com/kare/kkn.fi-srv/issues) to discuss or propose a new feature.
1. Open a new Pull Request for a new feature that has beed already discussed.

### Execute compiler, tests and tools
Use [`Makefile`](https://github.com/kare/kkn.fi-srv/blob/main/Makefile) to execute Go compiler, tests and tools.

Run all tests
```bash
make test
```
Execute short (unit) running tests
```bash
make test-unit
```
Execute long (integration) running tests
```bash
make test-integration
```
Run `goimports`
```bash
make goimports
```
Run `staticcheck`
```bash
make staticcheck
```
Run gofmt with simplify
```bash
make fmt
```
