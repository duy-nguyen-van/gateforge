# \OIDCAPI

All URIs are relative to *http://localhost:3000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Authorize**](OIDCAPI.md#Authorize) | **Get** /authorize | Authorize endpoint (authorization code + PKCE)
[**CreateToken**](OIDCAPI.md#CreateToken) | **Post** /token | Token endpoint (authorization_code / refresh_token grant)
[**GetJwks**](OIDCAPI.md#GetJwks) | **Get** /.well-known/jwks.json | JWKS endpoint (JSON Web Key Set)
[**GetOpenIdConfiguration**](OIDCAPI.md#GetOpenIdConfiguration) | **Get** /.well-known/openid-configuration | OpenID Configuration endpoint
[**GetUserInfo**](OIDCAPI.md#GetUserInfo) | **Get** /userinfo | Userinfo endpoint (Bearer access token from token endpoint, RS256)
[**LoginOidc**](OIDCAPI.md#LoginOidc) | **Post** /oidc/login | Login (OIDC browser flow) and continue /authorize



## Authorize

> Authorize(ctx).ResponseType(responseType).ClientId(clientId).RedirectUri(redirectUri).Scope(scope).State(state).Nonce(nonce).CodeChallenge(codeChallenge).CodeChallengeMethod(codeChallengeMethod).Execute()

Authorize endpoint (authorization code + PKCE)

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
	responseType := "responseType_example" // string | Response type (code) (optional)
	clientId := "clientId_example" // string | Client ID (optional)
	redirectUri := "redirectUri_example" // string | Redirect URI (optional)
	scope := "scope_example" // string | Scope (optional)
	state := "state_example" // string | State parameter (optional)
	nonce := "nonce_example" // string | Nonce (optional)
	codeChallenge := "codeChallenge_example" // string | PKCE code challenge (optional)
	codeChallengeMethod := "codeChallengeMethod_example" // string | PKCE method (S256) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OIDCAPI.Authorize(context.Background()).ResponseType(responseType).ClientId(clientId).RedirectUri(redirectUri).Scope(scope).State(state).Nonce(nonce).CodeChallenge(codeChallenge).CodeChallengeMethod(codeChallengeMethod).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OIDCAPI.Authorize``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthorizeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **responseType** | **string** | Response type (code) | 
 **clientId** | **string** | Client ID | 
 **redirectUri** | **string** | Redirect URI | 
 **scope** | **string** | Scope | 
 **state** | **string** | State parameter | 
 **nonce** | **string** | Nonce | 
 **codeChallenge** | **string** | PKCE code challenge | 
 **codeChallengeMethod** | **string** | PKCE method (S256) | 

### Return type

 (empty response body)

### Authorization

[CookieAuth](../README.md#CookieAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateToken

> OIDCTokenResponse CreateToken(ctx).GrantType(grantType).Code(code).RedirectUri(redirectUri).ClientId(clientId).ClientSecret(clientSecret).CodeVerifier(codeVerifier).RefreshToken(refreshToken).Scope(scope).Execute()

Token endpoint (authorization_code / refresh_token grant)

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
	grantType := "grantType_example" // string | authorization_code or refresh_token (optional)
	code := "code_example" // string | Authorization code (authorization_code grant) (optional)
	redirectUri := "redirectUri_example" // string |  (optional)
	clientId := "clientId_example" // string |  (optional)
	clientSecret := "clientSecret_example" // string |  (optional)
	codeVerifier := "codeVerifier_example" // string | PKCE code verifier (optional)
	refreshToken := "refreshToken_example" // string | Refresh token (refresh_token grant) (optional)
	scope := "scope_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OIDCAPI.CreateToken(context.Background()).GrantType(grantType).Code(code).RedirectUri(redirectUri).ClientId(clientId).ClientSecret(clientSecret).CodeVerifier(codeVerifier).RefreshToken(refreshToken).Scope(scope).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OIDCAPI.CreateToken``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateToken`: OIDCTokenResponse
	fmt.Fprintf(os.Stdout, "Response from `OIDCAPI.CreateToken`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateTokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **grantType** | **string** | authorization_code or refresh_token | 
 **code** | **string** | Authorization code (authorization_code grant) | 
 **redirectUri** | **string** |  | 
 **clientId** | **string** |  | 
 **clientSecret** | **string** |  | 
 **codeVerifier** | **string** | PKCE code verifier | 
 **refreshToken** | **string** | Refresh token (refresh_token grant) | 
 **scope** | **string** |  | 

### Return type

[**OIDCTokenResponse**](OIDCTokenResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetJwks

> JWKS GetJwks(ctx).Execute()

JWKS endpoint (JSON Web Key Set)

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
	resp, r, err := apiClient.OIDCAPI.GetJwks(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OIDCAPI.GetJwks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetJwks`: JWKS
	fmt.Fprintf(os.Stdout, "Response from `OIDCAPI.GetJwks`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetJwksRequest struct via the builder pattern


### Return type

[**JWKS**](JWKS.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOpenIdConfiguration

> OpenIDConfigurationResponse GetOpenIdConfiguration(ctx).Execute()

OpenID Configuration endpoint

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
	resp, r, err := apiClient.OIDCAPI.GetOpenIdConfiguration(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OIDCAPI.GetOpenIdConfiguration``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOpenIdConfiguration`: OpenIDConfigurationResponse
	fmt.Fprintf(os.Stdout, "Response from `OIDCAPI.GetOpenIdConfiguration`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetOpenIdConfigurationRequest struct via the builder pattern


### Return type

[**OpenIDConfigurationResponse**](OpenIDConfigurationResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetUserInfo

> UserInfoResponse GetUserInfo(ctx).Execute()

Userinfo endpoint (Bearer access token from token endpoint, RS256)

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
	resp, r, err := apiClient.OIDCAPI.GetUserInfo(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OIDCAPI.GetUserInfo``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetUserInfo`: UserInfoResponse
	fmt.Fprintf(os.Stdout, "Response from `OIDCAPI.GetUserInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetUserInfoRequest struct via the builder pattern


### Return type

[**UserInfoResponse**](UserInfoResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LoginOidc

> LoginOidc(ctx).LoginRequest(loginRequest).Execute()

Login (OIDC browser flow) and continue /authorize

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
	loginRequest := *openapiclient.NewLoginRequest("Email_example", "Password_example") // LoginRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OIDCAPI.LoginOidc(context.Background()).LoginRequest(loginRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OIDCAPI.LoginOidc``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLoginOidcRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginRequest** | [**LoginRequest**](LoginRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[CSRFToken](../README.md#CSRFToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

