SOURCE_FILES ?= ./...
TEST_PATTERN ?= .
TEST_OPTIONS ?=

export GO111MODULE := on
export GOBIN       := $(shell pwd)/bin
export PATH        := $(GOBIN):$(PATH)

export GOLANGCI_LINT_VERSION := v1.58.1

BINARY  := razorpay
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Install Go module dependencies
setup:
	go mod download
.PHONY: setup

# Run the end-to-end test suite. The suite is gated by the `e2e` build tag
# and runs the compiled CLI against the live Razorpay API. Export
# RAZORPAY_KEY_ID and RAZORPAY_KEY_SECRET before invoking — the suite
# exits with a non-zero status if either is missing. See tests/README.md
# for details.
test:
	go test -tags=e2e -v -timeout=10m ./tests/...
.PHONY: test

# Build the razorpay binary for the current platform
build:
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BINARY) .
.PHONY: build

# Cross-compile for all supported platforms
build-all-platforms:
	env GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BINARY)-darwin-amd64 .
	env GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BINARY)-linux-arm64 .
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BINARY)-windows-amd64.exe .
.PHONY: build-all-platforms

# Format all Go source files
fmt:
	find . -name '*.go' | xargs gofmt -w -s
.PHONY: fmt

# Run golangci-lint
lint: bin/golangci-lint
	./bin/golangci-lint run ./...
.PHONY: lint

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_LINT_VERSION)

# Tidy go.mod and verify no uncommitted changes
go-mod-tidy:
	@go mod tidy -v
	@git diff HEAD
	@git diff-index --quiet HEAD
.PHONY: go-mod-tidy

# Full CI check: build, test, lint, tidy
ci: build test lint go-mod-tidy
.PHONY: ci

# Tag and push a new release — actual binaries are built by CI (see .github/workflows/release.yml)
# Opens $EDITOR for changelog entry, writes CHANGELOG.md, commits, tags, and pushes.
release:
	@bash scripts/release.sh
.PHONY: release

# Clean all build artefacts
clean:
	go clean ./...
	rm -f $(BINARY) \
	      $(BINARY)-darwin-amd64 $(BINARY)-darwin-arm64 \
	      $(BINARY)-linux-amd64  $(BINARY)-linux-arm64  \
	      $(BINARY)-windows-amd64.exe
	rm -f coverage.txt
	rm -rf dist/ bin/
.PHONY: clean

.DEFAULT_GOAL := build
