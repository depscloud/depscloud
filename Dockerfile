FROM depscloud/base:latest

ARG VERSION=0.0.1

RUN install-depscloud-binary dis ${VERSION}

ENTRYPOINT [ "dis" ]
