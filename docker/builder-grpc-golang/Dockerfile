FROM golang:1.14-buster

RUN apt-get update && apt-get install -y unzip

ENV PROTOC_VERSION=3.9.1

RUN curl -sSL -o protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    unzip -d /usr/local/protoc protoc.zip && \
    rm -rf protoc.zip

ENV PATH="${PATH}:/usr/local/protoc/bin"

RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/gogo/protobuf/protoc-gen-gogo
RUN go get -u github.com/go-bindata/go-bindata/...
