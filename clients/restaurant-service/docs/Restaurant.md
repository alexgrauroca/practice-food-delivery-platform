# Restaurant

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique restaurant identifier | [optional] 
**VatCode** | **string** | Restaurant&#39;s VAT code | 
**Name** | **string** | Restaurant&#39;s commercial name | 
**LegalName** | **string** | Restaurant&#39;s legal name | 
**TaxId** | Pointer to **string** | Restaurant&#39;s tax identification number | [optional] 
**TimezoneId** | **string** | IANA Restaurant&#39;s timezone identifier | 
**Contact** | [**RestaurantContact**](RestaurantContact.md) |  | 
**CreatedAt** | Pointer to **time.Time** | Restaurant creation timestamp | [optional] [readonly] 
**UpdatedAt** | Pointer to **time.Time** | Restaurant last update timestamp | [optional] [readonly] 

## Methods

### NewRestaurant

`func NewRestaurant(vatCode string, name string, legalName string, timezoneId string, contact RestaurantContact, ) *Restaurant`

NewRestaurant instantiates a new Restaurant object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRestaurantWithDefaults

`func NewRestaurantWithDefaults() *Restaurant`

NewRestaurantWithDefaults instantiates a new Restaurant object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Restaurant) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Restaurant) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Restaurant) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Restaurant) HasId() bool`

HasId returns a boolean if a field has been set.

### GetVatCode

`func (o *Restaurant) GetVatCode() string`

GetVatCode returns the VatCode field if non-nil, zero value otherwise.

### GetVatCodeOk

`func (o *Restaurant) GetVatCodeOk() (*string, bool)`

GetVatCodeOk returns a tuple with the VatCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVatCode

`func (o *Restaurant) SetVatCode(v string)`

SetVatCode sets VatCode field to given value.


### GetName

`func (o *Restaurant) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Restaurant) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Restaurant) SetName(v string)`

SetName sets Name field to given value.


### GetLegalName

`func (o *Restaurant) GetLegalName() string`

GetLegalName returns the LegalName field if non-nil, zero value otherwise.

### GetLegalNameOk

`func (o *Restaurant) GetLegalNameOk() (*string, bool)`

GetLegalNameOk returns a tuple with the LegalName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLegalName

`func (o *Restaurant) SetLegalName(v string)`

SetLegalName sets LegalName field to given value.


### GetTaxId

`func (o *Restaurant) GetTaxId() string`

GetTaxId returns the TaxId field if non-nil, zero value otherwise.

### GetTaxIdOk

`func (o *Restaurant) GetTaxIdOk() (*string, bool)`

GetTaxIdOk returns a tuple with the TaxId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTaxId

`func (o *Restaurant) SetTaxId(v string)`

SetTaxId sets TaxId field to given value.

### HasTaxId

`func (o *Restaurant) HasTaxId() bool`

HasTaxId returns a boolean if a field has been set.

### GetTimezoneId

`func (o *Restaurant) GetTimezoneId() string`

GetTimezoneId returns the TimezoneId field if non-nil, zero value otherwise.

### GetTimezoneIdOk

`func (o *Restaurant) GetTimezoneIdOk() (*string, bool)`

GetTimezoneIdOk returns a tuple with the TimezoneId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimezoneId

`func (o *Restaurant) SetTimezoneId(v string)`

SetTimezoneId sets TimezoneId field to given value.


### GetContact

`func (o *Restaurant) GetContact() RestaurantContact`

GetContact returns the Contact field if non-nil, zero value otherwise.

### GetContactOk

`func (o *Restaurant) GetContactOk() (*RestaurantContact, bool)`

GetContactOk returns a tuple with the Contact field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContact

`func (o *Restaurant) SetContact(v RestaurantContact)`

SetContact sets Contact field to given value.


### GetCreatedAt

`func (o *Restaurant) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Restaurant) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Restaurant) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Restaurant) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *Restaurant) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Restaurant) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Restaurant) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *Restaurant) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


