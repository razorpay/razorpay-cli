# Service Architecture

This is a **Go microservice template** using the Razorpay Foundation framework. The example User service demonstrates canonical DDD layering.

## Layer Structure

```
cmd/<service>/main.go          → wires dependencies, starts Foundation server
internal/<domain>/server.go    → gRPC server struct, implements proto service interface
internal/<domain>/handler.go   → registers gRPC + HTTP handlers with Foundation
internal/<domain>/service/     → business logic layer
  service.go                   → use-case implementations
  repo.go                      → Repo interface (for testability)
internal/<domain>/repo/repo.go → data access layer (spine/GORM)
internal/<domain>/model/       → domain model, proto conversion
```

## Dependency Injection (cmd/user/main.go)

```go
// 1. Bootstrap Foundation server
server, err := foundation.NewServer("./config/user",
    providers.WithCustomConfig(config),  // inject typed AppConfig
)

// 2. Pull built-in deps from Container
tel := server.Container().Telemetry()
dbCollection := server.Container().GetDatabaseCollection()
primaryDB, err := dbCollection.Get("primary_db")  // must match config key

// 3. Wire your layers bottom-up
userRepo    := repo.New(tel, primaryDB)
userService := userservice.New(tel, userRepo, config.AWS)
userServer, err := user.New(tel, userService)

// 4. Start with handler registration
server.Start(
    foundation.WithGRPCHandlers(userServer.GRPCHandler),
    foundation.WithHTTPHandlers(userServer.HTTPHandler(server.Context())),
    foundation.WithHealthChecks(
        dbCollection.NewHealthCheck("primary_db", false),  // false = required
        dbCollection.NewHealthCheck("replica_db", true),   // true = optional
    ),
    foundation.WithShutdownSignals(syscall.SIGHUP),
)
```

## Server Struct Pattern (internal/<domain>/server.go)

```go
type Server struct {
    validator protovalidate.Validator   // proto request validation
    telemetry *deps.Telemetry           // logging + metrics
    service   Service                   // interface to service layer
    userv1.UnimplementedUserServiceServer  // forward compat embedding
}

func New(telemetry *deps.Telemetry, service Service) (*Server, error) {
    validator, err := protovalidate.New()
    // ...
}
```

## Handler Registration (internal/<domain>/handler.go)

```go
// gRPC
func (s *Server) GRPCHandler(server *grpc.Server) error {
    userv1.RegisterUserServiceServer(server, s)
    return nil
}

// HTTP via grpc-gateway
func (s *Server) HTTPHandler(ctx context.Context) func(mux *runtime.ServeMux, address string) error {
    return func(mux *runtime.ServeMux, address string) error {
        return userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, address,
            []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
    }
}
```

## Ports

| Port | Purpose |
|------|---------|
| 8080 | gRPC |
| 8081 | HTTP gateway (REST) |
| 8082 | Internal metrics (Prometheus) |
| 2345 | Delve debugger (devstack only) |

## Key Constraint

`Server.New()` returns an error if `protovalidate.New()` fails — always check it. The service layer (`Service` interface) is injected — the server never imports the concrete service struct, enabling mock injection in tests.
