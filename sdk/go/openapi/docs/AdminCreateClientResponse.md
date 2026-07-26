# AdminCreateClientResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**ClientId** | Pointer to **string** |  | [optional] 
**ClientSecret** | Pointer to **string** |  | [optional] 
**ClientSecretSet** | Pointer to **bool** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**TenantId** | Pointer to **string** |  | [optional] 
**IsPublic** | Pointer to **bool** |  | [optional] 
**RedirectUris** | Pointer to **[]string** |  | [optional] 
**GrantTypes** | Pointer to **[]string** |  | [optional] 
**Scopes** | Pointer to **[]string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewAdminCreateClientResponse

`func NewAdminCreateClientResponse() *AdminCreateClientResponse`

NewAdminCreateClientResponse instantiates a new AdminCreateClientResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminCreateClientResponseWithDefaults

`func NewAdminCreateClientResponseWithDefaults() *AdminCreateClientResponse`

NewAdminCreateClientResponseWithDefaults instantiates a new AdminCreateClientResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AdminCreateClientResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AdminCreateClientResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AdminCreateClientResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *AdminCreateClientResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetClientId

`func (o *AdminCreateClientResponse) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *AdminCreateClientResponse) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientId

`func (o *AdminCreateClientResponse) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *AdminCreateClientResponse) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetClientSecret

`func (o *AdminCreateClientResponse) GetClientSecret() string`

GetClientSecret returns the ClientSecret field if non-nil, zero value otherwise.

### GetClientSecretOk

`func (o *AdminCreateClientResponse) GetClientSecretOk() (*string, bool)`

GetClientSecretOk returns a tuple with the ClientSecret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientSecret

`func (o *AdminCreateClientResponse) SetClientSecret(v string)`

SetClientSecret sets ClientSecret field to given value.

### HasClientSecret

`func (o *AdminCreateClientResponse) HasClientSecret() bool`

HasClientSecret returns a boolean if a field has been set.

### GetClientSecretSet

`func (o *AdminCreateClientResponse) GetClientSecretSet() bool`

GetClientSecretSet returns the ClientSecretSet field if non-nil, zero value otherwise.

### GetClientSecretSetOk

`func (o *AdminCreateClientResponse) GetClientSecretSetOk() (*bool, bool)`

GetClientSecretSetOk returns a tuple with the ClientSecretSet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientSecretSet

`func (o *AdminCreateClientResponse) SetClientSecretSet(v bool)`

SetClientSecretSet sets ClientSecretSet field to given value.

### HasClientSecretSet

`func (o *AdminCreateClientResponse) HasClientSecretSet() bool`

HasClientSecretSet returns a boolean if a field has been set.

### GetName

`func (o *AdminCreateClientResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AdminCreateClientResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AdminCreateClientResponse) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *AdminCreateClientResponse) HasName() bool`

HasName returns a boolean if a field has been set.

### GetTenantId

`func (o *AdminCreateClientResponse) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *AdminCreateClientResponse) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *AdminCreateClientResponse) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *AdminCreateClientResponse) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetIsPublic

`func (o *AdminCreateClientResponse) GetIsPublic() bool`

GetIsPublic returns the IsPublic field if non-nil, zero value otherwise.

### GetIsPublicOk

`func (o *AdminCreateClientResponse) GetIsPublicOk() (*bool, bool)`

GetIsPublicOk returns a tuple with the IsPublic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPublic

`func (o *AdminCreateClientResponse) SetIsPublic(v bool)`

SetIsPublic sets IsPublic field to given value.

### HasIsPublic

`func (o *AdminCreateClientResponse) HasIsPublic() bool`

HasIsPublic returns a boolean if a field has been set.

### GetRedirectUris

`func (o *AdminCreateClientResponse) GetRedirectUris() []string`

GetRedirectUris returns the RedirectUris field if non-nil, zero value otherwise.

### GetRedirectUrisOk

`func (o *AdminCreateClientResponse) GetRedirectUrisOk() (*[]string, bool)`

GetRedirectUrisOk returns a tuple with the RedirectUris field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectUris

`func (o *AdminCreateClientResponse) SetRedirectUris(v []string)`

SetRedirectUris sets RedirectUris field to given value.

### HasRedirectUris

`func (o *AdminCreateClientResponse) HasRedirectUris() bool`

HasRedirectUris returns a boolean if a field has been set.

### GetGrantTypes

`func (o *AdminCreateClientResponse) GetGrantTypes() []string`

GetGrantTypes returns the GrantTypes field if non-nil, zero value otherwise.

### GetGrantTypesOk

`func (o *AdminCreateClientResponse) GetGrantTypesOk() (*[]string, bool)`

GetGrantTypesOk returns a tuple with the GrantTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantTypes

`func (o *AdminCreateClientResponse) SetGrantTypes(v []string)`

SetGrantTypes sets GrantTypes field to given value.

### HasGrantTypes

`func (o *AdminCreateClientResponse) HasGrantTypes() bool`

HasGrantTypes returns a boolean if a field has been set.

### GetScopes

`func (o *AdminCreateClientResponse) GetScopes() []string`

GetScopes returns the Scopes field if non-nil, zero value otherwise.

### GetScopesOk

`func (o *AdminCreateClientResponse) GetScopesOk() (*[]string, bool)`

GetScopesOk returns a tuple with the Scopes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScopes

`func (o *AdminCreateClientResponse) SetScopes(v []string)`

SetScopes sets Scopes field to given value.

### HasScopes

`func (o *AdminCreateClientResponse) HasScopes() bool`

HasScopes returns a boolean if a field has been set.

### GetCreatedAt

`func (o *AdminCreateClientResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *AdminCreateClientResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *AdminCreateClientResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *AdminCreateClientResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


