package db

import (
	"os"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	testutil.InitLogger()
	os.Exit(m.Run())
}

func testManagerConfig() *config.Config {
	return &config.Config{
		DatabaseMaxOpenConns:    5,
		DatabaseMaxIdleConns:    2,
		DatabaseConnMaxLifetime: time.Minute,
		DatabaseConnMaxIdleTime: time.Minute,
		DatabaseConnectTimeout:  2 * time.Second,
		DatabaseHealthTimeout:   2 * time.Second,
		DatabaseRetryAttempts:   1,
		DatabaseRetryDelay:      10 * time.Millisecond,
		DatabaseSSLMode:         "disable",
		DatabaseTimezone:        "UTC",
	}
}

func newSQLiteManager(t *testing.T) *DatabaseManager {
	t.Helper()

	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := gormDB.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)

	return &DatabaseManager{
		db:     gormDB,
		config: testManagerConfig(),
		metrics: &ConnectionMetrics{
			MaxOpenConnections: 5,
			MaxIdleConnections: 2,
		},
		healthStatus: HealthStatus{
			IsHealthy: true,
			LastCheck: time.Now(),
		},
	}
}

func TestDatabaseManager_HealthCheck_Success(t *testing.T) {
	t.Parallel()
	dm := newSQLiteManager(t)

	status := dm.HealthCheck()
	require.True(t, status.IsHealthy)
	require.Empty(t, status.LastError)
	require.True(t, status.ResponseTime >= 0)
}

func TestDatabaseManager_HealthCheck_PingFailure(t *testing.T) {
	t.Parallel()
	dm := newSQLiteManager(t)

	sqlDB, err := dm.db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	status := dm.HealthCheck()
	require.False(t, status.IsHealthy)
	require.NotEmpty(t, status.LastError)
}

func TestDatabaseManager_FastHealthCheck_Recent(t *testing.T) {
	t.Parallel()
	dm := newSQLiteManager(t)
	dm.healthStatus.LastCheck = time.Now()

	status := dm.FastHealthCheck()
	require.True(t, status.IsHealthy)
}

func TestDatabaseManager_FastHealthCheck_StaleTriggersAsyncCheck(t *testing.T) {
	t.Parallel()
	dm := newSQLiteManager(t)
	dm.healthStatus.LastCheck = time.Now().Add(-11 * time.Second)

	status := dm.FastHealthCheck()
	require.True(t, status.IsHealthy)

	require.Eventually(t, func() bool {
		return time.Since(dm.GetHealthStatus().LastCheck) < 2*time.Second
	}, 2*time.Second, 20*time.Millisecond)
}

func TestDatabaseManager_GetMetricsAndUpdateMetrics(t *testing.T) {
	t.Parallel()
	dm := newSQLiteManager(t)

	dm.updateMetrics()
	metrics := dm.GetMetrics()
	require.Equal(t, 5, metrics.MaxOpenConnections)
	require.Equal(t, 2, metrics.MaxIdleConnections)
	require.GreaterOrEqual(t, metrics.OpenConnections, 0)
}

func TestDatabaseManager_IsHealthyAndGetHealthStatus(t *testing.T) {
	t.Parallel()
	dm := newSQLiteManager(t)
	dm.updateHealthStatus(true, "", 5*time.Millisecond)

	require.True(t, dm.IsHealthy())
	status := dm.GetHealthStatus()
	require.True(t, status.IsHealthy)
	require.Equal(t, "", status.LastError)
}

func TestDatabaseManager_Close(t *testing.T) {
	t.Parallel()
	dm := newSQLiteManager(t)
	require.NoError(t, dm.Close())
}

func TestDatabaseManager_Close_NilDB(t *testing.T) {
	t.Parallel()
	dm := &DatabaseManager{config: testManagerConfig(), metrics: &ConnectionMetrics{}}
	require.NoError(t, dm.Close())
}

func TestDatabaseManager_connectWithRetry_MultipleAttempts(t *testing.T) {
	t.Parallel()
	cfg := testManagerConfig()
	cfg.DatabaseHost = "127.0.0.1"
	cfg.DatabasePort = "1"
	cfg.DatabaseUsername = "user"
	cfg.DatabasePassword = "pass"
	cfg.DatabaseName = "db"
	cfg.DatabaseRetryAttempts = 3
	cfg.DatabaseRetryDelay = 5 * time.Millisecond

	dm := &DatabaseManager{
		config:  cfg,
		metrics: &ConnectionMetrics{},
	}
	err := dm.connectWithRetry()
	require.Error(t, err)
	require.Equal(t, 3, dm.healthStatus.RetryCount)
}

func TestDatabaseManager_connect_InvalidConnectionString(t *testing.T) {
	t.Parallel()
	dm := &DatabaseManager{
		config: &config.Config{
			DatabaseHost:           "invalid-host",
			DatabasePort:           "1",
			DatabaseUsername:       "user",
			DatabasePassword:       "pass",
			DatabaseName:           "db",
			DatabaseSSLMode:        "disable",
			DatabaseTimezone:       "UTC",
			DatabaseConnectTimeout: 100 * time.Millisecond,
		},
		metrics: &ConnectionMetrics{},
	}

	err := dm.connect()
	require.Error(t, err)
}

// Postgres testcontainers coverage lives in postgres_integration_test.go (integration build tag).

func TestNewDatabaseManager_ConnectFailure(t *testing.T) {
	t.Parallel()
	cfg := testManagerConfig()
	cfg.DatabaseHost = "127.0.0.1"
	cfg.DatabasePort = "1"
	cfg.DatabaseUsername = "user"
	cfg.DatabasePassword = "pass"
	cfg.DatabaseName = "db"
	cfg.DatabaseRetryAttempts = 2
	cfg.DatabaseRetryDelay = 10 * time.Millisecond

	manager, err := NewDatabaseManager(cfg)
	require.Error(t, err)
	require.Nil(t, manager)
}

func TestDatabaseManager_startHealthCheckAndMetricsCollection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping background ticker coverage in short mode")
	}

	dm := newSQLiteManager(t)
	go dm.startHealthCheck()
	go dm.startMetricsCollection()

	require.Eventually(t, func() bool {
		metrics := dm.GetMetrics()
		return metrics.MaxOpenConnections == dm.config.DatabaseMaxOpenConns
	}, 12*time.Second, 100*time.Millisecond)

	require.Eventually(t, func() bool {
		return dm.GetHealthStatus().LastCheck.After(time.Now().Add(-35 * time.Second))
	}, 35*time.Second, 200*time.Millisecond)
}
