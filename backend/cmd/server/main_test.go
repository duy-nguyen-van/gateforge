package main

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/handlers"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

type noopLifecycle struct{}

func (noopLifecycle) Append(fx.Hook) {}

func TestProvideValidator(t *testing.T) {
	v := ProvideValidator()
	require.NotNil(t, v)
}

func TestNewHTTPServer(t *testing.T) {
	testutil.InitLogger()

	cfg := testutil.TestConfig()
	cfg.AppHTTPServer = ":18080"
	cfg.AppRequestTimeout = 30
	cfg.AppVersion = "test"

	tokenService, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)

	healthHandler := handlers.ProvideHealthHandler(cfg, nil)

	srv := NewHTTPServer(
		noopLifecycle{},
		healthHandler,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		tokenService,
		nil,
		cfg,
		&db.PostgresDB{},
	)

	require.NotNil(t, srv)
	require.Equal(t, ":18080", srv.Addr)
	require.NotNil(t, srv.Handler)
	require.Equal(t, 30*time.Second, srv.ReadHeaderTimeout)
}

func TestNewHTTPServer_RegistersLifecycleHook(t *testing.T) {
	testutil.InitLogger()

	cfg := testutil.TestConfig()
	cfg.AppHTTPServer = ":0"
	cfg.AppRequestTimeout = 5

	tokenService, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)

	healthHandler := handlers.ProvideHealthHandler(cfg, nil)

	var hooks []fx.Hook
	lc := hookLifecycle{appendHook: func(h fx.Hook) { hooks = append(hooks, h) }}

	srv := NewHTTPServer(lc, healthHandler, nil, nil, nil, nil, nil, nil, tokenService, nil, cfg, &db.PostgresDB{})
	require.NotNil(t, srv)
	require.Len(t, hooks, 1)
}

type hookLifecycle struct {
	appendHook func(fx.Hook)
}

func (h hookLifecycle) Append(hook fx.Hook) {
	h.appendHook(hook)
}

func TestProvideGormPostgres_InvalidConfigFatals(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		testutil.InitLogger()
		cfg := testutil.TestConfig()
		cfg.DatabaseHost = "invalid-host-should-not-resolve"
		cfg.DatabasePort = "1"
		cfg.DatabaseConnectTimeout = 1 * time.Second
		ProvideGormPostgres(cfg)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestProvideGormPostgres_InvalidConfigFatals", "-test.v")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	err := cmd.Run()
	require.Error(t, err)
}
