# Testing Patterns

## Test File Convention

Tests are co-located with the code they test, in a separate `_test` package:
```
internal/user/model/model.go        → source
internal/user/model/model_test.go   → tests (package model_test)
```

Using `package model_test` (not `package model`) enforces testing through the public API.

## Basic Test Structure

```go
package model_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/razorpay/go-foundation-v2/internal/user/model"
)

func TestNew(t *testing.T) {
    user := model.New()
    assert.NotNil(t, user)
    assert.IsType(t, &model.User{}, user)
}
```

## What to Test in Models

- `New()` — returns non-nil pointer of correct type
- `TableName()`, `EntityName()`, `EntityPrefix()` — return expected constant values
- `FromPublicID()` — test both 14-char internal format AND `"user_<14chars>"` public format; test invalid inputs
- `ToProto()` — verify all fields map correctly from model to proto
- `FromProto()` — verify all fields map correctly from proto to model

## Repository Interface Enables Mock Injection

The `Repo` interface in `internal/user/service/repo.go` is specifically for testability:
```go
type Repo interface {
    Create(ctx context.Context, user *model.User) goutilsError.IError
    Get(ctx context.Context, id string) (*model.User, goutilsError.IError)
}
```

In service tests, inject a mock implementing this interface — no real DB needed.

Similarly, `internal/user/server.go` defines the `Service` interface for injecting mock services in handler tests.

## Running Tests

```bash
make test ENV=local   # runs all tests with local Go toolchain
go test ./...         # equivalent direct command
go test ./internal/user/model/... -v  # specific package, verbose
```

## Dependencies

```go
github.com/stretchr/testify v1.11.1  // assert, require, mock packages
github.com/stretchr/objx v0.5.2      // used by testify/mock
```

## Test Naming Convention

Use `Test<FunctionName>[_<scenario>]` format:
- `TestNew` — basic constructor test
- `TestFromPublicID_InternalFormat` — specific scenario
- `TestFromPublicID_PublicFormat` — another scenario
- `TestFromPublicID_Invalid` — error case
