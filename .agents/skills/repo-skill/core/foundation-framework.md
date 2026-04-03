# Foundation Framework Wiring

`github.com/razorpay/foundation v1.0.0-alpha.8` — bootstraps the gRPC+HTTP server, telemetry, databases, interceptors, and health checks.

## What Foundation Manages Automatically

Reading `config/user/default.toml`, Foundation handles:
- gRPC server on `:8080`, HTTP gateway on `:8081`, internal metrics on `:8082`
- Structured logging (zap), Prometheus metrics, OpenTelemetry traces
- Access log interceptor (request/response logging)
- Optional basic auth interceptor (`[grpc_server.interceptor.basic_auth]`)
- Optional passport JWT auth (`[grpc_server.interceptor.passport]`)
- Database connection pools per `[databases.*]` entry (primary + replica)
- Health checks: goroutine count, GC max pause + your registered checks
- Graceful shutdown (30s timeout, SIGTERM/SIGINT + any extra signals you add)
- Config loading: `default.toml` base + env override (e.g. `dev.toml`)

## Server Bootstrap

```go
config := &config.AppConfig{}  // your typed config struct (uses mapstructure tags)
server, err := foundation.NewServer(
    "./config/user",                     // config directory path
    providers.WithCustomConfig(config),  // inject typed config
)
```

After `NewServer()` returns, `config` is already populated from TOML files.

## Accessing Built-in Dependencies

```go
// Telemetry: logger + metrics + tracer
tel := server.Container().Telemetry()  // *deps.Telemetry from foundation/deps

// All configured databases
dbCollection := server.Container().GetDatabaseCollection()
primaryDB, err := dbCollection.Get("primary_db")  // key must match [databases.primary_db] in config
```

## Start Options

```go
server.Start(
    foundation.WithGRPCHandlers(fn1, fn2),    // each: func(*grpc.Server) error
    foundation.WithHTTPHandlers(fn1, fn2),    // each: func(*runtime.ServeMux, address string) error
    foundation.WithHealthChecks(
        dbCollection.NewHealthCheck("primary_db", false),  // false = required (fails readiness)
        dbCollection.NewHealthCheck("replica_db", true),   // true = optional (warning only)
    ),
    foundation.WithShutdownSignals(syscall.SIGHUP),  // adds to SIGTERM+SIGINT
)
```

## Config Structure (default.toml sections)

| Section | What it configures |
|---------|-------------------|
| `[databases.<name>]` | Connection string, pool settings per DB |
| `[grpc_server]` | Addresses, shutdown timeout, interceptors |
| `[grpc_server.interceptor.basic_auth]` | Enable/disable, credentials map |
| `[grpc_server.interceptor.passport]` | JWT auth, JWKS host, retries |
| `[telemetry.logger]` | Log level, secure fields, redaction rules |
| `[telemetry.metrics]` | Exporter (prometheus), namespace, service name |
| `[health_check]` | Goroutine count limit, GC pause threshold |
| `[errors.mappings]` | Error mapping service names |

## Custom Config Pattern

```go
// internal/config/config.go
type AppConfig struct {
    AWS AWSConfig `mapstructure:"aws"`
}
type AWSConfig struct {
    Region   string `mapstructure:"region"`
    S3Bucket string `mapstructure:"s3_bucket"`
}
func (c *AppConfig) Validate() error { ... }  // Foundation calls this on boot

// In TOML:
// [aws]
//   region = "ap-south-1"
//   s3_bucket = "random_bucket"
```

Foundation calls `Validate()` on the custom config — if it returns an error, server startup fails.
