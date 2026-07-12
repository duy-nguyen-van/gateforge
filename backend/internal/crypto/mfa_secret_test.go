package crypto

import (
	"strings"
	"testing"
)

func TestEncryptMFASecret_EmptyPlaintext(t *testing.T) {
	_, err := EncryptMFASecret("key-material", "")
	if err == nil || !strings.Contains(err.Error(), "empty plaintext") {
		t.Fatalf("expected empty plaintext error, got %v", err)
	}
}

func TestDecryptMFASecret_InvalidCiphertext(t *testing.T) {
	key := "test-key-material-at-least-thirty-two-bytes-long"

	_, err := DecryptMFASecret(key, "not-valid-base64!!!")
	if err == nil {
		t.Fatal("expected base64 decode error")
	}

	_, err = DecryptMFASecret(key, "YWJj")
	if err == nil || !strings.Contains(err.Error(), "ciphertext too short") {
		t.Fatalf("expected ciphertext too short, got %v", err)
	}

	enc, err := EncryptMFASecret(key, "secret-value")
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecryptMFASecret("different-key-material-thirty-two-bytes", enc)
	if err == nil {
		t.Fatal("expected GCM authentication failure")
	}
}

func TestDeriveMFAKey(t *testing.T) {
	key := DeriveMFAKey("material")
	if len(key) != 32 {
		t.Fatalf("expected 32-byte key, got %d", len(key))
	}
}
