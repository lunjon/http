all: format build test

format:
	@go fmt ./...

build: format
	@go build ./...

test: build
	@go test ./... | grep -v 'no test files'

cover:
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out

install:
	go install
