package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashAndVerifyClientSecret(t *testing.T) {
	pepper := "test-pepper-at-least-long-enough"
	hashed, err := HashClientSecret(pepper, "super-secret-value")
	require.NoError(t, err)
	require.True(t, IsHashedClientSecret(hashed))

	ok, upgrade := VerifyClientSecret(pepper, hashed, "super-secret-value")
	require.True(t, ok)
	require.Empty(t, upgrade)

	ok, upgrade = VerifyClientSecret(pepper, hashed, "wrong")
	require.False(t, ok)
	require.Empty(t, upgrade)
}

func TestVerifyClientSecret_LegacyPlaintextUpgrades(t *testing.T) {
	pepper := "test-pepper-at-least-long-enough"
	ok, upgrade := VerifyClientSecret(pepper, "legacy-plain", "legacy-plain")
	require.True(t, ok)
	require.True(t, IsHashedClientSecret(upgrade))

	ok, _ = VerifyClientSecret(pepper, upgrade, "legacy-plain")
	require.True(t, ok)

	ok, upgrade = VerifyClientSecret(pepper, "legacy-plain", "other")
	require.False(t, ok)
	require.Empty(t, upgrade)
}

func TestHashClientSecret_Empty(t *testing.T) {
	_, err := HashClientSecret("", "secret")
	require.Error(t, err)
	_, err = HashClientSecret("pepper", "")
	require.Error(t, err)
}
