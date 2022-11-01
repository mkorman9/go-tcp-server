OUTPUT ?= go-tcp-server

.DEFAULT_GOAL := all

build:
	CGO_ENABLED=0 go build -o $(OUTPUT)

test:
	go test -v ./...

all: build test
