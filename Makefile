# =============================================================================
# PART A: SERVICE-SPECIFIC CONFIGURATION
# =============================================================================
# Edit this section to customize for your service.

# Binary names - these map to directories under cmd/
BINS := user user_migration

# Repository prefix for container images
# Used by: Dockerfile, docker-compose.yml, CI workflow
# Final image: $(REGISTRY)/$(REPO_PREFIX)-{binary}:{tag}
REPO_PREFIX := fnd

# Container registry for local development
# CI uses c.rzp.io/razorpay (configured in workflow)
REGISTRY := harbor.razorpay.com/razorpay

# Proto configuration
# Proto modules to fetch (space-separated list of directories from proto repo)
PROTO_MODULES := go_foundation_v2
PROTO_GIT_URL := https://github.com/razorpay/proto.git
PROTO_BRANCH  := master

# =============================================================================
# PART B: CENTRALIZED COMMANDS
# =============================================================================
# Do not edit below this line.

# -----------------------------------------------------------------------------
# Variables (auto-detected)
# -----------------------------------------------------------------------------
SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c

# Disable verbose output unless V=1
ifeq ($(V),1)
    Q :=
else
    Q := @
    MAKEFLAGS += -s
endif

# OS and architecture detection
OS   := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')

# Go configuration
GOFLAGS := -mod=mod
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Proto paths
PROTO_ROOT := proto
RPC_ROOT   := rpc

# Build output
BIN_DIR := bin

# Local tools directory (project-specific, not global)
TOOLS_DIR := $(shell pwd)/.tools

# Tool versions
GOLANGCI_LINT_VERSION := v2.7.1
MOCKGEN_VERSION       := v1.6.0
AIR_VERSION           := v1.63.4

# Protoc plugin versions
BUF_VERSION           := latest 
PROTOC_GEN_GO_VERSION         := latest
PROTOC_GEN_GO_GRPC_VERSION    := latest
PROTOC_GEN_GRPC_GATEWAY_VERSION := latest
PROTOC_GEN_OPENAPIV2_VERSION  := latest

# -----------------------------------------------------------------------------
# Help
# -----------------------------------------------------------------------------
.PHONY: help
help: ## Show this help message
	$(Q)echo "Usage: make [target]"
	$(Q)echo ""
	$(Q)echo "Service: $(REPO_PREFIX) | Binaries: $(BINS)"
	$(Q)echo ""
	$(Q)echo "Variables:"
	$(Q)echo "  REGISTRY         = $(REGISTRY)"
	$(Q)echo "  REPO_PREFIX      = $(REPO_PREFIX)"
	$(Q)echo "  BINS             = $(BINS)"
	$(Q)echo "  PROTO_MODULES    = $(PROTO_MODULES)"
	$(Q)echo "  PROTO_GIT_URL    = $(PROTO_GIT_URL)"
	$(Q)echo "  PROTO_BRANCH     = $(PROTO_BRANCH)"
	$(Q)echo "  PROTO_ROOT       = $(PROTO_ROOT)"
	$(Q)echo "  RPC_ROOT         = $(RPC_ROOT)"
	$(Q)echo "  BIN_DIR          = $(BIN_DIR)"
	$(Q)echo "  TOOLS_DIR        = $(TOOLS_DIR)"
	$(Q)echo "  OS               = $(OS)"
	$(Q)echo "  ARCH             = $(ARCH)"
	$(Q)echo "  VERSION          = $(VERSION)"
	$(Q)echo "  COMMIT           = $(COMMIT)"
	$(Q)echo "  GOFLAGS          = $(GOFLAGS)"
	$(Q)echo ""
	$(Q)echo "Available targets:"
	$(Q)awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# -----------------------------------------------------------------------------
# Tools
# -----------------------------------------------------------------------------
.PHONY: tools
tools: tools-lint tools-mock tools-proto tools-air ## Install all development tools

.PHONY: tools-lint
tools-lint: ## Install golangci-lint
	$(Q)if [ ! -x "$(TOOLS_DIR)/golangci-lint" ]; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION) to $(TOOLS_DIR)..."; \
		mkdir -p $(TOOLS_DIR); \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(TOOLS_DIR) $(GOLANGCI_LINT_VERSION); \
	else \
		echo "golangci-lint already installed in $(TOOLS_DIR)"; \
	fi

.PHONY: tools-mock
tools-mock: ## Install mockgen
	$(Q)if [ ! -x "$(TOOLS_DIR)/mockgen" ]; then \
		echo "Installing mockgen $(MOCKGEN_VERSION) to $(TOOLS_DIR)..."; \
		mkdir -p $(TOOLS_DIR); \
		GOBIN=$(TOOLS_DIR) go install github.com/golang/mock/mockgen@$(MOCKGEN_VERSION); \
	else \
		echo "mockgen already installed in $(TOOLS_DIR)"; \
	fi

.PHONY: tools-proto
tools-proto: ## Install buf and protoc plugins
	$(Q)mkdir -p $(TOOLS_DIR)
	$(Q)if [ ! -x "$(TOOLS_DIR)/buf" ]; then \
		echo "Installing buf $(BUF_VERSION) to $(TOOLS_DIR)..."; \
		GOBIN=$(TOOLS_DIR) go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION); \
	else \
		echo "buf already installed in $(TOOLS_DIR)"; \
	fi
	$(Q)if [ ! -x "$(TOOLS_DIR)/protoc-gen-go" ]; then \
		echo "Installing protoc-gen-go $(PROTOC_GEN_GO_VERSION) to $(TOOLS_DIR)..."; \
		GOBIN=$(TOOLS_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION); \
	else \
		echo "protoc-gen-go already installed in $(TOOLS_DIR)"; \
	fi
	$(Q)if [ ! -x "$(TOOLS_DIR)/protoc-gen-go-grpc" ]; then \
		echo "Installing protoc-gen-go-grpc $(PROTOC_GEN_GO_GRPC_VERSION) to $(TOOLS_DIR)..."; \
		GOBIN=$(TOOLS_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION); \
	else \
		echo "protoc-gen-go-grpc already installed in $(TOOLS_DIR)"; \
	fi
	$(Q)if [ ! -x "$(TOOLS_DIR)/protoc-gen-grpc-gateway" ]; then \
		echo "Installing protoc-gen-grpc-gateway $(PROTOC_GEN_GRPC_GATEWAY_VERSION) to $(TOOLS_DIR)..."; \
		GOBIN=$(TOOLS_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@$(PROTOC_GEN_GRPC_GATEWAY_VERSION); \
	else \
		echo "protoc-gen-grpc-gateway already installed in $(TOOLS_DIR)"; \
	fi
	$(Q)if [ ! -x "$(TOOLS_DIR)/protoc-gen-openapiv2" ]; then \
		echo "Installing protoc-gen-openapiv2 $(PROTOC_GEN_OPENAPIV2_VERSION) to $(TOOLS_DIR)..."; \
		GOBIN=$(TOOLS_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@$(PROTOC_GEN_OPENAPIV2_VERSION); \
	else \
		echo "protoc-gen-openapiv2 already installed in $(TOOLS_DIR)"; \
	fi

.PHONY: tools-air
tools-air: ## Install Air live reload tool
	$(Q)if [ ! -x "$(TOOLS_DIR)/air" ]; then \
		echo "Installing air $(AIR_VERSION) to $(TOOLS_DIR)..."; \
		mkdir -p $(TOOLS_DIR); \
		GOBIN=$(TOOLS_DIR) go install github.com/air-verse/air@$(AIR_VERSION); \
	else \
		echo "air already installed in $(TOOLS_DIR)"; \
	fi

# -----------------------------------------------------------------------------
# Build
# -----------------------------------------------------------------------------
.PHONY: build
build: ## Build all binaries
	$(Q)echo "Building for $(OS)/$(ARCH)..."
	$(Q)for bin in $(BINS); do \
		echo "  → $$bin"; \
		CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build \
			-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT)" \
			-o $(BIN_DIR)/$$bin \
			./cmd/$$bin; \
	done
	$(Q)echo "✅ Build complete: $(BIN_DIR)/"

.PHONY: build-%
build-%: ## Build a specific binary (e.g., make build-user)
	$(Q)echo "Building $* for $(OS)/$(ARCH)..."
	$(Q)CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build \
		-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT)" \
		-o $(BIN_DIR)/$* \
		./cmd/$*
	$(Q)echo "✅ Built: $(BIN_DIR)/$*"

# -----------------------------------------------------------------------------
# Test
# -----------------------------------------------------------------------------
.PHONY: test
test: ## Run tests
	$(Q)echo "Running tests..."
	$(Q)go test -race -cover ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	$(Q)echo "Running tests with coverage..."
	$(Q)go test -race -coverprofile=coverage.out ./...
	$(Q)go tool cover -html=coverage.out -o coverage.html
	$(Q)echo "✅ Coverage report: coverage.html"

# -----------------------------------------------------------------------------
# Code Quality
# -----------------------------------------------------------------------------
.PHONY: lint
lint: tools-lint ## Run linter
	$(Q)echo "Running linter..."
	$(Q)$(TOOLS_DIR)/golangci-lint run ./...

.PHONY: fmt
fmt: tools-lint ## Format code
	$(Q)echo "Formatting code..."
	$(Q)$(TOOLS_DIR)/golangci-lint run --fix ./...
	$(Q)go fmt ./...

.PHONY: mocks
mocks: tools-mock ## Generate mocks
	$(Q)echo "Generating mocks..."
	$(Q)PATH=$(TOOLS_DIR):$$PATH go generate ./...

# -----------------------------------------------------------------------------
# Proto
# -----------------------------------------------------------------------------
.PHONY: proto-fetch
proto-fetch: ## Fetch proto files from remote repository (sparse checkout)
	$(Q)echo "Fetching proto files from $(PROTO_GIT_URL) (branch: $(PROTO_BRANCH))..."
	$(Q)echo "Modules: $(PROTO_MODULES)"
	$(Q)rm -rf $(PROTO_ROOT)
	$(Q)mkdir -p $(PROTO_ROOT)
	$(Q)cd $(PROTO_ROOT) && \
		git init --quiet && \
		git config core.sparseCheckout true && \
		printf '%s\n' $(PROTO_MODULES) > .git/info/sparse-checkout && \
		git remote add origin $(PROTO_GIT_URL) && \
		git fetch origin $(PROTO_BRANCH) --quiet && \
		git checkout $(PROTO_BRANCH) --quiet && \
		rm -rf .git
	$(Q)echo "✅ Proto files fetched to $(PROTO_ROOT)/"

.PHONY: proto-generate
proto-generate: tools-proto ## Generate Go code from proto files
	$(Q)echo "Generating RPC code..."
	$(Q)PATH=$(TOOLS_DIR):$$PATH $(TOOLS_DIR)/buf generate --timeout 5m
	$(Q)echo "✅ Generated code in $(RPC_ROOT)/"

.PHONY: proto-lint
proto-lint: tools-proto ## Lint proto files
	$(Q)echo "Linting proto files..."
	$(Q)$(TOOLS_DIR)/buf lint

.PHONY: proto-refresh
proto-refresh: proto-clean proto-fetch proto-generate ## Clean, fetch, and regenerate proto files
	$(Q)echo "✅ Proto refresh complete"

.PHONY: proto-clean
proto-clean: ## Remove generated proto files
	$(Q)rm -rf $(RPC_ROOT) $(PROTO_ROOT)

# -----------------------------------------------------------------------------
# Clean
# -----------------------------------------------------------------------------
.PHONY: clean
clean: ## Remove build artifacts
	$(Q)echo "Cleaning..."
	$(Q)rm -rf $(BIN_DIR) coverage.out coverage.html
	$(Q)echo "✅ Clean complete"

.PHONY: tools-clean
tools-clean: ## Remove local tools
	$(Q)echo "Removing local tools..."
	$(Q)rm -rf $(TOOLS_DIR)
	$(Q)echo "✅ Tools removed"

.PHONY: clean-all
clean-all: clean proto-clean tools-clean ## Remove all artifacts including proto and tools
	$(Q)echo "✅ Full cleanup complete"

# -----------------------------------------------------------------------------
# Version
# -----------------------------------------------------------------------------
.PHONY: version
version: ## Show version information
	$(Q)echo "Version: $(VERSION)"
	$(Q)echo "Commit:  $(COMMIT)"
	$(Q)echo "OS/Arch: $(OS)/$(ARCH)"
