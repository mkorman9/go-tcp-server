OUTPUT ?= go-tcp-server

.DEFAULT_GOAL := all

compile:
	CGO_ENABLED=0 go build -ldflags "-X main.AppVersion=$(shell make version)" -o $(OUTPUT)

package:
	VERSION=$(shell make version) build/package.sh

version:
	@echo "1.0.0"

all: compile
