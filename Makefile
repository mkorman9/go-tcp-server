OUTPUT ?= go-tcp-server

.DEFAULT_GOAL := all

build:
	CGO_ENABLED=0 go build -o $(OUTPUT)

package:
	VERSION=$(shell make version) .package/build.sh

version:
	@echo "1.0.0"

all: build package
