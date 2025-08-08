# RegisterCustomerResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique customer identifier | 
**Email** | **string** | Customer&#39;s email address | 
**Name** | **string** | Customer&#39;s full name | 
**Address** | **string** | Customer&#39;s address | 
**City** | **string** | Customer&#39;s city | 
**PostalCode** | **string** | Customer&#39;s postal code | 
**CountryCode** | **string** | Customer&#39;s country code in ISO 3166-1 alpha-2 format | 
**CreatedAt** | **time.Time** | Account creation timestamp | 

## Methods

### NewRegisterCustomerResponse

`func NewRegisterCustomerResponse(id string, email string, name string, address string, city string, postalCode string, countryCode string, createdAt time.Time, ) *RegisterCustomerResponse`

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


### GetName

`func (o *RegisterCustomerResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *RegisterCustomerResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *RegisterCustomerResponse) SetName(v string)`

SetName sets Name field to given value.


### GetAddress

`func (o *RegisterCustomerResponse) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *RegisterCustomerResponse) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *RegisterCustomerResponse) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetCity

`func (o *RegisterCustomerResponse) GetCity() string`

GetCity returns the City field if non-nil, zero value otherwise.

### GetCityOk

`func (o *RegisterCustomerResponse) GetCityOk() (*string, bool)`

GetCityOk returns a tuple with the City field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCity

`func (o *RegisterCustomerResponse) SetCity(v string)`

SetCity sets City field to given value.


### GetPostalCode

`func (o *RegisterCustomerResponse) GetPostalCode() string`

GetPostalCode returns the PostalCode field if non-nil, zero value otherwise.

### GetPostalCodeOk

`func (o *RegisterCustomerResponse) GetPostalCodeOk() (*string, bool)`

GetPostalCodeOk returns a tuple with the PostalCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostalCode

`func (o *RegisterCustomerResponse) SetPostalCode(v string)`

SetPostalCode sets PostalCode field to given value.


### GetCountryCode

`func (o *RegisterCustomerResponse) GetCountryCode() string`

GetCountryCode returns the CountryCode field if non-nil, zero value otherwise.

### GetCountryCodeOk

`func (o *RegisterCustomerResponse) GetCountryCodeOk() (*string, bool)`

GetCountryCodeOk returns a tuple with the CountryCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryCode

`func (o *RegisterCustomerResponse) SetCountryCode(v string)`

SetCountryCode sets CountryCode field to given value.


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


