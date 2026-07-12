package dtos

// MFALoginChallengeResponse is returned when password or passkey succeeded but MFA is required.
type MFALoginChallengeResponse struct {
	MfaRequired bool   `json:"mfa_required" example:"true"`
	MfaTicket   string `json:"mfa_ticket" example:"550e8400-e29b-41d4-a716-446655440000"`
	ExpiresIn   int64  `json:"expires_in" example:"600"`
}

// MFATOTPSetupResponse is returned from POST /mfa/totp/setup.
type MFATOTPSetupResponse struct {
	Secret     string `json:"secret" example:"JBSWY3DPEHPK3PXP"`
	OtpauthURI string `json:"otpauth_uri" example:"otpauth://totp/..."`
}

// MFATOTPVerifyRequest confirms TOTP enrollment with a 6-digit code.
type MFATOTPVerifyRequest struct {
	Code string `json:"code" validate:"required,len=6" example:"123456"`
}

// MFAChallengeVerifyRequest exchanges an MFA ticket for tokens after password/passkey + MFA.
// Code is a 6-digit TOTP code, or a recovery code (from POST /mfa/recovery-codes) when not using TOTP.
type MFAChallengeVerifyRequest struct {
	MfaTicket string `json:"mfa_ticket" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Code      string `json:"code" validate:"required" example:"123456"`
}

// MFARecoveryCodesResponse returns one-time recovery codes (plain text once).
type MFARecoveryCodesResponse struct {
	Codes []string `json:"codes"`
}
