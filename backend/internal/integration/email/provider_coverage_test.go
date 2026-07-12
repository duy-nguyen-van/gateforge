package email

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestProvideEmailSender_SESInitFailure(t *testing.T) {
	testutil.InitLogger()
	cfg := config.Config{
		EmailProvider: constants.EmailProviderSES,
		AWSSESRegion:  "",
	}
	sender, err := ProvideEmailSender(cfg)
	if err != nil {
		require.Nil(t, sender)
		require.Contains(t, err.Error(), "Failed to initialize SES email sender")
		return
	}
	require.NotNil(t, sender)
}
