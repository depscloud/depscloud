default: install

# moved out of deps to decrease build time
build-deps:
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups
	GO111MODULE=off go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	GO111MODULE=off go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go
	GO111MODULE=off go get -u github.com/gogo/protobuf/protoc-gen-gogo

fmt:
	go-groups -w .
	gofmt -s -w .

deps:
	go get -v ./...

test:
	go vet ./...
	golint -set_exit_status ./...
	go test -v ./...

install:
	go install

deploy:
	mkdir -p bin
	gox -os="windows darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
	gox -os="linux" -arch="amd64 386 arm arm64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"

docker:
	docker build -t depscloud/tracker:latest -f Dockerfile.dev .

dockerx:
	docker buildx rm depscloud--tracker || echo "depscloud--tracker does not exist"
	docker buildx create --name depscloud--tracker --use
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t depscloud/tracker:latest .
