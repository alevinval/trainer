.PHONY: all build test vet cover update

all: update test vet build

build:
	go build -o bin/trainer ./cmd/trainer

test:
	go test ./...

vet:
	go vet ./...

cover:
	go test -coverprofile coverage.out ./...
	go tool cover -html coverage.out
	rm coverage.out

update:
	go get -v -u ./...
