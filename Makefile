VERSION := $(shell git describe --tags)

test:
	go test ./...

build:
	go build -ldflags="-X 'gorepo-cli/internal/commands/version.Version=$(VERSION)'" -o bin/gorepo_ ./cmd/gorepo/main.go
