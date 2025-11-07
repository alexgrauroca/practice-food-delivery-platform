# RegisterRestaurantRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Restaurant** | [**RegisterRestaurantRequestRestaurant**](RegisterRestaurantRequestRestaurant.md) |  | 
**StaffOwner** | [**RegisterRestaurantRequestStaffOwner**](RegisterRestaurantRequestStaffOwner.md) |  | 

## Methods

### NewRegisterRestaurantRequest

`func NewRegisterRestaurantRequest(restaurant RegisterRestaurantRequestRestaurant, staffOwner RegisterRestaurantRequestStaffOwner, ) *RegisterRestaurantRequest`

NewRegisterRestaurantRequest instantiates a new RegisterRestaurantRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterRestaurantRequestWithDefaults

`func NewRegisterRestaurantRequestWithDefaults() *RegisterRestaurantRequest`

NewRegisterRestaurantRequestWithDefaults instantiates a new RegisterRestaurantRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRestaurant

`func (o *RegisterRestaurantRequest) GetRestaurant() RegisterRestaurantRequestRestaurant`

GetRestaurant returns the Restaurant field if non-nil, zero value otherwise.

### GetRestaurantOk

`func (o *RegisterRestaurantRequest) GetRestaurantOk() (*RegisterRestaurantRequestRestaurant, bool)`

GetRestaurantOk returns a tuple with the Restaurant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRestaurant

`func (o *RegisterRestaurantRequest) SetRestaurant(v RegisterRestaurantRequestRestaurant)`

SetRestaurant sets Restaurant field to given value.


### GetStaffOwner

`func (o *RegisterRestaurantRequest) GetStaffOwner() RegisterRestaurantRequestStaffOwner`

GetStaffOwner returns the StaffOwner field if non-nil, zero value otherwise.

### GetStaffOwnerOk

`func (o *RegisterRestaurantRequest) GetStaffOwnerOk() (*RegisterRestaurantRequestStaffOwner, bool)`

GetStaffOwnerOk returns a tuple with the StaffOwner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStaffOwner

`func (o *RegisterRestaurantRequest) SetStaffOwner(v RegisterRestaurantRequestStaffOwner)`

SetStaffOwner sets StaffOwner field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


