package monitoring

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestInitSentry_WithDSNAndLogger(t *testing.T) {
	testutil.InitLogger()
	require.NotNil(t, logger.Log)

	cfg := config.Config{
		SentryDSN:  "https://examplePublicKey@o0.ingest.sentry.io/0",
		AppEnv:     config.EnvironmentDevelopment,
		AppName:    "gateforge-iam-test",
		AppVersion: "test",
	}
	InitSentry(cfg)
	FlushSentry()
}

func TestInitNewRelic_EmptyLicense(t *testing.T) {
	testutil.InitLogger()
	InitNewRelic(config.Config{NewRelicLicense: "  "})
}

func TestInitNewRelic_WithLicense(t *testing.T) {
	testutil.InitLogger()
	InitNewRelic(config.Config{
		NewRelicLicense: "test-license-key",
		AppName:         "gateforge-iam-test",
		AppEnv:          config.EnvironmentTest,
	})
}
