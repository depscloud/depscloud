FROM depscloud/base:latest

ARG VERSION=0.0.2

RUN install-depscloud-binary dis ${VERSION}

RUN useradd -ms /bin/sh dis
WORKDIR /home/dis
USER dis

ENTRYPOINT [ "dis" ]
