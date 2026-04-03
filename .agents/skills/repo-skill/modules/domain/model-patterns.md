# Model Patterns

## Model Struct

```go
type User struct {
    spine.SoftDeletableModel  // provides: ID string, CreatedAt/UpdatedAt/DeletedAt int64
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
}
```

Always use `spine.SoftDeletableModel` for soft-delete support. The `ID` field is a string (not int64/UUID), managed by spine with a prefix system.

## Required Model Methods

Every model MUST implement these methods for spine and Foundation to work correctly:

```go
const (
    tableName    = "users"   // PostgreSQL table name
    entityName   = "users"   // used for logging/identification
    entityPrefix = "user_"   // prefix for public IDs
)

func (*User) TableName() string    { return tableName }    // used by GORM
func (*User) EntityName() string   { return entityName }   // used for observability
func (*User) EntityPrefix() string { return entityPrefix } // used for ID generation

func (u *User) SetDefaults() errors.IError { return nil }  // hook for defaults
```

## Constructor Pattern

```go
func New() *User { return &User{} }
// Always use New(), never &User{} directly in application code
```

## Public ID Format

The template uses a **dual-format ID** for the `GetRequest`:
- **Internal**: 14-char string (e.g., `"abc123xyz01234"`) — stored in DB
- **Public**: `"user_" + internal` (e.g., `"user_abc123xyz01234"`) — exposed in API

```go
func (u *User) FromPublicID(id string) error {
    if len(id) == 14 {
        u.ID = id  // already internal format
        return nil
    }
    parts := strings.Split(id, entityPrefix)  // split on "user_"
    if len(parts) != 2 { return errors.New("invalid public ID") }
    u.ID = parts[1]
    return nil
}
```

Proto validation enforces `min_len: 14, max_len: 19` to accept both formats.

## Proto Conversion Methods

```go
// ToProto: internal model → proto response (called in gRPC handler)
func (u *User) ToProto() *userv1.User {
    return &userv1.User{
        Id:       u.ID,
        Username: u.Username,
        Password: u.Password,
        Email:    u.Email,
    }
}

// FromProto: proto request → internal model (called in service layer)
func (u *User) FromProto(pb *userv1.CreateRequest) {
    u.Username = pb.GetUsername()
    u.Password = pb.GetPassword()
    u.Email = pb.GetEmail()
}
```

**Pattern:** `FromProto` takes the specific request type (not the response type). For create operations pass `*CreateRequest`; for updates pass `*UpdateRequest`. The model's `ID` is NOT set from `CreateRequest` — spine assigns it on `Create()`.

## Adding a New Domain Entity

1. Create `internal/<domain>/model/model.go` with `spine.SoftDeletableModel` embedding
2. Implement `TableName()`, `EntityName()`, `EntityPrefix()`, `SetDefaults()`, `New()`
3. Implement `ToProto()` returning your response proto message
4. Implement `FromProto()` for each mutating request type
5. Create the matching SQL migration in `migrations/`
