FROM dock0/arch
MAINTAINER akerl <me@lesaker.org>
RUN pacman -Syu --needed --noconfirm go base-devel openssh
WORKDIR /opt/ducktape
CMD ["make", "local"]
