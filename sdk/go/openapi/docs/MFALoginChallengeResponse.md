# MFALoginChallengeResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MfaRequired** | **bool** |  | 
**MfaTicket** | **string** |  | 
**ExpiresIn** | **int32** |  | 

## Methods

### NewMFALoginChallengeResponse

`func NewMFALoginChallengeResponse(mfaRequired bool, mfaTicket string, expiresIn int32, ) *MFALoginChallengeResponse`

NewMFALoginChallengeResponse instantiates a new MFALoginChallengeResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMFALoginChallengeResponseWithDefaults

`func NewMFALoginChallengeResponseWithDefaults() *MFALoginChallengeResponse`

NewMFALoginChallengeResponseWithDefaults instantiates a new MFALoginChallengeResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMfaRequired

`func (o *MFALoginChallengeResponse) GetMfaRequired() bool`

GetMfaRequired returns the MfaRequired field if non-nil, zero value otherwise.

### GetMfaRequiredOk

`func (o *MFALoginChallengeResponse) GetMfaRequiredOk() (*bool, bool)`

GetMfaRequiredOk returns a tuple with the MfaRequired field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaRequired

`func (o *MFALoginChallengeResponse) SetMfaRequired(v bool)`

SetMfaRequired sets MfaRequired field to given value.


### GetMfaTicket

`func (o *MFALoginChallengeResponse) GetMfaTicket() string`

GetMfaTicket returns the MfaTicket field if non-nil, zero value otherwise.

### GetMfaTicketOk

`func (o *MFALoginChallengeResponse) GetMfaTicketOk() (*string, bool)`

GetMfaTicketOk returns a tuple with the MfaTicket field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaTicket

`func (o *MFALoginChallengeResponse) SetMfaTicket(v string)`

SetMfaTicket sets MfaTicket field to given value.


### GetExpiresIn

`func (o *MFALoginChallengeResponse) GetExpiresIn() int32`

GetExpiresIn returns the ExpiresIn field if non-nil, zero value otherwise.

### GetExpiresInOk

`func (o *MFALoginChallengeResponse) GetExpiresInOk() (*int32, bool)`

GetExpiresInOk returns a tuple with the ExpiresIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresIn

`func (o *MFALoginChallengeResponse) SetExpiresIn(v int32)`

SetExpiresIn sets ExpiresIn field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


