package request

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuditContextFromContext(t *testing.T) {
	ctx := context.Background()
	_, ok := AuditContextFromContext(ctx)
	require.False(t, ok)

	ac := AuditContext{
		IPAddress:     "127.0.0.1",
		UserAgent:     "test-agent",
		RequestID:     "req-1",
		CorrelationID: "corr-1",
		ActorType:     "user",
		ActorID:       "user-1",
		TenantID:      "tenant-1",
	}
	ctx = NewAuditContextContext(ctx, ac)
	got, ok := AuditContextFromContext(ctx)
	require.True(t, ok)
	require.Equal(t, ac, got)
}
