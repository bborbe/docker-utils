default: test

install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/monitoring_check/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/monitoring_cron/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/monitoring_server/*.go

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

goimports:
	go get golang.org/x/tools/cmd/goimports

format: goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +
