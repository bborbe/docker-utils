default: test

install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/docker-remote-repositories/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/docker-remote-sha-for-tag/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/docker-remote-size/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/docker-remote-tag-delete/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/docker-remote-tag-exists/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/docker-remote-tags/*.go

precommit: ensure format test check
	@echo "ready to commit"

ensure:
	dep ensure

format:
	@go get golang.org/x/tools/cmd/goimports
	@find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	@find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

check: lint vet errcheck

vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

lint:
	@go get github.com/golang/lint/golint
	@golint -min_confidence 1 $(shell go list ./... | grep -v /vendor/)

errcheck:
	@go get github.com/kisielk/errcheck
	@errcheck -ignore '(Close|Write|Fprint)' $(shell go list ./... | grep -v /vendor/)
