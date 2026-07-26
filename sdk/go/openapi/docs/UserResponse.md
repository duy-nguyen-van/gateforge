# UserResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 
**FirstName** | Pointer to **string** |  | [optional] 
**LastName** | Pointer to **string** |  | [optional] 
**EmailVerified** | Pointer to **bool** |  | [optional] 
**IsPlatformAdmin** | Pointer to **bool** |  | [optional] 
**MfaEnabled** | Pointer to **bool** |  | [optional] 
**ActiveTenantId** | Pointer to **string** |  | [optional] 
**Tenants** | Pointer to [**[]TenantSummary**](TenantSummary.md) |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewUserResponse

`func NewUserResponse() *UserResponse`

NewUserResponse instantiates a new UserResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserResponseWithDefaults

`func NewUserResponseWithDefaults() *UserResponse`

NewUserResponseWithDefaults instantiates a new UserResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UserResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UserResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UserResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *UserResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetEmail

`func (o *UserResponse) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *UserResponse) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *UserResponse) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *UserResponse) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetFirstName

`func (o *UserResponse) GetFirstName() string`

GetFirstName returns the FirstName field if non-nil, zero value otherwise.

### GetFirstNameOk

`func (o *UserResponse) GetFirstNameOk() (*string, bool)`

GetFirstNameOk returns a tuple with the FirstName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstName

`func (o *UserResponse) SetFirstName(v string)`

SetFirstName sets FirstName field to given value.

### HasFirstName

`func (o *UserResponse) HasFirstName() bool`

HasFirstName returns a boolean if a field has been set.

### GetLastName

`func (o *UserResponse) GetLastName() string`

GetLastName returns the LastName field if non-nil, zero value otherwise.

### GetLastNameOk

`func (o *UserResponse) GetLastNameOk() (*string, bool)`

GetLastNameOk returns a tuple with the LastName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastName

`func (o *UserResponse) SetLastName(v string)`

SetLastName sets LastName field to given value.

### HasLastName

`func (o *UserResponse) HasLastName() bool`

HasLastName returns a boolean if a field has been set.

### GetEmailVerified

`func (o *UserResponse) GetEmailVerified() bool`

GetEmailVerified returns the EmailVerified field if non-nil, zero value otherwise.

### GetEmailVerifiedOk

`func (o *UserResponse) GetEmailVerifiedOk() (*bool, bool)`

GetEmailVerifiedOk returns a tuple with the EmailVerified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailVerified

`func (o *UserResponse) SetEmailVerified(v bool)`

SetEmailVerified sets EmailVerified field to given value.

### HasEmailVerified

`func (o *UserResponse) HasEmailVerified() bool`

HasEmailVerified returns a boolean if a field has been set.

### GetIsPlatformAdmin

`func (o *UserResponse) GetIsPlatformAdmin() bool`

GetIsPlatformAdmin returns the IsPlatformAdmin field if non-nil, zero value otherwise.

### GetIsPlatformAdminOk

`func (o *UserResponse) GetIsPlatformAdminOk() (*bool, bool)`

GetIsPlatformAdminOk returns a tuple with the IsPlatformAdmin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPlatformAdmin

`func (o *UserResponse) SetIsPlatformAdmin(v bool)`

SetIsPlatformAdmin sets IsPlatformAdmin field to given value.

### HasIsPlatformAdmin

`func (o *UserResponse) HasIsPlatformAdmin() bool`

HasIsPlatformAdmin returns a boolean if a field has been set.

### GetMfaEnabled

`func (o *UserResponse) GetMfaEnabled() bool`

GetMfaEnabled returns the MfaEnabled field if non-nil, zero value otherwise.

### GetMfaEnabledOk

`func (o *UserResponse) GetMfaEnabledOk() (*bool, bool)`

GetMfaEnabledOk returns a tuple with the MfaEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaEnabled

`func (o *UserResponse) SetMfaEnabled(v bool)`

SetMfaEnabled sets MfaEnabled field to given value.

### HasMfaEnabled

`func (o *UserResponse) HasMfaEnabled() bool`

HasMfaEnabled returns a boolean if a field has been set.

### GetActiveTenantId

`func (o *UserResponse) GetActiveTenantId() string`

GetActiveTenantId returns the ActiveTenantId field if non-nil, zero value otherwise.

### GetActiveTenantIdOk

`func (o *UserResponse) GetActiveTenantIdOk() (*string, bool)`

GetActiveTenantIdOk returns a tuple with the ActiveTenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActiveTenantId

`func (o *UserResponse) SetActiveTenantId(v string)`

SetActiveTenantId sets ActiveTenantId field to given value.

### HasActiveTenantId

`func (o *UserResponse) HasActiveTenantId() bool`

HasActiveTenantId returns a boolean if a field has been set.

### GetTenants

`func (o *UserResponse) GetTenants() []TenantSummary`

GetTenants returns the Tenants field if non-nil, zero value otherwise.

### GetTenantsOk

`func (o *UserResponse) GetTenantsOk() (*[]TenantSummary, bool)`

GetTenantsOk returns a tuple with the Tenants field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenants

`func (o *UserResponse) SetTenants(v []TenantSummary)`

SetTenants sets Tenants field to given value.

### HasTenants

`func (o *UserResponse) HasTenants() bool`

HasTenants returns a boolean if a field has been set.

### GetCreatedAt

`func (o *UserResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *UserResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *UserResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *UserResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *UserResponse) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *UserResponse) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *UserResponse) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *UserResponse) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


