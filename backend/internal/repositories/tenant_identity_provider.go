package repositories

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// TenantIdentityProviderPatch updates tenant IdP settings.
type TenantIdentityProviderPatch struct {
	Enabled                    *bool
	OAuthClientID              string
	OAuthClientSecretPlaintext string // empty means keep existing secret
	OAuthClientSecretEncrypted string // set by service layer after encryption
}

// TenantIdentityProviderRepository reads/writes per-tenant IdP settings.
type TenantIdentityProviderRepository interface {
	IsProviderEnabled(ctx context.Context, tenantID, provider string) (bool, error)
	IsProviderConfigured(ctx context.Context, tenantID, provider string) (bool, error)
	GetByTenantAndProvider(ctx context.Context, tenantID, provider string) (*models.TenantIdentityProvider, error)
	SetProviderEnabled(ctx context.Context, tenantID, provider string, enabled bool) error
	UpdateProvider(ctx context.Context, tenantID, provider string, patch TenantIdentityProviderPatch) (*models.TenantIdentityProvider, error)
	ListByTenant(ctx context.Context, tenantID string) ([]models.TenantIdentityProvider, error)
}

type tenantIdentityProviderRepository struct {
	db *db.PostgresDB
}

// ProvideTenantIdentityProviderRepository wires tenant IdP settings.
func ProvideTenantIdentityProviderRepository(db *db.PostgresDB) TenantIdentityProviderRepository {
	return &tenantIdentityProviderRepository{db: db}
}

func (r *tenantIdentityProviderRepository) IsProviderEnabled(ctx context.Context, tenantID, provider string) (bool, error) {
	tip, err := r.GetByTenantAndProvider(ctx, tenantID, provider)
	if err != nil {
		if federationRepoNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return tip.Enabled, nil
}

func (r *tenantIdentityProviderRepository) IsProviderConfigured(ctx context.Context, tenantID, provider string) (bool, error) {
	tip, err := r.GetByTenantAndProvider(ctx, tenantID, provider)
	if err != nil {
		if federationRepoNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(tip.OAuthClientID) != "" && strings.TrimSpace(tip.OAuthClientSecretEncrypted) != "", nil
}

func (r *tenantIdentityProviderRepository) GetByTenantAndProvider(ctx context.Context, tenantID, provider string) (*models.TenantIdentityProvider, error) {
	var tip models.TenantIdentityProvider
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND provider = ?", tenantID, provider).
		First(&tip).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("Tenant identity provider not found", err).
				WithOperation("get_tenant_identity_provider").
				WithResource("tenant_identity_provider")
		}
		return nil, errors.DatabaseError("Failed to load tenant identity provider", err).
			WithOperation("get_tenant_identity_provider").
			WithResource("tenant_identity_provider")
	}
	return &tip, nil
}

func (r *tenantIdentityProviderRepository) SetProviderEnabled(ctx context.Context, tenantID, provider string, enabled bool) error {
	_, err := r.UpdateProvider(ctx, tenantID, provider, TenantIdentityProviderPatch{Enabled: &enabled})
	return err
}

func (r *tenantIdentityProviderRepository) UpdateProvider(ctx context.Context, tenantID, provider string, patch TenantIdentityProviderPatch) (*models.TenantIdentityProvider, error) {
	var out *models.TenantIdentityProvider
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var tip models.TenantIdentityProvider
		err := tx.Where("tenant_id = ? AND provider = ?", tenantID, provider).First(&tip).Error
		if err != nil {
			if stderrors.Is(err, gorm.ErrRecordNotFound) {
				tip = models.TenantIdentityProvider{
					BaseModel: models.NewBaseModel(),
					TenantID:  tenantID,
					Provider:  provider,
				}
			} else {
				return errors.DatabaseError("Failed to load tenant identity provider", err).
					WithOperation("update_tenant_identity_provider").
					WithResource("tenant_identity_provider")
			}
		}

		if patch.Enabled != nil {
			tip.Enabled = *patch.Enabled
		}
		if id := strings.TrimSpace(patch.OAuthClientID); id != "" {
			tip.OAuthClientID = id
		}
		if enc := strings.TrimSpace(patch.OAuthClientSecretEncrypted); enc != "" {
			tip.OAuthClientSecretEncrypted = enc
		}

		if tip.ID == "" {
			if err := tx.Create(&tip).Error; err != nil {
				return errors.DatabaseError("Failed to create tenant identity provider", err).
					WithOperation("update_tenant_identity_provider").
					WithResource("tenant_identity_provider")
			}
		} else if err := tx.Save(&tip).Error; err != nil {
			return errors.DatabaseError("Failed to update tenant identity provider", err).
				WithOperation("update_tenant_identity_provider").
				WithResource("tenant_identity_provider")
		}

		out = &tip
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *tenantIdentityProviderRepository) ListByTenant(ctx context.Context, tenantID string) ([]models.TenantIdentityProvider, error) {
	var rows []models.TenantIdentityProvider
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&rows).Error; err != nil {
		return nil, errors.DatabaseError("Failed to list tenant identity providers", err).
			WithOperation("list_tenant_identity_providers").
			WithResource("tenant_identity_provider")
	}
	return rows, nil
}

func federationRepoNotFound(err error) bool {
	appErr := errors.GetAppError(err)
	return appErr != nil && appErr.Type == errors.ErrorTypeNotFound
}
