# LoginResultEnvelope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Meta** | [**Meta**](Meta.md) |  | 
**Data** | [**LoginResult**](LoginResult.md) |  | 

## Methods

### NewLoginResultEnvelope

`func NewLoginResultEnvelope(meta Meta, data LoginResult, ) *LoginResultEnvelope`

NewLoginResultEnvelope instantiates a new LoginResultEnvelope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginResultEnvelopeWithDefaults

`func NewLoginResultEnvelopeWithDefaults() *LoginResultEnvelope`

NewLoginResultEnvelopeWithDefaults instantiates a new LoginResultEnvelope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMeta

`func (o *LoginResultEnvelope) GetMeta() Meta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *LoginResultEnvelope) GetMetaOk() (*Meta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *LoginResultEnvelope) SetMeta(v Meta)`

SetMeta sets Meta field to given value.


### GetData

`func (o *LoginResultEnvelope) GetData() LoginResult`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *LoginResultEnvelope) GetDataOk() (*LoginResult, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *LoginResultEnvelope) SetData(v LoginResult)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


