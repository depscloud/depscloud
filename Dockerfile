FROM depscloud/base:latest

ARG VERSION=0.0.1
ARG HEALTH_PROBE_VERSION=0.3.1

COPY download-health-probe.sh /usr/bin/download-health-probe

RUN download-health-probe ${HEALTH_PROBE_VERSION} && \
    rm -rf /usr/bin/download-health-probe && \
    install-depscloud-binary tracker ${VERSION}

RUN useradd -ms /bin/sh tracker
WORKDIR /home/tracker
USER tracker

ENTRYPOINT [ "tracker" ]
