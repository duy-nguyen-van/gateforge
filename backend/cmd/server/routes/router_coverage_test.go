package routes

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/handlers"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func testRouter(t *testing.T, cfg *config.Config) *echo.Echo {
	t.Helper()
	testutil.InitLogger()
	tokenService, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	return Router(nil, handlers.ProvideHealthHandler(cfg, nil), nil, nil, nil, nil, nil, tokenService, nil, cfg)
}

func TestRouter_ProductionOmitsSwagger(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AppEnv = config.EnvironmentProduction
	cfg.ServeEmbeddedFrontend = false

	e := testRouter(t, cfg)
	for _, route := range e.Routes() {
		require.NotEqual(t, "/swagger/*", route.Path)
	}
}

func TestRouter_ServeEmbeddedFrontend(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AppEnv = config.EnvironmentTest
	cfg.ServeEmbeddedFrontend = true
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>spa</html>"), 0o644))
	cfg.FrontendDistPath = dir

	e := testRouter(t, cfg)
	found := false
	for _, route := range e.Routes() {
		if route.Path == "/*" {
			found = true
			break
		}
	}
	require.True(t, found)
}

func TestRouter_AdminRoutesRegistered(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AppEnv = config.EnvironmentTest
	cfg.ServeEmbeddedFrontend = false

	e := testRouter(t, cfg)
	paths := map[string]bool{}
	for _, route := range e.Routes() {
		paths[route.Path] = true
	}
	require.True(t, paths["/api/v1/admin/stats"])
	require.True(t, paths["/api/v1/admin/users"])
	require.True(t, paths["/api/v1/admin/tenants"])
}

func TestSentryCaptureMiddleware_SkipsNotFound(t *testing.T) {
	cfg := testutil.TestConfig()
	e := echo.New()
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	e.Use(sentryCaptureMiddleware(cfg))
	e.GET("/missing", func(c echo.Context) error {
		return echo.ErrNotFound
	})

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestSentryCaptureMiddleware_SkipsMethodNotAllowed(t *testing.T) {
	cfg := testutil.TestConfig()
	e := echo.New()
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	e.Use(sentryCaptureMiddleware(cfg))
	e.GET("/only-get", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	})

	req := httptest.NewRequest(http.MethodPost, "/only-get", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestSentryCaptureMiddleware_CapturesInternalError(t *testing.T) {
	cfg := testutil.TestConfig()
	e := echo.New()
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	e.Use(sentryCaptureMiddleware(cfg))
	e.GET("/boom", func(c echo.Context) error {
		c.Set("organization_id", "org-1")
		c.Set("log_body", "request-body")
		return errors.New("internal failure")
	})

	req := httptest.NewRequest(http.MethodGet, "/boom?x=1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSentryCaptureMiddleware_CapturesHTTPError(t *testing.T) {
	cfg := testutil.TestConfig()
	e := echo.New()
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	e.Use(sentryCaptureMiddleware(cfg))
	e.GET("/bad", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "bad input")
	})

	req := httptest.NewRequest(http.MethodGet, "/bad", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSentryCaptureMiddleware_NoHubReturnsError(t *testing.T) {
	cfg := testutil.TestConfig()
	e := echo.New()
	e.Use(sentryCaptureMiddleware(cfg))
	e.GET("/err", func(c echo.Context) error {
		return errors.New("no sentry hub")
	})

	req := httptest.NewRequest(http.MethodGet, "/err", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRegisterGlobalMiddlewareChain(t *testing.T) {
	cfg := testutil.TestConfig()
	e := echo.New()
	registerGlobalMiddleware(e, cfg)
	e.GET("/ping", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}
