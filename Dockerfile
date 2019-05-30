FROM depscloud/base:latest

ARG VERSION=0.0.1

RUN install-depscloud-binary dts ${VERSION}

RUN useradd -ms /bin/sh dts
WORKDIR /home/dts
USER dts

ENTRYPOINT [ "dts" ]
