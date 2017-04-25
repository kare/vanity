
# Vanity [![Build Status](https://travis-ci.org/kare/vanity.svg?branch=master)](https://travis-ci.org/kare/vanity) [![GoDoc](https://godoc.org/kkn.fi/cmd/vanity?status.svg)](https://godoc.org/kkn.fi/cmd/vanity)

## Features
- Redirects browsers to godoc.org
- Redirects Go tool to VCS

## Installation
```
go get kkn.fi/cmd/vanity
```

## Configuration
Configuration example for Gorilla project:

```
/context  git https://github.com/gorilla/context
/mux  git https://github.com/gorilla/mux
```

## Running
Script to run Gorilla toolkit's vanity domain server:
```
vanity -d gorillatoolkit.org -c vanity.conf
```

## Specification
- [Go 1.4 Custom Import Path Checking](https://docs.google.com/document/d/1jVFkZTcYbNLaTxXD9OcGfn7vYv5hWtPx9--lTx1gPMs/edit)
