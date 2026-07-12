package constants

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIdentityProviderByID(t *testing.T) {
	spec, ok := IdentityProviderByID("google")
	require.True(t, ok)
	require.Equal(t, IdentityProviderGoogle, spec.ID)
	require.Equal(t, "Google", spec.DisplayName)

	_, ok = IdentityProviderByID("unknown")
	require.False(t, ok)

	spec, ok = IdentityProviderByID("  GOOGLE  ")
	require.True(t, ok)
	require.Equal(t, IdentityProviderGoogle, spec.ID)
}

func TestIsSupportedIdentityProvider(t *testing.T) {
	require.True(t, IsSupportedIdentityProvider("google"))
	require.False(t, IsSupportedIdentityProvider("facebook"))
}
