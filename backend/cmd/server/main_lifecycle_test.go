package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/handlers"
	"github.com/gateforge-iam/gateforge-iam/internal/services"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

type stubBootstrap struct {
	err error
}

func (s stubBootstrap) Run(context.Context) error {
	return s.err
}

var _ services.PlatformAdminBootstrap = stubBootstrap{}

func TestRunPlatformAdminBootstrap(t *testing.T) {
	require.NoError(t, runPlatformAdminBootstrap(stubBootstrap{}))
	require.Error(t, runPlatformAdminBootstrap(stubBootstrap{err: errors.New("bootstrap failed")}))
}

func TestNewHTTPServer_LifecycleOnStopWithoutStart(t *testing.T) {
	testutil.InitLogger()

	cfg := testutil.TestConfig()
	cfg.AppHTTPServer = ":0"
	cfg.AppRequestTimeout = 5

	tokenService, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)

	healthHandler := handlers.ProvideHealthHandler(cfg, nil)

	var hooks []fx.Hook
	lc := hookLifecycle{appendHook: func(h fx.Hook) { hooks = append(hooks, h) }}

	_ = NewHTTPServer(lc, healthHandler, nil, nil, nil, nil, nil, nil, tokenService, nil, cfg, &db.PostgresDB{})
	require.Len(t, hooks, 1)
	require.NoError(t, hooks[0].OnStop(context.Background()))
}

func TestNewHTTPServer_LifecycleOnStart(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, hooks[0].OnStart(ctx))
}

func TestNewHTTPServer_LifecycleOnStartListenError(t *testing.T) {
	testutil.InitLogger()

	cfg := testutil.TestConfig()
	cfg.AppHTTPServer = "invalid-listen-addr"
	cfg.AppRequestTimeout = 5

	tokenService, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)

	healthHandler := handlers.ProvideHealthHandler(cfg, nil)

	var hooks []fx.Hook
	lc := hookLifecycle{appendHook: func(h fx.Hook) { hooks = append(hooks, h) }}

	_ = NewHTTPServer(lc, healthHandler, nil, nil, nil, nil, nil, nil, tokenService, nil, cfg, &db.PostgresDB{})
	require.Len(t, hooks, 1)
	require.Error(t, hooks[0].OnStart(context.Background()))
}
