# AdminClientUsageResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientId** | Pointer to **string** |  | [optional] 
**ActiveRefreshTokens** | Pointer to **int32** |  | [optional] 
**TotalRefreshTokens** | Pointer to **int32** |  | [optional] 
**AuthorizeEvents30d** | Pointer to **int32** |  | [optional] 
**TokenIssueEvents30d** | Pointer to **int32** |  | [optional] 
**LastTokenIssuedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewAdminClientUsageResponse

`func NewAdminClientUsageResponse() *AdminClientUsageResponse`

NewAdminClientUsageResponse instantiates a new AdminClientUsageResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminClientUsageResponseWithDefaults

`func NewAdminClientUsageResponseWithDefaults() *AdminClientUsageResponse`

NewAdminClientUsageResponseWithDefaults instantiates a new AdminClientUsageResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClientId

`func (o *AdminClientUsageResponse) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *AdminClientUsageResponse) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientId

`func (o *AdminClientUsageResponse) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *AdminClientUsageResponse) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetActiveRefreshTokens

`func (o *AdminClientUsageResponse) GetActiveRefreshTokens() int32`

GetActiveRefreshTokens returns the ActiveRefreshTokens field if non-nil, zero value otherwise.

### GetActiveRefreshTokensOk

`func (o *AdminClientUsageResponse) GetActiveRefreshTokensOk() (*int32, bool)`

GetActiveRefreshTokensOk returns a tuple with the ActiveRefreshTokens field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActiveRefreshTokens

`func (o *AdminClientUsageResponse) SetActiveRefreshTokens(v int32)`

SetActiveRefreshTokens sets ActiveRefreshTokens field to given value.

### HasActiveRefreshTokens

`func (o *AdminClientUsageResponse) HasActiveRefreshTokens() bool`

HasActiveRefreshTokens returns a boolean if a field has been set.

### GetTotalRefreshTokens

`func (o *AdminClientUsageResponse) GetTotalRefreshTokens() int32`

GetTotalRefreshTokens returns the TotalRefreshTokens field if non-nil, zero value otherwise.

### GetTotalRefreshTokensOk

`func (o *AdminClientUsageResponse) GetTotalRefreshTokensOk() (*int32, bool)`

GetTotalRefreshTokensOk returns a tuple with the TotalRefreshTokens field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalRefreshTokens

`func (o *AdminClientUsageResponse) SetTotalRefreshTokens(v int32)`

SetTotalRefreshTokens sets TotalRefreshTokens field to given value.

### HasTotalRefreshTokens

`func (o *AdminClientUsageResponse) HasTotalRefreshTokens() bool`

HasTotalRefreshTokens returns a boolean if a field has been set.

### GetAuthorizeEvents30d

`func (o *AdminClientUsageResponse) GetAuthorizeEvents30d() int32`

GetAuthorizeEvents30d returns the AuthorizeEvents30d field if non-nil, zero value otherwise.

### GetAuthorizeEvents30dOk

`func (o *AdminClientUsageResponse) GetAuthorizeEvents30dOk() (*int32, bool)`

GetAuthorizeEvents30dOk returns a tuple with the AuthorizeEvents30d field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorizeEvents30d

`func (o *AdminClientUsageResponse) SetAuthorizeEvents30d(v int32)`

SetAuthorizeEvents30d sets AuthorizeEvents30d field to given value.

### HasAuthorizeEvents30d

`func (o *AdminClientUsageResponse) HasAuthorizeEvents30d() bool`

HasAuthorizeEvents30d returns a boolean if a field has been set.

### GetTokenIssueEvents30d

`func (o *AdminClientUsageResponse) GetTokenIssueEvents30d() int32`

GetTokenIssueEvents30d returns the TokenIssueEvents30d field if non-nil, zero value otherwise.

### GetTokenIssueEvents30dOk

`func (o *AdminClientUsageResponse) GetTokenIssueEvents30dOk() (*int32, bool)`

GetTokenIssueEvents30dOk returns a tuple with the TokenIssueEvents30d field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenIssueEvents30d

`func (o *AdminClientUsageResponse) SetTokenIssueEvents30d(v int32)`

SetTokenIssueEvents30d sets TokenIssueEvents30d field to given value.

### HasTokenIssueEvents30d

`func (o *AdminClientUsageResponse) HasTokenIssueEvents30d() bool`

HasTokenIssueEvents30d returns a boolean if a field has been set.

### GetLastTokenIssuedAt

`func (o *AdminClientUsageResponse) GetLastTokenIssuedAt() time.Time`

GetLastTokenIssuedAt returns the LastTokenIssuedAt field if non-nil, zero value otherwise.

### GetLastTokenIssuedAtOk

`func (o *AdminClientUsageResponse) GetLastTokenIssuedAtOk() (*time.Time, bool)`

GetLastTokenIssuedAtOk returns a tuple with the LastTokenIssuedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastTokenIssuedAt

`func (o *AdminClientUsageResponse) SetLastTokenIssuedAt(v time.Time)`

SetLastTokenIssuedAt sets LastTokenIssuedAt field to given value.

### HasLastTokenIssuedAt

`func (o *AdminClientUsageResponse) HasLastTokenIssuedAt() bool`

HasLastTokenIssuedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


