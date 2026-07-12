package email

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestProvideEmailSender_InvalidProvider(t *testing.T) {
	testutil.InitLogger()
	cfg := config.Config{EmailProvider: "unknown"}
	sender, err := ProvideEmailSender(cfg)
	require.Error(t, err)
	require.Nil(t, sender)
	require.Contains(t, err.Error(), "invalid email provider")
}

func TestProvideEmailSender_SESProvider(t *testing.T) {
	testutil.InitLogger()
	cfg := config.Config{
		EmailProvider:   constants.EmailProviderSES,
		AWSSESRegion:    "us-east-1",
		AWSSESAccessKey: "AKIATESTKEY",
		AWSSESSecretKey: "secret",
	}
	sender, err := ProvideEmailSender(cfg)
	require.NoError(t, err)
	require.NotNil(t, sender)
}

func TestSESSender_SendEmail_MissingBody(t *testing.T) {
	testutil.InitLogger()
	s := &SESSender{config: config.Config{}}
	resp, err := s.SendEmail(context.Background(), EmailRequest{
		To:      []string{"user@example.com"},
		Subject: "Hello",
	})
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "either HTML or text content must be provided")
}

func TestNewSESSender_LoadsClient(t *testing.T) {
	testutil.InitLogger()
	cfg := config.Config{
		AWSSESRegion:    "us-east-1",
		AWSSESAccessKey: "AKIATESTKEY",
		AWSSESSecretKey: "secret",
	}
	sender, err := NewSESSender(cfg)
	require.NoError(t, err)
	require.NotNil(t, sender)
	require.NotNil(t, sender.client)
}
