.PHONY: init build run test

all:
		make init
		make build
		make test
		make run
init:
		go mod tidy
build:
		make init
#		go build -v ./cmd/app
		go build -o build/server -v ./cmd/app
run:
		go run ./cmd/app

test:
		go test -v ./...

.DEFAULT_GOAL := build