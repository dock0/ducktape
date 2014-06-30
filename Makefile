DIR=$(shell pwd)
MUSL_VERSION="1.1.3"

.PHONY : container setup build install local

container:
	docker build --no-cache -t ducktape .
	docker run --rm -t -i -v $(DIR):/opt/ducktape ducktape

musl:
	rm -rf musl.tar.gz musl-$(MUSL_VERSION)
	curl -o musl.tar.gz http://www.musl-libc.org/releases/musl-$(MUSL_VERSION).tar.gz
	tar xf musl.tar.gz
	cd musl-$(MUSL_VERSION) ; CFLAGS=-fno-toplevel-reorder ./configure
	$(MAKE) -C musl-$(MUSL_VERSION) install

setup:
	mkdir -p gopath/{src,pkg,bin}

build:
	GOPATH=$(DIR)/gopath go get -d
	CC=/usr/local/musl/bin/musl-gcc GOPATH=$(DIR)/gopath go build -o build/ducktape -ldflags '-extldflags "-static"'

local: musl setup build
