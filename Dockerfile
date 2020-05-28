FROM depscloud/download:latest AS BUILDER

ARG VERSION=0.0.1

RUN install-depscloud-binary gateway ${VERSION}

FROM depscloud/base:latest

COPY --from=BUILDER /usr/bin/gateway /usr/bin/gateway

WORKDIR /home/depscloud
USER 13490:13490

ENTRYPOINT [ "/usr/bin/gateway" ]
