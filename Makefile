default: install

build-deps:
	GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go
	GO111MODULE=off go get -u github.com/gogo/protobuf/protoc-gen-gogo
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups

fmt:
	go-groups -w .
	gofmt -s -w .

deps:
	go mod download
	go mod verify

test:
	go vet ./...
	golint -set_exit_status ./...
	go test -v ./...

install:
	go install

deploy:
	mkdir -p bin
	CGO_ENABLED=0 gox -ldflags='-w -s -extldflags "-static"' -os="windows darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
	CGO_ENABLED=0 gox -ldflags='-w -s -extldflags "-static"' -os="linux" -arch="amd64 386 arm arm64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"

docker:
	docker build -t depscloud/indexer:latest -f Dockerfile.dev .
