# Build System & Service Rename

## Makefile Service Variables (Part A — Edit These)

```makefile
BINS        := user user_migration   # binary names → cmd/<name>/main.go directories
REPO_PREFIX := fnd                   # image prefix: harbor.../fnd-user:<tag>
REGISTRY    := harbor.razorpay.com/razorpay  # local dev registry

PROTO_MODULES := go_foundation_v2    # proto module dirs to fetch from proto repo
PROTO_GIT_URL := https://github.com/razorpay/proto.git
PROTO_BRANCH  := master
```

## Creating a New Service From This Template

```bash
# 1. Clone
git clone <this-repo> my-service && cd my-service

# 2. Rename (updates go.mod module path + all import statements + git remote)
make rename NAME=my-service

# 3. Update Makefile
#    BINS := <your-binary-names>
#    REPO_PREFIX := <your-prefix>

# 4. Manually update:
#    .github/workflows/*.yaml  — image names, registry references
#    Dockerfile.*              — hardcoded references
```

## Build Commands

| Command | What it does |
|---------|-------------|
| `make build ENV=local` | Build binaries locally (fast, no Docker, uses local Go) |
| `make build` | Build via Docker (default, multi-arch, used in CI) |
| `make test ENV=local` | Run tests locally |
| `make lint ENV=local` | Run golangci-lint locally |
| `make clean` | Remove build artifacts |
| `make help` | Show all targets with descriptions |

## Proto Commands

| Command | When to use |
|---------|------------|
| `make proto-fetch` | Pull updated `.proto` files from `github.com/razorpay/proto` into `proto/` |
| `make proto-generate` | Generate Go code from `proto/` into `rpc/` |
| `make proto-lint` | Lint proto files |
| `make proto-refresh` | `proto-fetch` + `proto-generate` |

**Docker builds always run `proto-generate` automatically — no manual step needed for CI.**

After `make proto-fetch`, commit the updated `proto/` files:
```bash
git add proto/
git commit -m "chore: update proto files"
```

## Docker Compose (local development)

```bash
docker-compose up        # starts: postgres (5432), postgres_replica (5433),
                         # user-service (8080/8081/8082/2345), prometheus (9090), grafana (3000)
```

The `user-service` container mounts `.:/src` and exposes port `2345` for Delve debugging.

## Tool Versions

```makefile
GOLANGCI_LINT_VERSION := v2.7.1
AIR_VERSION           := v1.63.4     # hot reload for devstack
```

Tools are installed to `.tools/` (project-local, not global). `.air.toml` configures hot-reload.

## CI/CD

GitHub Actions workflows handle:
- Multi-architecture Docker builds
- Image push to `c.rzp.io/razorpay` (production registry)
- Lint checks
- Security scanning (semgrep)

After `make rename`, update `.github/workflows/*.yaml` with your service name and image references.
