# Config Management

## File Layout

```
config/<service>/
  default.toml      # base config (all envs)
  dev.toml          # env override (only changed keys)
  route_config.toml # per-route auth/behavior overrides
```

Foundation loads `default.toml` first, then merges the env-specific file (e.g., `dev.toml`).

## Custom Config Pattern

```go
// internal/config/config.go
type AppConfig struct {
    AWS AWSConfig `mapstructure:"aws"`  // maps to [aws] section in TOML
}
type AWSConfig struct {
    Region   string `mapstructure:"region"`
    S3Bucket string `mapstructure:"s3_bucket"`
}

// Validation is called by Foundation at startup
func (c *AppConfig) Validate() error {
    return c.AWS.Validate()
}
func (c AWSConfig) Validate() error {
    return validation.ValidateStruct(&c,
        validation.Field(&c.Region, validation.Required, validation.In("us-east-1", "us-west-2", "ap-south-1")),
        validation.Field(&c.S3Bucket, validation.Required, validation.Length(3, 63)),
    )
}
```

**Rules:**
- Use `mapstructure` tags (not `toml`) — Foundation uses viper/mapstructure internally
- Implement `Validate() error` on root config — Foundation calls it; startup fails if it errors
- Validate nested structs from the parent's `Validate()` method

## Injecting Custom Config

```go
config := &config.AppConfig{}
server, err := foundation.NewServer("./config/user",
    providers.WithCustomConfig(config),
)
// After NewServer() returns, config is fully populated
tel := server.Container().Telemetry()
// Pass config fields down as needed:
userservice.New(tel, repo, config.AWS)
```

## Foundation-Managed Config Sections

Do NOT define these in your custom struct — Foundation manages them:

| Section | Purpose |
|---------|---------|
| `[databases.<name>]` | DB connections; access via `Container().GetDatabaseCollection()` |
| `[grpc_server]` | Server addresses, shutdown timeout, interceptors |
| `[grpc_server.interceptor.basic_auth]` | Enable/credentials for basic auth |
| `[grpc_server.interceptor.passport]` | JWT auth via JWKS |
| `[grpc_server.interceptor.accesslog]` | Request logging settings |
| `[telemetry]` | Logger level, metrics exporter, secure fields |
| `[health_check]` | Max goroutine count, GC pause thresholds |
| `[errors.mappings]` | Error mapping module service names |

## Route Config (config/user/route_config.toml)

Overrides interceptor behavior per gRPC route path:

```toml
"/common.health.v2.HealthService/ReadinessCheck" = { basic_auth = false }
"/common.health.v2.HealthService/LivenessCheck"  = { basic_auth = true }
```

## Key Config Values (default.toml)

```toml
[databases.primary_db.connection_pool]
  maxopenconnections = 25
  maxidleconnections = 10
  connectionmaxlifetime = 300  # seconds
  connectionmaxidletime = 60

[grpc_server.server_addresses]
  grpc = ":8080"
  http = ":8081"
  internal = ":8082"

[telemetry.logger]
  level = "DEBUG"
  secure_fields = ["password", "token", "authorization", "secret", "key"]

[telemetry.metrics]
  exporter = "prometheus"
  namespace = "foundation_app"
  service_name = "foundation_service"  # update this for your service
```
