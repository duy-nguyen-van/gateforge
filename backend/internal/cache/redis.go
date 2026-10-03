package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	appErrors "github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache is the Fx provider for RedisCache
func NewRedisCache(cfg *config.Config) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password:        cfg.RedisPassword,
		DB:              cfg.RedisDB,
		PoolSize:        cfg.PoolSize,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		PoolTimeout:     cfg.PoolTimeout,
		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: cfg.MinRetryBackoff,
		MaxRetryBackoff: cfg.MaxRetryBackoff,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, appErrors.CacheError("Failed to connect to Redis", err).
			WithOperation("connect_redis").
			WithResource("cache")
	}

	if monitoring.IsOTelEnabled(*cfg) {
		var instrumentErr error
		if cfg.OTelTracesEnabled {
			instrumentErr = errors.Join(instrumentErr, redisotel.InstrumentTracing(client))
		}
		if cfg.OTelMetricsEnabled {
			instrumentErr = errors.Join(instrumentErr, redisotel.InstrumentMetrics(client))
		}
		if instrumentErr != nil {
			logger.Sugar.Warnf("Failed to instrument Redis with OpenTelemetry: %v", instrumentErr)
		}
	}

	return &RedisCache{
		client: client,
	}, nil
}

// Get retrieves a value from Redis
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result := r.client.Get(ctx, key)
	if result.Err() != nil {
		if errors.Is(result.Err(), redis.Nil) {
			return "", appErrors.NotFoundError("Cache key", fmt.Errorf("key not found")).
				WithOperation("get_cache").
				WithResource("cache").
				WithContext("key", key)
		}
		return "", appErrors.CacheError("Failed to get from cache", result.Err()).
			WithOperation("get_cache").
			WithResource("cache").
			WithContext("key", key)
	}
	return result.Val(), nil
}

// Set stores a value in Redis with expiration
func (r *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return appErrors.CacheError("Failed to set cache", err).
			WithOperation("set_cache").
			WithResource("cache").
			WithContext("key", key)
	}
	return nil
}

// Delete removes a value from Redis
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return appErrors.CacheError("Failed to delete from cache", err).
			WithOperation("delete_cache").
			WithResource("cache").
			WithContext("key", key)
	}
	return nil
}

// Increment adds one to a counter and sets the window on the first hit.
func (r *RedisCache) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	n, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, appErrors.CacheError("Failed to increment cache key", err).
			WithOperation("increment_cache").
			WithResource("cache").
			WithContext("key", key)
	}
	if n == 1 && window > 0 {
		if err := r.client.PExpire(ctx, key, window).Err(); err != nil {
			return 0, appErrors.CacheError("Failed to expire cache key", err).
				WithOperation("increment_cache").
				WithResource("cache").
				WithContext("key", key)
		}
	}
	return n, nil
}

// Ping checks that Redis is reachable.
func (r *RedisCache) Ping(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return appErrors.CacheError("Failed to ping cache", err).
			WithOperation("ping_cache").
			WithResource("cache")
	}
	return nil
}

// Exists checks if a key exists in Redis
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result := r.client.Exists(ctx, key)
	if result.Err() != nil {
		return false, appErrors.CacheError("Failed to check cache existence", result.Err()).
			WithOperation("exists_cache").
			WithResource("cache").
			WithContext("key", key)
	}
	return result.Val() > 0, nil
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	err := r.client.Close()
	if err != nil {
		return appErrors.CacheError("Failed to close Redis connection", err).
			WithOperation("close_cache").
			WithResource("cache")
	}
	return nil
}
