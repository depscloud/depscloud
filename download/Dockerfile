ARG BASE_TAG=latest

FROM depscloud/base:${BASE_TAG}

RUN apk add -u curl

COPY ./install-depscloud-binary /usr/bin/install-depscloud-binary
COPY ./install-grpc-probe /usr/bin/install-grpc-probe
