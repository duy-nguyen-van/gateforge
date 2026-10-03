package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"
)

func lockoutKey(emailLower string) string {
	sum := sha256.Sum256([]byte(emailLower))
	return "iam:lockout:" + hex.EncodeToString(sum[:])
}

func (s *userService) loginLocked(ctx context.Context, emailLower string) bool {
	if s.loginCache == nil || s.cfg == nil || s.cfg.LockoutMaxFailures <= 0 {
		return false
	}
	raw, err := s.loginCache.Get(ctx, lockoutKey(emailLower))
	if err != nil {
		return false
	}
	n, convErr := strconv.Atoi(raw)
	if convErr != nil {
		return false
	}
	return n >= s.cfg.LockoutMaxFailures
}

func (s *userService) noteLoginFailure(ctx context.Context, emailLower, userID string) {
	if s.loginCache == nil || s.cfg == nil || s.cfg.LockoutMaxFailures <= 0 {
		return
	}
	window := s.cfg.LockoutWindow
	n, err := s.loginCache.Increment(ctx, lockoutKey(emailLower), window)
	if err != nil {
		return
	}
	if int(n) < s.cfg.LockoutMaxFailures {
		return
	}
	monitoring.AddAuthLockout(ctx)
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:    constants.AuditActionAuthLoginLocked,
		Result:    constants.AuditResultDenied,
		ActorType: constants.AuditActorTypeUser,
		ActorID:   userID,
	})
}

func (s *userService) clearLoginLock(ctx context.Context, emailLower string) {
	if s.loginCache == nil {
		return
	}
	_ = s.loginCache.Delete(ctx, lockoutKey(emailLower))
}
