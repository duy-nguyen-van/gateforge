# AdminClientResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**ClientId** | Pointer to **string** |  | [optional] 
**ClientSecretSet** | Pointer to **bool** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**TenantId** | Pointer to **string** |  | [optional] 
**IsPublic** | Pointer to **bool** |  | [optional] 
**RedirectUris** | Pointer to **[]string** |  | [optional] 
**GrantTypes** | Pointer to **[]string** |  | [optional] 
**Scopes** | Pointer to **[]string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewAdminClientResponse

`func NewAdminClientResponse() *AdminClientResponse`

NewAdminClientResponse instantiates a new AdminClientResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminClientResponseWithDefaults

`func NewAdminClientResponseWithDefaults() *AdminClientResponse`

NewAdminClientResponseWithDefaults instantiates a new AdminClientResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AdminClientResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AdminClientResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AdminClientResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *AdminClientResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetClientId

`func (o *AdminClientResponse) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *AdminClientResponse) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientId

`func (o *AdminClientResponse) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *AdminClientResponse) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetClientSecretSet

`func (o *AdminClientResponse) GetClientSecretSet() bool`

GetClientSecretSet returns the ClientSecretSet field if non-nil, zero value otherwise.

### GetClientSecretSetOk

`func (o *AdminClientResponse) GetClientSecretSetOk() (*bool, bool)`

GetClientSecretSetOk returns a tuple with the ClientSecretSet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientSecretSet

`func (o *AdminClientResponse) SetClientSecretSet(v bool)`

SetClientSecretSet sets ClientSecretSet field to given value.

### HasClientSecretSet

`func (o *AdminClientResponse) HasClientSecretSet() bool`

HasClientSecretSet returns a boolean if a field has been set.

### GetName

`func (o *AdminClientResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AdminClientResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AdminClientResponse) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *AdminClientResponse) HasName() bool`

HasName returns a boolean if a field has been set.

### GetTenantId

`func (o *AdminClientResponse) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *AdminClientResponse) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *AdminClientResponse) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *AdminClientResponse) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetIsPublic

`func (o *AdminClientResponse) GetIsPublic() bool`

GetIsPublic returns the IsPublic field if non-nil, zero value otherwise.

### GetIsPublicOk

`func (o *AdminClientResponse) GetIsPublicOk() (*bool, bool)`

GetIsPublicOk returns a tuple with the IsPublic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPublic

`func (o *AdminClientResponse) SetIsPublic(v bool)`

SetIsPublic sets IsPublic field to given value.

### HasIsPublic

`func (o *AdminClientResponse) HasIsPublic() bool`

HasIsPublic returns a boolean if a field has been set.

### GetRedirectUris

`func (o *AdminClientResponse) GetRedirectUris() []string`

GetRedirectUris returns the RedirectUris field if non-nil, zero value otherwise.

### GetRedirectUrisOk

`func (o *AdminClientResponse) GetRedirectUrisOk() (*[]string, bool)`

GetRedirectUrisOk returns a tuple with the RedirectUris field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectUris

`func (o *AdminClientResponse) SetRedirectUris(v []string)`

SetRedirectUris sets RedirectUris field to given value.

### HasRedirectUris

`func (o *AdminClientResponse) HasRedirectUris() bool`

HasRedirectUris returns a boolean if a field has been set.

### GetGrantTypes

`func (o *AdminClientResponse) GetGrantTypes() []string`

GetGrantTypes returns the GrantTypes field if non-nil, zero value otherwise.

### GetGrantTypesOk

`func (o *AdminClientResponse) GetGrantTypesOk() (*[]string, bool)`

GetGrantTypesOk returns a tuple with the GrantTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantTypes

`func (o *AdminClientResponse) SetGrantTypes(v []string)`

SetGrantTypes sets GrantTypes field to given value.

### HasGrantTypes

`func (o *AdminClientResponse) HasGrantTypes() bool`

HasGrantTypes returns a boolean if a field has been set.

### GetScopes

`func (o *AdminClientResponse) GetScopes() []string`

GetScopes returns the Scopes field if non-nil, zero value otherwise.

### GetScopesOk

`func (o *AdminClientResponse) GetScopesOk() (*[]string, bool)`

GetScopesOk returns a tuple with the Scopes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScopes

`func (o *AdminClientResponse) SetScopes(v []string)`

SetScopes sets Scopes field to given value.

### HasScopes

`func (o *AdminClientResponse) HasScopes() bool`

HasScopes returns a boolean if a field has been set.

### GetCreatedAt

`func (o *AdminClientResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *AdminClientResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *AdminClientResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *AdminClientResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


