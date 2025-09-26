# \StaffAPI

All URIs are relative to *http://localhost:80*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RegisterStaff**](StaffAPI.md#RegisterStaff) | **Post** /v1.0/auth/staff | Register a new staff user



## RegisterStaff

> RegisterStaffResponse RegisterStaff(ctx).RegisterStaffRequest(registerStaffRequest).Execute()

Register a new staff user



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
	registerStaffRequest := *openapiclient.NewRegisterStaffRequest("507f1f77bcf86cd799439011", "user@example.com", "strongpassword123", "John Doe") // RegisterStaffRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StaffAPI.RegisterStaff(context.Background()).RegisterStaffRequest(registerStaffRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StaffAPI.RegisterStaff``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RegisterStaff`: RegisterStaffResponse
	fmt.Fprintf(os.Stdout, "Response from `StaffAPI.RegisterStaff`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRegisterStaffRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registerStaffRequest** | [**RegisterStaffRequest**](RegisterStaffRequest.md) |  | 

### Return type

[**RegisterStaffResponse**](RegisterStaffResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

