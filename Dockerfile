FROM depscloud/download:latest AS BUILDER

ARG BINARY=""
ARG VERSION=0.0.1
ARG HEALTH_PROBE_VERSION=0.3.1
ENV RELEASE_CHAIN=goreleaser

RUN install-grpc-probe ${HEALTH_PROBE_VERSION}
RUN install-depscloud-binary depscloud ${VERSION} ${BINARY}

FROM depscloud/base:latest

WORKDIR /home/depscloud
USER 13490:13490

COPY --from=BUILDER /usr/bin/grpc_health_probe /usr/bin/grpc_health_probe
COPY --from=BUILDER /usr/bin/${BINARY} /usr/bin/${BINARY}

ENTRYPOINT [ "${BINARY}" ]
