package auth

import (
	"fmt"

	"github.com/gateforge-iam/gateforge-iam/internal/config"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// ProvideWebAuthn builds the go-webauthn instance from application config.
func ProvideWebAuthn(cfg *config.Config) (*webauthn.WebAuthn, error) {
	if len(cfg.WebauthnRPOrigins) == 0 {
		return nil, fmt.Errorf("WEBAUTHN_RP_ORIGINS must contain at least one origin (comma-separated)")
	}
	return webauthn.New(&webauthn.Config{
		RPID:          cfg.WebauthnRPID,
		RPDisplayName: cfg.WebauthnRPDisplayName,
		RPOrigins:     cfg.WebauthnRPOrigins,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			UserVerification: protocol.VerificationPreferred,
		},
		AttestationPreference: protocol.PreferNoAttestation,
	})
}
