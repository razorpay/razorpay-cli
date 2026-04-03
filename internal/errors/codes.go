package errors

import "github.com/razorpay/goutils/errors"

var (
	// TODO: change all these codes based on the mapping

	CodeRequestValidationFailure errors.IdentifierCode = "GOFND0000001"
	CodeDuplicateRequest         errors.IdentifierCode = "GOFND0000002"
	CodeNoRecordFound            errors.IdentifierCode = "GOFND0000003"
	CodeInternalServerError      errors.IdentifierCode = "GOFND0000004"
	CodeBadRequestError          errors.IdentifierCode = "GOFND0000005"
)
