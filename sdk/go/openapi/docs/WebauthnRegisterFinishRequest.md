# WebauthnRegisterFinishRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SessionToken** | **string** |  | 
**Credential** | Pointer to **map[string]interface{}** | PublicKeyCredential JSON from navigator.credentials.create() | [optional] 

## Methods

### NewWebauthnRegisterFinishRequest

`func NewWebauthnRegisterFinishRequest(sessionToken string, ) *WebauthnRegisterFinishRequest`

NewWebauthnRegisterFinishRequest instantiates a new WebauthnRegisterFinishRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebauthnRegisterFinishRequestWithDefaults

`func NewWebauthnRegisterFinishRequestWithDefaults() *WebauthnRegisterFinishRequest`

NewWebauthnRegisterFinishRequestWithDefaults instantiates a new WebauthnRegisterFinishRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSessionToken

`func (o *WebauthnRegisterFinishRequest) GetSessionToken() string`

GetSessionToken returns the SessionToken field if non-nil, zero value otherwise.

### GetSessionTokenOk

`func (o *WebauthnRegisterFinishRequest) GetSessionTokenOk() (*string, bool)`

GetSessionTokenOk returns a tuple with the SessionToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionToken

`func (o *WebauthnRegisterFinishRequest) SetSessionToken(v string)`

SetSessionToken sets SessionToken field to given value.


### GetCredential

`func (o *WebauthnRegisterFinishRequest) GetCredential() map[string]interface{}`

GetCredential returns the Credential field if non-nil, zero value otherwise.

### GetCredentialOk

`func (o *WebauthnRegisterFinishRequest) GetCredentialOk() (*map[string]interface{}, bool)`

GetCredentialOk returns a tuple with the Credential field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredential

`func (o *WebauthnRegisterFinishRequest) SetCredential(v map[string]interface{})`

SetCredential sets Credential field to given value.

### HasCredential

`func (o *WebauthnRegisterFinishRequest) HasCredential() bool`

HasCredential returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


