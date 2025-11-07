# \CustomersAPI

All URIs are relative to *http://localhost:80*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RegisterRestaurant**](CustomersAPI.md#RegisterRestaurant) | **Post** /v1.0/restaurants | Register a new restaurant



## RegisterRestaurant

> RegisterRestaurantResponse RegisterRestaurant(ctx).RegisterRestaurantRequest(registerRestaurantRequest).Execute()

Register a new restaurant



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
	registerRestaurantRequest := *openapiclient.NewRegisterRestaurantRequest(*openapiclient.NewRegisterRestaurantRequestRestaurant("GB123456789", "Acme Pizza", "Acme Pizza LLC", "America/New_York", *openapiclient.NewRestaurantContact("+1", "1234567890", "restaurant@example.com", "123 Main St", "New York", "10001", "US")), *openapiclient.NewRegisterRestaurantRequestStaffOwner("user@example.com", "strongpassword123", "John Doe", "123 Main St", "New York", "10001", "US")) // RegisterRestaurantRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CustomersAPI.RegisterRestaurant(context.Background()).RegisterRestaurantRequest(registerRestaurantRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CustomersAPI.RegisterRestaurant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RegisterRestaurant`: RegisterRestaurantResponse
	fmt.Fprintf(os.Stdout, "Response from `CustomersAPI.RegisterRestaurant`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRegisterRestaurantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registerRestaurantRequest** | [**RegisterRestaurantRequest**](RegisterRestaurantRequest.md) |  | 

### Return type

[**RegisterRestaurantResponse**](RegisterRestaurantResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

