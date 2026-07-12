package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptMFASecret_RoundTrip(t *testing.T) {
	key := "test-key-material-at-least-thirty-two-bytes-long"
	enc, err := EncryptMFASecret(key, "my-totp-secret")
	require.NoError(t, err)
	require.NotEmpty(t, enc)

	plain, err := DecryptMFASecret(key, enc)
	require.NoError(t, err)
	require.Equal(t, "my-totp-secret", plain)
}
