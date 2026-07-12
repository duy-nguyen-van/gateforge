package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// DeriveMFAKey normalizes arbitrary-length key material to 32 bytes for AES-256.
func DeriveMFAKey(material string) []byte {
	sum := sha256.Sum256([]byte(material))
	return sum[:]
}

// EncryptMFASecret returns base64(nonce|ciphertext) using AES-GCM.
func EncryptMFASecret(keyMaterial, plaintext string) (string, error) {
	if plaintext == "" {
		return "", fmt.Errorf("empty plaintext")
	}
	key := DeriveMFAKey(keyMaterial)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	out := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.RawStdEncoding.EncodeToString(out), nil
}

// DecryptMFASecret reverses EncryptMFASecret.
func DecryptMFASecret(keyMaterial, encoded string) (string, error) {
	raw, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	key := DeriveMFAKey(keyMaterial)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(raw) < ns {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ct := raw[:ns], raw[ns:]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}
