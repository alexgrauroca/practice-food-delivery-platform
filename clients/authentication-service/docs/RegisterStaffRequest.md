# RegisterStaffRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StaffId** | **string** | Unique staff identifier | 
**Email** | **string** | Staff&#39;s email address | 
**Password** | **string** | Password must be at least 8 characters long | 

## Methods

### NewRegisterStaffRequest

`func NewRegisterStaffRequest(staffId string, email string, password string, ) *RegisterStaffRequest`

NewRegisterStaffRequest instantiates a new RegisterStaffRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterStaffRequestWithDefaults

`func NewRegisterStaffRequestWithDefaults() *RegisterStaffRequest`

NewRegisterStaffRequestWithDefaults instantiates a new RegisterStaffRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStaffId

`func (o *RegisterStaffRequest) GetStaffId() string`

GetStaffId returns the StaffId field if non-nil, zero value otherwise.

### GetStaffIdOk

`func (o *RegisterStaffRequest) GetStaffIdOk() (*string, bool)`

GetStaffIdOk returns a tuple with the StaffId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStaffId

`func (o *RegisterStaffRequest) SetStaffId(v string)`

SetStaffId sets StaffId field to given value.


### GetEmail

`func (o *RegisterStaffRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegisterStaffRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegisterStaffRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetPassword

`func (o *RegisterStaffRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *RegisterStaffRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *RegisterStaffRequest) SetPassword(v string)`

SetPassword sets Password field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


