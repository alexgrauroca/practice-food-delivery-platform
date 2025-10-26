# RegisterCustomerResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique customer identifier in the auth service | 
**Email** | **string** | Customer&#39;s email address | 
**CreatedAt** | **time.Time** | Account creation timestamp | 

## Methods

### NewRegisterCustomerResponse

`func NewRegisterCustomerResponse(id string, email string, createdAt time.Time, ) *RegisterCustomerResponse`

NewRegisterCustomerResponse instantiates a new RegisterCustomerResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterCustomerResponseWithDefaults

`func NewRegisterCustomerResponseWithDefaults() *RegisterCustomerResponse`

NewRegisterCustomerResponseWithDefaults instantiates a new RegisterCustomerResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *RegisterCustomerResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *RegisterCustomerResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *RegisterCustomerResponse) SetId(v string)`

SetId sets Id field to given value.


### GetEmail

`func (o *RegisterCustomerResponse) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegisterCustomerResponse) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegisterCustomerResponse) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetCreatedAt

`func (o *RegisterCustomerResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *RegisterCustomerResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *RegisterCustomerResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


