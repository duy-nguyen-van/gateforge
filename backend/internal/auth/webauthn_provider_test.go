package auth

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestProvideWebAuthn_Success(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.WebauthnRPOrigins = []string{"http://localhost:3000"}
	cfg.WebauthnRPID = "localhost"
	cfg.WebauthnRPDisplayName = "IAM Test"

	w, err := ProvideWebAuthn(cfg)
	require.NoError(t, err)
	require.NotNil(t, w)
}

func TestProvideWebAuthn_MissingOrigins(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.WebauthnRPOrigins = nil

	_, err := ProvideWebAuthn(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "WEBAUTHN_RP_ORIGINS")
}
