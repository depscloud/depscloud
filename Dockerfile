FROM depscloud/base:latest

ARG VERSION=0.0.1

RUN install-depscloud-binary tracker ${VERSION}

RUN useradd -ms /bin/sh tracker
WORKDIR /home/tracker
USER tracker

ENTRYPOINT [ "tracker" ]
