FROM depscloud/download:latest AS BUILDER

ARG VERSION=0.0.1
ENV RELEASE_CHAIN=goreleaser

RUN install-depscloud-binary depscloud ${VERSION} gateway

FROM depscloud/runtime:latest

COPY --from=BUILDER /usr/bin/gateway /usr/bin/gateway

ENTRYPOINT [ "/usr/bin/gateway" ]
