# \FederationAPI

All URIs are relative to *http://localhost:3000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CallbackFederationOAuth**](FederationAPI.md#CallbackFederationOAuth) | **Get** /oidc/federation/{provider}/callback | Federated OAuth callback
[**ListFederationProviders**](FederationAPI.md#ListFederationProviders) | **Get** /api/v1/federation/providers | List enabled federation sign-in methods for a tenant
[**StartFederationOAuth**](FederationAPI.md#StartFederationOAuth) | **Get** /oidc/federation/{provider}/start | Start federated sign-in (OIDC browser flow)



## CallbackFederationOAuth

> CallbackFederationOAuth(ctx, provider).Code(code).State(state).Error_(error_).ErrorDescription(errorDescription).Execute()

Federated OAuth callback

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
	provider := "provider_example" // string | Provider id (e.g. google)
	code := "code_example" // string | Authorization code
	state := "state_example" // string | State
	error_ := "error__example" // string | OAuth error code from provider (optional)
	errorDescription := "errorDescription_example" // string | OAuth error description from provider (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.FederationAPI.CallbackFederationOAuth(context.Background(), provider).Code(code).State(state).Error_(error_).ErrorDescription(errorDescription).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FederationAPI.CallbackFederationOAuth``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**provider** | **string** | Provider id (e.g. google) | 

### Other Parameters

Other parameters are passed through a pointer to a apiCallbackFederationOAuthRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **code** | **string** | Authorization code | 
 **state** | **string** | State | 
 **error_** | **string** | OAuth error code from provider | 
 **errorDescription** | **string** | OAuth error description from provider | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListFederationProviders

> PublicFederationProviderListEnvelope ListFederationProviders(ctx).TenantId(tenantId).Page(page).PageSize(pageSize).Execute()

List enabled federation sign-in methods for a tenant

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
	tenantId := "tenantId_example" // string | Tenant UUID (defaults to DEFAULT_TENANT_ID) (optional)
	page := int32(56) // int32 | Page number (optional)
	pageSize := int32(56) // int32 | Page size (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.FederationAPI.ListFederationProviders(context.Background()).TenantId(tenantId).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FederationAPI.ListFederationProviders``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListFederationProviders`: PublicFederationProviderListEnvelope
	fmt.Fprintf(os.Stdout, "Response from `FederationAPI.ListFederationProviders`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListFederationProvidersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tenantId** | **string** | Tenant UUID (defaults to DEFAULT_TENANT_ID) | 
 **page** | **int32** | Page number | 
 **pageSize** | **int32** | Page size | 

### Return type

[**PublicFederationProviderListEnvelope**](PublicFederationProviderListEnvelope.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StartFederationOAuth

> StartFederationOAuth(ctx, provider).ReturnTo(returnTo).Execute()

Start federated sign-in (OIDC browser flow)

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
	provider := "provider_example" // string | Provider id (e.g. google)
	returnTo := "returnTo_example" // string | URL to return to after login (must be /authorize on this app)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.FederationAPI.StartFederationOAuth(context.Background(), provider).ReturnTo(returnTo).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FederationAPI.StartFederationOAuth``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**provider** | **string** | Provider id (e.g. google) | 

### Other Parameters

Other parameters are passed through a pointer to a apiStartFederationOAuthRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **returnTo** | **string** | URL to return to after login (must be /authorize on this app) | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

