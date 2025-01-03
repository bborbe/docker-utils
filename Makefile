
default: precommit

precommit: ensure format generate test check addlicense
	@echo "ready to commit"

ensure:
	go mod verify
	go mod vendor

format:
	go run -mod=vendor github.com/incu6us/goimports-reviser/v3 -project-name github.com/bborbe/docker-utils -format -excludes vendor ./...

generate:
	rm -rf mocks avro
	go generate -mod=vendor ./...

test:
	go test -mod=vendor -p=1 -cover -race $(shell go list -mod=vendor ./... | grep -v /vendor/)

check: vet errcheck vulncheck

vet:
	go vet -mod=vendor $(shell go list -mod=vendor ./... | grep -v /vendor/)

errcheck:
	go run -mod=vendor github.com/kisielk/errcheck -ignore '(Close|Write|Fprint)' $(shell go list -mod=vendor ./... | grep -v /vendor/)

addlicense:
	go run -mod=vendor github.com/google/addlicense -c "Benjamin Borbe" -y $$(date +'%Y') -l bsd $$(find . -name "*.go" -not -path './vendor/*')

vulncheck:
	go run -mod=vendor golang.org/x/vuln/cmd/govulncheck $(shell go list -mod=vendor ./... | grep -v /vendor/)

install:
	go install github.com/bborbe/docker-utils/cmd/docker-remote-repositories
	go install github.com/bborbe/docker-utils/cmd/docker-remote-sha-for-tag
	go install github.com/bborbe/docker-utils/cmd/docker-remote-size-repositories
	go install github.com/bborbe/docker-utils/cmd/docker-remote-size-tags
	go install github.com/bborbe/docker-utils/cmd/docker-remote-tag-delete
	go install github.com/bborbe/docker-utils/cmd/docker-remote-tag-exists
	go install github.com/bborbe/docker-utils/cmd/docker-remote-tags
	go install github.com/bborbe/docker-utils/cmd/dockerhub-cleaner

