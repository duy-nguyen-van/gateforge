package dtos

import (
	"encoding/json"
	"time"
)

// WebauthnCredentialResponse is a registered passkey visible to the owning user.
type WebauthnCredentialResponse struct {
	ID         string    `json:"id" example:"00000000-0000-0000-0000-000000000001"`
	DeviceName string    `json:"device_name" example:"Work MacBook"`
	CreatedAt  time.Time `json:"created_at"`
}

// WebauthnRegisterStartRequest is the body for POST /webauthn/register/start.
type WebauthnRegisterStartRequest struct {
	// DeviceName optional label stored with the passkey (e.g. "Work MacBook").
	DeviceName string `json:"device_name,omitempty" validate:"omitempty,max=255" example:"Personal iPhone"`
}

// WebauthnRegisterStartResponse returns creation options plus an opaque session token for finish.
// options follows W3C PublicKeyCredentialCreationOptions (see WebAuthn Level 2).
type WebauthnRegisterStartResponse struct {
	Options      json.RawMessage `json:"options" swaggertype:"object"`
	SessionToken string          `json:"session_token" example:"9b2c4e8a-1d3f-4a5b-8c7e-0f1a2b3c4d5e"`
}

// WebauthnRegisterFinishRequest completes passkey registration.
type WebauthnRegisterFinishRequest struct {
	SessionToken string `json:"session_token" validate:"required" example:"9b2c4e8a-1d3f-4a5b-8c7e-0f1a2b3c4d5e"`
	// Credential is the browser PublicKeyCredential JSON from navigator.credentials.create().
	Credential json.RawMessage `json:"credential" swaggertype:"object"`
}

// WebauthnLoginStartRequest begins passkey login for a user identified by email.
type WebauthnLoginStartRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	TenantID string `json:"tenant_id,omitempty" validate:"omitempty,uuid" example:"00000000-0000-0000-0000-000000000001"`
}

// WebauthnLoginStartResponse returns assertion options and session token.
// options follows W3C PublicKeyCredentialRequestOptions (see WebAuthn Level 2).
type WebauthnLoginStartResponse struct {
	Options      json.RawMessage `json:"options" swaggertype:"object"`
	SessionToken string          `json:"session_token" example:"9b2c4e8a-1d3f-4a5b-8c7e-0f1a2b3c4d5e"`
}

// WebauthnLoginFinishRequest completes passkey login.
type WebauthnLoginFinishRequest struct {
	Email        string `json:"email" validate:"required,email" example:"user@example.com"`
	TenantID     string `json:"tenant_id,omitempty" validate:"omitempty,uuid" example:"00000000-0000-0000-0000-000000000001"`
	SessionToken string `json:"session_token" validate:"required" example:"9b2c4e8a-1d3f-4a5b-8c7e-0f1a2b3c4d5e"`
	// Credential is the browser PublicKeyCredential JSON from navigator.credentials.get().
	Credential json.RawMessage `json:"credential" swaggertype:"object"`
	RememberMe bool            `json:"remember_me,omitempty" example:"false"`
	ReturnTo   string          `json:"return_to,omitempty" validate:"omitempty" example:""`
}
