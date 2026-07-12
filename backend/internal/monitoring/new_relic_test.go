package monitoring

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestInitNewRelic_EmptyLicenseIsNoOp(t *testing.T) {
	testutil.InitLogger()
	cfg := config.Config{NewRelicLicense: "", NewRelicAppName: ""}
	app := InitNewRelic(cfg)
	require.Nil(t, app)
}

func TestInitNewRelic_SuccessInitializesApp(t *testing.T) {
	testutil.InitLogger()
	cfg := config.Config{
		NewRelicLicense: "0000000000000000000000000000000000000000",
		NewRelicAppName: "gateforge-iam-test",
	}
	app := InitNewRelic(cfg)
	require.NotNil(t, app)
}

func TestInitNewRelic_DefaultAppName(t *testing.T) {
	testutil.InitLogger()
	cfg := config.Config{
		NewRelicLicense: "0000000000000000000000000000000000000000",
		NewRelicAppName: "",
	}
	app := InitNewRelic(cfg)
	require.NotNil(t, app)
}

func TestInitNewRelic_InvalidLicenseReturnsNil(t *testing.T) {
	testutil.InitLogger()
	cfg := config.Config{
		NewRelicLicense: "not-a-valid-license",
		NewRelicAppName: "",
	}
	app := InitNewRelic(cfg)
	require.Nil(t, app)
}
