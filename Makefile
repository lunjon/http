all: format build test

format:
	go fmt ./...

build: format
	go build ./...

test: build
	go test ./...

install:
	go install
