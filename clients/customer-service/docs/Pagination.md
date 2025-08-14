# Pagination

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TotalItems** | Pointer to **int32** | Total number of elements across all pages | [optional] 
**TotalPages** | Pointer to **int32** | Total number of pages | [optional] 
**CurrentPage** | Pointer to **int32** | Current page number | [optional] 
**PageSize** | Pointer to **int32** | Number of elements per page | [optional] 

## Methods

### NewPagination

`func NewPagination() *Pagination`

NewPagination instantiates a new Pagination object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPaginationWithDefaults

`func NewPaginationWithDefaults() *Pagination`

NewPaginationWithDefaults instantiates a new Pagination object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTotalItems

`func (o *Pagination) GetTotalItems() int32`

GetTotalItems returns the TotalItems field if non-nil, zero value otherwise.

### GetTotalItemsOk

`func (o *Pagination) GetTotalItemsOk() (*int32, bool)`

GetTotalItemsOk returns a tuple with the TotalItems field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalItems

`func (o *Pagination) SetTotalItems(v int32)`

SetTotalItems sets TotalItems field to given value.

### HasTotalItems

`func (o *Pagination) HasTotalItems() bool`

HasTotalItems returns a boolean if a field has been set.

### GetTotalPages

`func (o *Pagination) GetTotalPages() int32`

GetTotalPages returns the TotalPages field if non-nil, zero value otherwise.

### GetTotalPagesOk

`func (o *Pagination) GetTotalPagesOk() (*int32, bool)`

GetTotalPagesOk returns a tuple with the TotalPages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalPages

`func (o *Pagination) SetTotalPages(v int32)`

SetTotalPages sets TotalPages field to given value.

### HasTotalPages

`func (o *Pagination) HasTotalPages() bool`

HasTotalPages returns a boolean if a field has been set.

### GetCurrentPage

`func (o *Pagination) GetCurrentPage() int32`

GetCurrentPage returns the CurrentPage field if non-nil, zero value otherwise.

### GetCurrentPageOk

`func (o *Pagination) GetCurrentPageOk() (*int32, bool)`

GetCurrentPageOk returns a tuple with the CurrentPage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentPage

`func (o *Pagination) SetCurrentPage(v int32)`

SetCurrentPage sets CurrentPage field to given value.

### HasCurrentPage

`func (o *Pagination) HasCurrentPage() bool`

HasCurrentPage returns a boolean if a field has been set.

### GetPageSize

`func (o *Pagination) GetPageSize() int32`

GetPageSize returns the PageSize field if non-nil, zero value otherwise.

### GetPageSizeOk

`func (o *Pagination) GetPageSizeOk() (*int32, bool)`

GetPageSizeOk returns a tuple with the PageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageSize

`func (o *Pagination) SetPageSize(v int32)`

SetPageSize sets PageSize field to given value.

### HasPageSize

`func (o *Pagination) HasPageSize() bool`

HasPageSize returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


