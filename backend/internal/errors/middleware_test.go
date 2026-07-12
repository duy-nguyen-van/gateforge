package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/request"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"
	"github.com/gateforge-iam/gateforge-iam/internal/utils/i18n"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func init() {
	testutil.InitLogger()
	i18n.Init()
}

func TestRouteMatchers(t *testing.T) {
	require.True(t, isSwaggerRoute("/swagger/index.html"))
	require.True(t, isSwaggerRoute("/api/docs/openapi"))
	require.True(t, isFaviconRoute("/favicon.ico"))
	require.True(t, isFaviconRoute("/assets/favicon.png"))
	require.False(t, isSwaggerRoute("/api/v1/users"))
	require.False(t, isFaviconRoute("/api/v1/users"))
}

func TestRecoveryMiddleware(t *testing.T) {
	t.Run("recovers string panic in production", func(t *testing.T) {
		e := echo.New()
		cfg := &config.Config{AppEnv: config.EnvironmentProduction}
		e.Use(RecoveryMiddleware(cfg))
		e.GET("/panic", func(echo.Context) error {
			panic("boom")
		})

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp dtos.BaseResponse[map[string]interface{}]
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, constants.InternalError, resp.Meta.ErrorCode)
		require.Nil(t, resp.Data)
	})

	t.Run("recovers error panic in development", func(t *testing.T) {
		e := echo.New()
		cfg := &config.Config{AppEnv: config.EnvironmentDevelopment}
		e.Use(RecoveryMiddleware(cfg))
		e.GET("/panic", func(echo.Context) error {
			panic(errors.New("dev panic"))
		})

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp dtos.BaseResponse[map[string]interface{}]
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.NotNil(t, resp.Data)
		require.Contains(t, resp.Data, "panic_value")
		require.Contains(t, resp.Data, "stack_trace")
	})

	t.Run("recovers with nil config", func(t *testing.T) {
		e := echo.New()
		e.Use(RecoveryMiddleware(nil))
		e.GET("/panic", func(echo.Context) error {
			panic("nil config panic")
		})

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("passes through successful request", func(t *testing.T) {
		e := echo.New()
		e.Use(RecoveryMiddleware(testutil.TestConfig()))
		e.GET("/ok", func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "ok", rec.Body.String())
	})
}

func TestErrorMiddleware(t *testing.T) {
	t.Run("handles handler errors", func(t *testing.T) {
		e := echo.New()
		e.Use(ErrorMiddleware())
		e.GET("/fail", func(echo.Context) error {
			return ValidationError("bad input", nil)
		})

		req := httptest.NewRequest(http.MethodGet, "/fail", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var resp dtos.BaseResponse[any]
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, constants.ValidationError, resp.Meta.ErrorCode)
	})

	t.Run("skips swagger routes", func(t *testing.T) {
		e := echo.New()
		e.Use(ErrorMiddleware())
		e.GET("/swagger/*", func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusNotFound, "swagger missing")
		})

		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("skips favicon routes", func(t *testing.T) {
		e := echo.New()
		e.Use(ErrorMiddleware())
		e.GET("/favicon.ico", func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusNotFound, "missing icon")
		})

		req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("passes through successful request", func(t *testing.T) {
		e := echo.New()
		e.Use(ErrorMiddleware())
		e.GET("/ok", func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestRecoveryMiddleware_RequestContext(t *testing.T) {
	e := echo.New()
	e.Use(RecoveryMiddleware(&config.Config{AppEnv: config.EnvironmentProduction}))
	e.GET("/panic", func(c echo.Context) error {
		ctx := request.NewCorrelationIDContext(c.Request().Context(), "corr-123")
		ctx = request.NewLanguageCodeContext(ctx, "en")
		c.SetRequest(c.Request().WithContext(ctx))
		panic("with context")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRecoveryMiddleware_WithSentryHub(t *testing.T) {
	e := echo.New()
	e.Use(RecoveryMiddleware(&config.Config{AppEnv: config.EnvironmentProduction}))
	e.GET("/panic", func(c echo.Context) error {
		hub := sentry.CurrentHub().Clone()
		ctx := sentry.SetHubOnContext(c.Request().Context(), hub)
		ctx = request.NewCorrelationIDContext(ctx, "corr-456")
		ctx = request.NewLanguageCodeContext(ctx, "en")
		c.SetRequest(c.Request().WithContext(ctx))
		panic(errors.New("sentry panic"))
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
