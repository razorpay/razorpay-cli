# Proto & gRPC Workflow

## Directory Layout

```
proto/                          # COMMITTED to git — source of truth
  go_foundation_v2/
    user/v1/user.proto
    error/v1/error.proto
rpc/                            # NOT committed (.gitignore) — auto-generated
  go_foundation_v2/
    user/v1/                    # generated Go structs + gRPC server/client
buf.yaml                        # buf module config + lint/breaking rules
buf.gen.yaml                    # code generation plugins
buf.lock                        # pinned dependency versions
```

**Critical rule:** `proto/` is committed for reproducibility. `rpc/` is in `.gitignore` and regenerated automatically during Docker builds.

## Proto File Conventions

```proto
syntax = "proto3";
package go_foundation_v2.users.v1;
option go_package = "go_foundation_v2/users/v1;usersv1";

// Required imports for this template
import "google/api/annotations.proto";  // HTTP mapping
import "buf/validate/validate.proto";   // request validation
import "go_foundation_v2/error/v1/error.proto";  // embedded error

// Response messages embed an error field
message User {
  string id = 1; string username = 2; ...
  go_foundation_v2.error.v1.Error error = 5;  // ALWAYS embed this in responses
}

// Request validation via buf.validate annotations
message CreateRequest {
  string username = 1 [
    (buf.validate.field).string = {min_len: 4, max_len: 32},
    (buf.validate.field).required = true
  ];
  optional string email = 3 [(buf.validate.field).string.email = true];
}

// HTTP mapping on each RPC
service UserService {
  rpc Create(CreateRequest) returns (User) {
    option (google.api.http) = { post: "/v1/users" body: "*" };
  }
  rpc Get(GetRequest) returns (User) {
    option (google.api.http) = { get: "/v1/users/{id}" };
  }
}
```

## Make Targets

| Target | When to use |
|--------|------------|
| `make proto-fetch` | Pull updated `.proto` source files from `github.com/razorpay/proto` into `proto/` |
| `make proto-generate` | Run buf to generate Go code from `proto/` → `rpc/` |
| `make proto-lint` | Lint proto files (style + breaking change check) |
| `make proto-refresh` | `proto-fetch` + `proto-generate` in one step |

**You only need proto-fetch when updating to a newer proto version. For normal builds, Docker runs proto-generate automatically.**

## Makefile Proto Variables

```makefile
PROTO_MODULES := go_foundation_v2       # directory names from the proto repo
PROTO_GIT_URL := https://github.com/razorpay/proto.git
PROTO_BRANCH  := master
PROTO_ROOT    := proto                  # where .proto files live
RPC_ROOT      := rpc                   # where generated code goes
```

## Generated Go Import Path

```go
import userv1 "github.com/razorpay/go-foundation-v2/rpc/go_foundation_v2/user/v1"
// Note: this import fails locally until you run `make proto-generate`
// CI/Docker builds do this automatically
```

## buf.gen.yaml Plugins

The template uses 4 buf plugins:
1. `protoc-gen-go` — Go structs
2. `protoc-gen-go-grpc` — gRPC server/client interfaces
3. `protoc-gen-grpc-gateway` — HTTP gateway registration
4. `protoc-gen-openapiv2` — OpenAPI v2 spec (optional)

## Why Proto Files Are Committed

Proto source files are committed so:
1. The build is reproducible without network access to the proto repo
2. Git history tracks API contract changes
3. Breaking change detection works in CI (`buf breaking`)
