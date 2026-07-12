package domains

import (
	"strings"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWebAuthnUser_WebAuthnID_UUID(t *testing.T) {
	id := uuid.NewString()
	user := &WebAuthnUser{User: &models.User{BaseModel: models.BaseModel{ID: id}}}

	got := user.WebAuthnID()
	parsed, err := uuid.FromBytes(got)
	require.NoError(t, err)
	require.Equal(t, id, parsed.String())
}

func TestWebAuthnUser_WebAuthnID_NonUUID(t *testing.T) {
	user := &WebAuthnUser{User: &models.User{BaseModel: models.BaseModel{ID: "legacy-user-id"}}}
	require.Equal(t, []byte("legacy-user-id"), user.WebAuthnID())

	longID := strings.Repeat("x", 80)
	user = &WebAuthnUser{User: &models.User{BaseModel: models.BaseModel{ID: longID}}}
	got := user.WebAuthnID()
	require.Len(t, got, 64)
	require.Equal(t, []byte(longID[:64]), got)
}

func TestWebAuthnUser_WebAuthnName(t *testing.T) {
	user := &WebAuthnUser{User: &models.User{Email: "user@example.com"}}
	require.Equal(t, "user@example.com", user.WebAuthnName())
}

func TestWebAuthnUser_WebAuthnDisplayName(t *testing.T) {
	user := &WebAuthnUser{User: &models.User{
		Email:     "user@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
	}}
	require.Equal(t, "Jane Doe", user.WebAuthnDisplayName())

	user = &WebAuthnUser{User: &models.User{Email: "user@example.com"}}
	require.Equal(t, "user@example.com", user.WebAuthnDisplayName())
}

func TestWebAuthnUser_WebAuthnCredentials(t *testing.T) {
	validCred := webauthn.Credential{ID: []byte("cred-id")}
	validJSON, err := repositories.MarshalWebauthnCredential(&validCred)
	require.NoError(t, err)

	user := &WebAuthnUser{User: &models.User{
		WebauthnCredentials: []models.WebauthnCredential{
			{PublicKey: validJSON},
			{PublicKey: "not-json"},
		},
	}}

	creds := user.WebAuthnCredentials()
	require.Len(t, creds, 1)
	require.Equal(t, validCred.ID, creds[0].ID)
}
