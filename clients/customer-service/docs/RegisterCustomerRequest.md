# RegisterCustomerRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | Customer&#39;s email address | 
**Password** | **string** | Password must be at least 8 characters long | 
**Name** | **string** | Customer&#39;s full name | 
**Address** | **string** | Customer&#39;s address | 
**City** | **string** | Customer&#39;s city | 
**PostalCode** | **string** | Customer&#39;s postal code | 
**CountryCode** | **string** | Customer&#39;s country code in ISO 3166-1 alpha-2 format | 

## Methods

### NewRegisterCustomerRequest

`func NewRegisterCustomerRequest(email string, password string, name string, address string, city string, postalCode string, countryCode string, ) *RegisterCustomerRequest`

NewRegisterCustomerRequest instantiates a new RegisterCustomerRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterCustomerRequestWithDefaults

`func NewRegisterCustomerRequestWithDefaults() *RegisterCustomerRequest`

NewRegisterCustomerRequestWithDefaults instantiates a new RegisterCustomerRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *RegisterCustomerRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegisterCustomerRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegisterCustomerRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetPassword

`func (o *RegisterCustomerRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *RegisterCustomerRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *RegisterCustomerRequest) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetName

`func (o *RegisterCustomerRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *RegisterCustomerRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *RegisterCustomerRequest) SetName(v string)`

SetName sets Name field to given value.


### GetAddress

`func (o *RegisterCustomerRequest) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *RegisterCustomerRequest) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *RegisterCustomerRequest) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetCity

`func (o *RegisterCustomerRequest) GetCity() string`

GetCity returns the City field if non-nil, zero value otherwise.

### GetCityOk

`func (o *RegisterCustomerRequest) GetCityOk() (*string, bool)`

GetCityOk returns a tuple with the City field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCity

`func (o *RegisterCustomerRequest) SetCity(v string)`

SetCity sets City field to given value.


### GetPostalCode

`func (o *RegisterCustomerRequest) GetPostalCode() string`

GetPostalCode returns the PostalCode field if non-nil, zero value otherwise.

### GetPostalCodeOk

`func (o *RegisterCustomerRequest) GetPostalCodeOk() (*string, bool)`

GetPostalCodeOk returns a tuple with the PostalCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostalCode

`func (o *RegisterCustomerRequest) SetPostalCode(v string)`

SetPostalCode sets PostalCode field to given value.


### GetCountryCode

`func (o *RegisterCustomerRequest) GetCountryCode() string`

GetCountryCode returns the CountryCode field if non-nil, zero value otherwise.

### GetCountryCodeOk

`func (o *RegisterCustomerRequest) GetCountryCodeOk() (*string, bool)`

GetCountryCodeOk returns a tuple with the CountryCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryCode

`func (o *RegisterCustomerRequest) SetCountryCode(v string)`

SetCountryCode sets CountryCode field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


