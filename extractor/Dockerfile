FROM depscloud/download:latest AS installer

ARG HEALTH_PROBE_VERSION=0.3.1

RUN install-grpc-probe ${HEALTH_PROBE_VERSION}

FROM node:10 AS builder

WORKDIR /home/depscloud

ARG VERSION=0.2.10
ADD https://github.com/depscloud/depscloud/releases/download/v${VERSION}/extractor-${VERSION}.tar.gz extractor.tar.gz
RUN tar -zxvf extractor.tar.gz && rm extractor.tar.gz

RUN npm install --production

FROM node:10

COPY --from=installer /usr/bin/grpc_health_probe /usr/bin/grpc_health_probe
COPY --from=builder /home/depscloud /home/depscloud

WORKDIR /home/depscloud
USER 13490:13490

ENTRYPOINT [ "npm", "start", "--" ]
