ducktape
=======

[![MIT Licensed](http://img.shields.io/badge/license-MIT-green.svg)](https://tldrlegal.com/license/mit-license)

Minimal bootstrapping script for downloading and extracting a root tarball

## Usage

To bootstrap a system with this, put the `ducktape` binary on the system and run `./ducktape https://example.org/root.tar.bz2`

It supports uncompressed, gzipped, and bzipped archives.

To use this for a Docker container, start your Dockerfile with the following:

```
FROM scratch
MAINTAINER {you}
ADD ducktape /.ducktape
RUN ["/.ducktape", "https://example.org/path/to/root.tar.bz2"]
```

## License

ducktape is released under the MIT License. See the bundled LICENSE file for details.
