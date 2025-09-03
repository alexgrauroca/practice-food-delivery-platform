// Package api provides a set of tools and utilities for making HTTP requests to the authentication service API.
// It includes endpoints and methods specifically designed for customer-related operations such as registration and login.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

const (
	// BaseURL defines the default base URL for the authentication service API endpoints.
	// It is used as a prefix for all API requests in the testing environment.
	BaseURL = "http://localhost:80"

	contentTypeJSON = "application/json"
)

// DoPost sends a POST request to the specified endpoint with the given payload,
// decoding the response into the generic type T.
// It returns a pointer to the decoded response of type T or an error if the request or decoding fails.
func DoPost[P, R any](endpoint string, params P, config *RequestConfig) (*R, error) {
	url := buildURL(endpoint, params)

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Content-Type", contentTypeJSON)
	if config != nil && config.BearerToken != nil {
		req.Header.Add("Authorization", "Bearer "+*config.BearerToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
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
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil

}

// DoGet sends a GET request to the specified endpoint with the given parameters,
// decoding the response into the generic type T.
// It returns a pointer to the decoded response of type T or an error if the request or decoding fails.
func DoGet[P, R any](endpoint string, params P, config *RequestConfig) (*R, error) {
	url := buildURL(endpoint, params)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if config != nil && config.BearerToken != nil {
		req.Header.Add("Authorization", "Bearer "+*config.BearerToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if apiError, err := ParseErrorResponse(responseBody); err == nil {
			return nil, apiError
		}
		return nil, &ErrorResponse{
			Code:    "UNEXPECTED_ERROR",
			Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
			Details: []string{string(responseBody)},
		}
	}

	var result R
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil

}

func buildURL[P any](endpoint string, params P) string {
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	queryParams := make([]string, 0)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if pathTag := fieldType.Tag.Get("path"); pathTag != "" {
			placeholder := fmt.Sprintf(":%s", pathTag)
			endpoint = strings.Replace(endpoint, placeholder, toString(field), 1)
			continue
		}

		if queryTag := fieldType.Tag.Get("query"); queryTag != "" {
			if field.Kind() == reflect.Ptr && field.IsNil() {
				continue
			}

			value := field
			if field.Kind() == reflect.Ptr {
				value = field.Elem()
			}

			queryParams = append(queryParams, fmt.Sprintf("%s=%s", queryTag, toString(value)))
		}
	}

	if len(queryParams) > 0 {
		endpoint += "?" + strings.Join(queryParams, "&")
	}

	return endpoint
}

func toString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float())
	case reflect.Bool:
		return fmt.Sprintf("%v", v.Bool())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
