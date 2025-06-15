// Package errors provides custom errors types and handling functionality for the authentication service.
// It defines specific errors cases and standardizes errors responses across the customer-related operations.
package errors

import (
	"encoding/json"
	"errors"
)

// APIError represents an errors response from an API, containing an errors code, message, and optional details.
type APIError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

// ParseErrorResponse parses the errors response body into an APIError object and validates required fields.
// Returns an errors if unmarshalling fails or required fields are missing.
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
