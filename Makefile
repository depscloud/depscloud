DOCKER_PLATFORMS ?= linux/amd64,linux/arm64,linux/arm/v7
DOCKER_BUILDX_ARGS ?=

build-docker-base:
	docker build -t depscloud/base:latest base

build-dockerx-base:
	docker buildx create --name depscloud--base --use || echo "depscloud--base exists"
	docker buildx build --platform $(DOCKER_PLATFORMS) -t depscloud/base:latest base $(DOCKER_BUILDX_ARGS)

build-docker-download:
	docker build -t depscloud/download:latest download

build-dockerx-download:
	docker buildx create --name depscloud--download --use || echo "depscloud--download exists"
	docker buildx build --platform $(DOCKER_PLATFORMS) -t depscloud/download:latest download $(DOCKER_BUILDX_ARGS)
