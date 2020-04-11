#!/usr/bin/env bash

readonly SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

for file in $(find . -name *.proto | cut -c 3-); do
    echo "compiling ${file}"

    out=$(dirname ${file})

    protoc \
        -I=. \
        -I=$GOPATH/src \
        -I=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
        --gogo_out=plugins=grpc:$GOPATH/src \
        --grpc-gateway_out=logtostderr=true:$GOPATH/src \
        --swagger_out=logtostderr=true:. \
        ${file}
done

echo "generating swagger files"
go-bindata -fs -pkg swagger -o swagger/swagger.go $(find . -iname *.swagger.json)

echo "generating protoloader"
bash -s "${SCRIPT_DIR}/protoc-gen-protoloader.sh"
