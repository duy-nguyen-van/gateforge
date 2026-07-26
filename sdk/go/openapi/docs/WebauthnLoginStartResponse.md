# WebauthnLoginStartResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Options** | Pointer to **map[string]interface{}** |  | [optional] 
**SessionToken** | Pointer to **string** |  | [optional] 

## Methods

### NewWebauthnLoginStartResponse

`func NewWebauthnLoginStartResponse() *WebauthnLoginStartResponse`

NewWebauthnLoginStartResponse instantiates a new WebauthnLoginStartResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebauthnLoginStartResponseWithDefaults

`func NewWebauthnLoginStartResponseWithDefaults() *WebauthnLoginStartResponse`

NewWebauthnLoginStartResponseWithDefaults instantiates a new WebauthnLoginStartResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOptions

`func (o *WebauthnLoginStartResponse) GetOptions() map[string]interface{}`

GetOptions returns the Options field if non-nil, zero value otherwise.

### GetOptionsOk

`func (o *WebauthnLoginStartResponse) GetOptionsOk() (*map[string]interface{}, bool)`

GetOptionsOk returns a tuple with the Options field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOptions

`func (o *WebauthnLoginStartResponse) SetOptions(v map[string]interface{})`

SetOptions sets Options field to given value.

### HasOptions

`func (o *WebauthnLoginStartResponse) HasOptions() bool`

HasOptions returns a boolean if a field has been set.

### GetSessionToken

`func (o *WebauthnLoginStartResponse) GetSessionToken() string`

GetSessionToken returns the SessionToken field if non-nil, zero value otherwise.

### GetSessionTokenOk

`func (o *WebauthnLoginStartResponse) GetSessionTokenOk() (*string, bool)`

GetSessionTokenOk returns a tuple with the SessionToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionToken

`func (o *WebauthnLoginStartResponse) SetSessionToken(v string)`

SetSessionToken sets SessionToken field to given value.

### HasSessionToken

`func (o *WebauthnLoginStartResponse) HasSessionToken() bool`

HasSessionToken returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


