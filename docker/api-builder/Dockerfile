FROM debian:stretch-slim

ENV NODE_VERSION=12.8.0
ENV GO_VERSION=1.12.10
ENV PROTOC_VERSION=3.9.1
ENV GO111MODULE=off

ENV GRPC_PROTOC_PLUGINS_BUILD=5a57ed07-2c01-4973-a73d-eede65fb159b
ENV GRPC_PROTOC_PLUGINS_COMMIT=80f0c331f9c1b16f701a2d93e81770873459430f
ENV GRPC_PROTOC_PLUGINS_TAG=1.31.0-dev

RUN apt-get update && apt-get install -y unzip curl xz-utils git

RUN curl -sSL -o go.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xf go.tar.gz && \
    rm -rf go.tar.gz

RUN curl -sSL -o node.tar.xz https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.xz && \
    tar -C /usr/local -xf node.tar.xz && \
    rm -rf node.tar.xz && \
    mv /usr/local/node-v${NODE_VERSION}-linux-x64 /usr/local/node

RUN curl -sSL -o protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    unzip -d /usr/local/protoc protoc.zip && \
    rm -rf protoc.zip

RUN curl -sSL -o grpc-protoc.tar.gz https://packages.grpc.io/archive/2020/06/${GRPC_PROTOC_PLUGINS_COMMIT}-${GRPC_PROTOC_PLUGINS_BUILD}/protoc/grpc-protoc_linux_x64-${GRPC_PROTOC_PLUGINS_TAG}.tar.gz && \
    mkdir -p /usr/local/grpc-protoc/bin && \
    tar -C /usr/local/grpc-protoc/bin -xf grpc-protoc.tar.gz && \
    rm -rf protoc-plugins.tar.gz

ENV GOPATH=/go
ENV PATH="${PATH}:/usr/local/go/bin:/usr/local/node/bin:/usr/local/protoc/bin:/usr/local/grpc-protoc/bin:/go/bin"

RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/gogo/protobuf/protoc-gen-gogo
RUN go get -u github.com/go-bindata/go-bindata/...

WORKDIR /go/src/github.com/depscloud/api
