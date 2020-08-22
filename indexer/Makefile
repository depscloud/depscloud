default: install

build-deps:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups

generate:
	docker run --rm -it \
		-v $(PWD):/go/src/github.com/depscloud/depscloud/indexer \
		-w /go/src/github.com/depscloud/depscloud/indexer \
		depscloud/builder-grpc-golang \
		go generate ./...
	make fmt

fmt:
	go-groups -w .
	gofmt -s -w .

deps:
	go mod download
	go mod verify

test:
	go vet ./...
	golint -set_exit_status ./...
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

install:
	go install

docker:
	docker build -t depscloud/indexer:latest -f Dockerfile.dev .
