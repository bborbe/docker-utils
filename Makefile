include tools.env

default: precommit

precommit: ensure format generate test check addlicense
	@echo "ready to commit"

ensure:
	go mod tidy
	go mod verify
	rm -rf vendor

format:
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec go run github.com/incu6us/goimports-reviser/v3@$(GOIMPORTS_REVISER_VERSION) -project-name github.com/bborbe/docker-utils "{}" \;

generate:
	rm -rf mocks avro
	go generate -mod=mod ./...

test:
	go test -mod=mod -p=$${GO_TEST_PARALLEL:-1} -cover -race $(shell go list -mod=mod ./... | grep -v /vendor/)

check: vet errcheck vulncheck

vet:
	go vet -mod=mod $(shell go list -mod=mod ./... | grep -v /vendor/)

errcheck:
	go run github.com/kisielk/errcheck@$(ERRCHECK_VERSION) -ignore '(Close|Write|Fprint)' $(shell go list -mod=mod ./... | grep -v /vendor/)

addlicense:
	go run github.com/google/addlicense@$(ADDLICENSE_VERSION) -c "Benjamin Borbe" -y $$(date +'%Y') -l bsd $$(find . -name "*.go" -not -path './vendor/*')

VULNCHECK_IGNORE ?= GO-2026-6603 GO-2026-6605 GO-2026-6607 GO-2026-6608 GO-2026-6610 GO-2026-6611 GO-2026-6613 GO-2026-6617

.PHONY: vulncheck
vulncheck:
	@PKGS="$(shell go list -mod=mod ./... | grep -v /vendor/)"; \
	IGNORE_JSON=$$(printf '%s\n' $(VULNCHECK_IGNORE) | jq -R . | jq -s .); \
	REMAIN=$$(go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) -format json $$PKGS 2>/dev/null | \
		jq -rs --argjson ignore "$$IGNORE_JSON" \
			'(map(select(.osv != null)) | map({key: .osv.id, value: (.osv.summary // "")}) | from_entries) as $$sum | \
			 map(select(.finding != null) | .finding) | \
			 map(select(.osv as $$o | $$ignore | index($$o) | not)) | \
			 map("\(.osv)\t\(.trace[-1].module)@\(.trace[-1].version) -> \(.fixed_version)\t\($$sum[.osv] // "")") | \
			 unique | .[]'); \
	if [ -n "$$REMAIN" ]; then \
		echo "Unexpected vulnerabilities (ignored: $(VULNCHECK_IGNORE)):"; \
		printf '%s\n' "$$REMAIN" | column -t -s "$$(printf '\t')"; \
		exit 1; \
	else \
		echo "No unignored vulnerabilities found"; \
	fi

install:
	go install github.com/bborbe/docker-utils/cmd/docker-remote-repositories
	go install github.com/bborbe/docker-utils/cmd/docker-remote-sha-for-tag
	go install github.com/bborbe/docker-utils/cmd/docker-remote-size-repositories
	go install github.com/bborbe/docker-utils/cmd/docker-remote-size-tags
	go install github.com/bborbe/docker-utils/cmd/docker-remote-tag-delete
	go install github.com/bborbe/docker-utils/cmd/docker-remote-tag-exists
	go install github.com/bborbe/docker-utils/cmd/docker-remote-tags
	go install github.com/bborbe/docker-utils/cmd/dockerhub-cleaner

