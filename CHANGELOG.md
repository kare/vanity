
# Change Log

## Release v1.0.0-alpha.6 - 2020-TODO
- Add support for GET /robots.txt

## Release v1.0.0-alpha.5 - 2020-10-15
- Support GitHub as a vanity.ModuleServerURL() option. For example: Use https://github.com/kare/ as a module server url.
- Require Go 1.14
- Add IndexPageHandler() for serving index.html page from static file directory root.
- DefaultIndexPageHandler() uses static content directory path + "/index.html" as a default.

## Release v1.0.0-alpha.4 - 2020-09-10
- Support configurable static file directory with optional index.html file.

## Release v1.0.0-alpha.3 - 2020-02-05
- Support index HTML page on the root of the server

## Release v1.0.0-alpha.2 - 2020-01-23
- Add configurable option ModuleServerURL() for setting the used Go module server address.

## Release v1.0.0-alpha.1 - 2020-01-21
- Replace vanity.Redirect() with vanity.Handler()
	- Add SetLogger(), VCS() and VCSURL() functional options for configuring Handler
- Update docs
- Rewrite example
- Improve tests

## Release v0.2.2 - 2020-01-20
- Update docs for pkg.go.dev
- Improve tests
- Drop text/template. Just print out <meta> tag.
- Replace deprecated gometalinter with golint

## Release v0.2.1 - 2019-12-28
### Added
- Require Go 1.13
- Use pkg.go.dev instead of godoc.org

## Release v0.2.0 - 2018-11-06
### Added
- Rework internal structure
- Add support for configurable package logger
- Support Go 1.11 Modules

## Release v0.1.0 - 2016-12-12
### Added
- Add vanity.Path type
- Add vanity.NewServer and vanity.NewPackage functions

## Release v0.0.1 - 2016-10-20
### Added
- Supports Go tool and browsers for GoDoc
