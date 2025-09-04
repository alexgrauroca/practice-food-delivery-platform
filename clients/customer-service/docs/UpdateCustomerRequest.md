# UpdateCustomerRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Customer&#39;s full name | 
**Address** | **string** | Customer&#39;s address | 
**City** | **string** | Customer&#39;s city | 
**PostalCode** | **string** | Customer&#39;s postal code | 
**CountryCode** | **string** | Customer&#39;s country code in ISO 3166-1 alpha-2 format | 

## Methods

### NewUpdateCustomerRequest

`func NewUpdateCustomerRequest(name string, address string, city string, postalCode string, countryCode string, ) *UpdateCustomerRequest`

NewUpdateCustomerRequest instantiates a new UpdateCustomerRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateCustomerRequestWithDefaults

`func NewUpdateCustomerRequestWithDefaults() *UpdateCustomerRequest`

NewUpdateCustomerRequestWithDefaults instantiates a new UpdateCustomerRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *UpdateCustomerRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdateCustomerRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdateCustomerRequest) SetName(v string)`

SetName sets Name field to given value.


### GetAddress

`func (o *UpdateCustomerRequest) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *UpdateCustomerRequest) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *UpdateCustomerRequest) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetCity

`func (o *UpdateCustomerRequest) GetCity() string`

GetCity returns the City field if non-nil, zero value otherwise.

### GetCityOk

`func (o *UpdateCustomerRequest) GetCityOk() (*string, bool)`

GetCityOk returns a tuple with the City field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCity

`func (o *UpdateCustomerRequest) SetCity(v string)`

SetCity sets City field to given value.


### GetPostalCode

`func (o *UpdateCustomerRequest) GetPostalCode() string`

GetPostalCode returns the PostalCode field if non-nil, zero value otherwise.

### GetPostalCodeOk

`func (o *UpdateCustomerRequest) GetPostalCodeOk() (*string, bool)`

GetPostalCodeOk returns a tuple with the PostalCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostalCode

`func (o *UpdateCustomerRequest) SetPostalCode(v string)`

SetPostalCode sets PostalCode field to given value.


### GetCountryCode

`func (o *UpdateCustomerRequest) GetCountryCode() string`

GetCountryCode returns the CountryCode field if non-nil, zero value otherwise.

### GetCountryCodeOk

`func (o *UpdateCustomerRequest) GetCountryCodeOk() (*string, bool)`

GetCountryCodeOk returns a tuple with the CountryCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryCode

`func (o *UpdateCustomerRequest) SetCountryCode(v string)`

SetCountryCode sets CountryCode field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


