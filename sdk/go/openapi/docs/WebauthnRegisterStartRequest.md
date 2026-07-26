# WebauthnRegisterStartRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DeviceName** | Pointer to **string** | Optional label stored with the passkey | [optional] 

## Methods

### NewWebauthnRegisterStartRequest

`func NewWebauthnRegisterStartRequest() *WebauthnRegisterStartRequest`

NewWebauthnRegisterStartRequest instantiates a new WebauthnRegisterStartRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebauthnRegisterStartRequestWithDefaults

`func NewWebauthnRegisterStartRequestWithDefaults() *WebauthnRegisterStartRequest`

NewWebauthnRegisterStartRequestWithDefaults instantiates a new WebauthnRegisterStartRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDeviceName

`func (o *WebauthnRegisterStartRequest) GetDeviceName() string`

GetDeviceName returns the DeviceName field if non-nil, zero value otherwise.

### GetDeviceNameOk

`func (o *WebauthnRegisterStartRequest) GetDeviceNameOk() (*string, bool)`

GetDeviceNameOk returns a tuple with the DeviceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceName

`func (o *WebauthnRegisterStartRequest) SetDeviceName(v string)`

SetDeviceName sets DeviceName field to given value.

### HasDeviceName

`func (o *WebauthnRegisterStartRequest) HasDeviceName() bool`

HasDeviceName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


