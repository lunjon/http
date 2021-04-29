all: format build test

format:
	@go fmt ./...

build: format
	@go build ./...

test: build
	@go test ./... | grep -v 'no test files'

install:
	go install
