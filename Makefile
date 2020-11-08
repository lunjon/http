all: format build test

format:
	go fmt github.com/lunjon/httpreq/...

test:
	go test ./...

build:
	go build ./...

install:
	go install
