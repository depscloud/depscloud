FROM depscloud/download:latest AS BUILDER

ARG HEALTH_PROBE_VERSION=0.3.1

RUN install-grpc-probe ${HEALTH_PROBE_VERSION}

FROM node:10-alpine3.11 AS INSTALLER

ARG VERSION=0.2.10

WORKDIR /home/depscloud

RUN apk -U upgrade && apk add build-base git ca-certificates python2 python3

ADD https://github.com/deps-cloud/extractor/releases/download/v${VERSION}/extractor.tar.gz extractor.tar.gz
RUN tar -zxvf extractor.tar.gz && rm extractor.tar.gz

RUN npm install --production

FROM node:10-alpine3.11

RUN apk -U upgrade && apk add ca-certificates

COPY --from=BUILDER /usr/bin/grpc_health_probe /usr/bin/grpc_health_probe
COPY --from=INSTALLER /home/depscloud /home/depscloud

WORKDIR /home/depscloud
USER 13490:13490

ENTRYPOINT [ "npm", "start", "--" ]
