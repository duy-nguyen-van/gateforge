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

const retentionLockKey int64 = 742019

// RetentionCounts is how many rows one purge pass removed.
type RetentionCounts struct {
	AuthorizationCodes int64
	Sessions           int64
	RefreshTokens      int64
}

// RetentionRepository deletes expired auth data in bounded batches.
type RetentionRepository interface {
	Purge(ctx context.Context, refreshBefore time.Time, batch int) (locked bool, counts RetentionCounts, err error)
}

type retentionRepository struct {
	db *db.PostgresDB
}

func ProvideRetentionRepository(database *db.PostgresDB) RetentionRepository {
	return &retentionRepository{db: database}
}

func (r *retentionRepository) Purge(ctx context.Context, refreshBefore time.Time, batch int) (bool, RetentionCounts, error) {
	if batch <= 0 {
		batch = 1000
	}
	var counts RetentionCounts
	locked := true
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if r.db.Name() == "postgres" {
			var ok bool
			if err := tx.Raw("SELECT pg_try_advisory_xact_lock(?)", retentionLockKey).Scan(&ok).Error; err != nil {
				return err
			}
			if !ok {
				locked = false
				return nil
			}
		}
		var err error
		counts.AuthorizationCodes, err = deleteExpiredBatch(tx, &models.AuthorizationCode{}, "expires_at < ?", []any{time.Now().UTC()}, batch)
		if err != nil {
			return err
		}
		counts.Sessions, err = deleteExpiredBatch(tx, &models.Session{}, "expires_at IS NOT NULL AND expires_at < ?", []any{time.Now().UTC()}, batch)
		if err != nil {
			return err
		}
		counts.RefreshTokens, err = deleteExpiredBatch(tx, &models.RefreshToken{}, "revoked = ? OR expires_at < ?", []any{true, refreshBefore}, batch)
		return err
	})
	if err != nil {
		return false, RetentionCounts{}, errors.DatabaseError("Failed to purge expired rows", err).
			WithOperation("retention_purge").
			WithResource("database")
	}
	return locked, counts, nil
}

func deleteExpiredBatch(tx *gorm.DB, model any, where string, args []any, batch int) (int64, error) {
	var ids []string
	if err := tx.Model(model).Where(where, args...).Limit(batch).Pluck("id", &ids).Error; err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, nil
	}
	res := tx.Unscoped().Where("id IN ?", ids).Delete(model)
	if res.Error != nil && !stderrors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
