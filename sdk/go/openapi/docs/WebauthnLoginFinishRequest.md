# WebauthnLoginFinishRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** |  | 
**SessionToken** | **string** |  | 
**Credential** | Pointer to **map[string]interface{}** | PublicKeyCredential JSON from navigator.credentials.get() | [optional] 
**TenantId** | Pointer to **string** |  | [optional] 
**RememberMe** | Pointer to **bool** |  | [optional] 
**ReturnTo** | Pointer to **string** |  | [optional] 

## Methods

### NewWebauthnLoginFinishRequest

`func NewWebauthnLoginFinishRequest(email string, sessionToken string, ) *WebauthnLoginFinishRequest`

NewWebauthnLoginFinishRequest instantiates a new WebauthnLoginFinishRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebauthnLoginFinishRequestWithDefaults

`func NewWebauthnLoginFinishRequestWithDefaults() *WebauthnLoginFinishRequest`

NewWebauthnLoginFinishRequestWithDefaults instantiates a new WebauthnLoginFinishRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *WebauthnLoginFinishRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *WebauthnLoginFinishRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *WebauthnLoginFinishRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetSessionToken

`func (o *WebauthnLoginFinishRequest) GetSessionToken() string`

GetSessionToken returns the SessionToken field if non-nil, zero value otherwise.

### GetSessionTokenOk

`func (o *WebauthnLoginFinishRequest) GetSessionTokenOk() (*string, bool)`

GetSessionTokenOk returns a tuple with the SessionToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionToken

`func (o *WebauthnLoginFinishRequest) SetSessionToken(v string)`

SetSessionToken sets SessionToken field to given value.


### GetCredential

`func (o *WebauthnLoginFinishRequest) GetCredential() map[string]interface{}`

GetCredential returns the Credential field if non-nil, zero value otherwise.

### GetCredentialOk

`func (o *WebauthnLoginFinishRequest) GetCredentialOk() (*map[string]interface{}, bool)`

GetCredentialOk returns a tuple with the Credential field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredential

`func (o *WebauthnLoginFinishRequest) SetCredential(v map[string]interface{})`

SetCredential sets Credential field to given value.

### HasCredential

`func (o *WebauthnLoginFinishRequest) HasCredential() bool`

HasCredential returns a boolean if a field has been set.

### GetTenantId

`func (o *WebauthnLoginFinishRequest) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *WebauthnLoginFinishRequest) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *WebauthnLoginFinishRequest) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *WebauthnLoginFinishRequest) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetRememberMe

`func (o *WebauthnLoginFinishRequest) GetRememberMe() bool`

GetRememberMe returns the RememberMe field if non-nil, zero value otherwise.

### GetRememberMeOk

`func (o *WebauthnLoginFinishRequest) GetRememberMeOk() (*bool, bool)`

GetRememberMeOk returns a tuple with the RememberMe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRememberMe

`func (o *WebauthnLoginFinishRequest) SetRememberMe(v bool)`

SetRememberMe sets RememberMe field to given value.

### HasRememberMe

`func (o *WebauthnLoginFinishRequest) HasRememberMe() bool`

HasRememberMe returns a boolean if a field has been set.

### GetReturnTo

`func (o *WebauthnLoginFinishRequest) GetReturnTo() string`

GetReturnTo returns the ReturnTo field if non-nil, zero value otherwise.

### GetReturnToOk

`func (o *WebauthnLoginFinishRequest) GetReturnToOk() (*string, bool)`

GetReturnToOk returns a tuple with the ReturnTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReturnTo

`func (o *WebauthnLoginFinishRequest) SetReturnTo(v string)`

SetReturnTo sets ReturnTo field to given value.

### HasReturnTo

`func (o *WebauthnLoginFinishRequest) HasReturnTo() bool`

HasReturnTo returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


