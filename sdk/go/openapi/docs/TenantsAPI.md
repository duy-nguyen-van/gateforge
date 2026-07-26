# \TenantsAPI

All URIs are relative to *http://localhost:3000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ListMyTenants**](TenantsAPI.md#ListMyTenants) | **Get** /api/v1/me/tenants | List tenants the current user can access
[**SelectTenant**](TenantsAPI.md#SelectTenant) | **Post** /api/v1/tenants/select | Complete login by selecting a tenant
[**SwitchTenant**](TenantsAPI.md#SwitchTenant) | **Post** /api/v1/tenants/switch | Switch active tenant and re-issue tokens



## ListMyTenants

> TenantSummaryListEnvelope ListMyTenants(ctx).Page(page).PageSize(pageSize).Execute()

List tenants the current user can access

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
	resp, r, err := apiClient.TenantsAPI.ListMyTenants(context.Background()).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.ListMyTenants``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListMyTenants`: TenantSummaryListEnvelope
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.ListMyTenants`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListMyTenantsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** | Page number | 
 **pageSize** | **int32** | Page size | 

### Return type

[**TenantSummaryListEnvelope**](TenantSummaryListEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SelectTenant

> LoginResponseEnvelope SelectTenant(ctx).TenantSelectRequest(tenantSelectRequest).Execute()

Complete login by selecting a tenant

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
	tenantSelectRequest := *openapiclient.NewTenantSelectRequest("SelectionToken_example", "TenantId_example") // TenantSelectRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.SelectTenant(context.Background()).TenantSelectRequest(tenantSelectRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.SelectTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SelectTenant`: LoginResponseEnvelope
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.SelectTenant`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSelectTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tenantSelectRequest** | [**TenantSelectRequest**](TenantSelectRequest.md) |  | 

### Return type

[**LoginResponseEnvelope**](LoginResponseEnvelope.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SwitchTenant

> LoginResponseEnvelope SwitchTenant(ctx).TenantSwitchRequest(tenantSwitchRequest).Execute()

Switch active tenant and re-issue tokens

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
	tenantSwitchRequest := *openapiclient.NewTenantSwitchRequest("TenantId_example") // TenantSwitchRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.SwitchTenant(context.Background()).TenantSwitchRequest(tenantSwitchRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.SwitchTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SwitchTenant`: LoginResponseEnvelope
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.SwitchTenant`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tenantSwitchRequest** | [**TenantSwitchRequest**](TenantSwitchRequest.md) |  | 

### Return type

[**LoginResponseEnvelope**](LoginResponseEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

