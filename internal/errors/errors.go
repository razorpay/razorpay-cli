package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	goutilserrors "github.com/razorpay/goutils/errors"

	errrpc "github.com/razorpay/go-foundation-v2/rpc/go_foundation_v2/error/v1"
)

var (
	// ErrValidationFailure ...
	ErrValidationFailure = goutilserrors.NewClass(
		"validation_failure",
		goutilserrors.SeverityLow,
		goutilserrors.ValidationFailure,
	)

	// ErrInternalServerError ...
	ErrInternalServerError = goutilserrors.NewClass(
		"internal_error", goutilserrors.SeverityHigh, goutilserrors.Recoverable)

	// ErrBadRequestError ...
	ErrBadRequestError = goutilserrors.NewClass(
		"bad_request", goutilserrors.SeverityLow, goutilserrors.BadRequest)

	// ErrUnauthenticated ...
	ErrUnauthenticated = goutilserrors.NewClass(
		"authentication_failure",
		goutilserrors.SeverityLow,
		goutilserrors.Recoverable,
	)

	// ErrorInvalidData ...
	ErrorInvalidData = goutilserrors.NewClass(
		"invalid_data", goutilserrors.SeverityLow, goutilserrors.BadRequest)

	// ErrNotFound ...
	ErrNotFound = goutilserrors.NewClass(
		"not_found", goutilserrors.SeverityLow, goutilserrors.BadRequest)

	// ErrEntityAlreadyExists ...
	ErrEntityAlreadyExists = goutilserrors.NewClass(
		"already_exists", goutilserrors.SeverityLow, goutilserrors.BadRequest)

	// ErrServerUnavailable ...
	ErrServerUnavailable = goutilserrors.NewClass(
		"server_unavailable", goutilserrors.SeverityHigh, goutilserrors.Base)

	// ErrUnexpected is used when something that should never have occurred
	// theoretically happens like an invalid value of a constant.
	ErrUnexpected = goutilserrors.NewClass(
		"ErrUnexpected",
		goutilserrors.SeverityCritical, goutilserrors.Base,
	)

	// ErrContextCancelled represents an error class for context cancellation.
	// It is used to indicate that an operation was canceled, typically due to
	// a context timeout or cancellation.
	ErrContextCancelled = goutilserrors.NewClass(
		"CONTEXT_CANCELLED", goutilserrors.SeverityHigh, goutilserrors.Base)

	classToCodeMap = map[goutilserrors.IClass]codes.Code{
		ErrValidationFailure:   codes.InvalidArgument,
		ErrBadRequestError:     codes.InvalidArgument,
		ErrorInvalidData:       codes.InvalidArgument,
		ErrNotFound:            codes.NotFound,
		ErrInternalServerError: codes.Internal,
		ErrUnexpected:          codes.Internal,
		ErrServerUnavailable:   codes.Unavailable,
		ErrUnauthenticated:     codes.Unauthenticated,
		ErrContextCancelled:    codes.Canceled,
		ErrEntityAlreadyExists: codes.AlreadyExists,
	}
)

// WrapError wraps the goutils error package error into
// grpc compatible error
func WrapError(e goutilserrors.IError) error {
	s, _ := status.New(getGrpcCodeFromErrorClass(
		e.Class(),
	), e.Internal().Code().String()).WithDetails(ToProtoPublicError(e))
	return s.Err()
}

// ToProtoPublicError returns the error proto message
func ToProtoPublicError(err goutilserrors.IError) *errrpc.Error {
	return &errrpc.Error{
		Code:        err.Public().GetCode(),
		Field:       err.Public().GetField(),
		Source:      err.Public().GetSource(),
		Step:        err.Public().GetStep(),
		Reason:      err.Public().GetReason(),
		Metadata:    err.Public().GetMetadata(),
		Action:      err.Public().GetAction(),
		Description: err.Public().Error(),
	}
}

func getGrpcCodeFromErrorClass(errorClass goutilserrors.IClass) codes.Code {
	if code, exists := classToCodeMap[errorClass]; exists {
		return code
	}
	return codes.Internal
}
