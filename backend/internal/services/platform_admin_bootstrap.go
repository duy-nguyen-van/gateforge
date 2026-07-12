package services

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// PlatformAdminBootstrap ensures at least one platform admin exists on first startup.
type PlatformAdminBootstrap interface {
	Run(ctx context.Context) error
}

type platformAdminBootstrap struct {
	cfg            *config.Config
	users          repositories.UserRepository
	membershipRepo repositories.TenantMembershipRepository
}

// ProvidePlatformAdminBootstrap wires first-run platform admin creation when no admins exist.
func ProvidePlatformAdminBootstrap(
	cfg *config.Config,
	users repositories.UserRepository,
	membershipRepo repositories.TenantMembershipRepository,
) PlatformAdminBootstrap {
	return &platformAdminBootstrap{cfg: cfg, users: users, membershipRepo: membershipRepo}
}

func (b *platformAdminBootstrap) Run(ctx context.Context) error {
	count, err := b.users.CountPlatformAdmins(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	email := strings.TrimSpace(b.cfg.BootstrapAdminEmail)
	password := b.cfg.BootstrapAdminPassword
	if email != "" && password != "" {
		return b.bootstrapFromCredentials(ctx, email, password)
	}

	logger.Log.Warn("No platform admin configured; admin console APIs will be unavailable until a platform admin exists",
		zap.String("operation", "platform_admin_bootstrap"),
		zap.String("hint", "Set BOOTSTRAP_ADMIN_EMAIL and BOOTSTRAP_ADMIN_PASSWORD for first startup"))
	return nil
}

func (b *platformAdminBootstrap) bootstrapFromCredentials(ctx context.Context, email, password string) error {
	if len(password) < 8 {
		return errors.ValidationError("BOOTSTRAP_ADMIN_PASSWORD must be at least 8 characters", nil)
	}

	tenantID := b.cfg.DefaultTenantID
	emailLower := strings.ToLower(email)

	existing, err := b.users.GetByEmailLower(ctx, emailLower)
	if err == nil && existing != nil {
		if err := b.users.SetPlatformAdmin(ctx, existing.ID, true); err != nil {
			return err
		}
		if err := b.ensureMembership(ctx, existing.ID, tenantID); err != nil {
			return err
		}
		logger.Log.Info("platform admin bootstrapped (existing user promoted)",
			zap.String("operation", "platform_admin_bootstrap"),
			zap.String("user_id", existing.ID),
			zap.String("tenant_id", tenantID))
		return nil
	}
	var appErr *errors.AppError
	if !stderrors.As(err, &appErr) || appErr.Type != errors.ErrorTypeNotFound {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), constants.BcryptCost)
	if err != nil {
		return errors.InternalError("Failed to hash bootstrap admin password", err).
			WithOperation("platform_admin_bootstrap").
			WithResource("user")
	}

	u := &models.User{
		BaseModel:       models.NewBaseModel(),
		FirstName:       "Platform",
		LastName:        "Admin",
		Email:           email,
		EmailVerified:   false,
		Status:          constants.UserStatusActive,
		IsPlatformAdmin: true,
	}
	if err := b.users.CreateWithPasswordHash(ctx, u, string(hash)); err != nil {
		return err
	}
	if err := b.ensureMembership(ctx, u.ID, tenantID); err != nil {
		return err
	}

	logger.Log.Info("platform admin bootstrapped (new user created)",
		zap.String("operation", "platform_admin_bootstrap"),
		zap.String("user_id", u.ID),
		zap.String("tenant_id", tenantID))
	return nil
}

func (b *platformAdminBootstrap) ensureMembership(ctx context.Context, userID, tenantID string) error {
	ok, err := b.membershipRepo.ExistsActive(ctx, userID, tenantID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return b.membershipRepo.Create(ctx, &models.TenantMembership{
		BaseModel: models.NewBaseModel(),
		UserID:    userID,
		TenantID:  tenantID,
		Role:      constants.TenantMembershipRoleAdmin,
		Status:    constants.TenantMembershipStatusActive,
	})
}
