# RegisterRestaurantRequestRestaurant

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**VatCode** | **string** | Restaurant&#39;s VAT code | 
**Name** | **string** | Restaurant&#39;s commercial name | 
**LegalName** | **string** | Restaurant&#39;s legal name | 
**TaxId** | Pointer to **string** | Restaurant&#39;s tax identification number | [optional] [default to ""]
**TimezoneId** | **string** | IANA Restaurant&#39;s timezone identifier | 
**Contact** | [**RestaurantContact**](RestaurantContact.md) |  | 

## Methods

### NewRegisterRestaurantRequestRestaurant

`func NewRegisterRestaurantRequestRestaurant(vatCode string, name string, legalName string, timezoneId string, contact RestaurantContact, ) *RegisterRestaurantRequestRestaurant`

NewRegisterRestaurantRequestRestaurant instantiates a new RegisterRestaurantRequestRestaurant object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterRestaurantRequestRestaurantWithDefaults

`func NewRegisterRestaurantRequestRestaurantWithDefaults() *RegisterRestaurantRequestRestaurant`

NewRegisterRestaurantRequestRestaurantWithDefaults instantiates a new RegisterRestaurantRequestRestaurant object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVatCode

`func (o *RegisterRestaurantRequestRestaurant) GetVatCode() string`

GetVatCode returns the VatCode field if non-nil, zero value otherwise.

### GetVatCodeOk

`func (o *RegisterRestaurantRequestRestaurant) GetVatCodeOk() (*string, bool)`

GetVatCodeOk returns a tuple with the VatCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVatCode

`func (o *RegisterRestaurantRequestRestaurant) SetVatCode(v string)`

SetVatCode sets VatCode field to given value.


### GetName

`func (o *RegisterRestaurantRequestRestaurant) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *RegisterRestaurantRequestRestaurant) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *RegisterRestaurantRequestRestaurant) SetName(v string)`

SetName sets Name field to given value.


### GetLegalName

`func (o *RegisterRestaurantRequestRestaurant) GetLegalName() string`

GetLegalName returns the LegalName field if non-nil, zero value otherwise.

### GetLegalNameOk

`func (o *RegisterRestaurantRequestRestaurant) GetLegalNameOk() (*string, bool)`

GetLegalNameOk returns a tuple with the LegalName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLegalName

`func (o *RegisterRestaurantRequestRestaurant) SetLegalName(v string)`

SetLegalName sets LegalName field to given value.


### GetTaxId

`func (o *RegisterRestaurantRequestRestaurant) GetTaxId() string`

GetTaxId returns the TaxId field if non-nil, zero value otherwise.

### GetTaxIdOk

`func (o *RegisterRestaurantRequestRestaurant) GetTaxIdOk() (*string, bool)`

GetTaxIdOk returns a tuple with the TaxId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTaxId

`func (o *RegisterRestaurantRequestRestaurant) SetTaxId(v string)`

SetTaxId sets TaxId field to given value.

### HasTaxId

`func (o *RegisterRestaurantRequestRestaurant) HasTaxId() bool`

HasTaxId returns a boolean if a field has been set.

### GetTimezoneId

`func (o *RegisterRestaurantRequestRestaurant) GetTimezoneId() string`

GetTimezoneId returns the TimezoneId field if non-nil, zero value otherwise.

### GetTimezoneIdOk

`func (o *RegisterRestaurantRequestRestaurant) GetTimezoneIdOk() (*string, bool)`

GetTimezoneIdOk returns a tuple with the TimezoneId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimezoneId

`func (o *RegisterRestaurantRequestRestaurant) SetTimezoneId(v string)`

SetTimezoneId sets TimezoneId field to given value.


### GetContact

`func (o *RegisterRestaurantRequestRestaurant) GetContact() RestaurantContact`

GetContact returns the Contact field if non-nil, zero value otherwise.

### GetContactOk

`func (o *RegisterRestaurantRequestRestaurant) GetContactOk() (*RestaurantContact, bool)`

GetContactOk returns a tuple with the Contact field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContact

`func (o *RegisterRestaurantRequestRestaurant) SetContact(v RestaurantContact)`

SetContact sets Contact field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


