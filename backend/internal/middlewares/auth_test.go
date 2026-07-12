package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "test-jwt-secret-at-least-thirty-two-bytes"

func TestJWTBearerAuth(t *testing.T) {
	svc, err := auth.NewTokenService(testJWTSecret, "gateforge-iam", time.Hour)
	require.NoError(t, err)

	okHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"user_id":   c.Get(auth.EchoContextUserIDKey).(string),
			"tenant_id": c.Get(auth.EchoContextTenantIDKey).(string),
		})
	}

	t.Run("missing authorization header", func(t *testing.T) {
		e := echo.New()
		e.Use(JWTBearerAuth(svc))
		e.GET("/me", okHandler)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		var body map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		require.Equal(t, "Authorization header required", body["error"])
	})

	t.Run("invalid authorization format", func(t *testing.T) {
		e := echo.New()
		e.Use(JWTBearerAuth(svc))
		e.GET("/me", okHandler)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Basic abc")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("empty bearer token", func(t *testing.T) {
		e := echo.New()
		e.Use(JWTBearerAuth(svc))
		e.GET("/me", okHandler)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Bearer ")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		e := echo.New()
		e.Use(JWTBearerAuth(svc))
		e.GET("/me", okHandler)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Bearer not-a-jwt")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("valid token sets user and tenant", func(t *testing.T) {
		token, _, err := svc.SignAccessToken("user-1", "tenant-1")
		require.NoError(t, err)

		e := echo.New()
		e.Use(JWTBearerAuth(svc))
		e.GET("/me", okHandler)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var body map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		require.Equal(t, "user-1", body["user_id"])
		require.Equal(t, "tenant-1", body["tenant_id"])
	})
}
