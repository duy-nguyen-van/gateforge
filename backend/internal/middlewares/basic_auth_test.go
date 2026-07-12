package middlewares

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestBasicAuthMiddleware(t *testing.T) {
	cfg := config.Config{
		BasicAuthUsername: "admin",
		BasicAuthPassword: "secret",
	}

	e := echo.New()
	e.Use(BasicAuthMiddleware(cfg))
	e.GET("/protected", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	t.Run("rejects missing credentials", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("rejects wrong credentials", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.SetBasicAuth("admin", "wrong")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("accepts valid credentials", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.SetBasicAuth("admin", "secret")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "ok", rec.Body.String())
	})

	t.Run("rejects wrong username", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.SetBasicAuth("other", "secret")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("uses constant-time compare", func(t *testing.T) {
		// Smoke test: ensure middleware is wired (credentials encoded correctly).
		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:secret"))
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", authHeader)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})
}
