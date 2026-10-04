package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/request"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestRequestLogging(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AppVersion = "1.2.3"

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := request.NewCorrelationIDContext(c.Request().Context(), "corr-log")
			c.SetRequest(c.Request().WithContext(ctx))
			c.Set(auth.EchoContextUserIDKey, "user-log")
			return next(c)
		}
	})
	e.Use(RequestLogging(cfg))
	e.GET("/items", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/items?q=1", nil)
	req.Header.Set("User-Agent", "test-agent")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestLogBodyMiddleware(t *testing.T) {
	e := echo.New()
	e.Use(LogBodyMiddleware)
	e.POST("/echo", func(c echo.Context) error {
		body, _ := c.Get("log_body").(string)
		return c.String(http.StatusOK, body)
	})

	payload := `{"hello":"world"}`
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader([]byte(payload)))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, payload, rec.Body.String())
}

func TestLogBodyMiddleware_SkipsSensitivePaths(t *testing.T) {
	e := echo.New()
	e.Use(LogBodyMiddleware)
	e.POST("/api/v1/login", func(c echo.Context) error {
		_, ok := c.Get("log_body").(string)
		if ok {
			return c.String(http.StatusInternalServerError, "captured")
		}
		return c.String(http.StatusOK, "skipped")
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{"password":"secret"}`)))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "skipped", rec.Body.String())
}

func TestLogBodyMiddleware_SkipsMultipart(t *testing.T) {
	e := echo.New()
	e.Use(LogBodyMiddleware)
	e.POST("/upload", func(c echo.Context) error {
		if c.Get("log_body") != nil {
			return c.String(http.StatusInternalServerError, "captured")
		}
		return c.String(http.StatusOK, "skipped")
	})

	req := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewReader([]byte("file")))
	req.Header.Set(echo.HeaderContentType, "multipart/form-data; boundary=x")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "skipped", rec.Body.String())
}

func TestRequestLogging_NonStringUserID(t *testing.T) {
	cfg := testutil.TestConfig()
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(auth.EchoContextUserIDKey, 123)
			return next(c)
		}
	})
	e.Use(RequestLogging(cfg))
	e.GET("/items", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}
