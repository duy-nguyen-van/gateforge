package repositories

import (
	"context"
	stderrors "errors"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// TenantInviteRepository persists organization invites for emails that may not have an account yet.
type TenantInviteRepository interface {
	Create(ctx context.Context, invite *models.TenantInvite) error
	Save(ctx context.Context, invite *models.TenantInvite) error
	GetByID(ctx context.Context, id string) (*models.TenantInvite, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*models.TenantInvite, error)
	FindPendingByEmailAndTenant(ctx context.Context, emailLower, tenantID string) (*models.TenantInvite, error)
	ListPendingByTenant(ctx context.Context, tenantID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantInvite], error)
	ListPending(ctx context.Context, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantInvite], error)
}

type tenantInviteRepository struct {
	db *db.PostgresDB
}

// ProvideTenantInviteRepository wires invite persistence.
func ProvideTenantInviteRepository(db *db.PostgresDB) TenantInviteRepository {
	return &tenantInviteRepository{db: db}
}

func (r *tenantInviteRepository) Create(ctx context.Context, invite *models.TenantInvite) error {
	return Create(ctx, r.db, invite, DBOp{Operation: "create_tenant_invite", Resource: "tenant_invite"}, "Failed to create tenant invite")
}

func (r *tenantInviteRepository) Save(ctx context.Context, invite *models.TenantInvite) error {
	if err := r.db.WithContext(ctx).Save(invite).Error; err != nil {
		return errors.DatabaseError("Failed to save tenant invite", err).
			WithOperation("save_tenant_invite").
			WithResource("tenant_invite")
	}
	return nil
}

func (r *tenantInviteRepository) GetByID(ctx context.Context, id string) (*models.TenantInvite, error) {
	var invite models.TenantInvite
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&invite).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("Tenant invite", err).
				WithOperation("get_tenant_invite").
				WithResource("tenant_invite")
		}
		return nil, errors.DatabaseError("Failed to get tenant invite", err).
			WithOperation("get_tenant_invite").
			WithResource("tenant_invite")
	}
	return &invite, nil
}

func (r *tenantInviteRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.TenantInvite, error) {
	var invite models.TenantInvite
	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Where("token_hash = ?", tokenHash).
		First(&invite).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("Tenant invite", err).
				WithOperation("get_tenant_invite").
				WithResource("tenant_invite")
		}
		return nil, errors.DatabaseError("Failed to get tenant invite", err).
			WithOperation("get_tenant_invite").
			WithResource("tenant_invite")
	}
	return &invite, nil
}

func (r *tenantInviteRepository) FindPendingByEmailAndTenant(ctx context.Context, emailLower, tenantID string) (*models.TenantInvite, error) {
	var invite models.TenantInvite
	err := r.db.WithContext(ctx).
		Where("email_lower = ? AND tenant_id = ? AND status = ?", emailLower, tenantID, constants.TenantInviteStatusPending).
		First(&invite).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.DatabaseError("Failed to find tenant invite", err).
			WithOperation("find_tenant_invite").
			WithResource("tenant_invite")
	}
	return &invite, nil
}

func (r *tenantInviteRepository) ListPendingByTenant(ctx context.Context, tenantID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantInvite], error) {
	base := r.db.WithContext(ctx).Model(&models.TenantInvite{}).
		Where("tenant_id = ? AND status = ?", tenantID, constants.TenantInviteStatusPending)
	return paginatePendingInvites(ctx, base, pr)
}

func (r *tenantInviteRepository) ListPending(ctx context.Context, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantInvite], error) {
	base := r.db.WithContext(ctx).Model(&models.TenantInvite{}).
		Where("status = ?", constants.TenantInviteStatusPending)
	return paginatePendingInvites(ctx, base, pr)
}

func paginatePendingInvites(ctx context.Context, base *gorm.DB, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantInvite], error) {
	return PaginateWithFind[models.TenantInvite](ctx, base, base.Preload("Tenant"), pr, PaginateOptions{
		OrderBy:             "created_at DESC",
		CountFailureMessage: "Failed to count tenant invites",
		FindFailureMessage:  "Failed to list tenant invites",
		DBOp:                DBOp{Operation: "list_tenant_invites", Resource: "tenant_invite"},
	})
}

// PendingInviteFresh reports whether a pending invite can still be accepted.
func PendingInviteFresh(invite *models.TenantInvite, now time.Time) bool {
	return invite != nil && invite.Status == constants.TenantInviteStatusPending && now.Before(invite.ExpiresAt)
}
