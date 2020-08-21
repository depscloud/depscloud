GIT_SHA := $(shell git rev-parse HEAD)
VERSION ?= local
TIMESTAMP := $(shell date +%Y-%m-%dT%T)
LD_FLAGS := -X main.version=${VERSION} -X main.commit=${GIT_SHA} -X main.date=${TIMESTAMP}

build-deps:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups

deps:
	go mod download
	go mod verify

fmt:
	go-groups -w .
	gofmt -s -w .

test:
	go vet ./...
	# golint -set_exit_status ./...
	go test -v -cover -race ./...

install:
	go install -ldflags="${LD_FLAGS}" ./cmds/deps/
