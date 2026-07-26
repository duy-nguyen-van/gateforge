# AdminClientUsageEnvelope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Meta** | [**Meta**](Meta.md) |  | 
**Data** | [**AdminClientUsageResponse**](AdminClientUsageResponse.md) |  | 

## Methods

### NewAdminClientUsageEnvelope

`func NewAdminClientUsageEnvelope(meta Meta, data AdminClientUsageResponse, ) *AdminClientUsageEnvelope`

NewAdminClientUsageEnvelope instantiates a new AdminClientUsageEnvelope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminClientUsageEnvelopeWithDefaults

`func NewAdminClientUsageEnvelopeWithDefaults() *AdminClientUsageEnvelope`

NewAdminClientUsageEnvelopeWithDefaults instantiates a new AdminClientUsageEnvelope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMeta

`func (o *AdminClientUsageEnvelope) GetMeta() Meta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *AdminClientUsageEnvelope) GetMetaOk() (*Meta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *AdminClientUsageEnvelope) SetMeta(v Meta)`

SetMeta sets Meta field to given value.


### GetData

`func (o *AdminClientUsageEnvelope) GetData() AdminClientUsageResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *AdminClientUsageEnvelope) GetDataOk() (*AdminClientUsageResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *AdminClientUsageEnvelope) SetData(v AdminClientUsageResponse)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


