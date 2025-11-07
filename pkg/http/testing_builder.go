package http

// HttpResponseBuilder is a type used to build and return HTTP response payloads in JSON format.
type HttpResponseBuilder struct {
	json string
}

// Build returns the current JSON string.
func (b HttpResponseBuilder) Build() string {
	return b.json
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
