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

// AuthorizationCodeRepository persists OAuth2 authorization codes (PKCE).
type AuthorizationCodeRepository interface {
	Create(ctx context.Context, row *models.AuthorizationCode) error
	TakeByCode(ctx context.Context, code string) (*models.AuthorizationCode, error)
	DeleteByCode(ctx context.Context, code string) error
}

type authorizationCodeRepository struct {
	db *db.PostgresDB
}

// ProvideAuthorizationCodeRepository wires authorization code persistence.
func ProvideAuthorizationCodeRepository(db *db.PostgresDB) AuthorizationCodeRepository {
	return &authorizationCodeRepository{db: db}
}

func (r *authorizationCodeRepository) Create(ctx context.Context, row *models.AuthorizationCode) error {
	return Create(ctx, r.db, row, DBOp{Operation: "create_authorization_code", Resource: "authorization_code"}, "Failed to create authorization code")
}

func (r *authorizationCodeRepository) TakeByCode(ctx context.Context, code string) (*models.AuthorizationCode, error) {
	var row models.AuthorizationCode
	err := r.db.WithContext(ctx).
		Where("code = ? AND expires_at > ?", code, time.Now().UTC()).
		First(&row).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("Authorization code", err).
				WithOperation("get_authorization_code").
				WithResource("authorization_code")
		}
		return nil, errors.DatabaseError("Failed to load authorization code", err).
			WithOperation("get_authorization_code").
			WithResource("authorization_code")
	}
	return &row, nil
}

func (r *authorizationCodeRepository) DeleteByCode(ctx context.Context, code string) error {
	res := r.db.WithContext(ctx).Unscoped().Where("code = ?", code).Delete(&models.AuthorizationCode{})
	if res.Error != nil {
		return errors.DatabaseError("Failed to delete authorization code", res.Error).
			WithOperation("delete_authorization_code").
			WithResource("authorization_code")
	}
	return nil
}
