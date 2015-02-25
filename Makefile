DIR=$(shell pwd)

.PHONY : default build_container manual container build push local

default: container

build_container:
	docker build -t ducktape meta

manual: build_container
	./meta/launch /bin/bash || true

container: build_container
	./meta/launch

build:
	mkdir -p gopath/{src,pkg,bin}
	GOPATH=$(DIR)/gopath go get -d
	GOPATH=$(DIR)/gopath CC=musl-gcc go build -o build/ducktape -ldflags "-linkmode external -extldflags -static"
	strip build/ducktape

push:
	ssh -oStrictHostKeyChecking=no git@github.com &>/dev/null || true
	git tag -f "$$(./build/ducktape -v)"
	git push --tags origin master
	targit -a .github -c -f dock0/ducktape $$(./build/ducktape -v) build/ducktape

local: build push

