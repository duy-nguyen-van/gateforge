package repositories

import (
	"context"
	stderrors "errors"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// UserMFATOTPRepository persists TOTP enrollment.
type UserMFATOTPRepository interface {
	GetByUserID(ctx context.Context, userID string) (*models.UserMFATOTP, error)
	// GetActiveByUserID returns a row only when TOTP is enabled; nil, nil when MFA is off.
	GetActiveByUserID(ctx context.Context, userID string) (*models.UserMFATOTP, error)
	UpsertPending(ctx context.Context, row *models.UserMFATOTP) error
	MarkVerifiedAndEnabled(ctx context.Context, userID string) error
	Disable(ctx context.Context, userID string) error
	CountEnabled(ctx context.Context) (int64, error)
}

type userMFATOTPRepository struct {
	db *db.PostgresDB
}

func ProvideUserMFATOTPRepository(db *db.PostgresDB) UserMFATOTPRepository {
	return &userMFATOTPRepository{db: db}
}

func (r *userMFATOTPRepository) GetByUserID(ctx context.Context, userID string) (*models.UserMFATOTP, error) {
	var row models.UserMFATOTP
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&row).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("User MFA TOTP", err).
				WithOperation("get_user_mfa_totp").
				WithResource("user_mfa_totp")
		}
		return nil, errors.DatabaseError("Failed to load MFA TOTP", err).
			WithOperation("get_user_mfa_totp").
			WithResource("user_mfa_totp")
	}
	return &row, nil
}

func (r *userMFATOTPRepository) GetActiveByUserID(ctx context.Context, userID string) (*models.UserMFATOTP, error) {
	var row models.UserMFATOTP
	err := r.db.WithContext(ctx).Where("user_id = ? AND enabled = ?", userID, true).First(&row).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.DatabaseError("Failed to load active MFA TOTP", err).
			WithOperation("get_active_user_mfa_totp").
			WithResource("user_mfa_totp")
	}
	return &row, nil
}

func (r *userMFATOTPRepository) UpsertPending(ctx context.Context, row *models.UserMFATOTP) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("user_id = ?", row.UserID).Delete(&models.UserMFATOTP{}).Error; err != nil {
			return errors.DatabaseError("Failed to reset MFA TOTP", err).
				WithOperation("upsert_user_mfa_totp").
				WithResource("user_mfa_totp")
		}
		if err := tx.Create(row).Error; err != nil {
			return errors.DatabaseError("Failed to create MFA TOTP", err).
				WithOperation("upsert_user_mfa_totp").
				WithResource("user_mfa_totp")
		}
		return nil
	})
}

func (r *userMFATOTPRepository) MarkVerifiedAndEnabled(ctx context.Context, userID string) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&models.UserMFATOTP{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"enabled":     true,
			"verified_at": now,
		})
	if res.Error != nil {
		return errors.DatabaseError("Failed to enable MFA TOTP", res.Error).
			WithOperation("verify_user_mfa_totp").
			WithResource("user_mfa_totp")
	}
	return nil
}

func (r *userMFATOTPRepository) Disable(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.UserMFATOTP{}).Error; err != nil {
		return errors.DatabaseError("Failed to disable MFA TOTP", err).
			WithOperation("disable_user_mfa_totp").
			WithResource("user_mfa_totp")
	}
	return nil
}

func (r *userMFATOTPRepository) CountEnabled(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserMFATOTP{}).Where("enabled = ?", true).Count(&count).Error; err != nil {
		return 0, errors.DatabaseError("Failed to count MFA enrollments", err).
			WithOperation("count_mfa_enabled").
			WithResource("user_mfa_totp")
	}
	return count, nil
}
