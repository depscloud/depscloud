default: install

deps:
	go get -v ./...

test:
	go vet ./...
	golint -set_exit_status ./...
	go test -v ./...

install:
	go install

deploy:
	mkdir -p bin
	gox -os="windows linux darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
	gox -os="linux" -arch="arm" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
	GOOS=linux GOARCH=arm64 go build -o bin/gateway_linux_arm64

docker:
	GOOS=linux GOARCH=amd64 go build
	docker build -t depscloud/gateway:latest -f Dockerfile.dev .

dockerx:
	docker buildx rm depscloud--gateway || echo "depscloud--gateway does not exist"
	docker buildx create --name depscloud--gateway --use
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t depscloud/gateway:latest .
