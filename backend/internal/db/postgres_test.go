package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPostgresDB_NilManagerHealthAndMetrics(t *testing.T) {
	t.Parallel()
	var pg PostgresDB

	require.False(t, pg.HealthCheck().IsHealthy)
	require.False(t, pg.FastHealthCheck().IsHealthy)
	require.Equal(t, ConnectionMetrics{}, pg.GetMetrics())
	require.NoError(t, pg.Close())
}

func TestPostgresDB_WithStubManager(t *testing.T) {
	t.Parallel()
	pg := NewHealthyPostgresDBStub()

	fast := pg.FastHealthCheck()
	require.True(t, fast.IsHealthy)

	metrics := pg.GetMetrics()
	require.NotNil(t, metrics)
	require.NoError(t, pg.Close())
}

func TestPostgresDB_GetManager(t *testing.T) {
	t.Parallel()
	pg := NewHealthyPostgresDBStub()
	require.NotNil(t, pg.GetManager())
}

func TestPostgresDB_NewPostgresDB_InvalidConfig(t *testing.T) {
	t.Parallel()
	var pg PostgresDB
	cfg := testManagerConfig()
	cfg.DatabaseHost = "127.0.0.1"
	cfg.DatabasePort = "1"
	cfg.DatabaseUsername = "user"
	cfg.DatabasePassword = "pass"
	cfg.DatabaseName = "db"
	cfg.DatabaseRetryAttempts = 1
	cfg.DatabaseRetryDelay = 10 * time.Millisecond

	err := pg.NewPostgresDB(cfg)
	require.Error(t, err)
	require.Nil(t, pg.DB)
	require.Nil(t, pg.manager)
}

// Successful NewPostgresDB against real Postgres is covered in postgres_integration_test.go.
