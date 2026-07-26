# TenantSelectionResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SelectionRequired** | **bool** |  | 
**Tenants** | [**[]TenantSummary**](TenantSummary.md) |  | 
**SelectionToken** | **string** |  | 
**ExpiresIn** | **int32** |  | 

## Methods

### NewTenantSelectionResponse

`func NewTenantSelectionResponse(selectionRequired bool, tenants []TenantSummary, selectionToken string, expiresIn int32, ) *TenantSelectionResponse`

NewTenantSelectionResponse instantiates a new TenantSelectionResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTenantSelectionResponseWithDefaults

`func NewTenantSelectionResponseWithDefaults() *TenantSelectionResponse`

NewTenantSelectionResponseWithDefaults instantiates a new TenantSelectionResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSelectionRequired

`func (o *TenantSelectionResponse) GetSelectionRequired() bool`

GetSelectionRequired returns the SelectionRequired field if non-nil, zero value otherwise.

### GetSelectionRequiredOk

`func (o *TenantSelectionResponse) GetSelectionRequiredOk() (*bool, bool)`

GetSelectionRequiredOk returns a tuple with the SelectionRequired field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSelectionRequired

`func (o *TenantSelectionResponse) SetSelectionRequired(v bool)`

SetSelectionRequired sets SelectionRequired field to given value.


### GetTenants

`func (o *TenantSelectionResponse) GetTenants() []TenantSummary`

GetTenants returns the Tenants field if non-nil, zero value otherwise.

### GetTenantsOk

`func (o *TenantSelectionResponse) GetTenantsOk() (*[]TenantSummary, bool)`

GetTenantsOk returns a tuple with the Tenants field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenants

`func (o *TenantSelectionResponse) SetTenants(v []TenantSummary)`

SetTenants sets Tenants field to given value.


### GetSelectionToken

`func (o *TenantSelectionResponse) GetSelectionToken() string`

GetSelectionToken returns the SelectionToken field if non-nil, zero value otherwise.

### GetSelectionTokenOk

`func (o *TenantSelectionResponse) GetSelectionTokenOk() (*string, bool)`

GetSelectionTokenOk returns a tuple with the SelectionToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSelectionToken

`func (o *TenantSelectionResponse) SetSelectionToken(v string)`

SetSelectionToken sets SelectionToken field to given value.


### GetExpiresIn

`func (o *TenantSelectionResponse) GetExpiresIn() int32`

GetExpiresIn returns the ExpiresIn field if non-nil, zero value otherwise.

### GetExpiresInOk

`func (o *TenantSelectionResponse) GetExpiresInOk() (*int32, bool)`

GetExpiresInOk returns a tuple with the ExpiresIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresIn

`func (o *TenantSelectionResponse) SetExpiresIn(v int32)`

SetExpiresIn sets ExpiresIn field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


