#!/usr/bin/env bash

readonly home=$(pwd)
readonly protoc=${home}/build-deps/protoc/bin/protoc

for protofile in $(find ${home} -name *.proto | grep -v build-deps); do
    workdir=$(dirname ${protofile})
    file=$(basename ${protofile})

    echo "generating ${protofile}"

    cd "${workdir}"
    ${protoc} \
        -I=. \
        -I=$GOPATH/src \
        -I=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
        --gogo_out=plugins=grpc:. \
        --grpc-gateway_out=logtostderr=true:. \
        --swagger_out=logtostderr=true:. \
        ${file}
done
