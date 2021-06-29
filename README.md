# kkn.fi/vanity
[![Build Status](https://github.com/kare/vanity/workflows/CI/badge.svg)](https://github.com/kare/vanity/actions?query=workflow%3ACI)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/kkn.fi/vanity)
[![GoReportCard](https://goreportcard.com/badge/github.com/kare/vanity)](https://goreportcard.com/report/github.com/kare/vanity)
    

## Concepts
- VCS is Version Control System (such as 'git')
- Repo root is the root path the source code repository (such as 'https://github.com/kare')
- Domain is the internet address where the Go Vanity server is hosted (such as
  9fans.net or kkn.fi). Domain is deduced from HTTP request.
- Path is the path component of the Go package (such as /cmd/tcpproxy in
  kkn.fi/cmd/tcpproxy)

## Features
- Redirects Go tool to VCS
- Redirects browsers to pkg.go.dev module server by default. Module Server URL is configurable.
- Module server URL options are.
	- https://pkg.go.dev/
	- https://search.gocenter.io/
	- https://github.com/YOUR_USERNAME/
- Redirects HTTP to HTTPS
- Automatic configuration of packages:
	- All packages are redirected with full path to vcsroot.
	- Packages whose path is prefixed with "/cmd/" redirect automatically to
	  vcsroot by stripping the "/cmd" prefix from the package path.
	  Example: Redirect request "kkn.fi/cmd/tcpproxy" to "github.com/kare/tcpproxy"
- Configurable logger which is fully compatible with standard log package. Stdout is default.
- Supports index HTML file in the domain root and configurable static content directory (for images, CSS, and etc). 
- Supports [robots.txt file](https://www.robotstxt.org)

## Installation
```
go get kkn.fi/vanity
```

## Specification
- [Go 1.4 Custom Import Path Checking](https://docs.google.com/document/d/1jVFkZTcYbNLaTxXD9OcGfn7vYv5hWtPx9--lTx1gPMs/edit)
