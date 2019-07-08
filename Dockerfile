FROM depscloud/base:latest

ARG VERSION=0.0.2

RUN install-depscloud-binary indexer ${VERSION}

RUN useradd -ms /bin/sh indexer
WORKDIR /home/indexer
USER indexer

ENTRYPOINT [ "indexer" ]
