FROM depscloud/download:latest AS BUILDER

ARG VERSION=0.0.2
ENV RELEASE_CHAIN=goreleaser

RUN install-depscloud-binary depscloud ${VERSION} indexer

FROM depscloud/runtime:latest

COPY --from=BUILDER /usr/bin/indexer /usr/bin/indexer

ENTRYPOINT [ "/usr/bin/indexer" ]
