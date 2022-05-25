all: format build test
alias t := test

format:
	go fmt ./...

build: format
	go build ./...

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
