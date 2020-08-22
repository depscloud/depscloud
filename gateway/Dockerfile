FROM depscloud/download:latest AS BUILDER

ARG VERSION=0.0.1
ENV RELEASE_CHAIN=goreleaser

RUN install-depscloud-binary depscloud ${VERSION} gateway

FROM depscloud/base:latest

COPY --from=BUILDER /usr/bin/gateway /usr/bin/gateway

WORKDIR /home/depscloud
USER 13490:13490

ENTRYPOINT [ "/usr/bin/gateway" ]
