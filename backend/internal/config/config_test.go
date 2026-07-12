package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEnvironment_String(t *testing.T) {
	require.Equal(t, "production", EnvironmentProduction.String())
}

func TestEnvironment_IsDevelopment(t *testing.T) {
	require.True(t, EnvironmentDevelopment.IsDevelopment())
	require.True(t, EnvironmentTest.IsDevelopment())
	require.False(t, EnvironmentProduction.IsDevelopment())
	require.False(t, EnvironmentStaging.IsDevelopment())
}

func TestEnvironment_IsProduction(t *testing.T) {
	require.True(t, EnvironmentProduction.IsProduction())
	require.False(t, EnvironmentDevelopment.IsProduction())
	require.False(t, EnvironmentTest.IsProduction())
}

func TestGetEnv(t *testing.T) {
	t.Setenv("TEST_GET_ENV_KEY", "value")
	require.Equal(t, "value", getEnv("TEST_GET_ENV_KEY", "fallback"))
	require.Equal(t, "fallback", getEnv("TEST_GET_ENV_MISSING", "fallback"))
}

func TestGetEnvAsInt(t *testing.T) {
		t.Setenv("TEST_INT_KEY", "42")
	require.Equal(t, 42, getEnvAsInt("TEST_INT_KEY", 10))
	require.Equal(t, 7, getEnvAsInt("TEST_INT_MISSING", 7))

	t.Setenv("TEST_INT_BAD", "not-a-number")
	require.Equal(t, 3, getEnvAsInt("TEST_INT_BAD", 3))
}

func TestGetEnvAsBool(t *testing.T) {
	t.Setenv("TEST_BOOL_KEY", "true")
	require.True(t, getEnvAsBool("TEST_BOOL_KEY", false))
	require.False(t, getEnvAsBool("TEST_BOOL_MISSING", false))

	t.Setenv("TEST_BOOL_BAD", "maybe")
	require.True(t, getEnvAsBool("TEST_BOOL_BAD", true))
}

func TestGetEnvAsDuration(t *testing.T) {
	t.Setenv("TEST_DUR_KEY", "5s")
	require.Equal(t, 5*time.Second, getEnvAsDuration("TEST_DUR_KEY", time.Minute))
	require.Equal(t, time.Minute, getEnvAsDuration("TEST_DUR_MISSING", time.Minute))

	t.Setenv("TEST_DUR_BAD", "not-a-duration")
	require.Equal(t, 2*time.Second, getEnvAsDuration("TEST_DUR_BAD", 2*time.Second))
}

func TestSplitCommaTrim(t *testing.T) {
	require.Nil(t, splitCommaTrim(""))
	require.Equal(t, []string{"a", "b"}, splitCommaTrim("a,b"))
	require.Equal(t, []string{"a", "b"}, splitCommaTrim(" a , , b "))
}

func TestBaseConfigFromEnv(t *testing.T) {
	t.Setenv("APP_NAME", "test-app")
	t.Setenv("APP_VERSION", "2.0.0")
	t.Setenv("JWT_SECRET", "custom-secret-thirty-two-bytes-long!!")
	t.Setenv("WEBAUTHN_RP_ORIGINS", "http://a.com, http://b.com")
	t.Setenv("SERVE_EMBEDDED_FRONTEND", "false")

	cfg := baseConfigFromEnv(EnvironmentDevelopment)
	require.Equal(t, "test-app", cfg.AppName)
	require.Equal(t, "2.0.0", cfg.AppVersion)
	require.Equal(t, "custom-secret-thirty-two-bytes-long!!", cfg.JWTSecret)
	require.Equal(t, []string{"http://a.com", "http://b.com"}, cfg.WebauthnRPOrigins)
	require.False(t, cfg.ServeEmbeddedFrontend)
}

func TestBaseConfigFromEnv_ProductionDefaults(t *testing.T) {
	// Isolate from other tests that set SERVE_EMBEDDED_FRONTEND (Setenv is process-wide with parallel tests).
	t.Setenv("SERVE_EMBEDDED_FRONTEND", "")
	cfg := baseConfigFromEnv(EnvironmentProduction)
	require.True(t, cfg.ServeEmbeddedFrontend)
}

func TestApplyHTTPAndStorageEnv(t *testing.T) {
	t.Setenv("HTTP_CLIENT_TIMEOUT", "45s")
	t.Setenv("HTTP_CLIENT_RETRY_COUNT", "5")
	t.Setenv("STORAGE_PROVIDER", "s3")
	t.Setenv("S3_BUCKET", "my-bucket")
	t.Setenv("S3_REGION", "us-east-1")

	cfg := &Config{}
	applyHTTPAndStorageEnv(cfg)

	require.Equal(t, 45*time.Second, cfg.HTTPClientTimeout)
	require.Equal(t, 5, cfg.HTTPClientRetryCount)
	require.Equal(t, "s3", cfg.StorageProvider)
	require.Equal(t, "my-bucket", cfg.S3Bucket)
	require.Equal(t, "us-east-1", cfg.S3Region)
}

func TestLoad(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_NAME", "loaded-app")
	t.Setenv("RATE_LIMIT", "50")
	t.Setenv("HTTP_CLIENT_DEBUG", "true")

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, EnvironmentTest, cfg.AppEnv)
	require.Equal(t, "loaded-app", cfg.AppName)
	require.Equal(t, 50, cfg.RateLimit)
	require.True(t, cfg.HTTPClientDebug)
}

func TestConnectionString(t *testing.T) {
	cfg := &Config{
		DatabaseHost:           "db.example.com",
		DatabasePort:           "5433",
		DatabaseUsername:       "iam",
		DatabasePassword:       "secret",
		DatabaseName:           "iam_db",
		DatabaseSSLMode:        "require",
		DatabaseTimezone:       "UTC",
		DatabaseConnectTimeout: 15 * time.Second,
	}
	cs := cfg.ConnectionString()
	require.Contains(t, cs, "host=db.example.com")
	require.Contains(t, cs, "port=5433")
	require.Contains(t, cs, "user=iam")
	require.Contains(t, cs, "password=secret")
	require.Contains(t, cs, "dbname=iam_db")
	require.Contains(t, cs, "sslmode=require")
	require.Contains(t, cs, "connect_timeout=15")
}

func TestIsDebugMode(t *testing.T) {
	cfg := &Config{DatabaseEnableDebug: true}
	require.True(t, cfg.IsDebugMode())

	cfg.DatabaseEnableDebug = false
	require.False(t, cfg.IsDebugMode())
}

func TestPopulateFromJSONBytes(t *testing.T) {
	cfg := &Config{}
	data := []byte(`{"client_email":"svc@test.iam.gserviceaccount.com","private_key":"-----BEGIN","project_id":"proj-1"}`)
	require.NoError(t, cfg.PopulateFromJSONBytes(data))
}

func TestPopulateFromJSONBytes_InvalidJSON(t *testing.T) {
	cfg := &Config{}
	err := cfg.PopulateFromJSONBytes([]byte(`not json`))
	require.Error(t, err)
}

func TestPopulateFromJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "creds.json")
	data := []byte(`{"client_email":"svc@test.iam.gserviceaccount.com","private_key":"key","project_id":"proj-1"}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	cfg := &Config{}
	require.NoError(t, cfg.PopulateFromJSON(path))
}

func TestPopulateFromJSON_FileNotFound(t *testing.T) {
	cfg := &Config{}
	err := cfg.PopulateFromJSON("/nonexistent/path/creds.json")
	require.Error(t, err)
}
