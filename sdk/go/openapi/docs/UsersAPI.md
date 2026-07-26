# \UsersAPI

All URIs are relative to *http://localhost:3000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetMe**](UsersAPI.md#GetMe) | **Get** /api/v1/me | Current user profile
[**UpdateMe**](UsersAPI.md#UpdateMe) | **Patch** /api/v1/me | Update current user profile



## GetMe

> UserResponseEnvelope GetMe(ctx).Execute()

Current user profile

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
	resp, r, err := apiClient.UsersAPI.GetMe(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UsersAPI.GetMe``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMe`: UserResponseEnvelope
	fmt.Fprintf(os.Stdout, "Response from `UsersAPI.GetMe`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetMeRequest struct via the builder pattern


### Return type

[**UserResponseEnvelope**](UserResponseEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMe

> UserResponseEnvelope UpdateMe(ctx).UpdateProfileRequest(updateProfileRequest).Execute()

Update current user profile

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
	updateProfileRequest := *openapiclient.NewUpdateProfileRequest() // UpdateProfileRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UsersAPI.UpdateMe(context.Background()).UpdateProfileRequest(updateProfileRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UsersAPI.UpdateMe``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMe`: UserResponseEnvelope
	fmt.Fprintf(os.Stdout, "Response from `UsersAPI.UpdateMe`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateProfileRequest** | [**UpdateProfileRequest**](UpdateProfileRequest.md) |  | 

### Return type

[**UserResponseEnvelope**](UserResponseEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

