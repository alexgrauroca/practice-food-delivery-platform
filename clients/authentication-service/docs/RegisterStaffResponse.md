# RegisterStaffResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique staff identifier in the auth service | 
**Email** | **string** | Staff&#39;s email address | 
**CreatedAt** | **time.Time** | Account creation timestamp | 

## Methods

### NewRegisterStaffResponse

`func NewRegisterStaffResponse(id string, email string, createdAt time.Time, ) *RegisterStaffResponse`

NewRegisterStaffResponse instantiates a new RegisterStaffResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterStaffResponseWithDefaults

`func NewRegisterStaffResponseWithDefaults() *RegisterStaffResponse`

NewRegisterStaffResponseWithDefaults instantiates a new RegisterStaffResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *RegisterStaffResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *RegisterStaffResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *RegisterStaffResponse) SetId(v string)`

SetId sets Id field to given value.


### GetEmail

`func (o *RegisterStaffResponse) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegisterStaffResponse) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegisterStaffResponse) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetCreatedAt

`func (o *RegisterStaffResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *RegisterStaffResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *RegisterStaffResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


