package http

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

const (
	// CodeValidationError represents the error code for validation failures during input processing or validation checks.
	CodeValidationError = "VALIDATION_ERROR"
	// MsgValidationError represents the error message for validation failures during input validation checks.
	MsgValidationError = "validation failed"

	// CodeInvalidRequest represents the error code for an invalid or improper request made to the system.
	CodeInvalidRequest = "INVALID_REQUEST"
	// MsgInvalidRequest represents the error message for an invalid or improperly formed request.
	MsgInvalidRequest = "invalid request"

	// CodeInternalError represents the error code for an unspecified internal server error encountered in the system.
	CodeInternalError = "INTERNAL_ERROR"
	// MsgInternalError represents the error message returned when the system fails to log in a customer.
	MsgInternalError = "an unexpected error occurred"
)

// ErrorResponse represents a standardized structure for API error responses containing code, message, and optional details.
type ErrorResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

// NewErrorResponse creates and returns a new ErrorResponse with the provided code and message,
// and an empty details slice.
func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Details: make([]string, 0),
	}
}

// GetErrorResponseFromValidationErr gets the ErrorResponse based on the error type returned from the validation
func GetErrorResponseFromValidationErr(err error) ErrorResponse {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errResp := NewErrorResponse(CodeValidationError, MsgValidationError)
		details := make([]string, 0)

		for _, fe := range ve {
			details = append(details, getValidationErrorDetail(fe))
		}
		errResp.Details = details

		return errResp
	}
	return NewErrorResponse(CodeInvalidRequest, MsgInvalidRequest)
}

// getValidationErrorDetail returns a detailed error message based on the field error
func getValidationErrorDetail(fe validator.FieldError) string {
	field := strcase.ToSnake(fe.Field())
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		if field == "password" {
			return field + " must be a valid password with at least 8 characters long"
		}
		return field + " must be at least " + fe.Param() + " characters long"
	case "max":
		return field + " must not exceed " + fe.Param() + " characters long"
	default:
		return field + " is invalid"
	}
}
