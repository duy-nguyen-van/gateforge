package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"

	"golang.org/x/crypto/bcrypt"
)

type passwordResetMailer interface {
	SendPasswordResetEmail(ctx context.Context, userEmail, resetURL string) error
}

// ProvidePasswordResetMailer adapts EmailService for UserService.
// FX does not use a concrete type to fill an interface parameter on its own.
func ProvidePasswordResetMailer(svc EmailService) passwordResetMailer {
	return svc
}

func passwordResetKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return "iam:password-reset:" + hex.EncodeToString(sum[:])
}

func (s *userService) ForgotPassword(ctx context.Context, email string) error {
	return monitoring.ObserveErr(ctx, userTracer, "UserService.ForgotPassword", nil, func(ctx context.Context) error {
		return s.forgotPassword(ctx, email)
	})
}

func (s *userService) forgotPassword(ctx context.Context, email string) error {
	emailLower := strings.ToLower(strings.TrimSpace(email))
	if emailLower == "" || s.loginCache == nil {
		return nil
	}
	u, err := s.userRepo.GetByEmailLower(ctx, emailLower)
	if err != nil {
		return nil
	}
	token, _, err := auth.NewOpaqueRefreshToken()
	if err != nil {
		return errors.InternalError("Failed to create password reset token", err)
	}
	ttl := s.cfg.PasswordResetTTL
	if err := s.loginCache.Set(ctx, passwordResetKey(token), u.ID, ttl); err != nil {
		return err
	}
	if s.resetMailer == nil {
		return nil
	}
	base := strings.TrimRight(s.cfg.AppBaseURL, "/")
	if base == "" {
		base = "http://localhost:3000"
	}
	resetURL := base + "/reset-password?token=" + token
	if err := s.resetMailer.SendPasswordResetEmail(ctx, u.Email, resetURL); err != nil {
		return err
	}
	return nil
}

func (s *userService) ResetPassword(ctx context.Context, token, newPassword string) error {
	return monitoring.ObserveErr(ctx, userTracer, "UserService.ResetPassword", nil, func(ctx context.Context) error {
		return s.resetPassword(ctx, token, newPassword)
	})
}

func (s *userService) resetPassword(ctx context.Context, token, newPassword string) error {
	if err := ValidatePassword(newPassword); err != nil {
		return err
	}
	if s.loginCache == nil {
		return errors.ValidationError("Password reset is unavailable", nil)
	}
	key := passwordResetKey(strings.TrimSpace(token))
	userID, err := s.loginCache.Get(ctx, key)
	if err != nil || userID == "" {
		return errors.CodedUnauthorized(constants.InvalidRefreshToken, nil).
			WithOperation("reset_password").
			WithResource("user")
	}
	if err := s.loginCache.Delete(ctx, key); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), constants.BcryptCost)
	if err != nil {
		return errors.InternalError("Failed to hash password", err)
	}
	if err := s.userRepo.UpdatePasswordHash(ctx, userID, string(hash)); err != nil {
		return err
	}
	if err := s.refreshTokenRepo.RevokeAllValidForUser(ctx, userID); err != nil {
		return err
	}
	if s.sessions != nil {
		if err := s.sessions.InvalidateAllForUser(ctx, userID); err != nil {
			return err
		}
	}
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAuthPasswordReset,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      userID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   userID,
	})
	return nil
}
