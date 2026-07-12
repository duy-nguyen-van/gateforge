package monitoring

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
)

func TestInitSentry_EmptyDSNIsNoOp(t *testing.T) {
	testutil.InitLogger()
	cfg := config.Config{SentryDSN: "  "}
	InitSentry(cfg)
}

func TestGetSentryHub(t *testing.T) {
	hub := GetSentryHub(context.Background())
	require.NotNil(t, hub)

	ctxHub := sentry.SetHubOnContext(context.Background(), sentry.CurrentHub())
	require.Equal(t, sentry.GetHubFromContext(ctxHub), GetSentryHub(ctxHub))
}

func TestNewSentryWriter(t *testing.T) {
	//nolint:staticcheck // explicitly testing nil-context fallback to Background
	w := NewSentryWriter(nil)
	require.NotNil(t, w)
	n, err := w.Write([]byte("test log line"))
	require.NoError(t, err)
	require.Equal(t, len("test log line"), n)

	w2 := NewSentryWriter(context.Background())
	n2, err2 := w2.Write([]byte("another line"))
	require.NoError(t, err2)
	require.Equal(t, len("another line"), n2)
}

func TestFlushSentry(t *testing.T) {
	FlushSentry()
}
