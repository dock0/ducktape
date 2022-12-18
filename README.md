ducktape
=======

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/dock0/ducktape/build.yml?branch=main)](https://github.com/dock0/ducktape/actions))
[![GitHub release](https://img.shields.io/github/release/dock0/ducktape.svg)](https://github.com/dock0/ducktape/releases)
[![License](https://img.shields.io/github/license/dock0/ducktape)](https://github.com/dock0/ducktape/blob/master/LICENSE)

Minimal bootstrapping script for downloading and extracting a root tarball

## Usage

To bootstrap a system using ducktape, do the following:

1. Put the `ducktape` binary on the system
2. Put a `cert` file in the same directory as the binary
3. Run `./ducktape https://example.org/root.tar.bz2`

It supports uncompressed, gzipped, and bzipped archives.

To use this for a Docker container, start your Dockerfile with the following:

```
FROM scratch
MAINTAINER {you}
ADD ducktape /ducktape
ADD cert /cert
RUN ["/ducktape", "https://example.org/path/to/root.tar.bz2"]
```

## License

ducktape is released under the MIT License. See the bundled LICENSE file for details.
