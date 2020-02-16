FROM depscloud/download:latest AS BUILDER

ARG VERSION=0.0.2

RUN install-depscloud-binary indexer ${VERSION}

FROM depscloud/base:latest

COPY --from=BUILDER /usr/bin/indexer /usr/bin/indexer

USER deps-cloud

ENTRYPOINT [ "/usr/bin/indexer" ]
