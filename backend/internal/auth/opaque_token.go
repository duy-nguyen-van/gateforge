package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// NewOpaqueRefreshToken returns a URL-safe token for the client and a SHA-256 hex digest for storage.
func NewOpaqueRefreshToken() (raw string, hashHex string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("read random bytes: %w", err)
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	hashHex = hex.EncodeToString(sum[:])
	return raw, hashHex, nil
}

// HashOpaqueToken hashes a client-provided token the same way as issuance.
func HashOpaqueToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
