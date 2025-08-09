# \CustomersAPI

All URIs are relative to *http://localhost:80*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RegisterCustomer**](CustomersAPI.md#RegisterCustomer) | **Post** /v1.0/customers | Register a new customer



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

