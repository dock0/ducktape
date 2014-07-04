DIR=$(shell pwd)
MUSL_VERSION="1.1.3"

.PHONY : default build_container manual container musl build push local

default: container

build_container:
	docker build -t ducktape meta

manual: build_container
	./meta/launch /bin/bash || true

container: build_container
	./meta/launch

musl:
	rm -rf musl.tar.gz musl-$(MUSL_VERSION)
	curl -o musl.tar.gz http://www.musl-libc.org/releases/musl-$(MUSL_VERSION).tar.gz
	tar xf musl.tar.gz
	cd musl-$(MUSL_VERSION) ; CFLAGS=-fno-toplevel-reorder ./configure
	$(MAKE) -C musl-$(MUSL_VERSION) install

build: musl
	mkdir -p gopath/{src,pkg,bin}
	GOPATH=$(DIR)/gopath go get -d
	CC=/usr/local/musl/bin/musl-gcc GOPATH=$(DIR)/gopath go build -o build/ducktape -ldflags '-extldflags "-static"'

push:
	ssh -oStrictHostKeyChecking=no git@github.com &>/dev/null || true
	git tag -f "$$(./build/ducktape -v)"
	git push origin ":$$(./build/ducktape -v)"
	git push --tags origin master
	targit -a .github -c -f dock0/nsinit $$(./build/ducktape -v) build/ducktape
	targit -a .github -c -f dock0/nsinit master build/ducktape

local: build push

