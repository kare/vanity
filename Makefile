LANG = en_US.UTF-8
SHELL = /bin/bash
.SHELLFLAGS = -eu -o pipefail -c # run '/bin/bash ... -c /bin/cmd'
.DEFAULT_GOAL = build

name = kkn.fi/vanity

GOIMPORTS = $(GOPATH)/bin/goimports
STATICCHECK = $(GOPATH)/bin/staticcheck
GOLANGCI-LINT = $(GOPATH)/bin/golangci-lint

.PHONY: build
build:
	go build $(name)

.PHONY: test
test:
	go test $(name)

$(GOIMPORTS):
	go install golang.org/x/tools/cmd/goimports@latest

$(GOLANGCI-LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

$(STATICCHECK):
	go install honnef.co/go/tools/cmd/staticcheck@latest

fmt:
	gofmt -w -s .

goimports: fmt $(GOIMPORTS)
	goimports -w .

staticcheck: $(STATICCHECK)
	staticcheck -go 1.16 ./...

golangci-lint: $(GOLANGCI-LINT)
	golangci-lint run ./...

.PHONY: cover
cover:
	go test -coverprofile=coverage.out $(name)
	go tool cover -html=coverage.out
	@rm -f coverage.out

.PHONY: heat
heat:
	go test -covermode=count -coverprofile=count.out $(name)
	go tool cover -html=count.out
	@rm -f count.out
