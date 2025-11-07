# RestaurantContact

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PhonePrefix** | **string** | E.164 country/area prefix, leading &#39;+&#39; required | 
**PhoneNumber** | **string** | Restaurant&#39;s phone number | 
**Email** | **string** | Restaurant&#39;s email address | 
**Address** | **string** | Restaurant&#39;s address | 
**City** | **string** | Restaurant&#39;s city | 
**PostalCode** | **string** | Restaurant&#39;s postal code | 
**CountryCode** | **string** | Restaurant&#39;s country code in ISO 3166-1 alpha-2 format | 

## Methods

### NewRestaurantContact

`func NewRestaurantContact(phonePrefix string, phoneNumber string, email string, address string, city string, postalCode string, countryCode string, ) *RestaurantContact`

NewRestaurantContact instantiates a new RestaurantContact object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRestaurantContactWithDefaults

`func NewRestaurantContactWithDefaults() *RestaurantContact`

NewRestaurantContactWithDefaults instantiates a new RestaurantContact object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPhonePrefix

`func (o *RestaurantContact) GetPhonePrefix() string`

GetPhonePrefix returns the PhonePrefix field if non-nil, zero value otherwise.

### GetPhonePrefixOk

`func (o *RestaurantContact) GetPhonePrefixOk() (*string, bool)`

GetPhonePrefixOk returns a tuple with the PhonePrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhonePrefix

`func (o *RestaurantContact) SetPhonePrefix(v string)`

SetPhonePrefix sets PhonePrefix field to given value.


### GetPhoneNumber

`func (o *RestaurantContact) GetPhoneNumber() string`

GetPhoneNumber returns the PhoneNumber field if non-nil, zero value otherwise.

### GetPhoneNumberOk

`func (o *RestaurantContact) GetPhoneNumberOk() (*string, bool)`

GetPhoneNumberOk returns a tuple with the PhoneNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneNumber

`func (o *RestaurantContact) SetPhoneNumber(v string)`

SetPhoneNumber sets PhoneNumber field to given value.


### GetEmail

`func (o *RestaurantContact) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RestaurantContact) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RestaurantContact) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetAddress

`func (o *RestaurantContact) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *RestaurantContact) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *RestaurantContact) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetCity

`func (o *RestaurantContact) GetCity() string`

GetCity returns the City field if non-nil, zero value otherwise.

### GetCityOk

`func (o *RestaurantContact) GetCityOk() (*string, bool)`

GetCityOk returns a tuple with the City field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCity

`func (o *RestaurantContact) SetCity(v string)`

SetCity sets City field to given value.


### GetPostalCode

`func (o *RestaurantContact) GetPostalCode() string`

GetPostalCode returns the PostalCode field if non-nil, zero value otherwise.

### GetPostalCodeOk

`func (o *RestaurantContact) GetPostalCodeOk() (*string, bool)`

GetPostalCodeOk returns a tuple with the PostalCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostalCode

`func (o *RestaurantContact) SetPostalCode(v string)`

SetPostalCode sets PostalCode field to given value.


### GetCountryCode

`func (o *RestaurantContact) GetCountryCode() string`

GetCountryCode returns the CountryCode field if non-nil, zero value otherwise.

### GetCountryCodeOk

`func (o *RestaurantContact) GetCountryCodeOk() (*string, bool)`

GetCountryCodeOk returns a tuple with the CountryCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryCode

`func (o *RestaurantContact) SetCountryCode(v string)`

SetCountryCode sets CountryCode field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


