package http

import (
	"encoding/json"
	"strings"
)

// HttpResponseBuilder is a type used to build and return HTTP response payloads in JSON format.
type HttpResponseBuilder struct {
	json    string
	details []string
}

// WithDetails sets the details of the response.
func (b HttpResponseBuilder) WithDetails(details ...string) HttpResponseBuilder {
	b.details = details
	return b
}

// Build returns the JSON string.
func (b HttpResponseBuilder) Build() string {
	if len(b.details) == 0 {
		return b.json
	}

	details, err := json.Marshal(b.details)
	if err != nil {
		return b.json
	}
	return strings.ReplaceAll(b.json, `"details": []`, `"details": `+string(details))
}

// NewInternalErrorRespBuilder creates a builder for the INTERNAL_ERROR JSON used in tests.
func NewInternalErrorRespBuilder() HttpResponseBuilder {
	return HttpResponseBuilder{
		json: `{
			"code": "INTERNAL_ERROR",
			"message": "an unexpected error occurred",
			"details": []
		}`,
	}
}

// NewInvalidRequestRespBuilder creates a builder for the INVALID_REQUEST JSON used in tests.
func NewInvalidRequestRespBuilder() HttpResponseBuilder {
	return HttpResponseBuilder{
		json: `{
			"code": "INVALID_REQUEST",
			"message": "invalid request",
			"details": []
		}`,
	}
}

// NewValidationErrorRespBuilder creates a builder for the VALIDATION_ERROR JSON used in tests.
func NewValidationErrorRespBuilder() HttpResponseBuilder {
	return HttpResponseBuilder{
		json: `{
			"code": "VALIDATION_ERROR",
			"message": "validation failed",
			"details": []
		}`,
	}
}

// NewNotFoundRespBuilder creates a builder for the NOT_FOUND JSON used in tests.
func NewNotFoundRespBuilder() HttpResponseBuilder {
	return HttpResponseBuilder{
		json: `{
			"code": "NOT_FOUND",
			"message": "resource not found",
			"details": []
		}`,
	}
}
