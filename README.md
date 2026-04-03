# Go Foundation v2 - Golang Service Template

A production-ready Golang service template for Razorpay developers. This template provides a standardized foundation for building microservices following Domain-Driven Design (DDD) principles.

## What's Included

- Complete service architecture with gRPC and HTTP endpoints
- Protocol Buffers integration with automatic code generation
- Database integration with migrations and health checks
- Containerized build system using Docker and Make
- CI/CD workflows with GitHub Actions
- Production-ready patterns (logging, metrics, error handling)
- Example user service demonstrating best practices

## Prerequisites

- Go 1.25+
- Docker and Docker Compose
- Git with access to Razorpay repositories
- Make utility

## Quick Start

### 1. Review Documentation First

Before making any changes, it's highly recommended that you review the [documentation](docs/README.md) to understand the Foundation framework and how services should be modeled. Start with the [Example Service](docs/example-service.md) to see the architecture in action.

### 2. Clone and Rename

```bash
# Rename the service
make rename NAME=my-new-service
```

The `make rename` command updates:
- Go module path in `go.mod`
- Import statements in all `.go` files
- Git remote origin URL

### 3. Update Dependencies

```bash
# Update Go dependencies
go mod tidy

# Verify everything builds
make build
```

### 4. Configure Protocol Buffers

Proto source files are committed to the repository for version control and reproducibility. You only need to fetch and generate when updating to newer proto versions:

```bash
# Fetch latest proto files from central repository (when you need updates)
make proto-fetch

# Generate Go code from proto files (automatically done in Docker builds)
make proto-generate

# Commit the updated proto files
git add proto/
git commit -m "chore: update proto files"
```

**Note:** Docker builds automatically generate RPC code from committed proto files - no manual proto-fetch needed for building!

### 5. Post-Rename Configuration

Manually update these files with service-specific values:
- `.github/workflows/*.yaml` - Update image names and repository references
- `Dockerfile.*` - Update any hardcoded references
- `Makefile` - Update `BINS` and `IMAGE_PREFIX` variables

## Essential Commands

### Build and Test

```bash
# Build locally (faster for development)
make build ENV=local

# Build using Docker (default for CI/CD)
make build

# Run tests
make test ENV=local

# Lint code
make lint ENV=local
```

### Proto Management

```bash
# Fetch latest proto files (updates proto/ directory)
make proto-fetch

# Generate Go code from protos (creates rpc/ directory)
make proto-generate

# Lint proto files
make proto-lint

# Complete proto refresh (fetch + generate)
make proto-refresh

# Commit proto updates to version control
git add proto/
git commit -m "chore: update proto files"
```

**Important:** Proto source files (`proto/`) are committed to git for reproducibility. Generated code (`rpc/`) is excluded from git and auto-generated during builds.

### Docker Commands

```bash
# Run services locally
docker-compose up

# Build and push Docker images
make push
```

## Project Structure

```
cmd/                      # Service entry points
├── user/                # Example user service
└── user_migration/      # Database migration tool

internal/                # Private application code
└── user/               # Example domain logic
    ├── server.go       # gRPC server implementation
    ├── handler.go      # Handler registration
    ├── service/        # Business logic layer
    ├── repo/           # Data access layer
    ├── model/          # Domain models
    └── migrations/     # Database migrations

config/                  # Configuration files
proto/                   # Protocol buffer definitions (committed)
rpc/                     # Generated proto code (auto-generated, not committed)
scripts/                 # Build system makefiles
```

## Key Makefile Targets

Update these variables in the root `Makefile`:

```makefile
BINS ?= user user_migration    # Your service binaries
IMAGE_PREFIX ?= fnd-           # Docker image prefix
```

Common targets:
- `make build` - Build binaries
- `make test` - Run tests
- `make lint` - Lint code
- `make clean` - Clean artifacts
- `make help` - Show all targets

## Documentation

Detailed documentation is available in the `docs/` directory:

- **[Example Service](docs/example-service.md)** - Architecture and patterns of the included User Service
- **[Migration Guide](docs/migration-guide.md)** - Step-by-step guide to transition from example to your custom service
- **[Best Practices](docs/best-practices.md)** - Coding standards and architectural guidelines
- **[Build System](docs/build-system.md)** - Makefile structure and Docker build process
- **[Proto Management](docs/proto-management.md)** - Protocol Buffers workflow and central repository integration

## GitHub Actions Workflows

The template includes automated CI/CD workflows:

- **Build Workflow**: Multi-architecture builds, Docker image publishing
- **Lint Workflow**: Code formatting and linting checks
- **Security Workflows**: Automated security scanning

After renaming your service, update workflow files in `.github/workflows/` with your service name and image references.

## Additional Resources

- [Razorpay Go Style Guide](https://github.com/razorpay/styleguide/blob/master/go/style.md)
- [Foundation Library Documentation](https://github.com/razorpay/foundation)
- [Proto Repository](https://github.com/razorpay/proto)

## Acknowledgments

This template builds upon:
- [upi-switch](https://github.com/razorpay/upi-switch/) - Foundation patterns
- [thockin/go-build-template](https://github.com/thockin/go-build-template) - Build system inspiration

---

For questions or support, reach out to the Developer Experience Engineering team on slack [(#developer-experience)](https://razorpay.enterprise.slack.com/archives/C08DS8AE7T8)
