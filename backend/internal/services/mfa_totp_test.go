package services

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

func TestTOTPValidateKnownSecret(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if !totp.Validate(code, secret) {
		t.Fatal("expected valid TOTP")
	}
	if totp.Validate("000000", secret) {
		t.Fatal("expected invalid TOTP")
	}
}
