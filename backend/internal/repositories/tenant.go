package repositories

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// TenantPatch updates mutable tenant fields.
type TenantPatch struct {
	Name   *string
	Domain *string
}

// TenantRepository loads tenant records for admin APIs.
type TenantRepository interface {
	List(ctx context.Context, pr *dtos.PageableRequest) (*dtos.DataResponse[models.Tenant], error)
	GetByID(ctx context.Context, id string) (*models.Tenant, error)
	GetByDomain(ctx context.Context, domain string) (*models.Tenant, error)
	Create(ctx context.Context, tenant *models.Tenant) error
	Update(ctx context.Context, id string, patch TenantPatch) (*models.Tenant, error)
	Delete(ctx context.Context, id string) error
	DomainTaken(ctx context.Context, domain, excludeID string) (bool, error)
	CountUsersByTenantID(ctx context.Context, tenantID string) (int64, error)
}

type tenantRepository struct {
	db *db.PostgresDB
}

// ProvideTenantRepository wires tenant persistence.
func ProvideTenantRepository(db *db.PostgresDB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) List(ctx context.Context, pr *dtos.PageableRequest) (*dtos.DataResponse[models.Tenant], error) {
	query := r.db.WithContext(ctx).Model(&models.Tenant{})
	return Paginate[models.Tenant](ctx, query, pr, PaginateOptions{
		OrderBy:             "created_at DESC",
		CountFailureMessage: "Failed to count tenants",
		FindFailureMessage:  "Failed to list tenants",
		DBOp:                DBOp{Operation: "list_tenants", Resource: "tenant"},
	})
}

func (r *tenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	return Create(ctx, r.db, tenant, DBOp{Operation: "create_tenant", Resource: "tenant"}, "Failed to create tenant")
}

func (r *tenantRepository) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	var t models.Tenant
	if err := r.db.WithContext(ctx).First(&t, "id = ?", id).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("Tenant", err).
				WithOperation("get_tenant").
				WithResource("tenant")
		}
		return nil, errors.DatabaseError("Failed to get tenant", err).
			WithOperation("get_tenant").
			WithResource("tenant")
	}
	return &t, nil
}

func (r *tenantRepository) GetByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	var t models.Tenant
	if err := r.db.WithContext(ctx).Where("domain = ?", domain).First(&t).Error; err != nil {
		return nil, errors.NotFoundError("Tenant", err).
			WithOperation("get_tenant_by_domain").
			WithResource("tenant")
	}
	return &t, nil
}

func (r *tenantRepository) Update(ctx context.Context, id string, patch TenantPatch) (*models.Tenant, error) {
	updates := map[string]any{}
	if patch.Name != nil {
		updates["name"] = *patch.Name
	}
	if patch.Domain != nil {
		updates["domain"] = *patch.Domain
	}
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	res := r.db.WithContext(ctx).Model(&models.Tenant{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return nil, errors.DatabaseError("Failed to update tenant", res.Error).
			WithOperation("update_tenant").
			WithResource("tenant")
	}
	if res.RowsAffected == 0 {
		return nil, errors.NotFoundError("Tenant", nil).
			WithOperation("update_tenant").
			WithResource("tenant")
	}
	return r.GetByID(ctx, id)
}

func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&models.Tenant{}, "id = ?", id)
	if res.Error != nil {
		return errors.DatabaseError("Failed to delete tenant", res.Error).
			WithOperation("delete_tenant").
			WithResource("tenant")
	}
	if res.RowsAffected == 0 {
		return errors.NotFoundError("Tenant", nil).
			WithOperation("delete_tenant").
			WithResource("tenant")
	}
	return nil
}

func (r *tenantRepository) DomainTaken(ctx context.Context, domain, excludeID string) (bool, error) {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return false, nil
	}
	query := r.db.WithContext(ctx).Model(&models.Tenant{}).Where("domain = ?", domain)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, errors.DatabaseError("Failed to check tenant domain", err).
			WithOperation("check_tenant_domain").
			WithResource("tenant")
	}
	return count > 0, nil
}

func (r *tenantRepository) CountUsersByTenantID(ctx context.Context, tenantID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.TenantMembership{}).
		Where("tenant_id = ? AND status = ?", tenantID, "active").
		Count(&count).Error; err != nil {
		return 0, errors.DatabaseError("Failed to count tenant users", err).
			WithOperation("count_tenant_users").
			WithResource("user")
	}
	return count, nil
}
