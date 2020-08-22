FROM depscloud/download:latest AS BUILDER

ARG HEALTH_PROBE_VERSION=0.3.1

RUN install-grpc-probe ${HEALTH_PROBE_VERSION}

FROM node:10 AS INSTALLER

ARG VERSION=0.2.10

WORKDIR /home/depscloud

ADD https://github.com/depscloud/depscloud/releases/download/v${VERSION}/extractor-${VERSION}.tar.gz extractor.tar.gz
RUN tar -zxvf extractor.tar.gz && rm extractor.tar.gz

RUN npm install --production

FROM node:10

COPY --from=BUILDER /usr/bin/grpc_health_probe /usr/bin/grpc_health_probe
COPY --from=INSTALLER /home/depscloud /home/depscloud

WORKDIR /home/depscloud
USER 13490:13490

ENTRYPOINT [ "npm", "start", "--" ]
