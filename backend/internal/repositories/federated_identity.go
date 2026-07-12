package repositories

import (
	"context"
	stderrors "errors"

	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// FederatedIdentityRepository persists upstream IdP subject → user links.
type FederatedIdentityRepository interface {
	GetByProviderSubject(ctx context.Context, provider, subject string) (*models.FederatedIdentity, error)
	Create(ctx context.Context, row *models.FederatedIdentity) error
	ListByUserID(ctx context.Context, userID string) ([]models.FederatedIdentity, error)
}

type federatedIdentityRepository struct {
	db *db.PostgresDB
}

// ProvideFederatedIdentityRepository wires federated identity persistence.
func ProvideFederatedIdentityRepository(db *db.PostgresDB) FederatedIdentityRepository {
	return &federatedIdentityRepository{db: db}
}

func (r *federatedIdentityRepository) GetByProviderSubject(ctx context.Context, provider, subject string) (*models.FederatedIdentity, error) {
	var fi models.FederatedIdentity
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("User.PasswordCredential").
		Where("provider = ? AND subject = ?", provider, subject).
		First(&fi).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("FederatedIdentity", err).
				WithOperation("get_federated_identity").
				WithResource("federated_identity")
		}
		return nil, errors.DatabaseError("Failed to load federated identity", err).
			WithOperation("get_federated_identity").
			WithResource("federated_identity")
	}
	return &fi, nil
}

func (r *federatedIdentityRepository) Create(ctx context.Context, row *models.FederatedIdentity) error {
	return Create(ctx, r.db, row, DBOp{Operation: "create_federated_identity", Resource: "federated_identity"}, "Failed to create federated identity")
}

func (r *federatedIdentityRepository) ListByUserID(ctx context.Context, userID string) ([]models.FederatedIdentity, error) {
	var rows []models.FederatedIdentity
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&rows).Error
	if err != nil {
		return nil, errors.DatabaseError("Failed to list federated identities", err).
			WithOperation("list_federated_identities").
			WithResource("federated_identity")
	}
	return rows, nil
}
