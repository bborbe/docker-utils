
deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/maxbrunsfeld/counterfeiter

install:
	go install github.com/bborbe/docker-utils/cmd/docker-remote-repositories
	go install github.com/bborbe/docker-utils/cmd/docker-remote-sha-for-tag
	go install github.com/bborbe/docker-utils/cmd/docker-remote-size-repositories
	go install github.com/bborbe/docker-utils/cmd/docker-remote-size-tags
	go install github.com/bborbe/docker-utils/cmd/docker-remote-tag-delete
	go install github.com/bborbe/docker-utils/cmd/docker-remote-tag-exists
	go install github.com/bborbe/docker-utils/cmd/docker-remote-tags
	go install github.com/bborbe/docker-utils/cmd/dockerhub-cleaner

precommit: ensure format generate test check
	@echo "ready to commit"

ensure:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

format:
	@go get golang.org/x/tools/cmd/goimports
	@find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	@find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

generate:
	go get github.com/maxbrunsfeld/counterfeiter
	rm -rf mocks
	go generate ./...

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

check: lint vet errcheck

lint:
	@go get golang.org/x/lint/golint
	@golint -min_confidence 1 $(shell go list ./... | grep -v /vendor/)

vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

errcheck:
	@go get github.com/kisielk/errcheck
	@errcheck -ignore '(Close|Write|Fprint)' $(shell go list ./... | grep -v /vendor/)
