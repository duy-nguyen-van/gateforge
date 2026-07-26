# MFAChallengeVerifyRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | **string** |  | 
**MfaTicket** | **string** |  | 

## Methods

### NewMFAChallengeVerifyRequest

`func NewMFAChallengeVerifyRequest(code string, mfaTicket string, ) *MFAChallengeVerifyRequest`

NewMFAChallengeVerifyRequest instantiates a new MFAChallengeVerifyRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMFAChallengeVerifyRequestWithDefaults

`func NewMFAChallengeVerifyRequestWithDefaults() *MFAChallengeVerifyRequest`

NewMFAChallengeVerifyRequestWithDefaults instantiates a new MFAChallengeVerifyRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCode

`func (o *MFAChallengeVerifyRequest) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *MFAChallengeVerifyRequest) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *MFAChallengeVerifyRequest) SetCode(v string)`

SetCode sets Code field to given value.


### GetMfaTicket

`func (o *MFAChallengeVerifyRequest) GetMfaTicket() string`

GetMfaTicket returns the MfaTicket field if non-nil, zero value otherwise.

### GetMfaTicketOk

`func (o *MFAChallengeVerifyRequest) GetMfaTicketOk() (*string, bool)`

GetMfaTicketOk returns a tuple with the MfaTicket field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaTicket

`func (o *MFAChallengeVerifyRequest) SetMfaTicket(v string)`

SetMfaTicket sets MfaTicket field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


