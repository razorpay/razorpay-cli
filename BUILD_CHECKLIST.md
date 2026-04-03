# Agent-Readiness Extraction Checklist

**Repository:** go-foundation-v2 (Go service template)
**Module:** `github.com/razorpay/go-foundation-v2`
**Purpose:** Production-ready Golang microservice template built on the Razorpay Foundation framework.
**Domain:** Go service template — extract architectural patterns, conventions, and Foundation framework usage so AI agents can help developers build new services from this template.

## Extraction Tasks

Each task becomes a single file under `.agents/skills/repo-skill/`.

---

### [x] TASK-1: Service Architecture Overview
**Output:** `.agents/skills/repo-skill/core/architecture.md`

Extract:
- Layer structure: cmd → server → service → repo → model
- How each layer is wired in `cmd/user/main.go`
- Dependency injection pattern (telemetry, repo, service passed down)
- gRPC server + HTTP gateway dual-transport pattern via `foundation.WithGRPCHandlers` / `foundation.WithHTTPHandlers`
- Health check registration (`foundation.WithHealthChecks`)
- Shutdown signal handling (`foundation.WithShutdownSignals`)
- Primary + replica database pattern (`dbCollection.Get("primary_db")`)

Files to read: `cmd/user/main.go`, `internal/user/server.go`, `internal/user/handler.go`

---

### [x] TASK-2: Foundation Framework Wiring
**Output:** `.agents/skills/repo-skill/core/foundation-framework.md`

Extract:
- `foundation.NewServer()` usage pattern and config directory convention (`./config/user`)
- `providers.WithCustomConfig()` for injecting typed AppConfig structs
- `server.Container().Telemetry()` — accessing telemetry from the server container
- `server.Container().GetDatabaseCollection()` — database collection pattern
- `server.Start()` — all available `foundation.With*` options and their purposes
- `server.Context()` — context lifetime for HTTP handler registration
- What the Foundation library manages automatically (logging, metrics, tracing, interceptors)

Files to read: `cmd/user/main.go`, `config/user/default.toml`, `go.mod`

---

### [x] TASK-3: Proto & gRPC Workflow
**Output:** `.agents/skills/repo-skill/core/proto-grpc-workflow.md`

Extract:
- Proto file location: `proto/<module>/<entity>/v1/<entity>.proto`
- buf.build configuration: `buf.yaml`, `buf.gen.yaml`, `buf.lock`
- Proto conventions: package name, go_package option, HTTP annotations via `google.api.http`
- buf.validate field validation annotations (string constraints, required, email)
- Error embedding in response messages (embed `go_foundation_v2.error.v1.Error` in response protos)
- Generated code location: `rpc/` (excluded from git, auto-generated in Docker builds)
- `make proto-fetch` vs `make proto-generate` vs `make proto-refresh` — when to use each
- Why proto source files (`proto/`) ARE committed but generated code (`rpc/`) is NOT
- `grpc-gateway` HTTP mapping: how gRPC methods map to REST endpoints

Files to read: `proto/go_foundation_v2/user/v1/user.proto`, `proto/go_foundation_v2/error/v1/error.proto`, `buf.yaml`, `buf.gen.yaml`, `Makefile`

---

### [x] TASK-4: Error Handling System
**Output:** `.agents/skills/repo-skill/core/error-handling.md`

Extract:
- Error class declaration pattern using `goutilserrors.NewClass()` with severity + kind
- All defined error classes: `ErrValidationFailure`, `ErrInternalServerError`, `ErrBadRequestError`, `ErrUnauthenticated`, `ErrorInvalidData`, `ErrNotFound`, `ErrEntityAlreadyExists`, `ErrServerUnavailable`, `ErrUnexpected`, `ErrContextCancelled`
- Error code naming: `GOFND000000N` identifier codes from `internal/errors/codes.go`
- `classToCodeMap` — how error classes map to gRPC status codes
- `WrapError(ierr)` — converts IError to gRPC-compatible error with embedded proto details
- `ToProtoPublicError(err)` — converts IError to proto Error message for response embedding
- Error creation in service layer: `errpkg.ErrBadRequestError.New(code).Wrap(err)`
- Handling `spine.UniqueConstraintViolation` and `spine.RecordNotFound` sentinel errors
- Validation errors: `protovalidate.Validator`, `foundationErr.FromProtoValidationError(err)`, `WithPublicMetadata()`
- Pattern in server: return `&userv1.User{Error: errpkg.ToProtoPublicError(ierr)}, errpkg.WrapError(ierr)`

Files to read: `internal/errors/errors.go`, `internal/errors/codes.go`, `internal/user/server.go`, `internal/user/service/service.go`

---

### [x] TASK-5: Database & Repository Patterns
**Output:** `.agents/skills/repo-skill/modules/technical/database-patterns.md`

Extract:
- `spine.Repo` wrapper usage: `&spine.Repo{Db: database}`
- `r.database.DBInstance(ctx)` — context-aware DB instance retrieval
- GORM operations pattern: `.Create(model)`, `.Where().First()` returning `*gorm.DB`
- `spine.GetDBError(tx)` — extracting IError from GORM transaction result
- `spine.UniqueConstraintViolation` and `spine.RecordNotFound` — sentinel error detection
- Dual database setup: `primary_db` (read/write) + `replica_db` (read-only, optional health check)
- `dbCollection.NewHealthCheck("primary_db", false)` — required vs optional health checks
- Repository interface contract in service layer (`service/repo.go`) for testability

Files to read: `internal/user/repo/repo.go`, `internal/user/service/repo.go`, `cmd/user/main.go`, `migrations/000001_create_users.up.sql`, `config/user/default.toml`

---

### [x] TASK-6: Model Patterns
**Output:** `.agents/skills/repo-skill/modules/domain/model-patterns.md`

Extract:
- `spine.SoftDeletableModel` embedding — provides `id`, `created_at`, `updated_at`, `deleted_at`
- Required model methods: `TableName()`, `EntityName()`, `EntityPrefix()`
- Entity prefix pattern for public IDs: `"user_"` prefix, `entityPrefix = "user_"`
- `FromPublicID(id)` — two-format ID support: 14-char internal OR `user_<14chars>` public
- `SetDefaults()` — hook for default values, returns `errors.IError`
- `ToProto()` — converts internal model to protobuf response message
- `FromProto(req)` — populates model from protobuf request (only for create/update, not get)
- `New()` constructor pattern — always use this, not struct literal

Files to read: `internal/user/model/model.go`, `internal/user/model/model_test.go`, `migrations/000001_create_users.up.sql`

---

### [x] TASK-7: Config Management
**Output:** `.agents/skills/repo-skill/modules/technical/config-management.md`

Extract:
- TOML config file structure and location convention: `config/<service>/default.toml`, `dev.toml`, `route_config.toml`
- Custom config struct pattern: `AppConfig` with `mapstructure` tags, nested structs (e.g. `AWSConfig`)
- Config validation: implement `Validate() error` on both root and nested config structs
- `providers.WithCustomConfig(config)` — how to inject custom typed config into Foundation
- Accessing custom config after boot: `config` variable is populated after `foundation.NewServer()`
- Foundation-managed config sections: `[databases.*]`, `[grpc_server]`, `[telemetry]`, `[health_check]`, `[errors]`
- Config override pattern: `default.toml` → env-specific override (e.g. `dev.toml`)
- gRPC interceptors in config: basic auth, passport JWT, access logging, header filters
- AWS config example: region + s3_bucket with ozzo-validation

Files to read: `internal/config/config.go`, `config/user/default.toml`, `config/user/dev.toml`, `config/user/route_config.toml`, `cmd/user/main.go`

---

### [x] TASK-8: Build System & Service Rename
**Output:** `.agents/skills/repo-skill/core/build-system.md`

Extract:
- Makefile key variables: `BINS`, `IMAGE_PREFIX` — what to update when creating a new service
- `make build ENV=local` vs `make build` (Docker) — difference and when to use each
- `make test ENV=local`, `make lint ENV=local` — local development commands
- `make rename NAME=my-new-service` — what it updates (go.mod module path, import statements, git remote)
- `make proto-fetch`, `make proto-generate`, `make proto-refresh`, `make proto-lint`
- Files to manually update after `make rename`: `.github/workflows/*.yaml`, `Dockerfile.*`
- Docker Compose setup: services, postgres, postgres_replica
- `.air.toml` — live reload configuration for local development

Files to read: `Makefile`, `docs/build-system.md`, `docs/migration-guide.md`, `docker-compose.yml`, `.air.toml`

---

### [x] TASK-9: Testing Patterns
**Output:** `.agents/skills/repo-skill/modules/technical/testing-patterns.md`

Extract:
- Test file location convention: `*_test.go` co-located with code
- `testify` usage patterns from `internal/user/model/model_test.go`
- How to test model methods (TableName, EntityPrefix, FromPublicID, ToProto, FromProto)
- Repository interface (`service/repo.go`) — how it enables mock injection in unit tests
- `make test ENV=local` — how to run tests
- Test struct pattern and setup

Files to read: `internal/user/model/model_test.go`, `internal/user/service/repo.go`

---

## Mandatory Skills Installed

| Skill | Purpose |
|-------|---------|
| `code-security` | Security review for Go code |
| `tech-spec-reviewer` | Technical spec review |
| `devstack` | Deploy/debug via helmfile on devstack |
| `log-volume-optimizer` | Optimize log volume in Go services |
| `go-code-reviewer` | 5-layer Go code review framework |

## Status

- [x] Phase 0: Pre-flight check (go.mod exists, clean state detected)
- [x] Phase 1: Directory structure initialized (`.agents/`, `.claude/`)
- [x] Phase 1: Agentfill installed (Claude Code, Cursor, Gemini CLI)
- [x] Phase 1: Mandatory skills installed (5 skills)
- [x] Phase 2: Domain extraction (TASK-1 through TASK-9) — all 9 files written
- [x] Phase 3: AGENTS.md written
