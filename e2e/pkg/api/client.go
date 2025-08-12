// Package api provides a set of tools and utilities for making HTTP requests to the authentication service API.
// It includes endpoints and methods specifically designed for customer-related operations such as registration and login.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// BaseURL defines the default base URL for the authentication service API endpoints.
	// It is used as a prefix for all API requests in the testing environment.
	BaseURL = "http://localhost:80"

	contentTypeJSON = "application/json"
)

// DoPost sends a POST request to the specified endpoint with the given payload, decoding the response into the generic type T.
// It returns a pointer to the decoded response of type T or an error if the request or decoding fails.
func DoPost[P, R any](endpoint string, payload P) (*R, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("errors marshaling payload: %w", err)
	}

	resp, err := http.Post(endpoint, contentTypeJSON, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("errors making POST request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("errors reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// Try to parse as API errors first
		if apiError, err := ParseErrorResponse(responseBody); err == nil {
			return nil, apiError
		}
		// Fallback to generic errors
		return nil, &ErrorResponse{
			Code:    "UNEXPECTED_ERROR",
			Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
			Details: []string{string(responseBody)},
		}
	}

	var result R
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("errors decoding response: %w", err)
	}

	return &result, nil

}
