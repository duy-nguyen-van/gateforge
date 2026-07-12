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

// RefreshTokenUsage aggregates refresh token metrics for an OAuth client record.
type RefreshTokenUsage struct {
	TotalIssued  int64
	ActiveCount  int64
	LastIssuedAt *time.Time
}

// RefreshTokenRepository persists opaque refresh tokens (hashed at rest).
type RefreshTokenRepository interface {
	Create(ctx context.Context, rt *models.RefreshToken) error
	FindValidByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeByID(ctx context.Context, id string) error
	// RevokeAllValidForUser marks every non-expired, non-revoked refresh token for the user as revoked (global logout).
	RevokeAllValidForUser(ctx context.Context, userID string) error
	// RevokeAndCreate rotates a refresh token in one transaction (old revoked, new persisted).
	RevokeAndCreate(ctx context.Context, oldID string, newRT *models.RefreshToken) error
	UsageByClientRecordID(ctx context.Context, clientRecordID string) (*RefreshTokenUsage, error)
}

type refreshTokenRepository struct {
	db *db.PostgresDB
}

func ProvideRefreshTokenRepository(db *db.PostgresDB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, rt *models.RefreshToken) error {
	return Create(ctx, r.db, rt, DBOp{Operation: "create_refresh_token", Resource: "refresh_token"}, "Failed to create refresh token")
}

func (r *refreshTokenRepository) FindValidByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked = ? AND expires_at > ?", tokenHash, false, time.Now().UTC()).
		First(&rt).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("RefreshToken", err).
				WithOperation("refresh_token_lookup").
				WithResource("refresh_token")
		}
		return nil, errors.DatabaseError("Failed to load refresh token", err).
			WithOperation("refresh_token_lookup").
			WithResource("refresh_token")
	}
	return &rt, nil
}

func (r *refreshTokenRepository) RevokeByID(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Model(&models.RefreshToken{}).Where("id = ?", id).Update("revoked", true)
	if res.Error != nil {
		return errors.DatabaseError("Failed to revoke refresh token", res.Error).
			WithOperation("revoke_refresh_token").
			WithResource("refresh_token")
	}
	return nil
}

func (r *refreshTokenRepository) RevokeAllValidForUser(ctx context.Context, userID string) error {
	res := r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = ? AND expires_at > ?", userID, false, time.Now().UTC()).
		Update("revoked", true)
	if res.Error != nil {
		return errors.DatabaseError("Failed to revoke refresh tokens", res.Error).
			WithOperation("revoke_all_refresh_tokens_user").
			WithResource("refresh_token")
	}
	return nil
}

func (r *refreshTokenRepository) RevokeAndCreate(ctx context.Context, oldID string, newRT *models.RefreshToken) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.RefreshToken{}).Where("id = ?", oldID).Update("revoked", true).Error; err != nil {
			return err
		}
		return tx.Create(newRT).Error
	})
	if err != nil {
		return errors.DatabaseError("Failed to rotate refresh token", err).
			WithOperation("rotate_refresh_token").
			WithResource("refresh_token")
	}
	return nil
}

func (r *refreshTokenRepository) UsageByClientRecordID(ctx context.Context, clientRecordID string) (*RefreshTokenUsage, error) {
	now := time.Now().UTC()
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("client_record_id = ?", clientRecordID).
		Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("Failed to count refresh tokens for client", err).
			WithOperation("refresh_token_usage").
			WithResource("refresh_token")
	}

	var active int64
	if err := r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("client_record_id = ? AND revoked = ? AND expires_at > ?", clientRecordID, false, now).
		Count(&active).Error; err != nil {
		return nil, errors.DatabaseError("Failed to count active refresh tokens for client", err).
			WithOperation("refresh_token_usage").
			WithResource("refresh_token")
	}

	var lastIssued *time.Time
	var row models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("client_record_id = ?", clientRecordID).
		Order("created_at DESC").
		Limit(1).
		First(&row).Error
	if err == nil {
		lastIssued = &row.CreatedAt
	} else if !stderrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.DatabaseError("Failed to load last refresh token for client", err).
			WithOperation("refresh_token_usage").
			WithResource("refresh_token")
	}

	return &RefreshTokenUsage{
		TotalIssued:  total,
		ActiveCount:  active,
		LastIssuedAt: lastIssued,
	}, nil
}
