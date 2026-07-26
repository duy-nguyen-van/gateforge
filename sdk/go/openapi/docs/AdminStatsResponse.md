# AdminStatsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ActiveSessions** | Pointer to **int32** |  | [optional] 
**MfaEnabledCount** | Pointer to **int32** |  | [optional] 
**MfaEnabledPercent** | Pointer to **float32** |  | [optional] 
**TotalUsers** | Pointer to **int32** |  | [optional] 

## Methods

### NewAdminStatsResponse

`func NewAdminStatsResponse() *AdminStatsResponse`

NewAdminStatsResponse instantiates a new AdminStatsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminStatsResponseWithDefaults

`func NewAdminStatsResponseWithDefaults() *AdminStatsResponse`

NewAdminStatsResponseWithDefaults instantiates a new AdminStatsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActiveSessions

`func (o *AdminStatsResponse) GetActiveSessions() int32`

GetActiveSessions returns the ActiveSessions field if non-nil, zero value otherwise.

### GetActiveSessionsOk

`func (o *AdminStatsResponse) GetActiveSessionsOk() (*int32, bool)`

GetActiveSessionsOk returns a tuple with the ActiveSessions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActiveSessions

`func (o *AdminStatsResponse) SetActiveSessions(v int32)`

SetActiveSessions sets ActiveSessions field to given value.

### HasActiveSessions

`func (o *AdminStatsResponse) HasActiveSessions() bool`

HasActiveSessions returns a boolean if a field has been set.

### GetMfaEnabledCount

`func (o *AdminStatsResponse) GetMfaEnabledCount() int32`

GetMfaEnabledCount returns the MfaEnabledCount field if non-nil, zero value otherwise.

### GetMfaEnabledCountOk

`func (o *AdminStatsResponse) GetMfaEnabledCountOk() (*int32, bool)`

GetMfaEnabledCountOk returns a tuple with the MfaEnabledCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaEnabledCount

`func (o *AdminStatsResponse) SetMfaEnabledCount(v int32)`

SetMfaEnabledCount sets MfaEnabledCount field to given value.

### HasMfaEnabledCount

`func (o *AdminStatsResponse) HasMfaEnabledCount() bool`

HasMfaEnabledCount returns a boolean if a field has been set.

### GetMfaEnabledPercent

`func (o *AdminStatsResponse) GetMfaEnabledPercent() float32`

GetMfaEnabledPercent returns the MfaEnabledPercent field if non-nil, zero value otherwise.

### GetMfaEnabledPercentOk

`func (o *AdminStatsResponse) GetMfaEnabledPercentOk() (*float32, bool)`

GetMfaEnabledPercentOk returns a tuple with the MfaEnabledPercent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaEnabledPercent

`func (o *AdminStatsResponse) SetMfaEnabledPercent(v float32)`

SetMfaEnabledPercent sets MfaEnabledPercent field to given value.

### HasMfaEnabledPercent

`func (o *AdminStatsResponse) HasMfaEnabledPercent() bool`

HasMfaEnabledPercent returns a boolean if a field has been set.

### GetTotalUsers

`func (o *AdminStatsResponse) GetTotalUsers() int32`

GetTotalUsers returns the TotalUsers field if non-nil, zero value otherwise.

### GetTotalUsersOk

`func (o *AdminStatsResponse) GetTotalUsersOk() (*int32, bool)`

GetTotalUsersOk returns a tuple with the TotalUsers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalUsers

`func (o *AdminStatsResponse) SetTotalUsers(v int32)`

SetTotalUsers sets TotalUsers field to given value.

### HasTotalUsers

`func (o *AdminStatsResponse) HasTotalUsers() bool`

HasTotalUsers returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


