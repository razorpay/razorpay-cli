# =============================================================================
# Build Arguments
# =============================================================================
# REGISTRY: Container registry URL
#   - Default: c.rzp.io/razorpay (used in CI)
#   - Local override: harbor.razorpay.com/razorpay (pass via --build-arg)
ARG REGISTRY=c.rzp.io/razorpay

# ARG_BIN: Name of the binary to build (e.g., user, user_migration)
#   - This corresponds to a directory under cmd/
#   - Passed from Makefile's BINS variable, one at a time
ARG ARG_BIN=user

# =============================================================================
# Stage 1: Build
# =============================================================================
# This stage builds the Go binary.
FROM ${REGISTRY}/rzp-docker-image-inventory-multi-arch:rzp-golden-image-base-golang-1.25-alpine3.22 AS builder

ENV CGO_ENABLED=0
ENV GOPRIVATE=github.com/razorpay/*

# Install necessary tools for building
# bash is required by Makefile (SHELL := /bin/bash)
RUN apk add --no-cache git openssh-client make bash

# Add GitHub to known hosts for SSH connections
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

WORKDIR /src

# =============================================================================
# Authentication: Dual-mode support for private Go modules
# =============================================================================
# This RUN command configures Git authentication based on available credentials:
#   - CI (GIT_TOKEN available and non-empty): Uses HTTPS with token authentication
#   - Local (SSH available): Uses SSH agent forwarding
#
# The -s flag checks if file exists AND has content (size > 0).
# This handles the case where --secret is passed but the env var is empty.
RUN --mount=type=secret,id=git_token \
    --mount=type=ssh \
    if [ -s /run/secrets/git_token ]; then \
    echo "Using GIT_TOKEN for authentication (CI mode)"; \
    git config --global url."https://x-access-token:$(cat /run/secrets/git_token)@github.com/".insteadOf "https://github.com/"; \
    else \
    echo "Using SSH for authentication (local mode)"; \
    git config --global url."git@github.com:".insteadOf "https://github.com/"; \
    fi

# =============================================================================
# Proto Code Generation
# =============================================================================
# Proto source files are committed to the repository for version control.
# We only generate the Go code from these proto files during the build.
# This ensures reproducible builds and proper layer caching.

# Copy Makefile and buf configuration files (needed for proto generation)
# buf.lock pins the proto dependencies (e.g., google/api/annotations.proto)
COPY Makefile buf.yaml buf.gen.yaml buf.lock ./

# Copy proto source files (committed to repo, not fetched)
COPY proto/ ./proto/

# Install proto generation tools and generate Go code from proto files
# tools-proto: Installs buf and all protoc plugins (cached layer)
# proto-generate: Generates Go code in rpc/ directory (cached unless proto files change)
RUN make tools-proto && \
    make proto-generate

# =============================================================================
# Go Module Dependencies
# =============================================================================
# Pre-copy module files to leverage Docker's layer caching
COPY go.mod go.sum ./

# Download Go module dependencies
# This layer is cached unless go.mod or go.sum changes
RUN --mount=type=secret,id=git_token \
    --mount=type=ssh \
    go mod download

# =============================================================================
# Build Binary
# =============================================================================
# Copy the rest of the source code
COPY . .

# Re-declare ARG_BIN after source copy (placed AFTER COPY to preserve cache)
ARG ARG_BIN

# Build the Go binary for the specified service.
# TARGETARCH is an automatic platform ARG provided by Docker BuildKit.
# It will be 'amd64' or 'arm64' depending on the build target.
# This makes the build multi-arch compatible.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o /app-binary ./cmd/${ARG_BIN}

# =============================================================================
# Stage 2: Final
# =============================================================================
# This stage creates the final, lightweight image.
ARG REGISTRY
FROM ${REGISTRY}/rzp-docker-image-inventory-multi-arch:rzp-golden-image-base-golang-1.25-alpine3.22

WORKDIR /src

# It's good practice to run as a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the compiled binary from the 'builder' stage
COPY --from=builder /app-binary /app-binary

# Copy all configuration files needed by the application at runtime
# Copying the entire config directory ensures all binaries have access to their configs
COPY --chown=appuser:appgroup ./config ./config

# Set WORKDIR environment variable for foundation library
# Foundation checks this env var before trying 'go list' command
# This allows config loading without requiring Go toolchain in runtime image
ENV WORKDIR=/src

USER appuser

# Set the command to run the application
CMD ["/app-binary"]
