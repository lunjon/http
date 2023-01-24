default: format build test check

alias fmt := format
alias t := test

format:
	go fmt ./...

build: format
	go build ./...

check:
    staticcheck ./...

test: build
	go test ./... | grep -v 'no test files'

cover: build
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out

release version:
    git commit -am "Release {{version}}"
    git push
    ./build.sh
    gh release create {{version}} bin/*

# Get/install dependencies
deps:
    go get ./...
    go install honnef.co/go/tools/cmd/staticcheck@latest
