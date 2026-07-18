package middlewares

import (
	"net/http"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// RateLimit creates a rate limiting middleware with custom configuration
func RateLimit(config config.Config) echo.MiddlewareFunc {
	rateLimit := rate.Limit(float64(config.RateLimit) / config.RateLimitDuration.Seconds())
	store := echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
		echoMiddleware.RateLimiterMemoryStoreConfig{
			Rate:      rateLimit,
			Burst:     config.RateLimit,
			ExpiresIn: config.RateLimitDuration,
		},
	)

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
func DefaultRateLimit(cfg config.Config) echo.MiddlewareFunc {
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
	})
}

// StrictRateLimit creates a strict rate limiting middleware (5 requests per minute)
func StrictRateLimit() echo.MiddlewareFunc {
	return RateLimit(config.Config{
		RateLimit:         5,
		RateLimitDuration: time.Minute,
	})
}

// AuthRateLimit creates rate limiting for authentication endpoints from config
// (AUTH_RATE_LIMIT per minute; default 3/min).
func AuthRateLimit(cfg config.Config) echo.MiddlewareFunc {
	limit := cfg.AuthRateLimit
	if limit <= 0 {
		limit = 3
	}
	return RateLimit(config.Config{
		RateLimit:         limit,
		RateLimitDuration: time.Minute,
	})
}

// PublicRateLimit creates rate limiting for public endpoints from config
// (PUBLIC_RATE_LIMIT per minute; default 100/min).
func PublicRateLimit(cfg config.Config) echo.MiddlewareFunc {
	limit := cfg.PublicRateLimit
	if limit <= 0 {
		limit = 100
	}
	return RateLimit(config.Config{
		RateLimit:         limit,
		RateLimitDuration: time.Minute,
	})
}
