# kkn.fi/vanity
[![CI](https://github.com/kare/vanity/workflows/CI/badge.svg)](https://github.com/kare/vanity/actions?query=workflow%3ACI)
[![Go Reference](https://pkg.go.dev/badge/kkn.fi/vanity.svg)](https://pkg.go.dev/kkn.fi/vanity)
[![GoReportCard](https://goreportcard.com/badge/kkn.fi/vanity)](https://goreportcard.com/report/kkn.fi/vanity)
    

## Concepts
- VCS is Version Control System (such as 'git')
- Repo root is the root path the source code repository (such as 'https://github.com/kare')
- Domain is the internet address where the Go Vanity server is hosted (such as
  9fans.net or kkn.fi). Domain is deduced from HTTP request.
- Path is the path component of the Go package (such as /cmd/tcpproxy in
  kkn.fi/cmd/tcpproxy)

## Specification
- [Go 1.4 Custom Import Path Checking](https://docs.google.com/document/d/1jVFkZTcYbNLaTxXD9OcGfn7vYv5hWtPx9--lTx1gPMs/edit)

## Features
- Redirects Go tool to VCS
- Redirects browsers to pkg.go.dev module server by default. Module Server URL is configurable.
- Module server URL options are:
	- https://pkg.go.dev/
	- https://github.com/YOUR_USERNAME/
- Hostname defaults to request host, but it can also be configured.
- Automatic configuration of cmd packages:
	- All packages are redirected without sub-packages to vcsroot.
	- Packages whose path is prefixed with "/cmd/" redirect automatically to
	  vcsroot by stripping the "/cmd" prefix from the package path.
	  Example: Redirect request "kkn.fi/cmd/tcpproxy" to "github.com/kare/tcpproxy"
      Example: Redirect request "kkn.fi/project/sub/package" to
	  "github.com/kare/project"
- Configurable logger which is compatible with standard `log` package. Default output goes to `stderr`.
- Supports index HTML file in the domain root and configurable static content directory (for images, CSS, and etc). 
- Supports [robots.txt file](https://www.robotstxt.org)

## Installation
```
go get kkn.fi/vanity
```

## Development
Run all tests
```bash
make test
```
Run short (unit) tests
```bash
make test-unit
```
Run long (integration) tests
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
