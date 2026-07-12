package repositories

import (
	"context"
	stderrors "errors"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// TenantMembershipRepository persists user↔tenant memberships.
type TenantMembershipRepository interface {
	Create(ctx context.Context, m *models.TenantMembership) error
	GetActive(ctx context.Context, userID, tenantID string) (*models.TenantMembership, error)
	ExistsActive(ctx context.Context, userID, tenantID string) (bool, error)
	ListByUserID(ctx context.Context, userID string) ([]models.TenantMembership, error)
	ListByUserIDPaginated(ctx context.Context, userID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error)
	ListByTenantID(ctx context.Context, tenantID string) ([]models.TenantMembership, error)
	ListByTenantIDPaginated(ctx context.Context, tenantID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error)
	CountByTenantID(ctx context.Context, tenantID string) (int64, error)
	Delete(ctx context.Context, userID, tenantID string) error
}

type tenantMembershipRepository struct {
	db *db.PostgresDB
}

// ProvideTenantMembershipRepository wires membership persistence.
func ProvideTenantMembershipRepository(db *db.PostgresDB) TenantMembershipRepository {
	return &tenantMembershipRepository{db: db}
}

func (r *tenantMembershipRepository) Create(ctx context.Context, m *models.TenantMembership) error {
	return Create(ctx, r.db, m, DBOp{Operation: "create_tenant_membership", Resource: "tenant_membership"}, "Failed to create tenant membership")
}

func (r *tenantMembershipRepository) GetActive(ctx context.Context, userID, tenantID string) (*models.TenantMembership, error) {
	var m models.TenantMembership
	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Where("user_id = ? AND tenant_id = ? AND status = ?", userID, tenantID, constants.TenantMembershipStatusActive).
		First(&m).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("Tenant membership", err).
				WithOperation("get_tenant_membership").
				WithResource("tenant_membership")
		}
		return nil, errors.DatabaseError("Failed to get tenant membership", err).
			WithOperation("get_tenant_membership").
			WithResource("tenant_membership")
	}
	return &m, nil
}

func (r *tenantMembershipRepository) ExistsActive(ctx context.Context, userID, tenantID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.TenantMembership{}).
		Where("user_id = ? AND tenant_id = ? AND status = ?", userID, tenantID, constants.TenantMembershipStatusActive).
		Count(&count).Error
	if err != nil {
		return false, errors.DatabaseError("Failed to check tenant membership", err).
			WithOperation("exists_tenant_membership").
			WithResource("tenant_membership")
	}
	return count > 0, nil
}

func (r *tenantMembershipRepository) ListByUserID(ctx context.Context, userID string) ([]models.TenantMembership, error) {
	var rows []models.TenantMembership
	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Where("user_id = ? AND status = ?", userID, constants.TenantMembershipStatusActive).
		Order("created_at ASC").
		Find(&rows).Error
	if err != nil {
		return nil, errors.DatabaseError("Failed to list tenant memberships", err).
			WithOperation("list_tenant_memberships").
			WithResource("tenant_membership")
	}
	return rows, nil
}

func (r *tenantMembershipRepository) ListByUserIDPaginated(ctx context.Context, userID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	base := r.db.WithContext(ctx).Model(&models.TenantMembership{}).
		Where("user_id = ? AND status = ?", userID, constants.TenantMembershipStatusActive)
	return PaginateWithFind[models.TenantMembership](ctx, base, base.Preload("Tenant"), pr, PaginateOptions{
		OrderBy:             "created_at ASC",
		MaxPageSize:         constants.MaxPageSize,
		CountFailureMessage: "Failed to count tenant memberships",
		FindFailureMessage:  "Failed to list tenant memberships",
		DBOp:                DBOp{Operation: "list_tenant_memberships", Resource: "tenant_membership"},
	})
}

func (r *tenantMembershipRepository) ListByTenantID(ctx context.Context, tenantID string) ([]models.TenantMembership, error) {
	var rows []models.TenantMembership
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ?", tenantID, constants.TenantMembershipStatusActive).
		Find(&rows).Error
	if err != nil {
		return nil, errors.DatabaseError("Failed to list tenant memberships by tenant", err).
			WithOperation("list_tenant_memberships").
			WithResource("tenant_membership")
	}
	return rows, nil
}

func (r *tenantMembershipRepository) ListByTenantIDPaginated(ctx context.Context, tenantID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	base := r.db.WithContext(ctx).Model(&models.TenantMembership{}).
		Where("tenant_id = ? AND status = ?", tenantID, constants.TenantMembershipStatusActive)
	return PaginateWithFind[models.TenantMembership](ctx, base, base.Preload("User"), pr, PaginateOptions{
		OrderBy:             "created_at ASC",
		MaxPageSize:         constants.MaxPageSize,
		CountFailureMessage: "Failed to count tenant memberships",
		FindFailureMessage:  "Failed to list tenant memberships by tenant",
		DBOp:                DBOp{Operation: "list_tenant_memberships", Resource: "tenant_membership"},
	})
}

func (r *tenantMembershipRepository) CountByTenantID(ctx context.Context, tenantID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.TenantMembership{}).
		Where("tenant_id = ? AND status = ?", tenantID, constants.TenantMembershipStatusActive).
		Count(&count).Error
	if err != nil {
		return 0, errors.DatabaseError("Failed to count tenant memberships", err).
			WithOperation("count_tenant_memberships").
			WithResource("tenant_membership")
	}
	return count, nil
}

func (r *tenantMembershipRepository) Delete(ctx context.Context, userID, tenantID string) error {
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Delete(&models.TenantMembership{})
	if res.Error != nil {
		return errors.DatabaseError("Failed to delete tenant membership", res.Error).
			WithOperation("delete_tenant_membership").
			WithResource("tenant_membership")
	}
	if res.RowsAffected == 0 {
		return errors.NotFoundError("Tenant membership", nil).
			WithOperation("delete_tenant_membership").
			WithResource("tenant_membership")
	}
	return nil
}
