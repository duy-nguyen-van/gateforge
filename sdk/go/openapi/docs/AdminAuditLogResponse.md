# AdminAuditLogResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Action** | Pointer to **string** |  | [optional] 
**ActorId** | Pointer to **string** |  | [optional] 
**ActorType** | Pointer to **string** |  | [optional] 
**CorrelationId** | Pointer to **string** |  | [optional] 
**RequestId** | Pointer to **string** |  | [optional] 
**ResourceId** | Pointer to **string** |  | [optional] 
**ResourceName** | Pointer to **string** |  | [optional] 
**ResourceType** | Pointer to **string** |  | [optional] 
**Result** | Pointer to **string** |  | [optional] 
**TenantId** | Pointer to **string** |  | [optional] 
**IpAddress** | Pointer to **string** |  | [optional] 
**UserAgent** | Pointer to **string** |  | [optional] 
**OldValue** | Pointer to **map[string]interface{}** |  | [optional] 
**NewValue** | Pointer to **map[string]interface{}** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewAdminAuditLogResponse

`func NewAdminAuditLogResponse() *AdminAuditLogResponse`

NewAdminAuditLogResponse instantiates a new AdminAuditLogResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdminAuditLogResponseWithDefaults

`func NewAdminAuditLogResponseWithDefaults() *AdminAuditLogResponse`

NewAdminAuditLogResponseWithDefaults instantiates a new AdminAuditLogResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AdminAuditLogResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AdminAuditLogResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AdminAuditLogResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *AdminAuditLogResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetAction

`func (o *AdminAuditLogResponse) GetAction() string`

GetAction returns the Action field if non-nil, zero value otherwise.

### GetActionOk

`func (o *AdminAuditLogResponse) GetActionOk() (*string, bool)`

GetActionOk returns a tuple with the Action field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAction

`func (o *AdminAuditLogResponse) SetAction(v string)`

SetAction sets Action field to given value.

### HasAction

`func (o *AdminAuditLogResponse) HasAction() bool`

HasAction returns a boolean if a field has been set.

### GetActorId

`func (o *AdminAuditLogResponse) GetActorId() string`

GetActorId returns the ActorId field if non-nil, zero value otherwise.

### GetActorIdOk

`func (o *AdminAuditLogResponse) GetActorIdOk() (*string, bool)`

GetActorIdOk returns a tuple with the ActorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActorId

`func (o *AdminAuditLogResponse) SetActorId(v string)`

SetActorId sets ActorId field to given value.

### HasActorId

`func (o *AdminAuditLogResponse) HasActorId() bool`

HasActorId returns a boolean if a field has been set.

### GetActorType

`func (o *AdminAuditLogResponse) GetActorType() string`

GetActorType returns the ActorType field if non-nil, zero value otherwise.

### GetActorTypeOk

`func (o *AdminAuditLogResponse) GetActorTypeOk() (*string, bool)`

GetActorTypeOk returns a tuple with the ActorType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActorType

`func (o *AdminAuditLogResponse) SetActorType(v string)`

SetActorType sets ActorType field to given value.

### HasActorType

`func (o *AdminAuditLogResponse) HasActorType() bool`

HasActorType returns a boolean if a field has been set.

### GetCorrelationId

`func (o *AdminAuditLogResponse) GetCorrelationId() string`

GetCorrelationId returns the CorrelationId field if non-nil, zero value otherwise.

### GetCorrelationIdOk

`func (o *AdminAuditLogResponse) GetCorrelationIdOk() (*string, bool)`

GetCorrelationIdOk returns a tuple with the CorrelationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCorrelationId

`func (o *AdminAuditLogResponse) SetCorrelationId(v string)`

SetCorrelationId sets CorrelationId field to given value.

### HasCorrelationId

`func (o *AdminAuditLogResponse) HasCorrelationId() bool`

HasCorrelationId returns a boolean if a field has been set.

### GetRequestId

`func (o *AdminAuditLogResponse) GetRequestId() string`

GetRequestId returns the RequestId field if non-nil, zero value otherwise.

### GetRequestIdOk

`func (o *AdminAuditLogResponse) GetRequestIdOk() (*string, bool)`

GetRequestIdOk returns a tuple with the RequestId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestId

`func (o *AdminAuditLogResponse) SetRequestId(v string)`

SetRequestId sets RequestId field to given value.

### HasRequestId

`func (o *AdminAuditLogResponse) HasRequestId() bool`

HasRequestId returns a boolean if a field has been set.

### GetResourceId

`func (o *AdminAuditLogResponse) GetResourceId() string`

GetResourceId returns the ResourceId field if non-nil, zero value otherwise.

### GetResourceIdOk

`func (o *AdminAuditLogResponse) GetResourceIdOk() (*string, bool)`

GetResourceIdOk returns a tuple with the ResourceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceId

`func (o *AdminAuditLogResponse) SetResourceId(v string)`

SetResourceId sets ResourceId field to given value.

### HasResourceId

`func (o *AdminAuditLogResponse) HasResourceId() bool`

HasResourceId returns a boolean if a field has been set.

### GetResourceName

`func (o *AdminAuditLogResponse) GetResourceName() string`

GetResourceName returns the ResourceName field if non-nil, zero value otherwise.

### GetResourceNameOk

`func (o *AdminAuditLogResponse) GetResourceNameOk() (*string, bool)`

GetResourceNameOk returns a tuple with the ResourceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceName

`func (o *AdminAuditLogResponse) SetResourceName(v string)`

SetResourceName sets ResourceName field to given value.

### HasResourceName

`func (o *AdminAuditLogResponse) HasResourceName() bool`

HasResourceName returns a boolean if a field has been set.

### GetResourceType

`func (o *AdminAuditLogResponse) GetResourceType() string`

GetResourceType returns the ResourceType field if non-nil, zero value otherwise.

### GetResourceTypeOk

`func (o *AdminAuditLogResponse) GetResourceTypeOk() (*string, bool)`

GetResourceTypeOk returns a tuple with the ResourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceType

`func (o *AdminAuditLogResponse) SetResourceType(v string)`

SetResourceType sets ResourceType field to given value.

### HasResourceType

`func (o *AdminAuditLogResponse) HasResourceType() bool`

HasResourceType returns a boolean if a field has been set.

### GetResult

`func (o *AdminAuditLogResponse) GetResult() string`

GetResult returns the Result field if non-nil, zero value otherwise.

### GetResultOk

`func (o *AdminAuditLogResponse) GetResultOk() (*string, bool)`

GetResultOk returns a tuple with the Result field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResult

`func (o *AdminAuditLogResponse) SetResult(v string)`

SetResult sets Result field to given value.

### HasResult

`func (o *AdminAuditLogResponse) HasResult() bool`

HasResult returns a boolean if a field has been set.

### GetTenantId

`func (o *AdminAuditLogResponse) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *AdminAuditLogResponse) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *AdminAuditLogResponse) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *AdminAuditLogResponse) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetIpAddress

`func (o *AdminAuditLogResponse) GetIpAddress() string`

GetIpAddress returns the IpAddress field if non-nil, zero value otherwise.

### GetIpAddressOk

`func (o *AdminAuditLogResponse) GetIpAddressOk() (*string, bool)`

GetIpAddressOk returns a tuple with the IpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpAddress

`func (o *AdminAuditLogResponse) SetIpAddress(v string)`

SetIpAddress sets IpAddress field to given value.

### HasIpAddress

`func (o *AdminAuditLogResponse) HasIpAddress() bool`

HasIpAddress returns a boolean if a field has been set.

### GetUserAgent

`func (o *AdminAuditLogResponse) GetUserAgent() string`

GetUserAgent returns the UserAgent field if non-nil, zero value otherwise.

### GetUserAgentOk

`func (o *AdminAuditLogResponse) GetUserAgentOk() (*string, bool)`

GetUserAgentOk returns a tuple with the UserAgent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserAgent

`func (o *AdminAuditLogResponse) SetUserAgent(v string)`

SetUserAgent sets UserAgent field to given value.

### HasUserAgent

`func (o *AdminAuditLogResponse) HasUserAgent() bool`

HasUserAgent returns a boolean if a field has been set.

### GetOldValue

`func (o *AdminAuditLogResponse) GetOldValue() map[string]interface{}`

GetOldValue returns the OldValue field if non-nil, zero value otherwise.

### GetOldValueOk

`func (o *AdminAuditLogResponse) GetOldValueOk() (*map[string]interface{}, bool)`

GetOldValueOk returns a tuple with the OldValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOldValue

`func (o *AdminAuditLogResponse) SetOldValue(v map[string]interface{})`

SetOldValue sets OldValue field to given value.

### HasOldValue

`func (o *AdminAuditLogResponse) HasOldValue() bool`

HasOldValue returns a boolean if a field has been set.

### GetNewValue

`func (o *AdminAuditLogResponse) GetNewValue() map[string]interface{}`

GetNewValue returns the NewValue field if non-nil, zero value otherwise.

### GetNewValueOk

`func (o *AdminAuditLogResponse) GetNewValueOk() (*map[string]interface{}, bool)`

GetNewValueOk returns a tuple with the NewValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNewValue

`func (o *AdminAuditLogResponse) SetNewValue(v map[string]interface{})`

SetNewValue sets NewValue field to given value.

### HasNewValue

`func (o *AdminAuditLogResponse) HasNewValue() bool`

HasNewValue returns a boolean if a field has been set.

### GetCreatedAt

`func (o *AdminAuditLogResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *AdminAuditLogResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *AdminAuditLogResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *AdminAuditLogResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


