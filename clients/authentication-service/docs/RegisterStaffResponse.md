# RegisterStaffResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique staff identifier in the auth service | 
**Email** | **string** | Staff&#39;s email address | 
**RestaurantId** | Pointer to **string** | Unique restaurant identifier in the auth service | [optional] 
**CreatedAt** | **time.Time** | Staff creation timestamp | 
**UpdatedAt** | Pointer to **time.Time** | Staff update timestamp | [optional] 

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


### GetRestaurantId

`func (o *RegisterStaffResponse) GetRestaurantId() string`

GetRestaurantId returns the RestaurantId field if non-nil, zero value otherwise.

### GetRestaurantIdOk

`func (o *RegisterStaffResponse) GetRestaurantIdOk() (*string, bool)`

GetRestaurantIdOk returns a tuple with the RestaurantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRestaurantId

`func (o *RegisterStaffResponse) SetRestaurantId(v string)`

SetRestaurantId sets RestaurantId field to given value.

### HasRestaurantId

`func (o *RegisterStaffResponse) HasRestaurantId() bool`

HasRestaurantId returns a boolean if a field has been set.

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


### GetUpdatedAt

`func (o *RegisterStaffResponse) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *RegisterStaffResponse) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *RegisterStaffResponse) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *RegisterStaffResponse) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


