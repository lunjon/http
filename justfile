default: format build test check

alias fmt := format
alias t := test

format:
	go fmt ./...

build:
	go build ./...
	nix build

check:
    staticcheck ./...

test:
	go test ./... | grep -v 'no test files'

cover:
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out

# Get/install dependencies
deps:
    go get ./...
    go install honnef.co/go/tools/cmd/staticcheck@latest

release version: test
	./build.sh
	gh release create {{ version }} \
		--title {{ version }} \
		--generate-notes \
		bin/*
