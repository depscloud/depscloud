PROTOC_VERSION ?= 3.9.1
PROTOC_OS_ARCH ?= linux-x86_64

default:

build-deps-protoc:
	curl -sSL -o protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-${PROTOC_OS_ARCH}.zip
	unzip -d build-deps/protoc protoc.zip
	rm -rf protoc.zip

build-deps-go:
	GO111MODULE=off go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	GO111MODULE=off go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go
	GO111MODULE=off go get -u github.com/gogo/protobuf/protoc-gen-gogo

build-deps:
	mkdir build-deps
	make build-deps-protoc
	make build-deps-go

compile:
	bash compile.sh

clean:
	rm -rf build-deps
	rm -rf node_modules
