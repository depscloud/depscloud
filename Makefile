GIT_SHA := $(shell git rev-parse HEAD)
VERSION ?= local
TIMESTAMP := $(shell date +%Y-%m-%dT%T)
LD_FLAGS := -X main.version=${VERSION} -X main.commit=${GIT_SHA} -X main.timestamp=${TIMESTAMP}

build-deps:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups
	GO111MODULE=off go get -u github.com/mitchellh/gox

deps:
	go mod download
	go mod verify

fmt:
	go-groups -w .
	gofmt -s -w .

test:
	go vet ./...
	# golint -set_exit_status ./...
	go test -v ./...

install:
	go install -ldflags="${LD_FLAGS}" ./cmds/deps/

deploy:
	mkdir -p bin
	gox -ldflags="${LD_FLAGS}" -os="windows darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./cmds/deps
	gox -ldflags="${LD_FLAGS}" -os="linux" -arch="amd64 386 arm arm64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./cmds/deps
