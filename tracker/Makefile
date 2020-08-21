default: install

# moved out of deps to decrease build time
build-deps:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups
	GO111MODULE=off go get -u github.com/mitchellh/gox

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

deploy:
	mkdir -p bin
	gox -os="windows darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
	gox -os="linux" -arch="amd64 386 arm arm64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"

docker:
	docker build -t depscloud/tracker:latest -f Dockerfile.dev .
