package cache

import (
	"context"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/errors"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"
)

func TestNewRedisCache_ConnectionSuccess(t *testing.T) {
	mr := miniredis.RunT(t)
	rc, err := NewRedisCache(redisTestConfig(mr))
	require.NoError(t, err)
	require.NotNil(t, rc)
	require.NoError(t, rc.Close())
}

func TestNewRedisCache_ConnectionFailure(t *testing.T) {
	cfg := redisTestConfig(miniredis.RunT(t))
	cfg.RedisPort = "1"

	_, err := NewRedisCache(cfg)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeCache, appErr.Type)
}

func TestRedisCache_SetGetDeleteExists(t *testing.T) {
	mr := miniredis.RunT(t)
	rc, err := NewRedisCache(redisTestConfig(mr))
	require.NoError(t, err)
	t.Cleanup(func() { _ = rc.Close() })

	ctx := context.Background()
	key := "test:key"
	value := "hello"

	exists, err := rc.Exists(ctx, key)
	require.NoError(t, err)
	require.False(t, exists)

	require.NoError(t, rc.Set(ctx, key, value, time.Minute))

	got, err := rc.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, got)

	exists, err = rc.Exists(ctx, key)
	require.NoError(t, err)
	require.True(t, exists)

	require.NoError(t, rc.Delete(ctx, key))

	exists, err = rc.Exists(ctx, key)
	require.NoError(t, err)
	require.False(t, exists)
}

func TestRedisCache_GetNotFound(t *testing.T) {
	mr := miniredis.RunT(t)
	rc, err := NewRedisCache(redisTestConfig(mr))
	require.NoError(t, err)
	t.Cleanup(func() { _ = rc.Close() })

	_, err = rc.Get(context.Background(), "missing")
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeNotFound, appErr.Type)
}

func TestRedisCache_SetError(t *testing.T) {
	mr := miniredis.RunT(t)
	rc, err := NewRedisCache(redisTestConfig(mr))
	require.NoError(t, err)
	t.Cleanup(func() { _ = rc.Close() })

	mr.SetError("SET")
	err = rc.Set(context.Background(), "key", "value", time.Minute)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeCache, appErr.Type)
}

func TestRedisCache_GetError(t *testing.T) {
	mr := miniredis.RunT(t)
	rc, err := NewRedisCache(redisTestConfig(mr))
	require.NoError(t, err)
	t.Cleanup(func() { _ = rc.Close() })

	mr.SetError("GET")
	_, err = rc.Get(context.Background(), "key")
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeCache, appErr.Type)
}

func TestRedisCache_DeleteError(t *testing.T) {
	mr := miniredis.RunT(t)
	rc, err := NewRedisCache(redisTestConfig(mr))
	require.NoError(t, err)
	t.Cleanup(func() { _ = rc.Close() })

	mr.SetError("DEL")
	err = rc.Delete(context.Background(), "key")
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeCache, appErr.Type)
}

func TestRedisCache_ExistsError(t *testing.T) {
	mr := miniredis.RunT(t)
	rc, err := NewRedisCache(redisTestConfig(mr))
	require.NoError(t, err)
	t.Cleanup(func() { _ = rc.Close() })

	mr.SetError("EXISTS")
	_, err = rc.Exists(context.Background(), "key")
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeCache, appErr.Type)
}

func TestRedisCache_Close(t *testing.T) {
	mr := miniredis.RunT(t)
	rc, err := NewRedisCache(redisTestConfig(mr))
	require.NoError(t, err)
	require.NoError(t, rc.Close())
}
