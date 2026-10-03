package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
)

const clientSecretPrefix = "hmac-sha256:"

// HashClientSecret returns an HMAC-SHA256 digest of a high-entropy OAuth client secret.
func HashClientSecret(pepper, plaintext string) (string, error) {
	if strings.TrimSpace(pepper) == "" || plaintext == "" {
		return "", fmt.Errorf("client secret pepper and plaintext are required")
	}
	mac := hmac.New(sha256.New, []byte(pepper))
	_, _ = mac.Write([]byte(plaintext))
	return clientSecretPrefix + hex.EncodeToString(mac.Sum(nil)), nil
}

// VerifyClientSecret reports whether presented matches stored.
// upgrade is set when stored is a legacy plaintext secret and should be replaced with the HMAC.
func VerifyClientSecret(pepper, stored, presented string) (ok bool, upgrade string) {
	if stored == "" || presented == "" {
		return false, ""
	}
	if strings.HasPrefix(stored, clientSecretPrefix) {
		expected, err := HashClientSecret(pepper, presented)
		if err != nil {
			return false, ""
		}
		if subtle.ConstantTimeCompare([]byte(stored), []byte(expected)) == 1 {
			return true, ""
		}
		return false, ""
	}

	sumStored := sha256.Sum256([]byte(stored))
	sumPresented := sha256.Sum256([]byte(presented))
	if subtle.ConstantTimeCompare(sumStored[:], sumPresented[:]) != 1 {
		return false, ""
	}
	upgraded, err := HashClientSecret(pepper, presented)
	if err != nil {
		return true, ""
	}
	return true, upgraded
}

// IsHashedClientSecret reports whether stored is already an HMAC digest.
func IsHashedClientSecret(stored string) bool {
	return strings.HasPrefix(stored, clientSecretPrefix)
}
