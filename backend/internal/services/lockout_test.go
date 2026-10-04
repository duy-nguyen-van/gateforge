package services

import (
	"context"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"

	"github.com/stretchr/testify/require"
)

func TestUserService_LockoutThreshold(t *testing.T) {
	users := newUserTestRepo()
	users.seed("locked@example.com", "secret123")
	cfg := testConfig()
	cfg.LockoutMaxFailures = 2
	cfg.LockoutWindow = time.Minute
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	memberships := &stubMembershipRepo{active: map[string]map[string]bool{}}
	svc := ProvideUserService(
		users, memberships, newRefreshTokenTestRepo(),
		ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships),
		cfg, tokenSvc, &auditCapture{}, newMemCache(), nil, nil,
	)

	for i := 0; i < 2; i++ {
		_, err := svc.AuthenticateUser(context.Background(), &dtos.LoginRequest{
			Email: "locked@example.com", Password: "wrong-password-1",
		})
		require.Error(t, err)
	}
	_, err = svc.AuthenticateUser(context.Background(), &dtos.LoginRequest{
		Email: "locked@example.com", Password: "secret123",
	})
	require.Error(t, err)
}

func TestValidatePassword(t *testing.T) {
	require.Error(t, ValidatePassword("short"))
	require.Error(t, ValidatePassword("password123"))
	require.NoError(t, ValidatePassword("correct-horse-1"))
}
