# AdminStatsEnvelope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Meta** | [**Meta**](Meta.md) |  | 
**Data** | [**AdminStatsResponse**](AdminStatsResponse.md) |  | 

## Methods

### NewAdminStatsEnvelope

`func NewAdminStatsEnvelope(meta Meta, data AdminStatsResponse, ) *AdminStatsEnvelope`

NewAdminStatsEnvelope instantiates a new AdminStatsEnvelope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminStatsEnvelopeWithDefaults

`func NewAdminStatsEnvelopeWithDefaults() *AdminStatsEnvelope`

NewAdminStatsEnvelopeWithDefaults instantiates a new AdminStatsEnvelope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMeta

`func (o *AdminStatsEnvelope) GetMeta() Meta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *AdminStatsEnvelope) GetMetaOk() (*Meta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *AdminStatsEnvelope) SetMeta(v Meta)`

SetMeta sets Meta field to given value.


### GetData

`func (o *AdminStatsEnvelope) GetData() AdminStatsResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *AdminStatsEnvelope) GetDataOk() (*AdminStatsResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *AdminStatsEnvelope) SetData(v AdminStatsResponse)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


