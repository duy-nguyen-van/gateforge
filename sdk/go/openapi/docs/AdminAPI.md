# \AdminAPI

All URIs are relative to *http://localhost:3000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddAdminTenantMember**](AdminAPI.md#AddAdminTenantMember) | **Post** /api/v1/admin/tenants/{tenantId}/members | Add an existing user to a tenant
[**CreateAdminClient**](AdminAPI.md#CreateAdminClient) | **Post** /api/v1/admin/clients | Create an OAuth client (platform admin)
[**CreateAdminTenant**](AdminAPI.md#CreateAdminTenant) | **Post** /api/v1/admin/tenants | Create a tenant (platform admin)
[**DeleteAdminClient**](AdminAPI.md#DeleteAdminClient) | **Delete** /api/v1/admin/clients/{clientId} | Delete an OAuth client (platform admin)
[**DeleteAdminTenant**](AdminAPI.md#DeleteAdminTenant) | **Delete** /api/v1/admin/tenants/{tenantId} | Delete a tenant (platform admin)
[**DisableAdminUser**](AdminAPI.md#DisableAdminUser) | **Post** /api/v1/admin/users/{userId}/disable | Disable a user account (platform admin)
[**ForceLogoutAdminUser**](AdminAPI.md#ForceLogoutAdminUser) | **Post** /api/v1/admin/users/{userId}/force-logout | Force logout a user (platform admin)
[**GetAdminClient**](AdminAPI.md#GetAdminClient) | **Get** /api/v1/admin/clients/{clientId} | Get an OAuth client by ID (platform admin)
[**GetAdminClientUsage**](AdminAPI.md#GetAdminClientUsage) | **Get** /api/v1/admin/clients/{clientId}/usage | Get OAuth client usage metrics (platform admin)
[**GetAdminStats**](AdminAPI.md#GetAdminStats) | **Get** /api/v1/admin/stats | Platform admin dashboard stats
[**GetAdminTenant**](AdminAPI.md#GetAdminTenant) | **Get** /api/v1/admin/tenants/{tenantId} | Get a tenant by ID (platform admin)
[**GetAdminUser**](AdminAPI.md#GetAdminUser) | **Get** /api/v1/admin/users/{userId} | Get user details (platform admin)
[**ListAdminAuditLogs**](AdminAPI.md#ListAdminAuditLogs) | **Get** /api/v1/admin/audit-logs | List platform audit logs
[**ListAdminClients**](AdminAPI.md#ListAdminClients) | **Get** /api/v1/admin/clients | List OAuth clients (platform admin)
[**ListAdminIdentityProviders**](AdminAPI.md#ListAdminIdentityProviders) | **Get** /api/v1/admin/tenants/{tenantId}/identity-providers | List identity providers for a tenant
[**ListAdminLoginHistory**](AdminAPI.md#ListAdminLoginHistory) | **Get** /api/v1/admin/login-history | List login history events (platform admin)
[**ListAdminTenantMembers**](AdminAPI.md#ListAdminTenantMembers) | **Get** /api/v1/admin/tenants/{tenantId}/members | List members of a tenant (platform admin)
[**ListAdminTenants**](AdminAPI.md#ListAdminTenants) | **Get** /api/v1/admin/tenants | List tenants (platform admin)
[**ListAdminUsers**](AdminAPI.md#ListAdminUsers) | **Get** /api/v1/admin/users | List users (platform admin)
[**PatchAdminIdentityProvider**](AdminAPI.md#PatchAdminIdentityProvider) | **Patch** /api/v1/admin/tenants/{tenantId}/identity-providers/{provider} | Configure an upstream identity provider for a tenant (platform admin)
[**RemoveAdminTenantMember**](AdminAPI.md#RemoveAdminTenantMember) | **Delete** /api/v1/admin/tenants/{tenantId}/members/{userId} | Remove a user from a tenant
[**ResetAdminUserMfa**](AdminAPI.md#ResetAdminUserMfa) | **Post** /api/v1/admin/users/{userId}/reset-mfa | Reset MFA (TOTP + recovery codes) for a user (platform admin)
[**ResetAdminUserPasskeys**](AdminAPI.md#ResetAdminUserPasskeys) | **Post** /api/v1/admin/users/{userId}/reset-passkey | Reset all passkeys for a user (platform admin)
[**UpdateAdminClient**](AdminAPI.md#UpdateAdminClient) | **Patch** /api/v1/admin/clients/{clientId} | Update an OAuth client (platform admin)
[**UpdateAdminTenant**](AdminAPI.md#UpdateAdminTenant) | **Patch** /api/v1/admin/tenants/{tenantId} | Update a tenant (platform admin)



## AddAdminTenantMember

> AddAdminTenantMember(ctx, tenantId).AdminAddMemberRequest(adminAddMemberRequest).Execute()

Add an existing user to a tenant

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	tenantId := "tenantId_example" // string | Tenant UUID
	adminAddMemberRequest := *openapiclient.NewAdminAddMemberRequest("Email_example") // AdminAddMemberRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.AddAdminTenantMember(context.Background(), tenantId).AdminAddMemberRequest(adminAddMemberRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.AddAdminTenantMember``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**tenantId** | **string** | Tenant UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddAdminTenantMemberRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **adminAddMemberRequest** | [**AdminAddMemberRequest**](AdminAddMemberRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateAdminClient

> AdminCreateClientEnvelope CreateAdminClient(ctx).AdminCreateClientRequest(adminCreateClientRequest).Execute()

Create an OAuth client (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	adminCreateClientRequest := *openapiclient.NewAdminCreateClientRequest("Name_example", "TenantId_example", []string{"RedirectUris_example"}) // AdminCreateClientRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.CreateAdminClient(context.Background()).AdminCreateClientRequest(adminCreateClientRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.CreateAdminClient``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateAdminClient`: AdminCreateClientEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.CreateAdminClient`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAdminClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **adminCreateClientRequest** | [**AdminCreateClientRequest**](AdminCreateClientRequest.md) |  | 

### Return type

[**AdminCreateClientEnvelope**](AdminCreateClientEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateAdminTenant

> AdminTenantEnvelope CreateAdminTenant(ctx).AdminCreateTenantRequest(adminCreateTenantRequest).Execute()

Create a tenant (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	adminCreateTenantRequest := *openapiclient.NewAdminCreateTenantRequest("Name_example") // AdminCreateTenantRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.CreateAdminTenant(context.Background()).AdminCreateTenantRequest(adminCreateTenantRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.CreateAdminTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateAdminTenant`: AdminTenantEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.CreateAdminTenant`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAdminTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **adminCreateTenantRequest** | [**AdminCreateTenantRequest**](AdminCreateTenantRequest.md) |  | 

### Return type

[**AdminTenantEnvelope**](AdminTenantEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAdminClient

> DeleteAdminClient(ctx, clientId).Execute()

Delete an OAuth client (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	clientId := "clientId_example" // string | Client record UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.DeleteAdminClient(context.Background(), clientId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.DeleteAdminClient``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clientId** | **string** | Client record UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteAdminClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAdminTenant

> DeleteAdminTenant(ctx, tenantId).Execute()

Delete a tenant (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	tenantId := "tenantId_example" // string | Tenant UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.DeleteAdminTenant(context.Background(), tenantId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.DeleteAdminTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**tenantId** | **string** | Tenant UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteAdminTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DisableAdminUser

> DisableAdminUser(ctx, userId).Execute()

Disable a user account (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	userId := "userId_example" // string | User UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.DisableAdminUser(context.Background(), userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.DisableAdminUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**userId** | **string** | User UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDisableAdminUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ForceLogoutAdminUser

> ForceLogoutAdminUser(ctx, userId).Execute()

Force logout a user (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	userId := "userId_example" // string | User UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.ForceLogoutAdminUser(context.Background(), userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ForceLogoutAdminUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**userId** | **string** | User UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiForceLogoutAdminUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAdminClient

> AdminClientEnvelope GetAdminClient(ctx, clientId).Execute()

Get an OAuth client by ID (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	clientId := "clientId_example" // string | Client record UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.GetAdminClient(context.Background(), clientId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.GetAdminClient``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAdminClient`: AdminClientEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.GetAdminClient`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clientId** | **string** | Client record UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAdminClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**AdminClientEnvelope**](AdminClientEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAdminClientUsage

> AdminClientUsageEnvelope GetAdminClientUsage(ctx, clientId).Execute()

Get OAuth client usage metrics (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	clientId := "clientId_example" // string | Client record UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.GetAdminClientUsage(context.Background(), clientId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.GetAdminClientUsage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAdminClientUsage`: AdminClientUsageEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.GetAdminClientUsage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clientId** | **string** | Client record UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAdminClientUsageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**AdminClientUsageEnvelope**](AdminClientUsageEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAdminStats

> AdminStatsEnvelope GetAdminStats(ctx).Execute()

Platform admin dashboard stats

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.GetAdminStats(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.GetAdminStats``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAdminStats`: AdminStatsEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.GetAdminStats`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetAdminStatsRequest struct via the builder pattern


### Return type

[**AdminStatsEnvelope**](AdminStatsEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAdminTenant

> AdminTenantEnvelope GetAdminTenant(ctx, tenantId).Execute()

Get a tenant by ID (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	tenantId := "tenantId_example" // string | Tenant UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.GetAdminTenant(context.Background(), tenantId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.GetAdminTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAdminTenant`: AdminTenantEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.GetAdminTenant`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**tenantId** | **string** | Tenant UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAdminTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**AdminTenantEnvelope**](AdminTenantEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAdminUser

> AdminUserDetailEnvelope GetAdminUser(ctx, userId).Execute()

Get user details (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	userId := "userId_example" // string | User UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.GetAdminUser(context.Background(), userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.GetAdminUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAdminUser`: AdminUserDetailEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.GetAdminUser`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**userId** | **string** | User UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAdminUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**AdminUserDetailEnvelope**](AdminUserDetailEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAdminAuditLogs

> AdminAuditLogListEnvelope ListAdminAuditLogs(ctx).TenantId(tenantId).Action(action).Result(result).ActorId(actorId).From(from).To(to).Page(page).PageSize(pageSize).Execute()

List platform audit logs

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	tenantId := "tenantId_example" // string | Filter by tenant UUID (optional)
	action := "action_example" // string | Filter by action (exact or prefix with trailing dot) (optional)
	result := "result_example" // string | Filter by result (success, failure, denied) (optional)
	actorId := "actorId_example" // string | Filter by actor id (optional)
	from := "from_example" // string | ISO8601 lower bound on created_at (optional)
	to := "to_example" // string | ISO8601 upper bound on created_at (optional)
	page := int32(56) // int32 | Page number (optional)
	pageSize := int32(56) // int32 | Page size (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ListAdminAuditLogs(context.Background()).TenantId(tenantId).Action(action).Result(result).ActorId(actorId).From(from).To(to).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ListAdminAuditLogs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListAdminAuditLogs`: AdminAuditLogListEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ListAdminAuditLogs`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListAdminAuditLogsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tenantId** | **string** | Filter by tenant UUID | 
 **action** | **string** | Filter by action (exact or prefix with trailing dot) | 
 **result** | **string** | Filter by result (success, failure, denied) | 
 **actorId** | **string** | Filter by actor id | 
 **from** | **string** | ISO8601 lower bound on created_at | 
 **to** | **string** | ISO8601 upper bound on created_at | 
 **page** | **int32** | Page number | 
 **pageSize** | **int32** | Page size | 

### Return type

[**AdminAuditLogListEnvelope**](AdminAuditLogListEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAdminClients

> AdminClientListEnvelope ListAdminClients(ctx).Page(page).PageSize(pageSize).TenantId(tenantId).Execute()

List OAuth clients (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	page := int32(56) // int32 | Page number (optional)
	pageSize := int32(56) // int32 | Page size (optional)
	tenantId := "tenantId_example" // string | Filter by tenant UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ListAdminClients(context.Background()).Page(page).PageSize(pageSize).TenantId(tenantId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ListAdminClients``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListAdminClients`: AdminClientListEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ListAdminClients`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListAdminClientsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** | Page number | 
 **pageSize** | **int32** | Page size | 
 **tenantId** | **string** | Filter by tenant UUID | 

### Return type

[**AdminClientListEnvelope**](AdminClientListEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAdminIdentityProviders

> AdminIdentityProviderListEnvelope ListAdminIdentityProviders(ctx, tenantId).Page(page).PageSize(pageSize).Execute()

List identity providers for a tenant

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	tenantId := "tenantId_example" // string | Tenant UUID
	page := int32(56) // int32 | Page number (optional)
	pageSize := int32(56) // int32 | Page size (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ListAdminIdentityProviders(context.Background(), tenantId).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ListAdminIdentityProviders``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListAdminIdentityProviders`: AdminIdentityProviderListEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ListAdminIdentityProviders`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**tenantId** | **string** | Tenant UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListAdminIdentityProvidersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **int32** | Page number | 
 **pageSize** | **int32** | Page size | 

### Return type

[**AdminIdentityProviderListEnvelope**](AdminIdentityProviderListEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAdminLoginHistory

> AdminAuditLogListEnvelope ListAdminLoginHistory(ctx).TenantId(tenantId).Result(result).ActorId(actorId).From(from).To(to).Page(page).PageSize(pageSize).Execute()

List login history events (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	tenantId := "tenantId_example" // string | Filter by tenant UUID (optional)
	result := "result_example" // string | Filter by result (success, failure, denied) (optional)
	actorId := "actorId_example" // string | Filter by actor id (optional)
	from := "from_example" // string | ISO8601 lower bound on created_at (optional)
	to := "to_example" // string | ISO8601 upper bound on created_at (optional)
	page := int32(56) // int32 | Page number (optional)
	pageSize := int32(56) // int32 | Page size (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ListAdminLoginHistory(context.Background()).TenantId(tenantId).Result(result).ActorId(actorId).From(from).To(to).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ListAdminLoginHistory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListAdminLoginHistory`: AdminAuditLogListEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ListAdminLoginHistory`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListAdminLoginHistoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tenantId** | **string** | Filter by tenant UUID | 
 **result** | **string** | Filter by result (success, failure, denied) | 
 **actorId** | **string** | Filter by actor id | 
 **from** | **string** | ISO8601 lower bound on created_at | 
 **to** | **string** | ISO8601 upper bound on created_at | 
 **page** | **int32** | Page number | 
 **pageSize** | **int32** | Page size | 

### Return type

[**AdminAuditLogListEnvelope**](AdminAuditLogListEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAdminTenantMembers

> AdminTenantMemberListEnvelope ListAdminTenantMembers(ctx, tenantId).Page(page).PageSize(pageSize).Execute()

List members of a tenant (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	tenantId := "tenantId_example" // string | Tenant UUID
	page := int32(56) // int32 | Page number (optional)
	pageSize := int32(56) // int32 | Page size (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ListAdminTenantMembers(context.Background(), tenantId).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ListAdminTenantMembers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListAdminTenantMembers`: AdminTenantMemberListEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ListAdminTenantMembers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**tenantId** | **string** | Tenant UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListAdminTenantMembersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **int32** | Page number | 
 **pageSize** | **int32** | Page size | 

### Return type

[**AdminTenantMemberListEnvelope**](AdminTenantMemberListEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAdminTenants

> AdminTenantListEnvelope ListAdminTenants(ctx).Page(page).PageSize(pageSize).Execute()

List tenants (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	page := int32(56) // int32 | Page number (optional)
	pageSize := int32(56) // int32 | Page size (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ListAdminTenants(context.Background()).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ListAdminTenants``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListAdminTenants`: AdminTenantListEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ListAdminTenants`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListAdminTenantsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** | Page number | 
 **pageSize** | **int32** | Page size | 

### Return type

[**AdminTenantListEnvelope**](AdminTenantListEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAdminUsers

> AdminUserListEnvelope ListAdminUsers(ctx).Page(page).PageSize(pageSize).TenantId(tenantId).Search(search).Execute()

List users (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	page := int32(56) // int32 | Page number (optional)
	pageSize := int32(56) // int32 | Page size (optional)
	tenantId := "tenantId_example" // string | Filter by tenant UUID (optional)
	search := "search_example" // string | Search email or name (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ListAdminUsers(context.Background()).Page(page).PageSize(pageSize).TenantId(tenantId).Search(search).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ListAdminUsers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListAdminUsers`: AdminUserListEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ListAdminUsers`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListAdminUsersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** | Page number | 
 **pageSize** | **int32** | Page size | 
 **tenantId** | **string** | Filter by tenant UUID | 
 **search** | **string** | Search email or name | 

### Return type

[**AdminUserListEnvelope**](AdminUserListEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchAdminIdentityProvider

> PatchAdminIdentityProvider(ctx, tenantId, provider).PatchIdentityProviderRequest(patchIdentityProviderRequest).Execute()

Configure an upstream identity provider for a tenant (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	tenantId := "tenantId_example" // string | Tenant UUID
	provider := "provider_example" // string | Provider id (e.g. google)
	patchIdentityProviderRequest := *openapiclient.NewPatchIdentityProviderRequest() // PatchIdentityProviderRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.PatchAdminIdentityProvider(context.Background(), tenantId, provider).PatchIdentityProviderRequest(patchIdentityProviderRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.PatchAdminIdentityProvider``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**tenantId** | **string** | Tenant UUID | 
**provider** | **string** | Provider id (e.g. google) | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchAdminIdentityProviderRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **patchIdentityProviderRequest** | [**PatchIdentityProviderRequest**](PatchIdentityProviderRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveAdminTenantMember

> RemoveAdminTenantMember(ctx, tenantId, userId).Execute()

Remove a user from a tenant

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	tenantId := "tenantId_example" // string | Tenant UUID
	userId := "userId_example" // string | User UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.RemoveAdminTenantMember(context.Background(), tenantId, userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.RemoveAdminTenantMember``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**tenantId** | **string** | Tenant UUID | 
**userId** | **string** | User UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveAdminTenantMemberRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetAdminUserMfa

> ResetAdminUserMfa(ctx, userId).Execute()

Reset MFA (TOTP + recovery codes) for a user (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	userId := "userId_example" // string | User UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.ResetAdminUserMfa(context.Background(), userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ResetAdminUserMfa``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**userId** | **string** | User UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiResetAdminUserMfaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetAdminUserPasskeys

> ResetAdminUserPasskeys(ctx, userId).Execute()

Reset all passkeys for a user (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	userId := "userId_example" // string | User UUID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.ResetAdminUserPasskeys(context.Background(), userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ResetAdminUserPasskeys``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**userId** | **string** | User UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiResetAdminUserPasskeysRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateAdminClient

> AdminClientEnvelope UpdateAdminClient(ctx, clientId).AdminUpdateClientRequest(adminUpdateClientRequest).Execute()

Update an OAuth client (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	clientId := "clientId_example" // string | Client record UUID
	adminUpdateClientRequest := *openapiclient.NewAdminUpdateClientRequest() // AdminUpdateClientRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.UpdateAdminClient(context.Background(), clientId).AdminUpdateClientRequest(adminUpdateClientRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.UpdateAdminClient``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateAdminClient`: AdminClientEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.UpdateAdminClient`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clientId** | **string** | Client record UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateAdminClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **adminUpdateClientRequest** | [**AdminUpdateClientRequest**](AdminUpdateClientRequest.md) |  | 

### Return type

[**AdminClientEnvelope**](AdminClientEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateAdminTenant

> AdminTenantEnvelope UpdateAdminTenant(ctx, tenantId).AdminUpdateTenantRequest(adminUpdateTenantRequest).Execute()

Update a tenant (platform admin)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	tenantId := "tenantId_example" // string | Tenant UUID
	adminUpdateTenantRequest := *openapiclient.NewAdminUpdateTenantRequest() // AdminUpdateTenantRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.UpdateAdminTenant(context.Background(), tenantId).AdminUpdateTenantRequest(adminUpdateTenantRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.UpdateAdminTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateAdminTenant`: AdminTenantEnvelope
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.UpdateAdminTenant`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**tenantId** | **string** | Tenant UUID | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateAdminTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **adminUpdateTenantRequest** | [**AdminUpdateTenantRequest**](AdminUpdateTenantRequest.md) |  | 

### Return type

[**AdminTenantEnvelope**](AdminTenantEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

