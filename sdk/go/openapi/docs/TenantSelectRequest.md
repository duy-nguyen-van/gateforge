# TenantSelectRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SelectionToken** | **string** |  | 
**TenantId** | **string** |  | 
**RememberMe** | Pointer to **bool** |  | [optional] 

## Methods

### NewTenantSelectRequest

`func NewTenantSelectRequest(selectionToken string, tenantId string, ) *TenantSelectRequest`

NewTenantSelectRequest instantiates a new TenantSelectRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTenantSelectRequestWithDefaults

`func NewTenantSelectRequestWithDefaults() *TenantSelectRequest`

NewTenantSelectRequestWithDefaults instantiates a new TenantSelectRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSelectionToken

`func (o *TenantSelectRequest) GetSelectionToken() string`

GetSelectionToken returns the SelectionToken field if non-nil, zero value otherwise.

### GetSelectionTokenOk

`func (o *TenantSelectRequest) GetSelectionTokenOk() (*string, bool)`

GetSelectionTokenOk returns a tuple with the SelectionToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSelectionToken

`func (o *TenantSelectRequest) SetSelectionToken(v string)`

SetSelectionToken sets SelectionToken field to given value.


### GetTenantId

`func (o *TenantSelectRequest) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *TenantSelectRequest) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *TenantSelectRequest) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.


### GetRememberMe

`func (o *TenantSelectRequest) GetRememberMe() bool`

GetRememberMe returns the RememberMe field if non-nil, zero value otherwise.

### GetRememberMeOk

`func (o *TenantSelectRequest) GetRememberMeOk() (*bool, bool)`

GetRememberMeOk returns a tuple with the RememberMe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRememberMe

`func (o *TenantSelectRequest) SetRememberMe(v bool)`

SetRememberMe sets RememberMe field to given value.

### HasRememberMe

`func (o *TenantSelectRequest) HasRememberMe() bool`

HasRememberMe returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


