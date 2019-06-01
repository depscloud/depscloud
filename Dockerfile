FROM depscloud/base:latest

ARG VERSION=0.0.1

RUN install-depscloud-binary gateway ${VERSION}

RUN useradd -ms /bin/sh gateway
WORKDIR /home/gateway
USER gateway

ENTRYPOINT [ "gateway" ]
