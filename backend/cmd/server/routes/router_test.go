package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/handlers"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestRouter_HealthCheckSmoke(t *testing.T) {
	testutil.InitLogger()

	cfg := testutil.TestConfig()
	cfg.AppEnv = config.EnvironmentTest
	cfg.AppVersion = "test"
	cfg.ServeEmbeddedFrontend = false

	healthHandler := handlers.ProvideHealthHandler(cfg, nil)
	tokenService, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)

	e := Router(
		nil,
		healthHandler,
		nil,
		nil,
		nil,
		nil,
		nil,
		tokenService,
		nil,
		cfg,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "healthy")
}

func TestRouter_SwaggerRegisteredInNonProduction(t *testing.T) {
	testutil.InitLogger()

	cfg := testutil.TestConfig()
	cfg.AppEnv = config.EnvironmentDevelopment
	cfg.ServeEmbeddedFrontend = false

	healthHandler := handlers.ProvideHealthHandler(cfg, nil)
	tokenService, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, time.Hour)
	require.NoError(t, err)

	e := Router(nil, healthHandler, nil, nil, nil, nil, nil, tokenService, nil, cfg)

	found := false
	for _, route := range e.Routes() {
		if route.Path == "/swagger/*" {
			found = true
			break
		}
	}
	require.True(t, found)
}

func TestRouter_OIDCRoutesRegistered(t *testing.T) {
	testutil.InitLogger()

	cfg := testutil.TestConfig()
	cfg.AppEnv = config.EnvironmentTest
	cfg.ServeEmbeddedFrontend = false

	healthHandler := handlers.ProvideHealthHandler(cfg, nil)
	tokenService, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, time.Hour)
	require.NoError(t, err)

	e := Router(nil, healthHandler, nil, nil, nil, nil, nil, tokenService, nil, cfg)

	paths := map[string]bool{}
	for _, route := range e.Routes() {
		paths[route.Path] = true
	}
	require.True(t, paths["/.well-known/openid-configuration"])
	require.True(t, paths["/authorize"])
	require.True(t, paths["/token"])
}
