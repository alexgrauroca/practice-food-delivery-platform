// Package api provides a set of tools and utilities for making HTTP requests to the authentication service API.
// It includes endpoints and methods specifically designed for customer-related operations such as registration and login.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apierrors "github.com/alexgrauroca/practice-food-delivery-platform/e2e/authentication-service/customers/api/errors"
)

const (
	baseURL         = "http://localhost:80"
	contentTypeJSON = "application/json"
)

var (
	// RegisterEndpoint defines the API endpoint URL for customer registration under version 1.0 of the API.
	RegisterEndpoint = baseURL + "/v1.0/customers/register"
	// LoginEndpoint defines the API endpoint URL for customer login under version 1.0 of the API.
	LoginEndpoint = baseURL + "/v1.0/customers/login"
	// RefreshEndpoint defines the API endpoint URL for refreshing customer data under version 1.0 of the API.
	RefreshEndpoint = baseURL + "/v1.0/customers/refresh"
)

// DoPost sends a POST request to the specified endpoint with the given payload, decoding the response into the generic type T.
// It returns a pointer to the decoded response of type T or an error if the request or decoding fails.
func DoPost[T any](endpoint string, payload any) (*T, error) {
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
		if apiError, err := apierrors.ParseErrorResponse(responseBody); err == nil {
			return nil, apiError
		}
		// Fallback to generic errors
		return nil, &apierrors.APIError{
			Code:    "UNEXPECTED_ERROR",
			Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
			Details: []string{string(responseBody)},
		}
	}

	var result T
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("errors decoding response: %w", err)
	}

	return &result, nil

}
