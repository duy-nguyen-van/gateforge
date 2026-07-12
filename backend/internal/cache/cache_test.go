package cache

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"
)

func TestProvideCache_Redis(t *testing.T) {
	mr := miniredis.RunT(t)
	cfg := redisTestConfig(mr)

	c, err := ProvideCache(cfg)
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NoError(t, c.Close())
}

func TestProvideCache_InvalidProvider(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.CacheProvider = "memory"

	_, err := ProvideCache(cfg)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeInternal, appErr.Type)
}

func TestProvideCache_RedisConnectionFailure(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.CacheProvider = constants.CacheProviderRedis
	cfg.RedisHost = "127.0.0.1"
	cfg.RedisPort = "1"

	_, err := ProvideCache(cfg)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeCache, appErr.Type)
}

func redisTestConfig(mr *miniredis.Miniredis) *config.Config {
	cfg := testutil.TestConfig()
	cfg.RedisHost = mr.Host()
	cfg.RedisPort = mr.Port()
	return cfg
}
