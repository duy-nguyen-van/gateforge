# AdminUserListEnvelope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Meta** | [**Meta**](Meta.md) |  | 
**Data** | [**[]AdminUserResponse**](AdminUserResponse.md) |  | 

## Methods

### NewAdminUserListEnvelope

`func NewAdminUserListEnvelope(meta Meta, data []AdminUserResponse, ) *AdminUserListEnvelope`

NewAdminUserListEnvelope instantiates a new AdminUserListEnvelope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminUserListEnvelopeWithDefaults

`func NewAdminUserListEnvelopeWithDefaults() *AdminUserListEnvelope`

NewAdminUserListEnvelopeWithDefaults instantiates a new AdminUserListEnvelope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMeta

`func (o *AdminUserListEnvelope) GetMeta() Meta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *AdminUserListEnvelope) GetMetaOk() (*Meta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *AdminUserListEnvelope) SetMeta(v Meta)`

SetMeta sets Meta field to given value.


### GetData

`func (o *AdminUserListEnvelope) GetData() []AdminUserResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *AdminUserListEnvelope) GetDataOk() (*[]AdminUserResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *AdminUserListEnvelope) SetData(v []AdminUserResponse)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


