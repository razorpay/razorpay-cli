# Example User Service

This template includes a complete **User Service** that demonstrates the recommended architecture and patterns for building microservices.

## Architecture Overview

```
cmd/user/               # Service entry point
├── main.go             # Application bootstrap

internal/user/          # Business logic
├── server.go           # gRPC server implementation
├── handler.go          # HTTP/gRPC handler registration
├── service/            # Business logic layer
├── repo/               # Data access layer
├── model/              # Domain models
└── migrations/         # Database migrations

config/user/            # Configuration files
├── default.toml        # Default configuration
├── dev.toml            # Development overrides
└── route_config.toml   # API routing configuration

rpc/go_foundation_v2/   # Generated proto code
└── user/v1/            # User service proto definitions
```

## Key Components

### 1. Entry Point (`cmd/user/main.go`)
- Initializes the Foundation server with configuration
- Sets up database connections and health checks
- Wires together all service dependencies
- Registers gRPC and HTTP handlers

### 2. Server Layer (`internal/user/server.go`)
- Implements the gRPC service interface
- Handles request validation using proto validation
- Manages error handling and response formatting
- Delegates business logic to the service layer

### 3. Service Layer (`internal/user/service/`)
- Contains core business logic
- Implements domain rules and validations
- Orchestrates data access operations
- Returns domain models

### 4. Repository Layer (`internal/user/repo/`)
- Handles data persistence operations
- Abstracts database implementation details
- Provides clean interfaces for data access
- Manages database transactions

### 5. Model Layer (`internal/user/model/`)
- Defines domain entities and value objects
- Contains business logic specific to entities
- Provides conversion methods (e.g., `ToProto()`)
- Implements validation rules

## Request Flow

```
HTTP/gRPC Request
       ↓
   Handler Layer (server.go)
       ↓ [validation]
   Service Layer (service/)
       ↓ [business logic]
   Repository Layer (repo/)
       ↓ [data access]
   Database
```

1. Handler receives and validates the request
2. Service layer applies business rules
3. Repository layer handles database interactions
4. Results flow back up through the layers
5. Errors are properly formatted and returned

## Code Patterns

### Dependency Injection
```go
// Clean dependency injection in main.go
userServer, err := user.New(
    tel,                    // Telemetry
    userservice.New(        // Service layer
        tel,
        repo.New(tel, primaryDB), // Repository layer
    ),
)
```

### Interface-Based Design
```go
// Service interface defines business operations
type Service interface {
    CreateUser(ctx context.Context, req *userv1.CreateRequest) (*model.User, goutilsError.IError)
    GetUser(ctx context.Context, req *userv1.GetRequest) (*model.User, goutilsError.IError)
}
```

### Error Handling
```go
// Structured error handling with public/private error separation
if err != nil {
    return &userv1.User{
        Error: errpkg.ToProtoPublicError(err),
    }, errpkg.WrapError(err)
}
```

## Running the Example

```bash
# Run the example service
make build ENV=local
./bin/user --config ./config/user/dev.toml

# Or using Docker
docker-compose up
```

Examine the API endpoints, database migrations, and configuration files to understand the complete flow.

