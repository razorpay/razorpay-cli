# go-foundation-v2 — Agent Instructions

Production-ready Go microservice template using the Razorpay Foundation framework (gRPC + HTTP, PostgreSQL, buf/proto). Module: `github.com/razorpay/go-foundation-v2`.

## Skills Index

Load these skills **before** working on any matching task:

| Skill | Trigger |
|-------|---------|
| `.agents/skills/repo-skill/core/architecture.md` | Adding/modifying layers, wiring dependencies, server startup |
| `.agents/skills/repo-skill/core/foundation-framework.md` | Foundation server config, Container, telemetry, DB collection |
| `.agents/skills/repo-skill/core/proto-grpc-workflow.md` | Proto files, buf commands, gRPC/HTTP handler registration |
| `.agents/skills/repo-skill/core/error-handling.md` | Creating errors, error codes, gRPC status mapping, WrapError |
| `.agents/skills/repo-skill/core/build-system.md` | Makefile targets, service rename, Docker, CI |
| `.agents/skills/repo-skill/modules/domain/model-patterns.md` | Domain models, spine.SoftDeletableModel, ToProto/FromProto |
| `.agents/skills/repo-skill/modules/technical/database-patterns.md` | spine.Repo, GORM, Create/Get, primary/replica DB |
| `.agents/skills/repo-skill/modules/technical/config-management.md` | Custom config structs, TOML layout, mapstructure tags |
| `.agents/skills/repo-skill/modules/technical/testing-patterns.md` | Unit tests, mock injection, testify patterns |
| `.agents/skills/code-security/SKILL.md` | Writing or reviewing any code |
| `.agents/skills/go-code-reviewer/SKILL.md` | Go code review requests |
| `.agents/skills/log-volume-optimizer/SKILL.md` | Log optimization |
| `.agents/skills/tech-spec-reviewer/SKILL.md` | Tech spec review |
| `.agents/skills/devstack/SKILL.md` | Deploying to devstack |

## Critical Rules

1. **Proto-generated code** (`rpc/`) is NOT committed. Run `make proto-generate` locally or let Docker do it. If imports from `rpc/` fail, run `make proto-generate` first.

2. **Dual error return** — gRPC handlers ALWAYS return both: embed error in response proto (`Error: errpkg.ToProtoPublicError(err)`) AND return `errpkg.WrapError(err)` as the error return.

3. **Context propagation** — always use `r.database.DBInstance(ctx)` (not `r.database.Db`) for database ops.

4. **Service layer returns IError** — never return `error` from service/repo; always `errors.IError` from `github.com/razorpay/goutils/errors`.

5. **Spine sentinel errors** — check `spine.UniqueConstraintViolation` and `spine.RecordNotFound` in the service layer after repo calls, not in the repo itself.

6. **Repo interface** — `service/repo.go` defines the `Repo` interface; use it for mock injection in tests.

## Repository Structure

```
cmd/<service>/main.go           # entry point — Foundation bootstrap
internal/<domain>/
  server.go                     # gRPC server struct + RPC method implementations
  handler.go                    # GRPCHandler + HTTPHandler registration
  service/
    service.go                  # business logic
    repo.go                     # Repo interface (for mocking)
  repo/repo.go                  # spine/GORM data access
  model/model.go                # domain model + proto conversion
config/<service>/               # TOML config files
proto/<module>/<entity>/v1/     # .proto source files (COMMITTED)
rpc/                            # generated Go code (NOT committed)
migrations/                     # SQL migration files
```

## Common Tasks

**Add a new entity:** Follow the pattern in `internal/user/` — create model, repo, service, server, handler, then wire in `cmd/<service>/main.go`. Add proto file in `proto/`, run `make proto-generate`.

**Add a new RPC:** Define in `.proto` file, run `make proto-generate`, implement in `server.go`, add error handling with the dual return pattern.

**Create a new service from this template:** Run `make rename NAME=<your-service>`, update `Makefile` BINS and REPO_PREFIX, update `.github/workflows/*.yaml`.
