# AdminCreateTenantRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Domain** | Pointer to **string** |  | [optional] 

## Methods

### NewAdminCreateTenantRequest

`func NewAdminCreateTenantRequest(name string, ) *AdminCreateTenantRequest`

NewAdminCreateTenantRequest instantiates a new AdminCreateTenantRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminCreateTenantRequestWithDefaults

`func NewAdminCreateTenantRequestWithDefaults() *AdminCreateTenantRequest`

NewAdminCreateTenantRequestWithDefaults instantiates a new AdminCreateTenantRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *AdminCreateTenantRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AdminCreateTenantRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AdminCreateTenantRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDomain

`func (o *AdminCreateTenantRequest) GetDomain() string`

GetDomain returns the Domain field if non-nil, zero value otherwise.

### GetDomainOk

`func (o *AdminCreateTenantRequest) GetDomainOk() (*string, bool)`

GetDomainOk returns a tuple with the Domain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomain

`func (o *AdminCreateTenantRequest) SetDomain(v string)`

SetDomain sets Domain field to given value.

### HasDomain

`func (o *AdminCreateTenantRequest) HasDomain() bool`

HasDomain returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


