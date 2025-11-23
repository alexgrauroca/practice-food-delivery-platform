package staff

import "github.com/alexgrauroca/practice-food-delivery-platform/e2e/pkg/api"

var (
	// LoginEndpoint defines the API endpoint URL for staff login under version 1.0 of the API.
	LoginEndpoint = api.BaseURL + "/v1.0/staff/login"
	// RefreshEndpoint defines the API endpoint URL for refreshing staff data under version 1.0 of the API.
	RefreshEndpoint = api.BaseURL + "/v1.0/staff/refresh"
)
