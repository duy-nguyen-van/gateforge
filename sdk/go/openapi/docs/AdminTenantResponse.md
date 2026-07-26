# AdminTenantResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Domain** | Pointer to **string** |  | [optional] 
**UserCount** | Pointer to **int32** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewAdminTenantResponse

`func NewAdminTenantResponse() *AdminTenantResponse`

NewAdminTenantResponse instantiates a new AdminTenantResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminTenantResponseWithDefaults

`func NewAdminTenantResponseWithDefaults() *AdminTenantResponse`

NewAdminTenantResponseWithDefaults instantiates a new AdminTenantResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AdminTenantResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AdminTenantResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AdminTenantResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *AdminTenantResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *AdminTenantResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AdminTenantResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AdminTenantResponse) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *AdminTenantResponse) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDomain

`func (o *AdminTenantResponse) GetDomain() string`

GetDomain returns the Domain field if non-nil, zero value otherwise.

### GetDomainOk

`func (o *AdminTenantResponse) GetDomainOk() (*string, bool)`

GetDomainOk returns a tuple with the Domain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomain

`func (o *AdminTenantResponse) SetDomain(v string)`

SetDomain sets Domain field to given value.

### HasDomain

`func (o *AdminTenantResponse) HasDomain() bool`

HasDomain returns a boolean if a field has been set.

### GetUserCount

`func (o *AdminTenantResponse) GetUserCount() int32`

GetUserCount returns the UserCount field if non-nil, zero value otherwise.

### GetUserCountOk

`func (o *AdminTenantResponse) GetUserCountOk() (*int32, bool)`

GetUserCountOk returns a tuple with the UserCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserCount

`func (o *AdminTenantResponse) SetUserCount(v int32)`

SetUserCount sets UserCount field to given value.

### HasUserCount

`func (o *AdminTenantResponse) HasUserCount() bool`

HasUserCount returns a boolean if a field has been set.

### GetCreatedAt

`func (o *AdminTenantResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *AdminTenantResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *AdminTenantResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *AdminTenantResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


