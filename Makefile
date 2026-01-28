BINARY=cmd/server/main.go

.PHONY: build run tidy test

build:
	go build -o bin/agodrift ./cmd/server

run: build
	./bin/agodrift

tidy:
	go mod tidy

test:
	go test ./... 
