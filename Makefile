default: install

build-deps:
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups

deps:
	go get -v ./...

fmt:
	go-groups -w .
	gofmt -s -w .

test:
	go vet ./...
	golint -set_exit_status ./...
	go test -v ./...

install:
	go install

deploy:
	mkdir -p bin
	gox -os="windows darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
	gox -os="linux" -arch="amd64 386 arm arm64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"

docker:
	docker build -t depscloud/indexer:latest -f Dockerfile.dev .

dockerx:
	docker buildx rm depscloud--indexer || echo "depscloud--indexer does not exist"
	docker buildx create --name depscloud--indexer --use
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t depscloud/indexer:latest .
