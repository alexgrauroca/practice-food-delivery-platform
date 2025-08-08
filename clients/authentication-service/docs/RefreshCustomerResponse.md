# RefreshCustomerResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessToken** | **string** | JWT access token for API authentication | 
**RefreshToken** | **string** | Token used to obtain a new access token when it expires | 
**ExpiresIn** | **int32** | Access token expiration time in seconds | 
**TokenType** | **string** | Access token type | 

## Methods

### NewRefreshCustomerResponse

`func NewRefreshCustomerResponse(accessToken string, refreshToken string, expiresIn int32, tokenType string, ) *RefreshCustomerResponse`

NewRefreshCustomerResponse instantiates a new RefreshCustomerResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRefreshCustomerResponseWithDefaults

`func NewRefreshCustomerResponseWithDefaults() *RefreshCustomerResponse`

NewRefreshCustomerResponseWithDefaults instantiates a new RefreshCustomerResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessToken

`func (o *RefreshCustomerResponse) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *RefreshCustomerResponse) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *RefreshCustomerResponse) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.


### GetRefreshToken

`func (o *RefreshCustomerResponse) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *RefreshCustomerResponse) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *RefreshCustomerResponse) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.


### GetExpiresIn

`func (o *RefreshCustomerResponse) GetExpiresIn() int32`

GetExpiresIn returns the ExpiresIn field if non-nil, zero value otherwise.

### GetExpiresInOk

`func (o *RefreshCustomerResponse) GetExpiresInOk() (*int32, bool)`

GetExpiresInOk returns a tuple with the ExpiresIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresIn

`func (o *RefreshCustomerResponse) SetExpiresIn(v int32)`

SetExpiresIn sets ExpiresIn field to given value.


### GetTokenType

`func (o *RefreshCustomerResponse) GetTokenType() string`

GetTokenType returns the TokenType field if non-nil, zero value otherwise.

### GetTokenTypeOk

`func (o *RefreshCustomerResponse) GetTokenTypeOk() (*string, bool)`

GetTokenTypeOk returns a tuple with the TokenType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenType

`func (o *RefreshCustomerResponse) SetTokenType(v string)`

SetTokenType sets TokenType field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


