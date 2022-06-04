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
# install psql for wait-for-postgres.sh
	    apt-get update
		apt-get -y install postgresql-client
		chmod +x wait-for-postgres.sh
run:
		go run ./cmd/app

test:
		go test -v ./...

.DEFAULT_GOAL := build