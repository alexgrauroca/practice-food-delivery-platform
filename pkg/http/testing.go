package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// TestHandler defines an interface for registering routes within a Gin web engine.
type TestHandler interface {
	RegisterRoutes(*gin.Engine)
}

// SetupTestEnv initializes the test environment with default values common to all tests.
func SetupTestEnv() log.Logger {
	// Setting up the default values
	gin.SetMode(gin.TestMode)
	logger, _ := log.NewTest()
	return logger
}

// ServeTestHTTPRequest executes a test case for the given handler.
func ServeTestHTTPRequest(
	t *testing.T,
	h TestHandler,
	httpMethod string,
	route string,
	token string,
	queryParams map[string]string,
	jsonPayload string,
) *httptest.ResponseRecorder {
	// Initialize the Gin router and register the routes
	router := gin.New()
	h.RegisterRoutes(router)

	// Create a new HTTP request with the test case's params
	var req *http.Request

	switch httpMethod {
	case http.MethodGet:
		baseURL, err := url.Parse(route)
		if err != nil {
			t.Fatalf("failed to parse route: %v", err)
		}

		// Add query parameters
		if len(queryParams) > 0 {
			q := baseURL.Query()
			for k, v := range queryParams {
				q.Add(k, v)
			}
			baseURL.RawQuery = q.Encode()
		}

		req = httptest.NewRequest(http.MethodGet, baseURL.String(), nil)

	case http.MethodPost, http.MethodPut:
		req = httptest.NewRequest(httpMethod, route, strings.NewReader(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

	default:
		t.Fatalf("unsupported HTTP method: %s", httpMethod)
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	w := httptest.NewRecorder()

	// Make the request to the handler
	router.ServeHTTP(w, req)

	return w
}
