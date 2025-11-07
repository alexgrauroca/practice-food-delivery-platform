package restaurants_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	customhttp "github.com/alexgrauroca/practice-food-delivery-platform/pkg/http"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants"
)

var errUnexpected = errors.New("unexpected error")

type handlerTestCase struct {
	name        string
	token       string
	pathParams  map[string]string
	queryParams map[string]string
	jsonPayload string
	mocksSetup  func()
	wantJSON    string
	wantStatus  int
}

func TestHandler_RegisterRestaurant(t *testing.T) {
	logger := customhttp.SetupTestEnv()
	_ = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []handlerTestCase{
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			jsonPayload: `{"restaurant": 1.2, "staff_owner": true}`,
			wantJSON: `{
				"code": "INVALID_REQUEST",
				"message": "invalid request",
				"details": []
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "when empty payload is provided, then it should return a 400 with the validation error",
			jsonPayload: `{}`,
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"restaurant is required",
					"staff_owner is required"
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runHandlerTestCase(t, logger, http.MethodPost, "/v1.0/restaurants", tt, "")
		})
	}
}

func runHandlerTestCase(
	t *testing.T,
	logger log.Logger,
	httpMethod string,
	route string,
	tt handlerTestCase,
	token string,
) {
	if tt.mocksSetup != nil {
		tt.mocksSetup()
	}

	// Initialize the handler
	h := restaurants.NewHandler(logger)

	// Make HTTP request
	w := customhttp.ServeTestHTTPRequest(t, h, httpMethod, route, token, tt.queryParams, tt.jsonPayload)

	assert.Equal(t, tt.wantStatus, w.Code)
	assert.JSONEq(t, tt.wantJSON, w.Body.String())
}
