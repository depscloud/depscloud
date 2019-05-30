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
	GOOS=linux GOARCH=arm64 go build -o bin/dis_linux_arm64

docker:
	docker build -t depscloud/dis:latest -f Dockerfile.dev .

dockerx:
	docker buildx rm depscloud--dis || echo "depscloud--dis does not exist"
	docker buildx create --name depscloud--dis --use
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t depscloud/dis:latest .
