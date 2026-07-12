package services

import (
	"os"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"
)

func TestMain(m *testing.M) {
	testutil.InitLogger()
	os.Exit(m.Run())
}

func testConfig() *config.Config {
	cfg := testutil.TestConfig()
	cfg.JWTRefreshTTL = 24 * time.Hour
	cfg.NativeOAuthClientID = "native"
	cfg.OIDCAuthCodeTTL = 10 * time.Minute
	cfg.OIDCAccessTTL = time.Hour
	cfg.OIDCIDTokenTTL = time.Hour
	cfg.MFAPendingTicketTTL = 5 * time.Minute
	cfg.MFARecoveryCodeCount = 3
	cfg.WebauthnSessionTTL = 5 * time.Minute
	cfg.WebauthnRPID = "localhost"
	cfg.WebauthnRPDisplayName = "GateForge IAM Test"
	cfg.WebauthnRPOrigins = []string{"http://localhost:5173"}
	cfg.SSOSessionTTL = 8 * time.Hour
	cfg.SSOSessionRememberTTL = 720 * time.Hour
	cfg.MFAEncryptionKey = "01234567890123456789012345678901"
	return cfg
}
