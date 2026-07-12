package auth

import (
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

const testJWTSecret = "test-jwt-secret-at-least-thirty-two-bytes"

func TestNewTokenService_Validation(t *testing.T) {
	_, err := NewTokenService("short", "issuer", time.Hour)
	require.Error(t, err)

	_, err = NewTokenService(testJWTSecret, "issuer", 0)
	require.Error(t, err)
}

func TestTokenService_SignAndParseAccessToken(t *testing.T) {
	svc, err := NewTokenService(testJWTSecret, "gateforge-iam", time.Hour)
	require.NoError(t, err)

	token, expiresAt, err := svc.SignAccessToken("user-1", "tenant-1")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.True(t, expiresAt.After(time.Now()))

	userID, tenantID, err := svc.ParseAccessToken(token)
	require.NoError(t, err)
	require.Equal(t, "user-1", userID)
	require.Equal(t, "tenant-1", tenantID)
}

func TestTokenService_SignAndParseSelectionToken(t *testing.T) {
	svc, err := NewTokenService(testJWTSecret, "gateforge-iam", time.Hour)
	require.NoError(t, err)

	token, expiresIn, err := svc.SignSelectionToken("user-1")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Positive(t, expiresIn)

	userID, err := svc.ParseSelectionToken(token)
	require.NoError(t, err)
	require.Equal(t, "user-1", userID)
}

func TestTokenService_ParseAccessToken_Errors(t *testing.T) {
	svc, err := NewTokenService(testJWTSecret, "gateforge-iam", time.Hour)
	require.NoError(t, err)

	_, _, err = svc.ParseAccessToken("not-a-jwt")
	require.Error(t, err)

	other, err := NewTokenService("another-secret-thirty-two-bytes-xx", "issuer", time.Hour)
	require.NoError(t, err)
	token, _, err := other.SignAccessToken("user-1", "tenant-1")
	require.NoError(t, err)
	_, _, err = svc.ParseAccessToken(token)
	require.Error(t, err)
}

func TestProvideTokenService(t *testing.T) {
	cfg := testutil.TestConfig()
	svc, err := ProvideTokenService(cfg)
	require.NoError(t, err)
	require.NotNil(t, svc)
}

func TestOpaqueRefreshToken(t *testing.T) {
	raw, hash, err := NewOpaqueRefreshToken()
	require.NoError(t, err)
	require.NotEmpty(t, raw)
	require.NotEmpty(t, hash)
	require.Equal(t, hash, HashOpaqueToken(raw))
}
