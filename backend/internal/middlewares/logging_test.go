package middlewares

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/request"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/require"
)

func TestRequestLogging(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AppVersion = "1.2.3"

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := request.NewCorrelationIDContext(c.Request().Context(), "corr-log")
			ctx = request.NewLanguageCodeContext(ctx, "en")
			c.SetRequest(c.Request().WithContext(ctx))
			c.Set(auth.EchoContextUserIDKey, "user-log")
			c.Set("organization_id", "org-1")
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

func TestBuildRequestLogContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/42?sort=name", nil)
	req.Header.Set("Authorization", "Bearer "+makeTestJWT(t))
	req.Header.Set("Cookie", "session=secret")
	req.Header.Set("X-Custom", "value")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("42")
	c.Set(auth.EchoContextUserIDKey, "user-99")
	c.Set("organization_id", "org-55")
	c.Set("log_body", `{"k":"v"}`)

	ctx := request.NewCorrelationIDContext(c.Request().Context(), "corr-ctx")
	ctx = request.NewLanguageCodeContext(ctx, "fr")
	c.SetRequest(c.Request().WithContext(ctx))

	logCtx := buildRequestLogContext(c)
	require.Equal(t, "user-99", logCtx.userID)
	require.Equal(t, "org-55", logCtx.organizationID)
	require.Equal(t, "corr-ctx", logCtx.correlationID)
	require.Equal(t, "fr", logCtx.languageCode)
	require.Equal(t, `{"k":"v"}`, logCtx.body)
	require.Contains(t, logCtx.queryJSON, "sort")
	require.Contains(t, logCtx.pathParamsJSON, "42")
	require.Contains(t, logCtx.headersJSON, "X-Custom")
	require.NotContains(t, logCtx.headersJSON, "Authorization")
	require.NotContains(t, logCtx.headersJSON, "Cookie")
	require.NotEmpty(t, logCtx.userLoginInfo)
}

func TestLogRequestHelpers(t *testing.T) {
	cfg := testutil.TestConfig()
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	values := middleware.RequestLoggerValues{
		Method:    http.MethodPost,
		URI:       "/test",
		Status:    http.StatusCreated,
		Latency:   10 * time.Millisecond,
		RequestID: "req-1",
		RemoteIP:  "127.0.0.1",
		UserAgent: "agent",
	}
	logCtx := requestLogContext{correlationID: "cid", languageCode: "en"}

	require.NotPanics(t, func() {
		logRequestToSentry(c, values, cfg, logCtx)
		logRequestToZap(c, values, cfg, logCtx)
	})
}

func makeTestJWT(t *testing.T) string {
	t.Helper()
	header := base64.RawStdEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawStdEncoding.EncodeToString([]byte(`{"sub":"user-jwt"}`))
	return strings.Join([]string{header, payload, "sig"}, ".")
}

func TestBuildRequestLogContext_NoOptionalFields(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	logCtx := buildRequestLogContext(c)
	require.Empty(t, logCtx.userID)
	require.Empty(t, logCtx.organizationID)
	require.Equal(t, "", logCtx.body)
}

func TestBuildRequestLogContext_InvalidJWTHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Authorization", "Bearer not-three-parts")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	logCtx := buildRequestLogContext(c)
	require.Empty(t, logCtx.userLoginInfo)
}

func TestBuildRequestLogContext_JWTPayloadDecode(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	payload, _ := json.Marshal(map[string]string{"sub": "decoded-user"})
	encoded := base64.RawStdEncoding.EncodeToString(payload)
	req.Header.Set("Authorization", "Bearer header."+encoded+".sig")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	logCtx := buildRequestLogContext(c)
	require.Contains(t, logCtx.userLoginInfo, "decoded-user")
}
