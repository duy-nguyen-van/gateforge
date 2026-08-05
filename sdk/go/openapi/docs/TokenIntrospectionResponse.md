# TokenIntrospectionResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | **bool** | Whether the token is currently active | 
**Scope** | Pointer to **string** |  | [optional] 
**ClientId** | Pointer to **string** |  | [optional] 
**Username** | Pointer to **string** |  | [optional] 
**TokenType** | Pointer to **string** | access_token or refresh_token when active | [optional] 
**Exp** | Pointer to **int64** |  | [optional] 
**Iat** | Pointer to **int64** |  | [optional] 
**Nbf** | Pointer to **int64** |  | [optional] 
**Sub** | Pointer to **string** |  | [optional] 
**Aud** | Pointer to **string** |  | [optional] 
**Iss** | Pointer to **string** |  | [optional] 
**Jti** | Pointer to **string** |  | [optional] 

## Methods

### NewTokenIntrospectionResponse

`func NewTokenIntrospectionResponse(active bool, ) *TokenIntrospectionResponse`

NewTokenIntrospectionResponse instantiates a new TokenIntrospectionResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTokenIntrospectionResponseWithDefaults

`func NewTokenIntrospectionResponseWithDefaults() *TokenIntrospectionResponse`

NewTokenIntrospectionResponseWithDefaults instantiates a new TokenIntrospectionResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActive

`func (o *TokenIntrospectionResponse) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *TokenIntrospectionResponse) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *TokenIntrospectionResponse) SetActive(v bool)`

SetActive sets Active field to given value.


### GetScope

`func (o *TokenIntrospectionResponse) GetScope() string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *TokenIntrospectionResponse) GetScopeOk() (*string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *TokenIntrospectionResponse) SetScope(v string)`

SetScope sets Scope field to given value.

### HasScope

`func (o *TokenIntrospectionResponse) HasScope() bool`

HasScope returns a boolean if a field has been set.

### GetClientId

`func (o *TokenIntrospectionResponse) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *TokenIntrospectionResponse) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientId

`func (o *TokenIntrospectionResponse) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *TokenIntrospectionResponse) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetUsername

`func (o *TokenIntrospectionResponse) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *TokenIntrospectionResponse) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *TokenIntrospectionResponse) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *TokenIntrospectionResponse) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

### GetTokenType

`func (o *TokenIntrospectionResponse) GetTokenType() string`

GetTokenType returns the TokenType field if non-nil, zero value otherwise.

### GetTokenTypeOk

`func (o *TokenIntrospectionResponse) GetTokenTypeOk() (*string, bool)`

GetTokenTypeOk returns a tuple with the TokenType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenType

`func (o *TokenIntrospectionResponse) SetTokenType(v string)`

SetTokenType sets TokenType field to given value.

### HasTokenType

`func (o *TokenIntrospectionResponse) HasTokenType() bool`

HasTokenType returns a boolean if a field has been set.

### GetExp

`func (o *TokenIntrospectionResponse) GetExp() int64`

GetExp returns the Exp field if non-nil, zero value otherwise.

### GetExpOk

`func (o *TokenIntrospectionResponse) GetExpOk() (*int64, bool)`

GetExpOk returns a tuple with the Exp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExp

`func (o *TokenIntrospectionResponse) SetExp(v int64)`

SetExp sets Exp field to given value.

### HasExp

`func (o *TokenIntrospectionResponse) HasExp() bool`

HasExp returns a boolean if a field has been set.

### GetIat

`func (o *TokenIntrospectionResponse) GetIat() int64`

GetIat returns the Iat field if non-nil, zero value otherwise.

### GetIatOk

`func (o *TokenIntrospectionResponse) GetIatOk() (*int64, bool)`

GetIatOk returns a tuple with the Iat field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIat

`func (o *TokenIntrospectionResponse) SetIat(v int64)`

SetIat sets Iat field to given value.

### HasIat

`func (o *TokenIntrospectionResponse) HasIat() bool`

HasIat returns a boolean if a field has been set.

### GetNbf

`func (o *TokenIntrospectionResponse) GetNbf() int64`

GetNbf returns the Nbf field if non-nil, zero value otherwise.

### GetNbfOk

`func (o *TokenIntrospectionResponse) GetNbfOk() (*int64, bool)`

GetNbfOk returns a tuple with the Nbf field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNbf

`func (o *TokenIntrospectionResponse) SetNbf(v int64)`

SetNbf sets Nbf field to given value.

### HasNbf

`func (o *TokenIntrospectionResponse) HasNbf() bool`

HasNbf returns a boolean if a field has been set.

### GetSub

`func (o *TokenIntrospectionResponse) GetSub() string`

GetSub returns the Sub field if non-nil, zero value otherwise.

### GetSubOk

`func (o *TokenIntrospectionResponse) GetSubOk() (*string, bool)`

GetSubOk returns a tuple with the Sub field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSub

`func (o *TokenIntrospectionResponse) SetSub(v string)`

SetSub sets Sub field to given value.

### HasSub

`func (o *TokenIntrospectionResponse) HasSub() bool`

HasSub returns a boolean if a field has been set.

### GetAud

`func (o *TokenIntrospectionResponse) GetAud() string`

GetAud returns the Aud field if non-nil, zero value otherwise.

### GetAudOk

`func (o *TokenIntrospectionResponse) GetAudOk() (*string, bool)`

GetAudOk returns a tuple with the Aud field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAud

`func (o *TokenIntrospectionResponse) SetAud(v string)`

SetAud sets Aud field to given value.

### HasAud

`func (o *TokenIntrospectionResponse) HasAud() bool`

HasAud returns a boolean if a field has been set.

### GetIss

`func (o *TokenIntrospectionResponse) GetIss() string`

GetIss returns the Iss field if non-nil, zero value otherwise.

### GetIssOk

`func (o *TokenIntrospectionResponse) GetIssOk() (*string, bool)`

GetIssOk returns a tuple with the Iss field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIss

`func (o *TokenIntrospectionResponse) SetIss(v string)`

SetIss sets Iss field to given value.

### HasIss

`func (o *TokenIntrospectionResponse) HasIss() bool`

HasIss returns a boolean if a field has been set.

### GetJti

`func (o *TokenIntrospectionResponse) GetJti() string`

GetJti returns the Jti field if non-nil, zero value otherwise.

### GetJtiOk

`func (o *TokenIntrospectionResponse) GetJtiOk() (*string, bool)`

GetJtiOk returns a tuple with the Jti field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJti

`func (o *TokenIntrospectionResponse) SetJti(v string)`

SetJti sets Jti field to given value.

### HasJti

`func (o *TokenIntrospectionResponse) HasJti() bool`

HasJti returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


