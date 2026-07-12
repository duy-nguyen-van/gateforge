package db

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestDatabaseManager_GetDB(t *testing.T) {
	t.Parallel()
	dm := newSQLiteManager(t)
	require.NotNil(t, dm.GetDB())
}

func TestDatabaseManager_connect_DebugMode(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()
	cfg := testManagerConfig()
	cfg.AppEnv = config.EnvironmentDevelopment
	cfg.DatabaseHost = "127.0.0.1"
	cfg.DatabasePort = "1"
	cfg.DatabaseUsername = "user"
	cfg.DatabasePassword = "pass"
	cfg.DatabaseName = "db"
	cfg.DatabaseConnectTimeout = testManagerConfig().DatabaseConnectTimeout

	dm := &DatabaseManager{config: cfg, metrics: &ConnectionMetrics{}}
	err := dm.connect()
	require.Error(t, err)
}
