# TokenIntrospectionRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Token** | **string** | The access or refresh token to introspect | 
**TokenTypeHint** | Pointer to **string** | Optional hint — access_token or refresh_token | [optional] 
**ClientId** | Pointer to **string** | Confidential client id (when not using HTTP Basic) | [optional] 
**ClientSecret** | Pointer to **string** | Confidential client secret (when not using HTTP Basic) | [optional] 

## Methods

### NewTokenIntrospectionRequest

`func NewTokenIntrospectionRequest(token string, ) *TokenIntrospectionRequest`

NewTokenIntrospectionRequest instantiates a new TokenIntrospectionRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTokenIntrospectionRequestWithDefaults

`func NewTokenIntrospectionRequestWithDefaults() *TokenIntrospectionRequest`

NewTokenIntrospectionRequestWithDefaults instantiates a new TokenIntrospectionRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetToken

`func (o *TokenIntrospectionRequest) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *TokenIntrospectionRequest) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *TokenIntrospectionRequest) SetToken(v string)`

SetToken sets Token field to given value.


### GetTokenTypeHint

`func (o *TokenIntrospectionRequest) GetTokenTypeHint() string`

GetTokenTypeHint returns the TokenTypeHint field if non-nil, zero value otherwise.

### GetTokenTypeHintOk

`func (o *TokenIntrospectionRequest) GetTokenTypeHintOk() (*string, bool)`

GetTokenTypeHintOk returns a tuple with the TokenTypeHint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenTypeHint

`func (o *TokenIntrospectionRequest) SetTokenTypeHint(v string)`

SetTokenTypeHint sets TokenTypeHint field to given value.

### HasTokenTypeHint

`func (o *TokenIntrospectionRequest) HasTokenTypeHint() bool`

HasTokenTypeHint returns a boolean if a field has been set.

### GetClientId

`func (o *TokenIntrospectionRequest) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *TokenIntrospectionRequest) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientId

`func (o *TokenIntrospectionRequest) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *TokenIntrospectionRequest) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetClientSecret

`func (o *TokenIntrospectionRequest) GetClientSecret() string`

GetClientSecret returns the ClientSecret field if non-nil, zero value otherwise.

### GetClientSecretOk

`func (o *TokenIntrospectionRequest) GetClientSecretOk() (*string, bool)`

GetClientSecretOk returns a tuple with the ClientSecret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientSecret

`func (o *TokenIntrospectionRequest) SetClientSecret(v string)`

SetClientSecret sets ClientSecret field to given value.

### HasClientSecret

`func (o *TokenIntrospectionRequest) HasClientSecret() bool`

HasClientSecret returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


