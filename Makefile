all: build test

test:
	go test ./...

build:
	go build ./...

install:
	go install
