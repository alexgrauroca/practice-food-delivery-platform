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
	baseURL         = "http://localhost:80"
	apiV1           = baseURL + "/v1.0"
	contentTypeJSON = "application/json"
)

var (

	// RegisterEndpoint defines the API endpoint URL for customer registration under version 1.0 of the API.
	RegisterEndpoint = apiV1 + "/customers/register"

	// LoginEndpoint defines the API endpoint URL for customer login under version 1.0 of the API.
	LoginEndpoint = apiV1 + "/customers/login"
)

// DoPost sends a POST request to the specified endpoint with the given payload, decoding the response into the generic type T.
// It returns a pointer to the decoded response of type T or an error if the request or decoding fails.
func DoPost[T any](endpoint string, payload any) (*T, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	resp, err := http.Post(endpoint, contentTypeJSON, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}
