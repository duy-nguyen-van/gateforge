package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/require"
)

func TestCSRF(t *testing.T) {
	cfg := &config.Config{AppEnv: config.EnvironmentDevelopment}

	skippedPaths := []string{
		"/api/v1/register",
		"/api/v1/login",
		"/api/v1/login/session",
		"/api/v1/refresh",
		"/api/v1/logout",
		"/api/v1/webauthn/register/begin",
		"/api/v1/mfa/setup",
		"/token",
		"/authorize",
		"/userinfo",
		"/.well-known/openid-configuration",
		"/oidc/federation/google/start",
	}

	for _, path := range skippedPaths {
		t.Run("skips "+path, func(t *testing.T) {
			e := echo.New()
			e.Use(CSRF(cfg))
			e.Any(path, func(c echo.Context) error {
				return c.String(http.StatusOK, "ok")
			})

			req := httptest.NewRequest(http.MethodPost, path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)
		})
	}

	t.Run("skips when bearer token present", func(t *testing.T) {
		e := echo.New()
		e.Use(CSRF(cfg))
		e.POST("/oidc/login", func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodPost, "/oidc/login", nil)
		req.Header.Set("Authorization", "Bearer some-token")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("enforces CSRF on protected browser route", func(t *testing.T) {
		e := echo.New()
		e.Use(CSRF(cfg))
		e.POST("/oidc/login", func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodPost, "/oidc/login", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.NotEqual(t, http.StatusOK, rec.Code)
	})

	t.Run("production sets secure cookie", func(t *testing.T) {
		prodCfg := &config.Config{AppEnv: config.EnvironmentProduction}
		e := echo.New()
		e.Use(CSRF(prodCfg))
		e.GET("/csrf-check", func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/csrf-check", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		cookies := rec.Result().Cookies()
		var csrfCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "csrf_token" {
				csrfCookie = c
				break
			}
		}
		require.NotNil(t, csrfCookie)
		require.True(t, csrfCookie.Secure)
		require.True(t, csrfCookie.HttpOnly)
	})
}

func TestExposeCSRFToken(t *testing.T) {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(echoMiddleware.DefaultCSRFConfig.ContextKey, "csrf-token-value")
			return next(c)
		}
	})
	e.Use(ExposeCSRFToken())
	e.GET("/token", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/token", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "csrf-token-value", rec.Header().Get("X-CSRF-Token"))
}

func TestExposeCSRFToken_EmptyToken(t *testing.T) {
	e := echo.New()
	e.Use(ExposeCSRFToken())
	e.GET("/token", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/token", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Header().Get("X-CSRF-Token"))
}
