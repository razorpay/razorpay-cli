# Database Patterns

Libraries: `github.com/razorpay/goutils/spine v0.12.4` (wraps GORM), `gorm.io/gorm v1.31.1`, `gorm.io/driver/postgres v1.6.0`

## Repo Struct

```go
type Repo struct {
    database  *spine.Repo   // wraps DB with context-awareness
    telemetry *deps.Telemetry
}

func New(telemetry *deps.Telemetry, database *db.DB) *Repo {
    return &Repo{
        database:  &spine.Repo{Db: database},  // wrap *db.DB in spine.Repo
        telemetry: telemetry,
    }
}
```

## Database Operations

```go
// Create
func (r *Repo) Create(ctx context.Context, user *model.User) errors.IError {
    tx := r.database.DBInstance(ctx).Create(user)
    return spine.GetDBError(tx)  // converts *gorm.DB error to IError
}

// Get by ID
func (r *Repo) Get(ctx context.Context, id string) (*model.User, errors.IError) {
    user := &model.User{}
    // Convert public ID to internal ID first
    if err := user.FromPublicID(id); err != nil {
        return nil, errors.New("invalid public ID")
    }
    tx := r.database.DBInstance(ctx).Where("id = ?", user.ID).First(&user)
    return user, spine.GetDBError(tx)
}
```

## Key Patterns

**`r.database.DBInstance(ctx)`** — always use this, not `r.database.Db` directly. It propagates the context for tracing, timeouts, and transaction support.

**`spine.GetDBError(tx)`** — converts the GORM `*gorm.DB` result into `errors.IError`. Returns `nil` on success.

**Sentinel errors** (check these in the service layer, not repo):
```go
// In service.go, after calling repo:
if errors.Is(err, spine.UniqueConstraintViolation) { ... }
if errors.Is(err, spine.RecordNotFound) { ... }
```

## Dual Database Setup

```toml
# config/user/default.toml
[databases.primary_db]     # read-write DB
  [databases.primary_db.connection]
    dialect = "postgres"
    url = "postgres"       # Docker service name
    port = 5432
    username = "example"
    password = "example_pass"
    name = "example"

[databases.replica_db]     # read-only replica
  [databases.replica_db.connection]
    url = "postgres_replica"
    name = "example_replica"
```

Health check registration:
```go
dbCollection.NewHealthCheck("primary_db", false)  // false = required (fails readiness if down)
dbCollection.NewHealthCheck("replica_db", true)   // true = optional (warning only)
```

**Note:** The template only uses `primary_db` for writes/reads. In a real service, route reads to the replica by getting it with `dbCollection.Get("replica_db")`.

## Database Schema (migrations/000001_create_users.up.sql)

```sql
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE,
    created_at BIGINT DEFAULT EXTRACT(epoch FROM NOW())::BIGINT,
    updated_at BIGINT DEFAULT EXTRACT(epoch FROM NOW())::BIGINT,
    deleted_at BIGINT DEFAULT NULL  -- soft delete (NULL = not deleted)
);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE email IS NOT NULL;
```

`created_at`, `updated_at`, `deleted_at` are Unix epoch **bigints**, not timestamps — managed by `spine.SoftDeletableModel`.

## Repo Interface for Testability

The service layer consumes a `Repo` interface (defined in `internal/user/service/repo.go`), not the concrete `*repo.Repo`. This enables mock injection in unit tests without a real database.
