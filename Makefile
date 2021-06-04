define WELCOME

Welcome to the deps.cloud project's source repository!

  :service must be one of - deps, extractor, gateway, indexer, tracker
  :db must be one of - cockroachdb, mariadb, mysql, postgres, sqlite

Available Targets:

        build-deps    install some development tools you need for the project
              deps    install software dependencies for the project

              test    runs unit tests for all projects
     test-:service    run unit tests for :service

            docker    build all containers using docker (common)
   docker-:service    build the container for :service using docker

        run-docker    run the stack using docker and a sqlite backend
    run-docker-:db    run the stack using docker and a :db backend

           install    installs all application binaries locally (rarely used)
  install-:service    install a specific service locally (often used for deps)


endef
export WELCOME

GIT_SHA ?= $(shell git rev-parse HEAD)
VERSION ?= local
TIMESTAMP ?= $(shell date +%Y-%m-%dT%T)
LD_FLAGS := -X main.version=${VERSION} -X main.commit=${GIT_SHA} -X main.date=${TIMESTAMP}

# Use a registry prefix when building the docker images locally.
ifeq (${USE_REGISTRY},1)
	REGISTRY_PREFIX = ocr.sh/
endif

help:
	@echo "$$WELCOME"

build-deps: .build-deps
.build-deps:
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups
#	GO111MODULE=off go get -u github.com/google/addlicense

deps: .deps
.deps:
	[[ -e services/extractor ]] && { make services/extractor/node_modules; }
	[[ -e go.mod ]] && { go mod download; go mod verify; }

fmt: .fmt
.fmt:
	cd services/extractor && npm run lint

	ls -1 services | grep -v extractor | xargs -I{} echo -n "./services/{} " | \
		xargs -I++ bash -c '{ goimports -w ./internal ++ ; go-groups -w ./internal ++ ; gofmt -s -w ./internal ++ ; }'

docker: docker-deps docker-extractor docker-gateway docker-indexer docker-tracker

install: install-deps install-extractor install-gateway install-indexer install-tracker

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

test: internal/hack/ test-extractor
	@ls -1 services | grep -v extractor | xargs -I{} echo -n "./services/{}/... " | \
		xargs -I{} make .test PACKAGES="./internal/... {}"

##===
## Common
##===

APP_LANG ?= go

.docker:
	docker build . \
		--build-arg BINARY=${BINARY} \
		--build-arg VERSION=${VERSION} \
		--build-arg GIT_SHA=${GIT_SHA} \
		-t depscloud/${BINARY}:latest \
		-f dockerfiles/${APP_LANG}-branch/Dockerfile

.install:
	go install -ldflags="${LD_FLAGS}" ./services/${BINARY}/

# Build the dockerfiles
dockerfiles: docker-base docker-devbase

## Build the `depscloud/base:latest` development container
docker-base:
	docker build ./dockerfiles/base -t ${REGISTRY_PREFIX}depscloud/base:latest

## Build the `depscloud/devbase:latest` development container
docker-devbase:
	docker build ./dockerfiles/devbase -t ${REGISTRY_PREFIX}depscloud/devbase:latest

## Build `depscloud/deps:latest` development container
docker-deps:
	@make .docker BINARY=deps

install-deps:
	@make .install BINARY=deps

test-deps:
	@make .test PACKAGES="./services/deps/..."


## Build `depscloud/extractor:latest` development container
docker-extractor:
	@make .docker BINARY=extractor APP_LANG=node

install-extractor:
	@cd services/extractor && npm run build

package-extractor:
	@cd services/extractor && npm run build && npm run package

services/extractor/node_modules: services/extractor/package-lock.json
	@cd services/extractor && npm install

test-extractor: services/extractor/node_modules
	@cd services/extractor && npm run test


## Build `depscloud/gateway:latest` development container
docker-gateway:
	@make .docker BINARY=gateway

install-gateway:
	@make .install BINARY=gateway

test-gateway:
	@make .test PACKAGES="./services/gateway/..."


## Build `depscloud/deps:latest` development container
docker-indexer:
	@make .docker BINARY=indexer

install-indexer:
	@make .install BINARY=indexer

test-indexer:
	@make .test PACKAGES="./services/indexer/..."


## Build `depscloud/deps:latest` development container
docker-tracker:
	@make .docker BINARY=tracker

install-tracker:
	@make .install BINARY=tracker

test-tracker:
	@make .test PACKAGES="./services/tracker/..."

## helper docker-compose configurations

.run:
	@cd docker/$(PLATFORM) && docker-compose up

run-docker-cockroachdb:
	@make .run PLATFORM=cockroachdb

run-docker-mariadb:
	@make .run PLATFORM=mariadb

run-docker-mysql:
	@make .run PLATFORM=mysql

run-docker-postgres:
	@make .run PLATFORM=postgres

run-docker-sqlite:
	@make .run PLATFORM=sqlite

run-docker: run-docker-sqlite
