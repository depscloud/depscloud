#!/usr/bin/env bash

readonly home=$(pwd)
readonly protoc=${home}/build-deps/protoc/bin/protoc

for file in $(find . -name *.proto | grep -v build-deps | grep -v node_modules | cut -c 3-); do
    out=$(dirname ${file})

    ${protoc} \
        -I=${home} \
        -I=$GOPATH/src \
        -I=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
        --gogo_out=plugins=grpc:$GOPATH/src \
        --grpc-gateway_out=logtostderr=true:$GOPATH/src \
        --swagger_out=logtostderr=true:$GOPATH/src \
        ${file}
done
