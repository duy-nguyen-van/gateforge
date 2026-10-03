package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/cache"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type redisRateLimiterStore struct {
	cache  cache.Cache
	limit  int
	window time.Duration
	group  string
}

func (s *redisRateLimiterStore) Allow(identifier string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	n, err := s.cache.Increment(ctx, "iam:ratelimit:"+s.group+":"+identifier, s.window)
	if err != nil || int(n) > s.limit {
		monitoring.AddRateLimitDenied(ctx, s.group)
		return false, nil
	}
	return true, nil
}

// RateLimit creates a rate limiting middleware with custom configuration.
// A non-nil cache shares the limit across processes. Otherwise the limit is in-memory.
func RateLimit(config config.Config, shared cache.Cache, group string) echo.MiddlewareFunc {
	var store echoMiddleware.RateLimiterStore
	if shared != nil {
		store = &redisRateLimiterStore{
			cache:  shared,
			limit:  config.RateLimit,
			window: config.RateLimitDuration,
			group:  group,
		}
	} else {
		rateLimit := rate.Limit(float64(config.RateLimit) / config.RateLimitDuration.Seconds())
		store = echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
			echoMiddleware.RateLimiterMemoryStoreConfig{
				Rate:      rateLimit,
				Burst:     config.RateLimit,
				ExpiresIn: config.RateLimitDuration,
			},
		)
	}

	return echoMiddleware.RateLimiterWithConfig(echoMiddleware.RateLimiterConfig{
		Store: store,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message_code": constants.RateLimitExceeded,
				"message":      "Rate limit exceeded",
				"limit":        config.RateLimit,
				"window":       config.RateLimitDuration.String(),
			})
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			monitoring.AddRateLimitDenied(c.Request().Context(), group)
			return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
				"error_code": constants.RateLimitExceeded,
				"message":    "Rate limit exceeded",
				"limit":      config.RateLimit,
				"window":     config.RateLimitDuration.String(),
			})
		},
	})
}

// DefaultRateLimit creates the global rate limiting middleware from config
// (DEFAULT_RATE_LIMIT per RATE_LIMIT_DURATION; defaults 20/s).
func DefaultRateLimit(cfg config.Config, shared cache.Cache) echo.MiddlewareFunc {
	limit := cfg.DefaultRateLimit
	if limit <= 0 {
		limit = 20
	}
	window := cfg.RateLimitDuration
	if window <= 0 {
		window = time.Second
	}
	return RateLimit(config.Config{
		RateLimit:         limit,
		RateLimitDuration: window,
	}, shared, "default")
}

// StrictRateLimit creates a strict rate limiting middleware (5 requests per minute)
func StrictRateLimit(shared cache.Cache) echo.MiddlewareFunc {
	return RateLimit(config.Config{
		RateLimit:         5,
		RateLimitDuration: time.Minute,
	}, shared, "auth")
}

// AuthRateLimit creates rate limiting for authentication endpoints from config
// (AUTH_RATE_LIMIT per minute; default 3/min).
func AuthRateLimit(cfg config.Config, shared cache.Cache) echo.MiddlewareFunc {
	limit := cfg.AuthRateLimit
	if limit <= 0 {
		limit = 3
	}
	return RateLimit(config.Config{
		RateLimit:         limit,
		RateLimitDuration: time.Minute,
	}, shared, "auth")
}

// PublicRateLimit creates rate limiting for public endpoints from config
// (PUBLIC_RATE_LIMIT per minute; default 100/min).
func PublicRateLimit(cfg config.Config, shared cache.Cache) echo.MiddlewareFunc {
	limit := cfg.PublicRateLimit
	if limit <= 0 {
		limit = 100
	}
	return RateLimit(config.Config{
		RateLimit:         limit,
		RateLimitDuration: time.Minute,
	}, shared, "public")
}
