all: build test

test:
	go test ./...

build:
	go build ./cmd/httpreq/

install:
	go install ./cmd/httpreq
