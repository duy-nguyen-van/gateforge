# WebauthnCredentialListEnvelope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Meta** | [**Meta**](Meta.md) |  | 
**Data** | [**[]WebauthnCredentialResponse**](WebauthnCredentialResponse.md) |  | 

## Methods

### NewWebauthnCredentialListEnvelope

`func NewWebauthnCredentialListEnvelope(meta Meta, data []WebauthnCredentialResponse, ) *WebauthnCredentialListEnvelope`

NewWebauthnCredentialListEnvelope instantiates a new WebauthnCredentialListEnvelope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebauthnCredentialListEnvelopeWithDefaults

`func NewWebauthnCredentialListEnvelopeWithDefaults() *WebauthnCredentialListEnvelope`

NewWebauthnCredentialListEnvelopeWithDefaults instantiates a new WebauthnCredentialListEnvelope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMeta

`func (o *WebauthnCredentialListEnvelope) GetMeta() Meta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *WebauthnCredentialListEnvelope) GetMetaOk() (*Meta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *WebauthnCredentialListEnvelope) SetMeta(v Meta)`

SetMeta sets Meta field to given value.


### GetData

`func (o *WebauthnCredentialListEnvelope) GetData() []WebauthnCredentialResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *WebauthnCredentialListEnvelope) GetDataOk() (*[]WebauthnCredentialResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *WebauthnCredentialListEnvelope) SetData(v []WebauthnCredentialResponse)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


