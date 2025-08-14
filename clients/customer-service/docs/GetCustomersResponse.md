# GetCustomersResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Items** | [**[]Customer**](Customer.md) | List of customers | 
**Pagination** | [**Pagination**](Pagination.md) |  | 

## Methods

### NewGetCustomersResponse

`func NewGetCustomersResponse(items []Customer, pagination Pagination, ) *GetCustomersResponse`

NewGetCustomersResponse instantiates a new GetCustomersResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetCustomersResponseWithDefaults

`func NewGetCustomersResponseWithDefaults() *GetCustomersResponse`

NewGetCustomersResponseWithDefaults instantiates a new GetCustomersResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetItems

`func (o *GetCustomersResponse) GetItems() []Customer`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *GetCustomersResponse) GetItemsOk() (*[]Customer, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *GetCustomersResponse) SetItems(v []Customer)`

SetItems sets Items field to given value.


### GetPagination

`func (o *GetCustomersResponse) GetPagination() Pagination`

GetPagination returns the Pagination field if non-nil, zero value otherwise.

### GetPaginationOk

`func (o *GetCustomersResponse) GetPaginationOk() (*Pagination, bool)`

GetPaginationOk returns a tuple with the Pagination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPagination

`func (o *GetCustomersResponse) SetPagination(v Pagination)`

SetPagination sets Pagination field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


