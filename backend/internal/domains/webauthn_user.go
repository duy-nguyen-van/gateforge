package domains

import (
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

// WebAuthnUser adapts models.User for github.com/go-webauthn/webauthn.
type WebAuthnUser struct {
	User *models.User
}

func (u *WebAuthnUser) WebAuthnID() []byte {
	id, err := uuid.Parse(u.User.ID)
	if err != nil {
		b := []byte(u.User.ID)
		if len(b) > 64 {
			return b[:64]
		}
		return b
	}
	return id[:]
}

func (u *WebAuthnUser) WebAuthnName() string {
	return u.User.Email
}

func (u *WebAuthnUser) WebAuthnDisplayName() string {
	s := strings.TrimSpace(u.User.FirstName + " " + u.User.LastName)
	if s == "" {
		return u.User.Email
	}
	return s
}

func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	var out []webauthn.Credential
	for i := range u.User.WebauthnCredentials {
		c, err := repositories.UnmarshalWebauthnCredential(u.User.WebauthnCredentials[i].PublicKey)
		if err != nil {
			continue
		}
		out = append(out, *c)
	}
	return out
}
