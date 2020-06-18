#!/usr/bin/env bash

readonly SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "generating nodejs content"
cp -R proto/* nodejs/
pushd nodejs
npm run generate
popd

###///====

echo "generating golang and swagger content"
pushd proto
for file in $(find . -name *.proto | cut -c 3-); do
    echo "compiling ${file}"

    protoc \
        -I=. \
        -I=$GOPATH/src \
        -I=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
        --gogo_out=plugins=grpc:$GOPATH/src \
        --grpc-gateway_out=logtostderr=true:$GOPATH/src \
        --swagger_out=logtostderr=true:../swagger \
        ${file}
done
popd

###///====

echo "generating swagger files"
pushd swagger
go-bindata -fs -pkg swagger -o swagger.go $(find . -iname *.swagger.json)
popd
