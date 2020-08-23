FROM depscloud/download:latest AS BUILDER

ARG VERSION=0.0.1
ARG HEALTH_PROBE_VERSION=0.3.1
ENV RELEASE_CHAIN=goreleaser

RUN install-grpc-probe ${HEALTH_PROBE_VERSION}
RUN install-depscloud-binary depscloud ${VERSION} tracker

FROM depscloud/runtime:latest

COPY --from=BUILDER /usr/bin/grpc_health_probe /usr/bin/grpc_health_probe
COPY --from=BUILDER /usr/bin/tracker /usr/bin/tracker

ENTRYPOINT [ "tracker" ]
