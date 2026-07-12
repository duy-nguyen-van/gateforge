package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/cache"
	apperrors "github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestProvideTokenService_InvalidSecret(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.JWTSecret = "short"
	_, err := ProvideTokenService(cfg)
	require.Error(t, err)
}

func TestTokenService_ParseAccessToken_InvalidClaims(t *testing.T) {
	svc, err := NewTokenService(testJWTSecret, "issuer", time.Hour)
	require.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(testJWTSecret))
	require.NoError(t, err)

	_, _, err = svc.ParseAccessToken(signed)
	require.Error(t, err)
}

func TestTokenService_ParseAccessToken_WrongSigningMethod(t *testing.T) {
	svc, err := NewTokenService(testJWTSecret, "issuer", time.Hour)
	require.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodNone, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"},
	})
	unsigned, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, _, err = svc.ParseAccessToken(unsigned)
	require.Error(t, err)
}

func TestTokenService_ParseSelectionToken_Errors(t *testing.T) {
	svc, err := NewTokenService(testJWTSecret, "issuer", time.Hour)
	require.NoError(t, err)

	_, err = svc.ParseSelectionToken("not-a-jwt")
	require.Error(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(testJWTSecret))
	require.NoError(t, err)
	_, err = svc.ParseSelectionToken(signed)
	require.Error(t, err)
}


type errCache struct{}

func (errCache) Get(context.Context, string) (string, error) {
	return "", errors.New("cache down")
}

func (errCache) Set(context.Context, string, string, time.Duration) error {
	return errors.New("cache down")
}

func (errCache) Delete(context.Context, string) error { return errors.New("cache down") }

func (errCache) Exists(context.Context, string) (bool, error) { return false, errors.New("cache down") }

func (errCache) Close() error { return nil }

func TestEphemeralStore_putJSON_CacheError(t *testing.T) {
	cfg := testutil.TestConfig()
	store := NewEphemeralStore(errCache{}, cfg)

	_, err := store.PutWebauthnRegistrationSession(context.Background(), map[string]string{"k": "v"})
	require.Error(t, err)
}

func TestEphemeralStore_takeJSON_CacheError(t *testing.T) {
	cfg := testutil.TestConfig()
	store := NewEphemeralStore(errCache{}, cfg)

	var out map[string]string
	err := store.TakeWebauthnRegistrationSession(context.Background(), "token", &out)
	require.Error(t, err)
}

func TestEphemeralStore_takeJSON_InvalidPayload(t *testing.T) {
	cfg := testutil.TestConfig()
	c := newFakeCache()
	key := fmt.Sprintf(redisKeyWebauthnReg, "bad-token")
	require.NoError(t, c.Set(context.Background(), key, "not-json", time.Minute))

	store := NewEphemeralStore(c, cfg)
	var out map[string]string
	err := store.TakeWebauthnRegistrationSession(context.Background(), "bad-token", &out)
	require.Error(t, err)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, apperrors.ErrorTypeValidation, appErr.Type)
}

func TestEphemeralStore_PutMFAPending_CacheError(t *testing.T) {
	cfg := testutil.TestConfig()
	store := NewEphemeralStore(errCache{}, cfg)
	_, err := store.PutMFAPending(context.Background(), MFAPendingPayload{UserID: "u1"})
	require.Error(t, err)
}

func TestEphemeralStore_TakeMFAPending_CacheError(t *testing.T) {
	cfg := testutil.TestConfig()
	store := NewEphemeralStore(errCache{}, cfg)
	_, err := store.TakeMFAPending(context.Background(), "ticket")
	require.Error(t, err)
}

func TestEphemeralStore_TakeMFAPending_InvalidPayload(t *testing.T) {
	cfg := testutil.TestConfig()
	c := newFakeCache()
	key := fmt.Sprintf(redisKeyMFAPending, "ticket-1")
	require.NoError(t, c.Set(context.Background(), key, "{bad", time.Minute))

	store := NewEphemeralStore(c, cfg)
	_, err := store.TakeMFAPending(context.Background(), "ticket-1")
	require.Error(t, err)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, apperrors.ErrorTypeValidation, appErr.Type)
}

func TestEphemeralStore_putJSON_MarshalError(t *testing.T) {
	cfg := testutil.TestConfig()
	store := NewEphemeralStore(newFakeCache(), cfg)
	_, err := store.putJSON(context.Background(), redisKeyWebauthnReg, make(chan int), time.Minute)
	require.Error(t, err)
}

// Ensure errCache satisfies cache.Cache.
var _ cache.Cache = errCache{}
