package request

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLanguageCodeContext(t *testing.T) {
	ctx := NewLanguageCodeContext(context.Background(), "en")
	got, ok := LanguageCodeFromContext(ctx)
	require.True(t, ok)
	require.Equal(t, "en", got)
}

func TestRequestTimestampContext(t *testing.T) {
	ctx := NewRequestTimestampContext(context.Background(), 12345)
	got, ok := RequestTimestampFromContext(ctx)
	require.True(t, ok)
	require.Equal(t, int64(12345), got)
}

func TestCorrelationIDContext(t *testing.T) {
	ctx := NewCorrelationIDContext(context.Background(), "cid-1")
	got, ok := CorrelationIDFromContext(ctx)
	require.True(t, ok)
	require.Equal(t, "cid-1", got)
}

func TestRequestURLContext(t *testing.T) {
	ctx := NewRequestURLContext(context.Background(), "/api/v1/me")
	got, ok := RequestURLFromContext(ctx)
	require.True(t, ok)
	require.Equal(t, "/api/v1/me", got)
}

func TestCtxKeyString(t *testing.T) {
	require.Contains(t, ctxKeyLanguageCode.String(), "language_code")
}
