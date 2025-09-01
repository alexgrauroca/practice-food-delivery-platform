package customer

import "github.com/alexgrauroca/practice-food-delivery-platform/e2e/pkg/api"

var (
	// RegisterEndpoint defines the API endpoint URL for customer registration under version 1.0 of the API.
	RegisterEndpoint = api.BaseURL + "/v1.0/customers"
	// LoginEndpoint defines the API endpoint URL for customer login under version 1.0 of the API.
	LoginEndpoint = api.BaseURL + "/v1.0/customers/login"
	// RefreshEndpoint defines the API endpoint URL for refreshing customer data under version 1.0 of the API.
	RefreshEndpoint = api.BaseURL + "/v1.0/customers/refresh"
	// GetCustomerEndpoint defines the API endpoint URL for retrieving customer data under version 1.0 of the API.
	GetCustomerEndpoint = api.BaseURL + "/v1.0/customers/:customerID"
)
