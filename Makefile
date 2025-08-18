.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golangci-lint run ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	go build ./cmd/archivist/
.PHONY:build

run: vet
	go run ./cmd/archivist/
.PHONY:run
