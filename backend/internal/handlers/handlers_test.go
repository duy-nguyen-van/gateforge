package handlers

import (
	"os"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"
	"github.com/gateforge-iam/gateforge-iam/internal/utils/i18n"
)

const (
	testUserID   = "00000000-0000-4000-8000-000000000001"
	testTenantID = "00000000-0000-4000-8000-000000000002"
	testClientPK = "00000000-0000-4000-8000-000000000003"
)

func TestMain(m *testing.M) {
	testutil.InitLogger()
	i18n.Init()
	os.Exit(m.Run())
}

func handlerTestConfig() *config.Config {
	cfg := testutil.TestConfig()
	cfg.AppEnv = config.EnvironmentTest
	cfg.AppVersion = "1.0.0-test"
	cfg.OIDCLoginPageURL = "http://localhost:8080/login"
	return cfg
}
