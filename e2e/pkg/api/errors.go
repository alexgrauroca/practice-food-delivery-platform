package api

import (
	"encoding/json"
	"errors"
)

// APIError represents an error response from an API, containing an error code, message, and optional details.
type APIError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

// ParseErrorResponse parses the error response body into an APIError object and validates required fields.
// Returns an error if unmarshalling fails or required fields are missing.
func ParseErrorResponse(body []byte) (*APIError, error) {
	var apiError APIError
	if err := json.Unmarshal(body, &apiError); err != nil {
		return nil, err
	}
	// Validate that we got at least the required fields according to the OpenAPI spec
	if apiError.Code == "" || apiError.Message == "" {
		return nil, errors.New("invalid errors response: missing required fields")
	}
	return &apiError, nil
}
