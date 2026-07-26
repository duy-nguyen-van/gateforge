# LoginResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessToken** | Pointer to **string** |  | [optional] 
**RefreshToken** | Pointer to **string** |  | [optional] 
**TokenType** | Pointer to **string** |  | [optional] 
**ExpiresIn** | **int32** |  | 
**RefreshExpiresIn** | Pointer to **int32** | Refresh token lifetime in seconds | [optional] 
**ActiveTenantId** | Pointer to **string** |  | [optional] 
**MfaRequired** | **bool** |  | 
**MfaTicket** | **string** |  | 
**SelectionRequired** | **bool** |  | 
**Tenants** | [**[]TenantSummary**](TenantSummary.md) |  | 
**SelectionToken** | **string** |  | 

## Methods

### NewLoginResult

`func NewLoginResult(expiresIn int32, mfaRequired bool, mfaTicket string, selectionRequired bool, tenants []TenantSummary, selectionToken string, ) *LoginResult`

NewLoginResult instantiates a new LoginResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginResultWithDefaults

`func NewLoginResultWithDefaults() *LoginResult`

NewLoginResultWithDefaults instantiates a new LoginResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessToken

`func (o *LoginResult) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *LoginResult) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *LoginResult) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.

### HasAccessToken

`func (o *LoginResult) HasAccessToken() bool`

HasAccessToken returns a boolean if a field has been set.

### GetRefreshToken

`func (o *LoginResult) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *LoginResult) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *LoginResult) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.

### HasRefreshToken

`func (o *LoginResult) HasRefreshToken() bool`

HasRefreshToken returns a boolean if a field has been set.

### GetTokenType

`func (o *LoginResult) GetTokenType() string`

GetTokenType returns the TokenType field if non-nil, zero value otherwise.

### GetTokenTypeOk

`func (o *LoginResult) GetTokenTypeOk() (*string, bool)`

GetTokenTypeOk returns a tuple with the TokenType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenType

`func (o *LoginResult) SetTokenType(v string)`

SetTokenType sets TokenType field to given value.

### HasTokenType

`func (o *LoginResult) HasTokenType() bool`

HasTokenType returns a boolean if a field has been set.

### GetExpiresIn

`func (o *LoginResult) GetExpiresIn() int32`

GetExpiresIn returns the ExpiresIn field if non-nil, zero value otherwise.

### GetExpiresInOk

`func (o *LoginResult) GetExpiresInOk() (*int32, bool)`

GetExpiresInOk returns a tuple with the ExpiresIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresIn

`func (o *LoginResult) SetExpiresIn(v int32)`

SetExpiresIn sets ExpiresIn field to given value.


### GetRefreshExpiresIn

`func (o *LoginResult) GetRefreshExpiresIn() int32`

GetRefreshExpiresIn returns the RefreshExpiresIn field if non-nil, zero value otherwise.

### GetRefreshExpiresInOk

`func (o *LoginResult) GetRefreshExpiresInOk() (*int32, bool)`

GetRefreshExpiresInOk returns a tuple with the RefreshExpiresIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshExpiresIn

`func (o *LoginResult) SetRefreshExpiresIn(v int32)`

SetRefreshExpiresIn sets RefreshExpiresIn field to given value.

### HasRefreshExpiresIn

`func (o *LoginResult) HasRefreshExpiresIn() bool`

HasRefreshExpiresIn returns a boolean if a field has been set.

### GetActiveTenantId

`func (o *LoginResult) GetActiveTenantId() string`

GetActiveTenantId returns the ActiveTenantId field if non-nil, zero value otherwise.

### GetActiveTenantIdOk

`func (o *LoginResult) GetActiveTenantIdOk() (*string, bool)`

GetActiveTenantIdOk returns a tuple with the ActiveTenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActiveTenantId

`func (o *LoginResult) SetActiveTenantId(v string)`

SetActiveTenantId sets ActiveTenantId field to given value.

### HasActiveTenantId

`func (o *LoginResult) HasActiveTenantId() bool`

HasActiveTenantId returns a boolean if a field has been set.

### GetMfaRequired

`func (o *LoginResult) GetMfaRequired() bool`

GetMfaRequired returns the MfaRequired field if non-nil, zero value otherwise.

### GetMfaRequiredOk

`func (o *LoginResult) GetMfaRequiredOk() (*bool, bool)`

GetMfaRequiredOk returns a tuple with the MfaRequired field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaRequired

`func (o *LoginResult) SetMfaRequired(v bool)`

SetMfaRequired sets MfaRequired field to given value.


### GetMfaTicket

`func (o *LoginResult) GetMfaTicket() string`

GetMfaTicket returns the MfaTicket field if non-nil, zero value otherwise.

### GetMfaTicketOk

`func (o *LoginResult) GetMfaTicketOk() (*string, bool)`

GetMfaTicketOk returns a tuple with the MfaTicket field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaTicket

`func (o *LoginResult) SetMfaTicket(v string)`

SetMfaTicket sets MfaTicket field to given value.


### GetSelectionRequired

`func (o *LoginResult) GetSelectionRequired() bool`

GetSelectionRequired returns the SelectionRequired field if non-nil, zero value otherwise.

### GetSelectionRequiredOk

`func (o *LoginResult) GetSelectionRequiredOk() (*bool, bool)`

GetSelectionRequiredOk returns a tuple with the SelectionRequired field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSelectionRequired

`func (o *LoginResult) SetSelectionRequired(v bool)`

SetSelectionRequired sets SelectionRequired field to given value.


### GetTenants

`func (o *LoginResult) GetTenants() []TenantSummary`

GetTenants returns the Tenants field if non-nil, zero value otherwise.

### GetTenantsOk

`func (o *LoginResult) GetTenantsOk() (*[]TenantSummary, bool)`

GetTenantsOk returns a tuple with the Tenants field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenants

`func (o *LoginResult) SetTenants(v []TenantSummary)`

SetTenants sets Tenants field to given value.


### GetSelectionToken

`func (o *LoginResult) GetSelectionToken() string`

GetSelectionToken returns the SelectionToken field if non-nil, zero value otherwise.

### GetSelectionTokenOk

`func (o *LoginResult) GetSelectionTokenOk() (*string, bool)`

GetSelectionTokenOk returns a tuple with the SelectionToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSelectionToken

`func (o *LoginResult) SetSelectionToken(v string)`

SetSelectionToken sets SelectionToken field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


