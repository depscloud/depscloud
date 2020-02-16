FROM depscloud/download:latest AS BUILDER

ARG VERSION=0.0.1
ARG HEALTH_PROBE_VERSION=0.3.1

RUN install-grpc-probe ${HEALTH_PROBE_VERSION}
RUN install-depscloud-binary tracker ${VERSION}

FROM depscloud/base:latest

COPY --from=BUILDER /usr/bin/grpc_health_probe /usr/bin/grpc_health_probe
COPY --from=BUILDER /usr/bin/tracker /usr/bin/tracker

ENTRYPOINT [ "tracker" ]
