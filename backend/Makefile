.DEFAULT_GOAL := check

.PHONY: fmt lint test build audit vulncheck check

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

lint:
	go vet ./...

test:
	go test ./...

build:
	go build ./...

# Validates module integrity and lists the resolved dependency graph. Run the
# pinned govulncheck command in docs/development.md for vulnerability analysis.
audit:
	go mod verify
	go list -m all

# Network access is required the first time this pinned tool is downloaded.
vulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...

check: lint test build audit
