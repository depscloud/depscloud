build-docker:
	docker build . -t depscloud/api-builder

compile-docker:
	docker run --rm \
		-v $(PWD)/swagger:/go/src/github.com/deps-cloud/api/swagger \
		-v $(PWD)/v1alpha:/go/src/github.com/deps-cloud/api/v1alpha \
		depscloud/api-builder
