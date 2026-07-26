# \MFAAPI

All URIs are relative to *http://localhost:3000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateMfaRecoveryCodes**](MFAAPI.md#CreateMfaRecoveryCodes) | **Post** /api/v1/mfa/recovery-codes | Generate MFA recovery codes
[**SetupTotp**](MFAAPI.md#SetupTotp) | **Post** /api/v1/mfa/totp/setup | Start TOTP enrollment
[**VerifyMfaChallenge**](MFAAPI.md#VerifyMfaChallenge) | **Post** /api/v1/mfa/challenge/verify | Complete MFA after password or passkey login
[**VerifyTotpEnrollment**](MFAAPI.md#VerifyTotpEnrollment) | **Post** /api/v1/mfa/totp/verify | Confirm TOTP enrollment



## CreateMfaRecoveryCodes

> MFARecoveryCodesEnvelope CreateMfaRecoveryCodes(ctx).Execute()

Generate MFA recovery codes



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
	resp, r, err := apiClient.MFAAPI.CreateMfaRecoveryCodes(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MFAAPI.CreateMfaRecoveryCodes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateMfaRecoveryCodes`: MFARecoveryCodesEnvelope
	fmt.Fprintf(os.Stdout, "Response from `MFAAPI.CreateMfaRecoveryCodes`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiCreateMfaRecoveryCodesRequest struct via the builder pattern


### Return type

[**MFARecoveryCodesEnvelope**](MFARecoveryCodesEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetupTotp

> MFATOTPSetupEnvelope SetupTotp(ctx).Execute()

Start TOTP enrollment



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
	resp, r, err := apiClient.MFAAPI.SetupTotp(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MFAAPI.SetupTotp``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetupTotp`: MFATOTPSetupEnvelope
	fmt.Fprintf(os.Stdout, "Response from `MFAAPI.SetupTotp`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiSetupTotpRequest struct via the builder pattern


### Return type

[**MFATOTPSetupEnvelope**](MFATOTPSetupEnvelope.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## VerifyMfaChallenge

> LoginResponseEnvelope VerifyMfaChallenge(ctx).MFAChallengeVerifyRequest(mFAChallengeVerifyRequest).Execute()

Complete MFA after password or passkey login



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
	mFAChallengeVerifyRequest := *openapiclient.NewMFAChallengeVerifyRequest("123456", "550e8400-e29b-41d4-a716-446655440000") // MFAChallengeVerifyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MFAAPI.VerifyMfaChallenge(context.Background()).MFAChallengeVerifyRequest(mFAChallengeVerifyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MFAAPI.VerifyMfaChallenge``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `VerifyMfaChallenge`: LoginResponseEnvelope
	fmt.Fprintf(os.Stdout, "Response from `MFAAPI.VerifyMfaChallenge`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiVerifyMfaChallengeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mFAChallengeVerifyRequest** | [**MFAChallengeVerifyRequest**](MFAChallengeVerifyRequest.md) |  | 

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


## VerifyTotpEnrollment

> EmptyDataEnvelope VerifyTotpEnrollment(ctx).MFATOTPVerifyRequest(mFATOTPVerifyRequest).Execute()

Confirm TOTP enrollment



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
	mFATOTPVerifyRequest := *openapiclient.NewMFATOTPVerifyRequest("123456") // MFATOTPVerifyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MFAAPI.VerifyTotpEnrollment(context.Background()).MFATOTPVerifyRequest(mFATOTPVerifyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MFAAPI.VerifyTotpEnrollment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `VerifyTotpEnrollment`: EmptyDataEnvelope
	fmt.Fprintf(os.Stdout, "Response from `MFAAPI.VerifyTotpEnrollment`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiVerifyTotpEnrollmentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mFATOTPVerifyRequest** | [**MFATOTPVerifyRequest**](MFATOTPVerifyRequest.md) |  | 

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

