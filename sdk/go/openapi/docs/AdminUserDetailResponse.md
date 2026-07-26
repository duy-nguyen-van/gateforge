# AdminUserDetailResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 
**FirstName** | Pointer to **string** |  | [optional] 
**LastName** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**MfaEnabled** | Pointer to **bool** |  | [optional] 
**IsPlatformAdmin** | Pointer to **bool** |  | [optional] 
**PasskeyCount** | Pointer to **int32** |  | [optional] 
**ActiveSessions** | Pointer to **int32** |  | [optional] 
**Memberships** | Pointer to [**[]AdminUserMembershipResponse**](AdminUserMembershipResponse.md) |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewAdminUserDetailResponse

`func NewAdminUserDetailResponse() *AdminUserDetailResponse`

NewAdminUserDetailResponse instantiates a new AdminUserDetailResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminUserDetailResponseWithDefaults

`func NewAdminUserDetailResponseWithDefaults() *AdminUserDetailResponse`

NewAdminUserDetailResponseWithDefaults instantiates a new AdminUserDetailResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AdminUserDetailResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AdminUserDetailResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AdminUserDetailResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *AdminUserDetailResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetEmail

`func (o *AdminUserDetailResponse) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *AdminUserDetailResponse) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *AdminUserDetailResponse) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *AdminUserDetailResponse) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetFirstName

`func (o *AdminUserDetailResponse) GetFirstName() string`

GetFirstName returns the FirstName field if non-nil, zero value otherwise.

### GetFirstNameOk

`func (o *AdminUserDetailResponse) GetFirstNameOk() (*string, bool)`

GetFirstNameOk returns a tuple with the FirstName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstName

`func (o *AdminUserDetailResponse) SetFirstName(v string)`

SetFirstName sets FirstName field to given value.

### HasFirstName

`func (o *AdminUserDetailResponse) HasFirstName() bool`

HasFirstName returns a boolean if a field has been set.

### GetLastName

`func (o *AdminUserDetailResponse) GetLastName() string`

GetLastName returns the LastName field if non-nil, zero value otherwise.

### GetLastNameOk

`func (o *AdminUserDetailResponse) GetLastNameOk() (*string, bool)`

GetLastNameOk returns a tuple with the LastName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastName

`func (o *AdminUserDetailResponse) SetLastName(v string)`

SetLastName sets LastName field to given value.

### HasLastName

`func (o *AdminUserDetailResponse) HasLastName() bool`

HasLastName returns a boolean if a field has been set.

### GetStatus

`func (o *AdminUserDetailResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *AdminUserDetailResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *AdminUserDetailResponse) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *AdminUserDetailResponse) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetMfaEnabled

`func (o *AdminUserDetailResponse) GetMfaEnabled() bool`

GetMfaEnabled returns the MfaEnabled field if non-nil, zero value otherwise.

### GetMfaEnabledOk

`func (o *AdminUserDetailResponse) GetMfaEnabledOk() (*bool, bool)`

GetMfaEnabledOk returns a tuple with the MfaEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaEnabled

`func (o *AdminUserDetailResponse) SetMfaEnabled(v bool)`

SetMfaEnabled sets MfaEnabled field to given value.

### HasMfaEnabled

`func (o *AdminUserDetailResponse) HasMfaEnabled() bool`

HasMfaEnabled returns a boolean if a field has been set.

### GetIsPlatformAdmin

`func (o *AdminUserDetailResponse) GetIsPlatformAdmin() bool`

GetIsPlatformAdmin returns the IsPlatformAdmin field if non-nil, zero value otherwise.

### GetIsPlatformAdminOk

`func (o *AdminUserDetailResponse) GetIsPlatformAdminOk() (*bool, bool)`

GetIsPlatformAdminOk returns a tuple with the IsPlatformAdmin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPlatformAdmin

`func (o *AdminUserDetailResponse) SetIsPlatformAdmin(v bool)`

SetIsPlatformAdmin sets IsPlatformAdmin field to given value.

### HasIsPlatformAdmin

`func (o *AdminUserDetailResponse) HasIsPlatformAdmin() bool`

HasIsPlatformAdmin returns a boolean if a field has been set.

### GetPasskeyCount

`func (o *AdminUserDetailResponse) GetPasskeyCount() int32`

GetPasskeyCount returns the PasskeyCount field if non-nil, zero value otherwise.

### GetPasskeyCountOk

`func (o *AdminUserDetailResponse) GetPasskeyCountOk() (*int32, bool)`

GetPasskeyCountOk returns a tuple with the PasskeyCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPasskeyCount

`func (o *AdminUserDetailResponse) SetPasskeyCount(v int32)`

SetPasskeyCount sets PasskeyCount field to given value.

### HasPasskeyCount

`func (o *AdminUserDetailResponse) HasPasskeyCount() bool`

HasPasskeyCount returns a boolean if a field has been set.

### GetActiveSessions

`func (o *AdminUserDetailResponse) GetActiveSessions() int32`

GetActiveSessions returns the ActiveSessions field if non-nil, zero value otherwise.

### GetActiveSessionsOk

`func (o *AdminUserDetailResponse) GetActiveSessionsOk() (*int32, bool)`

GetActiveSessionsOk returns a tuple with the ActiveSessions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActiveSessions

`func (o *AdminUserDetailResponse) SetActiveSessions(v int32)`

SetActiveSessions sets ActiveSessions field to given value.

### HasActiveSessions

`func (o *AdminUserDetailResponse) HasActiveSessions() bool`

HasActiveSessions returns a boolean if a field has been set.

### GetMemberships

`func (o *AdminUserDetailResponse) GetMemberships() []AdminUserMembershipResponse`

GetMemberships returns the Memberships field if non-nil, zero value otherwise.

### GetMembershipsOk

`func (o *AdminUserDetailResponse) GetMembershipsOk() (*[]AdminUserMembershipResponse, bool)`

GetMembershipsOk returns a tuple with the Memberships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemberships

`func (o *AdminUserDetailResponse) SetMemberships(v []AdminUserMembershipResponse)`

SetMemberships sets Memberships field to given value.

### HasMemberships

`func (o *AdminUserDetailResponse) HasMemberships() bool`

HasMemberships returns a boolean if a field has been set.

### GetCreatedAt

`func (o *AdminUserDetailResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *AdminUserDetailResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *AdminUserDetailResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *AdminUserDetailResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


