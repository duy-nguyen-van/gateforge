# WebauthnLoginStartRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** |  | 
**TenantId** | Pointer to **string** |  | [optional] 

## Methods

### NewWebauthnLoginStartRequest

`func NewWebauthnLoginStartRequest(email string, ) *WebauthnLoginStartRequest`

NewWebauthnLoginStartRequest instantiates a new WebauthnLoginStartRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebauthnLoginStartRequestWithDefaults

`func NewWebauthnLoginStartRequestWithDefaults() *WebauthnLoginStartRequest`

NewWebauthnLoginStartRequestWithDefaults instantiates a new WebauthnLoginStartRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *WebauthnLoginStartRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *WebauthnLoginStartRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *WebauthnLoginStartRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetTenantId

`func (o *WebauthnLoginStartRequest) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *WebauthnLoginStartRequest) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *WebauthnLoginStartRequest) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *WebauthnLoginStartRequest) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


