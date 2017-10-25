.PHONY: test cover

test:
	go test

cover:
	go test -cover -coverprofile coverage.out
	go tool cover -html coverage.out
	rm coverage.out
