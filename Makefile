OUTPUT ?= go-tcp-server

.DEFAULT_GOAL := all

build:
	CGO_ENABLED=0 go build -o $(OUTPUT)

package:
	.package/build.sh

all: build package
