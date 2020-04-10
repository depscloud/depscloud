build-deps:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups

fmt:
	go-groups -w .
	gofmt -s -w .

install:
	go install ./cmds/depscloud-cli/
