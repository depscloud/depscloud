## ======
## Docker
## ======

builder-grpc-golang:
	docker build docker/builder-grpc-golang/ -t depscloud/builder-grpc-golang

builder-grpc-nodejs:
	docker build docker/builder-grpc-nodejs/ -t depscloud/builder-grpc-nodejs

builder-grpc-python:
	docker build docker/builder-grpc-python/ -t depscloud/builder-grpc-python

builder: builder-grpc-golang builder-grpc-nodejs builder-grpc-python

## ===================
## Compilation Targets
## ===================

compile-golang:
	docker run --rm -it \
		-v $(PWD):/go/src/github.com/depscloud/api \
		-w /go/src/github.com/depscloud/api \
		depscloud/builder-grpc-golang \
		bash scripts/compile-files.sh compile-golang

compile-nodejs:
	docker run --rm -it \
		-v $(PWD):/depscloud/api \
		-w /depscloud/api \
		depscloud/builder-grpc-nodejs \
		bash scripts/compile-nodejs.sh
	cp LICENSE packages/depscloud-api-nodejs

compile-python:
	# compile src
	docker run --rm -it \
		-v $(PWD):/depscloud/api \
		-w /depscloud/api \
		depscloud/builder-grpc-python \
		bash scripts/compile-files.sh compile-python
	cp LICENSE packages/depscloud-api-python/LICENSE.txt

compile-swagger:
	docker run --rm -it \
		-v $(PWD):/depscloud/api \
		-w /depscloud/api \
		depscloud/builder-grpc-golang \
		bash scripts/compile-swagger.sh

compile: compile-golang compile-nodejs compile-python compile-swagger

## ==================
## Versioning Helpers
## ==================

version-patch:
	bash scripts/version.sh patch

version-minor:
	bash scripts/version.sh minor

version-major:
	bash scripts/version.sh major
