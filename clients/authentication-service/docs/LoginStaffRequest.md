# LoginStaffRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** |  | 
**RestaurantId** | **string** |  | 
**Password** | **string** |  | 

## Methods

### NewLoginStaffRequest

`func NewLoginStaffRequest(email string, restaurantId string, password string, ) *LoginStaffRequest`

NewLoginStaffRequest instantiates a new LoginStaffRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginStaffRequestWithDefaults

`func NewLoginStaffRequestWithDefaults() *LoginStaffRequest`

NewLoginStaffRequestWithDefaults instantiates a new LoginStaffRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *LoginStaffRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *LoginStaffRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *LoginStaffRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetRestaurantId

`func (o *LoginStaffRequest) GetRestaurantId() string`

GetRestaurantId returns the RestaurantId field if non-nil, zero value otherwise.

### GetRestaurantIdOk

`func (o *LoginStaffRequest) GetRestaurantIdOk() (*string, bool)`

GetRestaurantIdOk returns a tuple with the RestaurantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRestaurantId

`func (o *LoginStaffRequest) SetRestaurantId(v string)`

SetRestaurantId sets RestaurantId field to given value.


### GetPassword

`func (o *LoginStaffRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *LoginStaffRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *LoginStaffRequest) SetPassword(v string)`

SetPassword sets Password field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


