
name := kkn.fi/vanity
golint := $(GOPATH)/bin/golint
goimports := $(GOPATH)/bin/goimports

.PHONY: build
build:
	go build $(name)

.PHONY: test
test:
	go test $(name)

.PHONY: lint
lint: $(golint)
	golangci-lint run --config .golangci.yml ./...

$(golint):
	GO111MODULE=off go get -u golang.org/x/lint/golint

.PHONY: fmt
fmt: $(goimports)
	goimports -w *.go

$(goimports):
	GO111MODULE=off go get -u golang.org/x/lint/goimports

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
