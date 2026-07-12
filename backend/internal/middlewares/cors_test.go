package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestCORS(t *testing.T) {
	e := echo.New()
	e.Use(CORS())
	e.GET("/api", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	t.Run("sets CORS headers on GET", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set("Origin", "https://example.com")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
		require.Contains(t, rec.Header().Get("Access-Control-Expose-Headers"), "X-CSRF-Token")
	})

	t.Run("handles preflight OPTIONS", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Code)
		require.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
		require.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), http.MethodPost)
	})
}
