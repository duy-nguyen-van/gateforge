# PublicFederationProviderListEnvelope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Meta** | [**Meta**](Meta.md) |  | 
**Data** | [**[]PublicFederationProviderResponse**](PublicFederationProviderResponse.md) |  | 

## Methods

### NewPublicFederationProviderListEnvelope

`func NewPublicFederationProviderListEnvelope(meta Meta, data []PublicFederationProviderResponse, ) *PublicFederationProviderListEnvelope`

NewPublicFederationProviderListEnvelope instantiates a new PublicFederationProviderListEnvelope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPublicFederationProviderListEnvelopeWithDefaults

`func NewPublicFederationProviderListEnvelopeWithDefaults() *PublicFederationProviderListEnvelope`

NewPublicFederationProviderListEnvelopeWithDefaults instantiates a new PublicFederationProviderListEnvelope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMeta

`func (o *PublicFederationProviderListEnvelope) GetMeta() Meta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *PublicFederationProviderListEnvelope) GetMetaOk() (*Meta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *PublicFederationProviderListEnvelope) SetMeta(v Meta)`

SetMeta sets Meta field to given value.


### GetData

`func (o *PublicFederationProviderListEnvelope) GetData() []PublicFederationProviderResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *PublicFederationProviderListEnvelope) GetDataOk() (*[]PublicFederationProviderResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *PublicFederationProviderListEnvelope) SetData(v []PublicFederationProviderResponse)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


