package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestCORS(t *testing.T) {
	e := echo.New()
	e.Use(CORS(&config.Config{CORSAllowedOrigins: []string{"http://localhost:5173"}}))
	e.GET("/api", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	t.Run("allows the configured origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "http://localhost:5173", rec.Header().Get("Access-Control-Allow-Origin"))
		require.Contains(t, rec.Header().Get("Access-Control-Expose-Headers"), "X-CSRF-Token")
	})

	t.Run("rejects a wildcard origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api", nil)
		req.Header.Set("Origin", "*")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	})
}
