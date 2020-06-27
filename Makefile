build-docker:
	docker build docker/api-builder/ -t depscloud/api-builder

compile-docker:
	docker run --rm \
		-v $(PWD):/go/src/github.com/depscloud/api \
		depscloud/api-builder \
		bash scripts/compile.sh

version-patch:
	bash scripts/version.sh patch

version-minor:
	bash scripts/version.sh minor

version-major:
	bash scripts/version.sh major
