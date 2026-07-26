# LoginResponseEnvelope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Meta** | [**Meta**](Meta.md) |  | 
**Data** | [**LoginResponse**](LoginResponse.md) |  | 

## Methods

### NewLoginResponseEnvelope

`func NewLoginResponseEnvelope(meta Meta, data LoginResponse, ) *LoginResponseEnvelope`

NewLoginResponseEnvelope instantiates a new LoginResponseEnvelope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginResponseEnvelopeWithDefaults

`func NewLoginResponseEnvelopeWithDefaults() *LoginResponseEnvelope`

NewLoginResponseEnvelopeWithDefaults instantiates a new LoginResponseEnvelope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMeta

`func (o *LoginResponseEnvelope) GetMeta() Meta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *LoginResponseEnvelope) GetMetaOk() (*Meta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *LoginResponseEnvelope) SetMeta(v Meta)`

SetMeta sets Meta field to given value.


### GetData

`func (o *LoginResponseEnvelope) GetData() LoginResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *LoginResponseEnvelope) GetDataOk() (*LoginResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *LoginResponseEnvelope) SetData(v LoginResponse)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


