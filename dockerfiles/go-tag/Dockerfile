FROM ocr.sh/depscloud/base:latest

WORKDIR /home/depscloud

ARG VERSION=0.2.10
ARG BINARY=""
ARG TARGETOS
ARG TARGETARCH

ADD https://github.com/depscloud/depscloud/releases/download/v${VERSION}/${BINARY}_${VERSION}_${TARGETOS}_${TARGETARCH}.tar.gz ${BINARY}.tar.gz
RUN tar zxf ${BINARY}.tar.gz && \
    rm ${BINARY}.tar.gz && \
    mv ${BINARY} /usr/bin/${BINARY} && \
    ln -s /usr/bin/${BINARY} /usr/bin/entrypoint

USER 13490:13490

ENTRYPOINT [ "/usr/bin/entrypoint" ]
