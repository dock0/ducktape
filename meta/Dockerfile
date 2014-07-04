FROM dock0/arch
MAINTAINER akerl <me@lesaker.org>
RUN pacman -Syu --needed --noconfirm go base-devel ruby openssh
RUN gem install --no-user-install --no-rdoc --no-ri targit
WORKDIR /opt/ducktape
CMD ["make", "local"]
