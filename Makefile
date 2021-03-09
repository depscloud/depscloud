GIT_SHA ?= $(shell git rev-parse HEAD)
VERSION ?= local
TIMESTAMP ?= $(shell date +%Y-%m-%dT%T)
LD_FLAGS := -X main.version=${VERSION} -X main.commit=${GIT_SHA} -X main.date=${TIMESTAMP}

# Use a registry prefix when building the docker images localy.
ifeq (${USE_REGISTRY},1)
	REGISTRY_PREFIX = ocr.sh/
endif

default: docker

build-deps: .build-deps
.build-deps:
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups
#	GO111MODULE=off go get -u github.com/google/addlicense

deps: .deps
.deps:
	@if [ -d extractor ]; then cd extractor && npm install; fi
	go mod download
	go mod verify

fmt: .fmt
.fmt:
	cd extractor && npm run lint
	goimports -w ./deps ./gateway ./indexer ./tracker ./internal
	go-groups -w ./deps ./gateway ./indexer ./tracker ./internal
	gofmt -s -w ./deps ./gateway ./indexer ./tracker ./internal
#	addlicense -c deps.cloud -l mit ./deps/**/*.go ./extractor/src/**/*.ts ./gateway/**/*.go ./indexer/**/*.go ./tracker/**/*.go ./internal/**/*.go

docker: deps/docker extractor/docker gateway/docker indexer/docker tracker/docker

install: deps/install extractor/install gateway/install indexer/install tracker/install

generate:
	docker run --rm -it \
		-v $(PWD)/indexer:/go/src/github.com/depscloud/depscloud/indexer \
		-w /go/src/github.com/depscloud/depscloud/indexer \
		ocr.sh/depscloud/builder-grpc-golang \
		go generate ./...
	make fmt

.test:
	go vet ${PACKAGES}
	#golint -set_exit_status ${PACKAGES}
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ${PACKAGES}

internal/hack/:
	mkdir -p internal/hack/
	openssl req -x509 -sha256 -newkey rsa:4096 \
		-keyout internal/hack/ca.key -out internal/hack/ca.crt \
		-nodes -subj '/CN=localhost'
	openssl req -new -newkey rsa:4096 \
		-keyout internal/hack/test.key -out internal/hack/test.csr \
		-nodes -subj "/CN=test"
	openssl x509 -req -sha256 -days 365 -in internal/hack/test.csr \
            -CA internal/hack/ca.crt -CAkey internal/hack/ca.key \
            -set_serial 01 -out internal/hack/test.crt

test: internal/hack/ extractor/test
	@make .test PACKAGES="./deps/... ./gateway/... ./indexer/... ./tracker/... ./internal/..."

##===
## Common
##===

.docker:
	docker build . \
		--build-arg BINARY=${BINARY} \
		--build-arg VERSION=${VERSION} \
		--build-arg GIT_SHA=${GIT_SHA} \
		-t depscloud/${BINARY}:latest \
		-f dockerfiles/go-branch/Dockerfile

.install:
	go install -ldflags="${LD_FLAGS}" ./${BINARY}/

# Build the dockerfiles
dockerfiles: base/docker devbase/docker

## Build the `depscloud/base:latest` development container
base/docker:
	docker build ./dockerfiles/base -t ${REGISTRY_PREFIX}depscloud/base:latest

## Build the `depscloud/devbase:latest` development container
devbase/docker:
	docker build ./dockerfiles/devbase -t ${REGISTRY_PREFIX}depscloud/devbase:latest

## Build `depscloud/deps:latest` development container
deps/docker:
	@make .docker BINARY=deps

deps/install:
	@make .install BINARY=deps

deps/test:
	@make .test PACKAGES="./deps/..."


## Build `depscloud/extractor:latest` development container
extractor/docker:
	@cd extractor && npm run docker

extractor/install:
	@cd extractor && npm run build

extractor/test:
	@cd extractor && npm run test


## Build `depscloud/gateway:latest` development container
gateway/docker:
	@make .docker BINARY=gateway

gateway/install:
	@make .install BINARY=gateway

gateway/test:
	@make .test PACKAGES="./gateway/..."


## Build `depscloud/deps:latest` development container
indexer/docker:
	@make .docker BINARY=indexer

indexer/install:
	@make .install BINARY=indexer

indexer/test:
	@make .test PACKAGES="./indexer/..."


## Build `depscloud/deps:latest` development container
tracker/docker:
	@make .docker BINARY=tracker

tracker/install:
	@make .install BINARY=tracker

tracker/test:
	@make .test PACKAGES="./tracker/..."

## helper docker-compose configurations

.run:
	@cd docker/$(PLATFORM) && docker-compose up

run/docker/cockroachdb:
	@make .run PLATFORM=cockroachdb

run/docker/mariadb:
	@make .run PLATFORM=mariadb

run/docker/mysql:
	@make .run PLATFORM=mysql

run/docker/postgres:
	@make .run PLATFORM=postgres

run/docker/sqlite:
	@make .run PLATFORM=sqlite

run/docker: run/docker/sqlite
