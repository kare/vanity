LANG := en_US.UTF-8
SHELL := /bin/bash
.SHELLFLAGS := --norc --noprofile -e -u -o pipefail -c
.DEFAULT_GOAL := build

name := kkn.fi/vanity

GOIMPORTS := $(GOPATH)/bin/goimports
STATICCHECK := $(GOPATH)/bin/staticcheck

.PHONY: build
build:
	go build $(name)

.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit:
	go test -v -short $(name)

.PHONY: test-integration
test-integration:
	go test -v $(name)

$(GOIMPORTS):
	go install golang.org/x/tools/cmd/goimports@latest

$(STATICCHECK):
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: fmt
fmt:
	gofmt -w -s .

.PHONY: goimports
goimports: fmt $(GOIMPORTS)
	$(GOIMPORTS) -w .

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./...

.PHONY: cover
cover:
	go test -coverprofile=coverage.out $(name)/...
	go tool cover -html=coverage.out
	@$(RM) -f coverage.out

.PHONY: heat
heat:
	go test -covermode=count -coverprofile=count.out $(name)/...
	go tool cover -html=count.out
	@$(RM) -f count.out

.PHONY: clean
clean:
	go clean
