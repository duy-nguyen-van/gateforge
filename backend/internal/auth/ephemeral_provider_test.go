package auth

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestProvideEphemeralStore(t *testing.T) {
	cfg := testutil.TestConfig()
	store := ProvideEphemeralStore(newFakeCache(), cfg)
	require.NotNil(t, store)
	require.NotNil(t, store.cache)
	require.Equal(t, cfg, store.cfg)
}
