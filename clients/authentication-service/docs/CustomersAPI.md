# \CustomersAPI

All URIs are relative to *http://localhost:80*

Method | HTTP request | Description
------------- | ------------- | -------------
[**LoginCustomer**](CustomersAPI.md#LoginCustomer) | **Post** /v1.0/customers/login | Login as a customer
[**RefreshCustomer**](CustomersAPI.md#RefreshCustomer) | **Post** /v1.0/customers/refresh | Refresh access token
[**RegisterCustomer**](CustomersAPI.md#RegisterCustomer) | **Post** /v1.0/customers/register | Register a new customer



## LoginCustomer

> LoginCustomerResponse LoginCustomer(ctx).LoginCustomerRequest(loginCustomerRequest).Execute()

Login as a customer



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/alexgrauroca/practice-food-delivery-platform/authclient"
)

func main() {
	loginCustomerRequest := *openapiclient.NewLoginCustomerRequest("user@example.com", "strongpassword123") // LoginCustomerRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CustomersAPI.LoginCustomer(context.Background()).LoginCustomerRequest(loginCustomerRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CustomersAPI.LoginCustomer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `LoginCustomer`: LoginCustomerResponse
	fmt.Fprintf(os.Stdout, "Response from `CustomersAPI.LoginCustomer`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLoginCustomerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginCustomerRequest** | [**LoginCustomerRequest**](LoginCustomerRequest.md) |  | 

### Return type

[**LoginCustomerResponse**](LoginCustomerResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RefreshCustomer

> RefreshCustomerResponse RefreshCustomer(ctx).RefreshCustomerRequest(refreshCustomerRequest).Execute()

Refresh access token



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/alexgrauroca/practice-food-delivery-platform/authclient"
)

func main() {
	refreshCustomerRequest := *openapiclient.NewRefreshCustomerRequest("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", "dGhpc2lzYXJlZnJlc2h0b2tlbg==") // RefreshCustomerRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CustomersAPI.RefreshCustomer(context.Background()).RefreshCustomerRequest(refreshCustomerRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CustomersAPI.RefreshCustomer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RefreshCustomer`: RefreshCustomerResponse
	fmt.Fprintf(os.Stdout, "Response from `CustomersAPI.RefreshCustomer`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRefreshCustomerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **refreshCustomerRequest** | [**RefreshCustomerRequest**](RefreshCustomerRequest.md) |  | 

### Return type

[**RefreshCustomerResponse**](RefreshCustomerResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
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
	openapiclient "github.com/alexgrauroca/practice-food-delivery-platform/authclient"
)

func main() {
	registerCustomerRequest := *openapiclient.NewRegisterCustomerRequest("507f1f77bcf86cd799439011", "user@example.com", "strongpassword123", "John Doe") // RegisterCustomerRequest | 

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

