VERSION := $(shell git describe --tags)

test:
	go test ./...

build:
	go build -ldflags="-X 'version._version=$(VERSION)'" -o bin/gorepo_ ./cmd/gorepo/main.go
