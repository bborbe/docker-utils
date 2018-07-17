default: test

install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/docker-remote-repositories/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/docker-remote-tag-delete/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/docker-remote-tag-exists/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/docker-remote-tags/*.go

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

goimports:
	go get golang.org/x/tools/cmd/goimports

format: goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +
