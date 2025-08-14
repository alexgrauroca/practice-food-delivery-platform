# \CustomersAPI

All URIs are relative to *http://localhost:80*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCustomers**](CustomersAPI.md#GetCustomers) | **Get** /v1.0/customers | Get the list of customers
[**RegisterCustomer**](CustomersAPI.md#RegisterCustomer) | **Post** /v1.0/customers | Register a new customer



## GetCustomers

> GetCustomersResponse GetCustomers(ctx).Page(page).PageSize(pageSize).Sort(sort).Execute()

Get the list of customers



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	page := int32(1) // int32 | Page number for pagination (optional) (default to 1)
	pageSize := int32(20) // int32 | Number of items per page (optional) (default to 10)
	sort := "name,-email" // string | Sort fields and directions, comma-separated. Prefix field with '-' for descending order. Multiple fields can be specified (e.g., 'name,-email'). Available sort fields:   - name: Customer's full name   - email: Customer's email address   - created-at: Account creation date  (optional) (default to "name")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CustomersAPI.GetCustomers(context.Background()).Page(page).PageSize(pageSize).Sort(sort).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CustomersAPI.GetCustomers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCustomers`: GetCustomersResponse
	fmt.Fprintf(os.Stdout, "Response from `CustomersAPI.GetCustomers`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetCustomersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** | Page number for pagination | [default to 1]
 **pageSize** | **int32** | Number of items per page | [default to 10]
 **sort** | **string** | Sort fields and directions, comma-separated. Prefix field with &#39;-&#39; for descending order. Multiple fields can be specified (e.g., &#39;name,-email&#39;). Available sort fields:   - name: Customer&#39;s full name   - email: Customer&#39;s email address   - created-at: Account creation date  | [default to &quot;name&quot;]

### Return type

[**GetCustomersResponse**](GetCustomersResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RegisterCustomer

> RegisterCustomerResponse RegisterCustomer(ctx).RegisterCustomerRequest(registerCustomerRequest).Execute()

Register a new customer



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	registerCustomerRequest := *openapiclient.NewRegisterCustomerRequest("user@example.com", "strongpassword123", "John Doe", "123 Main St", "New York", "10001", "US") // RegisterCustomerRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CustomersAPI.RegisterCustomer(context.Background()).RegisterCustomerRequest(registerCustomerRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CustomersAPI.RegisterCustomer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RegisterCustomer`: RegisterCustomerResponse
	fmt.Fprintf(os.Stdout, "Response from `CustomersAPI.RegisterCustomer`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRegisterCustomerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registerCustomerRequest** | [**RegisterCustomerRequest**](RegisterCustomerRequest.md) |  | 

### Return type

[**RegisterCustomerResponse**](RegisterCustomerResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

