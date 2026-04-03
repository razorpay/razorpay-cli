# Build System

Overview of the Makefile structure and Docker build process.

## Makefile Structure

The build system uses a modular approach with multiple makefiles:

```
Makefile                 # Main Makefile (service-specific variables)
scripts/
├── variables.mk        # Global variables and configuration
├── common.mk          # Common targets and rules
├── docker.mk          # Docker-based build targets
└── local.mk           # Local development targets
```

## Key Variables

Configure these variables in the root `Makefile` for your service:

```makefile
# Update these variables as per your application
BINS ?= user user_migration          # Your service binaries
IMAGE_PREFIX ?= fnd-                 # Docker image prefix
```

## Essential Make Targets

### Development Targets (Faster)
```bash
# Build binaries using local Go installation
make build ENV=local

# Build specific binary
make build ENV=local BIN=user

# Run tests locally
make test ENV=local

# Lint code locally
make lint ENV=local
```

### Docker Targets (Default for CI/CD)
```bash
# Build using Docker containers
make build

# Run tests in containers
make test

# Lint code in containers
make lint

# Build and push Docker images
make push
```

### Proto Targets
```bash
# Fetch proto files from central repo
make proto-fetch

# Generate Go code from protos
make proto-generate

# Lint proto files
make proto-lint

# Complete proto workflow (clean + fetch + generate + lint)
make proto-refresh
```

### Utility Targets
```bash
# Show all available targets and variables
make help

# Clean build artifacts
make clean

# Show version information
make version
```

## Environment-Specific Builds

The build system supports different environments:

```bash
# Local: Uses your local Go installation (faster for development)
make build ENV=local

# Docker: Uses Docker containers (default, recommended for CI/CD)
make build ENV=docker
```

## Extending the Makefile

To add custom targets for your service, add them to the root `Makefile` after the includes:

```makefile
# Custom targets for your service
run-local: build
	./bin/your-service --config ./config/your-service/dev.toml

migrate-up:
	./bin/your-service_migration up

migrate-down:
	./bin/your-service_migration down
```

Use existing patterns from `scripts/common.mk` for consistency.

## Build Dependencies

The build system automatically handles:
- Go module dependencies via `go mod download`
- Proto generation tools (buf, protoc-gen-*)
- Linting tools (golangci-lint)
- Docker image dependencies from the base image

## Common Issues

### Proto Generation Fails
```bash
# Ensure buf is installed and proto files are fetched
make proto-fetch
make proto-generate ENV=local
```

### Build Failures
```bash
# Clean and rebuild
make clean
make build ENV=local
```

### Dependency Issues
```bash
# Update Go dependencies
go mod tidy
go mod download
```

