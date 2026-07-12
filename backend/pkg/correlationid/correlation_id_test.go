package correlationid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCorrelationID_ContextRoundTrip(t *testing.T) {
	ctx := context.Background()
	const id = "corr-123"
	ctx = NewContext(ctx, id)

	got, ok := FromContext(ctx)
	require.True(t, ok)
	require.Equal(t, id, got)
}

func TestCorrelationID_FromEmptyContext(t *testing.T) {
	_, ok := FromContext(context.Background())
	require.False(t, ok)
}

func TestCorrelationID_Header(t *testing.T) {
	require.Equal(t, "X-Correlation-Id", Header)
}
