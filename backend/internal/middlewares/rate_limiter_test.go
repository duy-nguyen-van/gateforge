package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestRateLimit(t *testing.T) {
	cfg := config.Config{
		RateLimit:         2,
		RateLimitDuration: time.Second,
	}

	e := echo.New()
	e.Use(RateLimit(cfg))
	e.GET("/limited", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/limited", nil)
		req.RemoteAddr = "192.0.2.1:1234"
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code, "request %d should succeed", i+1)
	}

	req := httptest.NewRequest(http.MethodGet, "/limited", nil)
	req.RemoteAddr = "192.0.2.1:1234"
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusTooManyRequests, rec.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, constants.RateLimitExceeded, body["error_code"])
}

func TestDefaultRateLimit(t *testing.T) {
	e := echo.New()
	e.Use(DefaultRateLimit(config.Config{DefaultRateLimit: 20, RateLimitDuration: time.Second}))
	e.GET("/default", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/default", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestStrictRateLimit(t *testing.T) {
	e := echo.New()
	e.Use(StrictRateLimit())
	e.GET("/strict", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/strict", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthRateLimit(t *testing.T) {
	e := echo.New()
	e.Use(AuthRateLimit(config.Config{AuthRateLimit: 3}))
	e.POST("/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestPublicRateLimit(t *testing.T) {
	e := echo.New()
	e.Use(PublicRateLimit(config.Config{PublicRateLimit: 100}))
	e.GET("/public", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestRateLimit_ErrorHandler(t *testing.T) {
	cfg := config.Config{
		RateLimit:         1,
		RateLimitDuration: time.Minute,
	}

	e := echo.New()
	e.Use(RateLimit(cfg))
	e.GET("/err", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Exhaust burst
	req := httptest.NewRequest(http.MethodGet, "/err", nil)
	req.RemoteAddr = "198.51.100.10:4321"
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// Trigger deny handler
	req = httptest.NewRequest(http.MethodGet, "/err", nil)
	req.RemoteAddr = "198.51.100.10:4321"
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusTooManyRequests, rec.Code)
}
