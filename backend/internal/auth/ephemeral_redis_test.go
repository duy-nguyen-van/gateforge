package auth

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

type fakeCacheItem struct {
	value  string
	expiry time.Time
}

type fakeCache struct {
	mu    sync.Mutex
	items map[string]fakeCacheItem
}

func newFakeCache() *fakeCache {
	return &fakeCache{items: make(map[string]fakeCacheItem)}
}

func (f *fakeCache) Get(_ context.Context, key string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	item, ok := f.items[key]
	if !ok || (!item.expiry.IsZero() && time.Now().After(item.expiry)) {
		return "", errors.NotFoundError("Cache key", fmt.Errorf("key not found")).
			WithOperation("get_cache").
			WithResource("cache").
			WithContext("key", key)
	}
	return item.value, nil
}

func (f *fakeCache) Set(_ context.Context, key, value string, expiration time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	item := fakeCacheItem{value: value}
	if expiration > 0 {
		item.expiry = time.Now().Add(expiration)
	}
	f.items[key] = item
	return nil
}

func (f *fakeCache) Delete(_ context.Context, key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.items, key)
	return nil
}

func (f *fakeCache) Exists(ctx context.Context, key string) (bool, error) {
	_, err := f.Get(ctx, key)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil && appErr.Type == errors.ErrorTypeNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (f *fakeCache) Close() error { return nil }

func TestEphemeralStore_WebauthnRegistrationSession(t *testing.T) {
	cfg := testutil.TestConfig()
	store := NewEphemeralStore(newFakeCache(), cfg)
	ctx := context.Background()

	type sessionPayload struct {
		Challenge string `json:"challenge"`
		UserID    string `json:"user_id"`
	}
	payload := sessionPayload{Challenge: "abc", UserID: "user-1"}

	token, err := store.PutWebauthnRegistrationSession(ctx, payload)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	var out sessionPayload
	require.NoError(t, store.TakeWebauthnRegistrationSession(ctx, token, &out))
	require.Equal(t, payload, out)

	var again sessionPayload
	err = store.TakeWebauthnRegistrationSession(ctx, token, &again)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeUnauthorized, appErr.Type)
}

func TestEphemeralStore_WebauthnLoginSession(t *testing.T) {
	cfg := testutil.TestConfig()
	store := NewEphemeralStore(newFakeCache(), cfg)
	ctx := context.Background()

	type sessionPayload struct {
		SessionID string `json:"session_id"`
	}
	payload := sessionPayload{SessionID: "sess-1"}

	token, err := store.PutWebauthnLoginSession(ctx, payload)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	var out sessionPayload
	require.NoError(t, store.TakeWebauthnLoginSession(ctx, token, &out))
	require.Equal(t, payload, out)

	err = store.TakeWebauthnLoginSession(ctx, token, &out)
	require.Error(t, err)
}

func TestEphemeralStore_MFAPending(t *testing.T) {
	cfg := testutil.TestConfig()
	store := NewEphemeralStore(newFakeCache(), cfg)
	ctx := context.Background()

	payload := MFAPendingPayload{
		UserID:     "user-1",
		TenantID:   "tenant-1",
		RememberMe: true,
		ReturnTo:   "/dashboard",
	}

	ticket, err := store.PutMFAPending(ctx, payload)
	require.NoError(t, err)
	require.NotEmpty(t, ticket)

	got, err := store.TakeMFAPending(ctx, ticket)
	require.NoError(t, err)
	require.Equal(t, payload, *got)

	_, err = store.TakeMFAPending(ctx, ticket)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeUnauthorized, appErr.Type)
}

func TestEphemeralStore_MFAPending_InvalidTicket(t *testing.T) {
	cfg := testutil.TestConfig()
	store := NewEphemeralStore(newFakeCache(), cfg)

	_, err := store.TakeMFAPending(context.Background(), "missing-ticket")
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeUnauthorized, appErr.Type)
}
