# Error Handling System

## Error Class Declaration (internal/errors/errors.go)

```go
import goutilserrors "github.com/razorpay/goutils/errors"

var (
    ErrValidationFailure   = goutilserrors.NewClass("validation_failure",   SeverityLow,      ValidationFailure)
    ErrInternalServerError = goutilserrors.NewClass("internal_error",       SeverityHigh,     Recoverable)
    ErrBadRequestError     = goutilserrors.NewClass("bad_request",          SeverityLow,      BadRequest)
    ErrUnauthenticated     = goutilserrors.NewClass("authentication_failure",SeverityLow,     Recoverable)
    ErrorInvalidData       = goutilserrors.NewClass("invalid_data",         SeverityLow,      BadRequest)
    ErrNotFound            = goutilserrors.NewClass("not_found",            SeverityLow,      BadRequest)
    ErrEntityAlreadyExists = goutilserrors.NewClass("already_exists",       SeverityLow,      BadRequest)
    ErrServerUnavailable   = goutilserrors.NewClass("server_unavailable",   SeverityHigh,     Base)
    ErrUnexpected          = goutilserrors.NewClass("ErrUnexpected",        SeverityCritical, Base)
    ErrContextCancelled    = goutilserrors.NewClass("CONTEXT_CANCELLED",    SeverityHigh,     Base)
)
```

## Error Codes (internal/errors/codes.go)

Naming pattern: `GOFND` prefix + 7-digit zero-padded number. When creating a new service, replace `GOFND` with your service prefix.

```go
CodeRequestValidationFailure errors.IdentifierCode = "GOFND0000001"
CodeDuplicateRequest         errors.IdentifierCode = "GOFND0000002"
CodeNoRecordFound            errors.IdentifierCode = "GOFND0000003"
CodeInternalServerError      errors.IdentifierCode = "GOFND0000004"
CodeBadRequestError          errors.IdentifierCode = "GOFND0000005"
```

## Error Class → gRPC Status Code Mapping

| Error Class | gRPC Code |
|-------------|-----------|
| ErrValidationFailure, ErrBadRequestError, ErrorInvalidData | InvalidArgument |
| ErrNotFound | NotFound |
| ErrInternalServerError, ErrUnexpected | Internal |
| ErrServerUnavailable | Unavailable |
| ErrUnauthenticated | Unauthenticated |
| ErrContextCancelled | Canceled |
| ErrEntityAlreadyExists | AlreadyExists |

## Creating Errors — Service Layer

```go
// UniqueConstraintViolation from DB
if errors.Is(err, spine.UniqueConstraintViolation) {
    return nil, errpkg.ErrBadRequestError.New(errpkg.CodeBadRequestError).Wrap(err)
}
// RecordNotFound from DB
if errors.Is(err, spine.RecordNotFound) {
    return nil, errpkg.ErrNotFound.New(errpkg.CodeNoRecordFound).Wrap(err)
}
// Generic internal error
return nil, errpkg.ErrInternalServerError.New(errpkg.CodeInternalServerError).Wrap(err)
```

## Creating Errors — Server Layer (Validation)

```go
if err := s.validator.Validate(req); err != nil {
    ierr := errpkg.ErrValidationFailure.
        New(errpkg.CodeRequestValidationFailure).
        Wrap(err).
        WithPublicMetadata(foundationErr.FromProtoValidationError(err))  // field-level details
    return &userv1.User{
        Error: errpkg.ToProtoPublicError(ierr),  // embed in response proto
    }, errpkg.WrapError(ierr)  // also return as gRPC status error
}
```

## Dual Return Pattern (ALWAYS do both)

```go
// For any error in a gRPC handler:
return &EntityProto{
    Error: errpkg.ToProtoPublicError(err),  // embeds error in response body
}, errpkg.WrapError(err)                    // sets gRPC status code
```

**Why both?** gRPC clients can decode the status error. HTTP clients (via grpc-gateway) receive the embedded error in the JSON response body.

## WrapError and ToProtoPublicError

```go
// WrapError: creates gRPC status error with embedded proto details
func WrapError(e IError) error {
    s, _ := status.New(grpcCode, e.Internal().Code().String()).WithDetails(ToProtoPublicError(e))
    return s.Err()
}

// ToProtoPublicError: converts IError to the error.proto Error message
func ToProtoPublicError(err IError) *errrpc.Error {
    return &errrpc.Error{
        Code: err.Public().GetCode(), Field: ..., Source: ..., Step: ...,
        Reason: ..., Metadata: ..., Action: ..., Description: err.Public().Error(),
    }
}
```
