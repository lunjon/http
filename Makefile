
test:
	go test -v ./...

build:
	go build ./cmd/httpreq/

install:
	go install ./cmd/httpreq
