all: deps extractor gateway indexer tracker

##===
## Build `depscloud/deps:latest` development container
##==
.deps:
	@cd deps && make docker
deps: .deps

##===
## Build `depscloud/extractor:latest` development container
##==
.extractor:
	@cd extractor && npm run docker
extractor: .extractor

##===
## Build `depscloud/gateway:latest` development container
##==
.gateway:
	@cd gateway && make docker
gateway: .gateway

##===
## Build `depscloud/deps:latest` development container
##==
.indexer:
	@cd indexer && make docker
indexer: .indexer

##===
## Build `depscloud/deps:latest` development container
##==
.tracker:
	@cd tracker && make docker
tracker: .tracker
