# DatabaseHealthEnvelope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Meta** | [**Meta**](Meta.md) |  | 
**Data** | **map[string]interface{}** |  | 

## Methods

### NewDatabaseHealthEnvelope

`func NewDatabaseHealthEnvelope(meta Meta, data map[string]interface{}, ) *DatabaseHealthEnvelope`

NewDatabaseHealthEnvelope instantiates a new DatabaseHealthEnvelope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDatabaseHealthEnvelopeWithDefaults

`func NewDatabaseHealthEnvelopeWithDefaults() *DatabaseHealthEnvelope`

NewDatabaseHealthEnvelopeWithDefaults instantiates a new DatabaseHealthEnvelope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMeta

`func (o *DatabaseHealthEnvelope) GetMeta() Meta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *DatabaseHealthEnvelope) GetMetaOk() (*Meta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *DatabaseHealthEnvelope) SetMeta(v Meta)`

SetMeta sets Meta field to given value.


### GetData

`func (o *DatabaseHealthEnvelope) GetData() map[string]interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *DatabaseHealthEnvelope) GetDataOk() (*map[string]interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *DatabaseHealthEnvelope) SetData(v map[string]interface{})`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


