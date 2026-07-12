// Package testutil provides shared helpers for backend unit tests.
package testutil

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var initLoggerOnce sync.Once

// InitLogger sets a nop logger so packages that log during tests do not panic.
func InitLogger() {
	initLoggerOnce.Do(func() {
		logger.Log = zap.NewNop()
		logger.Sugar = logger.Log.Sugar()
	})
}

// TestConfig returns a minimal config suitable for unit tests.
func TestConfig() *config.Config {
	return &config.Config{
		AppName:         "gateforge-iam-test",
		AppBaseURL:      "http://localhost:8080",
		JWTSecret:       "test-jwt-secret-at-least-thirty-two-bytes",
		JWTAccessTTL:    time.Hour,
		DefaultTenantID: "default-tenant",
		CacheProvider:   constants.CacheProviderRedis,
		RedisHost:       "127.0.0.1",
		RedisPort:       "6379",
		Environment:     "test",
	}
}

// NewEchoContext builds an Echo context for handler tests.
func NewEchoContext(method, path, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
