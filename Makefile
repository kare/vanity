
NAME := kkn.fi/vanity

.PHONY: build test vet lint errcheck cover heat check

build:
	go build $(NAME)/...

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

