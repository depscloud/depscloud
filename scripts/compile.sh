#!/usr/bin/env bash

set -e -o pipefail

readonly ROOT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )" )"

readonly NODEJS_PACKAGE_DIR="${ROOT_DIR}/packages/depscloud-api-nodejs"

readonly PROTO_DIR="${ROOT_DIR}/proto"
readonly SWAGGER_DIR="${ROOT_DIR}/swagger"

function update_package_nodejs {
    echo "updating package nodejs"
    cp -R proto/* ${NODEJS_PACKAGE_DIR}

    pushd ${NODEJS_PACKAGE_DIR}
    npm run generate
    popd
}

function update_protoc {
    echo "update packages using protoc"

    pushd ${PROTO_DIR}
    for file in $(find . -name *.proto | cut -c 3-); do
        echo "compiling ${file}"

        protoc \
            -I=. \
            -I=$GOPATH/src \
            -I=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
            --gogo_out=plugins=grpc:$GOPATH/src \
            --grpc-gateway_out=logtostderr=true:$GOPATH/src \
            --swagger_out=logtostderr=true:${SWAGGER_DIR} \
            ${file}
    done
    popd
}

function update_swagger {
    echo "generating swagger fs"
    pushd ${SWAGGER_DIR}
    go-bindata -fs -pkg swagger -o swagger.go $(find . -iname *.swagger.json)
    popd
}

update_package_nodejs
update_protoc
update_swagger
