# WebauthnCredentialResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**DeviceName** | Pointer to **string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewWebauthnCredentialResponse

`func NewWebauthnCredentialResponse() *WebauthnCredentialResponse`

NewWebauthnCredentialResponse instantiates a new WebauthnCredentialResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebauthnCredentialResponseWithDefaults

`func NewWebauthnCredentialResponseWithDefaults() *WebauthnCredentialResponse`

NewWebauthnCredentialResponseWithDefaults instantiates a new WebauthnCredentialResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *WebauthnCredentialResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *WebauthnCredentialResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *WebauthnCredentialResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *WebauthnCredentialResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetDeviceName

`func (o *WebauthnCredentialResponse) GetDeviceName() string`

GetDeviceName returns the DeviceName field if non-nil, zero value otherwise.

### GetDeviceNameOk

`func (o *WebauthnCredentialResponse) GetDeviceNameOk() (*string, bool)`

GetDeviceNameOk returns a tuple with the DeviceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceName

`func (o *WebauthnCredentialResponse) SetDeviceName(v string)`

SetDeviceName sets DeviceName field to given value.

### HasDeviceName

`func (o *WebauthnCredentialResponse) HasDeviceName() bool`

HasDeviceName returns a boolean if a field has been set.

### GetCreatedAt

`func (o *WebauthnCredentialResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *WebauthnCredentialResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *WebauthnCredentialResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *WebauthnCredentialResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


