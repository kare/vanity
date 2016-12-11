
NAME := kkn.fi/vanity/cmd/vanity

.PHONY: build install clean test vet lint errcheck cover heat check dist

build:
	go build $(NAME)/...

install:
	go install $(NAME)

clean:
	@rm -rf vanity

test:
	go test -v $(NAME)/...

vet:
	go vet $(NAME)/...

lint:
	golint $(NAME)/...

errcheck:
	errcheck $(NAME)/...

cover:
	go test -coverprofile=coverage.out $(NAME)
	go tool cover -html=coverage.out
	@rm -f coverage.out

heat:
	go test -covermode=count -coverprofile=count.out $(NAME)
	go tool cover -html=count.out
	@rm -f count.out

check: vet lint errcheck test heat

dist: clean
	GOOS=plan9 GOARCH=amd64 go build $(NAME)

