# AdminIdentityProviderResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TenantId** | Pointer to **string** |  | [optional] 
**Provider** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Configured** | Pointer to **bool** |  | [optional] 
**OauthClientId** | Pointer to **string** |  | [optional] 
**OauthClientSecretSet** | Pointer to **bool** |  | [optional] 
**RedirectUri** | Pointer to **string** |  | [optional] 
**SetupConsoleUrl** | Pointer to **string** |  | [optional] 

## Methods

### NewAdminIdentityProviderResponse

`func NewAdminIdentityProviderResponse() *AdminIdentityProviderResponse`

NewAdminIdentityProviderResponse instantiates a new AdminIdentityProviderResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminIdentityProviderResponseWithDefaults

`func NewAdminIdentityProviderResponseWithDefaults() *AdminIdentityProviderResponse`

NewAdminIdentityProviderResponseWithDefaults instantiates a new AdminIdentityProviderResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTenantId

`func (o *AdminIdentityProviderResponse) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *AdminIdentityProviderResponse) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *AdminIdentityProviderResponse) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *AdminIdentityProviderResponse) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetProvider

`func (o *AdminIdentityProviderResponse) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *AdminIdentityProviderResponse) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *AdminIdentityProviderResponse) SetProvider(v string)`

SetProvider sets Provider field to given value.

### HasProvider

`func (o *AdminIdentityProviderResponse) HasProvider() bool`

HasProvider returns a boolean if a field has been set.

### GetName

`func (o *AdminIdentityProviderResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AdminIdentityProviderResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AdminIdentityProviderResponse) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *AdminIdentityProviderResponse) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnabled

`func (o *AdminIdentityProviderResponse) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *AdminIdentityProviderResponse) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *AdminIdentityProviderResponse) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *AdminIdentityProviderResponse) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetConfigured

`func (o *AdminIdentityProviderResponse) GetConfigured() bool`

GetConfigured returns the Configured field if non-nil, zero value otherwise.

### GetConfiguredOk

`func (o *AdminIdentityProviderResponse) GetConfiguredOk() (*bool, bool)`

GetConfiguredOk returns a tuple with the Configured field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfigured

`func (o *AdminIdentityProviderResponse) SetConfigured(v bool)`

SetConfigured sets Configured field to given value.

### HasConfigured

`func (o *AdminIdentityProviderResponse) HasConfigured() bool`

HasConfigured returns a boolean if a field has been set.

### GetOauthClientId

`func (o *AdminIdentityProviderResponse) GetOauthClientId() string`

GetOauthClientId returns the OauthClientId field if non-nil, zero value otherwise.

### GetOauthClientIdOk

`func (o *AdminIdentityProviderResponse) GetOauthClientIdOk() (*string, bool)`

GetOauthClientIdOk returns a tuple with the OauthClientId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOauthClientId

`func (o *AdminIdentityProviderResponse) SetOauthClientId(v string)`

SetOauthClientId sets OauthClientId field to given value.

### HasOauthClientId

`func (o *AdminIdentityProviderResponse) HasOauthClientId() bool`

HasOauthClientId returns a boolean if a field has been set.

### GetOauthClientSecretSet

`func (o *AdminIdentityProviderResponse) GetOauthClientSecretSet() bool`

GetOauthClientSecretSet returns the OauthClientSecretSet field if non-nil, zero value otherwise.

### GetOauthClientSecretSetOk

`func (o *AdminIdentityProviderResponse) GetOauthClientSecretSetOk() (*bool, bool)`

GetOauthClientSecretSetOk returns a tuple with the OauthClientSecretSet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOauthClientSecretSet

`func (o *AdminIdentityProviderResponse) SetOauthClientSecretSet(v bool)`

SetOauthClientSecretSet sets OauthClientSecretSet field to given value.

### HasOauthClientSecretSet

`func (o *AdminIdentityProviderResponse) HasOauthClientSecretSet() bool`

HasOauthClientSecretSet returns a boolean if a field has been set.

### GetRedirectUri

`func (o *AdminIdentityProviderResponse) GetRedirectUri() string`

GetRedirectUri returns the RedirectUri field if non-nil, zero value otherwise.

### GetRedirectUriOk

`func (o *AdminIdentityProviderResponse) GetRedirectUriOk() (*string, bool)`

GetRedirectUriOk returns a tuple with the RedirectUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectUri

`func (o *AdminIdentityProviderResponse) SetRedirectUri(v string)`

SetRedirectUri sets RedirectUri field to given value.

### HasRedirectUri

`func (o *AdminIdentityProviderResponse) HasRedirectUri() bool`

HasRedirectUri returns a boolean if a field has been set.

### GetSetupConsoleUrl

`func (o *AdminIdentityProviderResponse) GetSetupConsoleUrl() string`

GetSetupConsoleUrl returns the SetupConsoleUrl field if non-nil, zero value otherwise.

### GetSetupConsoleUrlOk

`func (o *AdminIdentityProviderResponse) GetSetupConsoleUrlOk() (*string, bool)`

GetSetupConsoleUrlOk returns a tuple with the SetupConsoleUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSetupConsoleUrl

`func (o *AdminIdentityProviderResponse) SetSetupConsoleUrl(v string)`

SetSetupConsoleUrl sets SetupConsoleUrl field to given value.

### HasSetupConsoleUrl

`func (o *AdminIdentityProviderResponse) HasSetupConsoleUrl() bool`

HasSetupConsoleUrl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


