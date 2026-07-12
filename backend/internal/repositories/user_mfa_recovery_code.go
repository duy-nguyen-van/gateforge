package repositories

import (
	"context"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// UserMFARecoveryCodeRepository stores hashed recovery codes.
type UserMFARecoveryCodeRepository interface {
	ReplaceAllForUser(ctx context.Context, userID string, rows []*models.UserMFARecoveryCode) error
	FindUnusedByUserID(ctx context.Context, userID string) ([]models.UserMFARecoveryCode, error)
	MarkUsed(ctx context.Context, id string) error
}

type userMFARecoveryCodeRepository struct {
	db *db.PostgresDB
}

func ProvideUserMFARecoveryCodeRepository(db *db.PostgresDB) UserMFARecoveryCodeRepository {
	return &userMFARecoveryCodeRepository{db: db}
}

func (r *userMFARecoveryCodeRepository) ReplaceAllForUser(ctx context.Context, userID string, rows []*models.UserMFARecoveryCode) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserMFARecoveryCode{}).Error; err != nil {
			return errors.DatabaseError("Failed to clear recovery codes", err).
				WithOperation("replace_recovery_codes").
				WithResource("user_mfa_recovery_code")
		}
		for _, row := range rows {
			if err := tx.Create(row).Error; err != nil {
				return errors.DatabaseError("Failed to insert recovery code", err).
					WithOperation("replace_recovery_codes").
					WithResource("user_mfa_recovery_code")
			}
		}
		return nil
	})
}

func (r *userMFARecoveryCodeRepository) FindUnusedByUserID(ctx context.Context, userID string) ([]models.UserMFARecoveryCode, error) {
	var rows []models.UserMFARecoveryCode
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND used_at IS NULL", userID).
		Find(&rows).Error; err != nil {
		return nil, errors.DatabaseError("Failed to list recovery codes", err).
			WithOperation("list_recovery_codes").
			WithResource("user_mfa_recovery_code")
	}
	return rows, nil
}

func (r *userMFARecoveryCodeRepository) MarkUsed(ctx context.Context, id string) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&models.UserMFARecoveryCode{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", now)
	if res.Error != nil {
		return errors.DatabaseError("Failed to mark recovery code used", res.Error).
			WithOperation("mark_recovery_code_used").
			WithResource("user_mfa_recovery_code")
	}
	return nil
}
