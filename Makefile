docker:
	docker build -t depscloud/base:latest .

dockerx:
	docker buildx rm depscloud--base || echo "depscloud--base does not exist"
	docker buildx create --name depscloud--base --use
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t depscloud/base:latest .
