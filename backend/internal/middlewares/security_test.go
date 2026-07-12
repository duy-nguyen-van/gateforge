package middlewares

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestSecurity(t *testing.T) {
	e := echo.New()
	e.Use(Security())
	e.GET("/secure", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.TLS = &tls.ConnectionState{}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "1; mode=block", rec.Header().Get("X-Xss-Protection"))
	require.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	require.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	require.Contains(t, rec.Header().Get("Strict-Transport-Security"), "max-age=31536000")
	require.Equal(t, "no-referrer", rec.Header().Get("Referrer-Policy"))
}
