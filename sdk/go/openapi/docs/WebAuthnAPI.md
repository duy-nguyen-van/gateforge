# \WebAuthnAPI

All URIs are relative to *http://localhost:3000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**FinishWebauthnLogin**](WebAuthnAPI.md#FinishWebauthnLogin) | **Post** /api/v1/webauthn/login/finish | Finish passkey login
[**FinishWebauthnRegister**](WebAuthnAPI.md#FinishWebauthnRegister) | **Post** /api/v1/webauthn/register/finish | Finish passkey registration
[**ListWebauthnCredentials**](WebAuthnAPI.md#ListWebauthnCredentials) | **Get** /api/v1/webauthn/credentials | List registered passkeys
[**StartWebauthnLogin**](WebAuthnAPI.md#StartWebauthnLogin) | **Post** /api/v1/webauthn/login/start | Begin passkey login
[**StartWebauthnRegister**](WebAuthnAPI.md#StartWebauthnRegister) | **Post** /api/v1/webauthn/register/start | Begin passkey registration



## FinishWebauthnLogin

> LoginResultEnvelope FinishWebauthnLogin(ctx).WebauthnLoginFinishRequest(webauthnLoginFinishRequest).Execute()

Finish passkey login



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
	webauthnLoginFinishRequest := *openapiclient.NewWebauthnLoginFinishRequest("user@example.com", "SessionToken_example") // WebauthnLoginFinishRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WebAuthnAPI.FinishWebauthnLogin(context.Background()).WebauthnLoginFinishRequest(webauthnLoginFinishRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnAPI.FinishWebauthnLogin``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `FinishWebauthnLogin`: LoginResultEnvelope
	fmt.Fprintf(os.Stdout, "Response from `WebAuthnAPI.FinishWebauthnLogin`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiFinishWebauthnLoginRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **webauthnLoginFinishRequest** | [**WebauthnLoginFinishRequest**](WebauthnLoginFinishRequest.md) |  | 

### Return type

[**LoginResultEnvelope**](LoginResultEnvelope.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## FinishWebauthnRegister

> EmptyDataEnvelope FinishWebauthnRegister(ctx).WebauthnRegisterFinishRequest(webauthnRegisterFinishRequest).Execute()

Finish passkey registration



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
	webauthnRegisterFinishRequest := *openapiclient.NewWebauthnRegisterFinishRequest("SessionToken_example") // WebauthnRegisterFinishRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WebAuthnAPI.FinishWebauthnRegister(context.Background()).WebauthnRegisterFinishRequest(webauthnRegisterFinishRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnAPI.FinishWebauthnRegister``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `FinishWebauthnRegister`: EmptyDataEnvelope
	fmt.Fprintf(os.Stdout, "Response from `WebAuthnAPI.FinishWebauthnRegister`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiFinishWebauthnRegisterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **webauthnRegisterFinishRequest** | [**WebauthnRegisterFinishRequest**](WebauthnRegisterFinishRequest.md) |  | 

### Return type

[**EmptyDataEnvelope**](EmptyDataEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListWebauthnCredentials

> WebauthnCredentialListEnvelope ListWebauthnCredentials(ctx).Page(page).PageSize(pageSize).Execute()

List registered passkeys



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
	resp, r, err := apiClient.WebAuthnAPI.ListWebauthnCredentials(context.Background()).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnAPI.ListWebauthnCredentials``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListWebauthnCredentials`: WebauthnCredentialListEnvelope
	fmt.Fprintf(os.Stdout, "Response from `WebAuthnAPI.ListWebauthnCredentials`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListWebauthnCredentialsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** | Page number | 
 **pageSize** | **int32** | Page size | 

### Return type

[**WebauthnCredentialListEnvelope**](WebauthnCredentialListEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StartWebauthnLogin

> WebauthnLoginStartEnvelope StartWebauthnLogin(ctx).WebauthnLoginStartRequest(webauthnLoginStartRequest).Execute()

Begin passkey login



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
	webauthnLoginStartRequest := *openapiclient.NewWebauthnLoginStartRequest("user@example.com") // WebauthnLoginStartRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WebAuthnAPI.StartWebauthnLogin(context.Background()).WebauthnLoginStartRequest(webauthnLoginStartRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnAPI.StartWebauthnLogin``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `StartWebauthnLogin`: WebauthnLoginStartEnvelope
	fmt.Fprintf(os.Stdout, "Response from `WebAuthnAPI.StartWebauthnLogin`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiStartWebauthnLoginRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **webauthnLoginStartRequest** | [**WebauthnLoginStartRequest**](WebauthnLoginStartRequest.md) |  | 

### Return type

[**WebauthnLoginStartEnvelope**](WebauthnLoginStartEnvelope.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StartWebauthnRegister

> WebauthnRegisterStartEnvelope StartWebauthnRegister(ctx).WebauthnRegisterStartRequest(webauthnRegisterStartRequest).Execute()

Begin passkey registration

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
	webauthnRegisterStartRequest := *openapiclient.NewWebauthnRegisterStartRequest() // WebauthnRegisterStartRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WebAuthnAPI.StartWebauthnRegister(context.Background()).WebauthnRegisterStartRequest(webauthnRegisterStartRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnAPI.StartWebauthnRegister``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `StartWebauthnRegister`: WebauthnRegisterStartEnvelope
	fmt.Fprintf(os.Stdout, "Response from `WebAuthnAPI.StartWebauthnRegister`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiStartWebauthnRegisterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **webauthnRegisterStartRequest** | [**WebauthnRegisterStartRequest**](WebauthnRegisterStartRequest.md) |  | 

### Return type

[**WebauthnRegisterStartEnvelope**](WebauthnRegisterStartEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

