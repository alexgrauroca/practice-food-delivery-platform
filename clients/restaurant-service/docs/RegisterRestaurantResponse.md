# RegisterRestaurantResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Restaurant** | [**Restaurant**](Restaurant.md) |  | 
**StaffOwner** | [**Staff**](Staff.md) |  | 

## Methods

### NewRegisterRestaurantResponse

`func NewRegisterRestaurantResponse(restaurant Restaurant, staffOwner Staff, ) *RegisterRestaurantResponse`

NewRegisterRestaurantResponse instantiates a new RegisterRestaurantResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterRestaurantResponseWithDefaults

`func NewRegisterRestaurantResponseWithDefaults() *RegisterRestaurantResponse`

NewRegisterRestaurantResponseWithDefaults instantiates a new RegisterRestaurantResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRestaurant

`func (o *RegisterRestaurantResponse) GetRestaurant() Restaurant`

GetRestaurant returns the Restaurant field if non-nil, zero value otherwise.

### GetRestaurantOk

`func (o *RegisterRestaurantResponse) GetRestaurantOk() (*Restaurant, bool)`

GetRestaurantOk returns a tuple with the Restaurant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRestaurant

`func (o *RegisterRestaurantResponse) SetRestaurant(v Restaurant)`

SetRestaurant sets Restaurant field to given value.


### GetStaffOwner

`func (o *RegisterRestaurantResponse) GetStaffOwner() Staff`

GetStaffOwner returns the StaffOwner field if non-nil, zero value otherwise.

### GetStaffOwnerOk

`func (o *RegisterRestaurantResponse) GetStaffOwnerOk() (*Staff, bool)`

GetStaffOwnerOk returns a tuple with the StaffOwner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStaffOwner

`func (o *RegisterRestaurantResponse) SetStaffOwner(v Staff)`

SetStaffOwner sets StaffOwner field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


