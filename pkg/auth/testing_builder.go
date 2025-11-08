package auth

import (
	"encoding/json"
	"strings"
)

// AuthResponseBuilder is a type used to build and return HTTP response payloads in JSON format.
type AuthResponseBuilder struct {
	json    string
	details []string
}

// WithDetails sets the details of the response.
func (b AuthResponseBuilder) WithDetails(details ...string) AuthResponseBuilder {
	b.details = details
	return b
}

// Build returns the JSON string.
func (b AuthResponseBuilder) Build() string {
	if len(b.details) == 0 {
		return b.json
	}

	details, err := json.Marshal(b.details)
	if err != nil {
		return b.json
	}
	return strings.ReplaceAll(b.json, `"details": []`, `"details": `+string(details))
}

// NewUnauthorizedRespBuilder creates a builder for the UNAUTHORIZED JSON used in tests.
func NewUnauthorizedRespBuilder() AuthResponseBuilder {
	return AuthResponseBuilder{
		json: `{
			"code": "UNAUTHORIZED",
			"message": "Authentication is required to access this resource",
			"details": []
		}`,
	}
}

// NewForbiddenRespBuilder creates a builder for the FORBIDDEN JSON used in tests.
func NewForbiddenRespBuilder() AuthResponseBuilder {
	return AuthResponseBuilder{
		json: `{
			"code": "FORBIDDEN",
			"message": "You do not have permission to access this resource",
			"details": []
		}`,
	}
}
