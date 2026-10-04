package services

import (
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/errors"
)

const (
	passwordMinLength = 12
	passwordMaxLength = 128
)

var commonPasswords = map[string]struct{}{
	"password":     {},
	"password123":  {},
	"password1234": {},
	"12345678":     {},
	"123456789":    {},
	"1234567890":   {},
	"qwertyuiop":   {},
	"letmein":      {},
	"welcome":      {},
	"changeme":     {},
	"admin123":     {},
	"iloveyou":     {},
}

// ValidatePassword enforces length and a small denylist for newly chosen passwords.
func ValidatePassword(password string) error {
	if len(password) < passwordMinLength || len(password) > passwordMaxLength {
		return errors.ValidationError("Password must be between 12 and 128 characters", nil)
	}
	if _, ok := commonPasswords[strings.ToLower(password)]; ok {
		return errors.ValidationError("Password is too common", nil)
	}
	return nil
}
