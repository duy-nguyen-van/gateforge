package repositories

import (
	"context"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// AuditLogListFilters narrows admin audit log queries.
type AuditLogListFilters struct {
	TenantID     string
	Action       string
	ActionsIn    []string
	Result       string
	ActorID      string
	ResourceType string
	ResourceName string
	From         *time.Time
	To           *time.Time
}

// AuditLogRepository persists append-only audit log rows.
type AuditLogRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	List(ctx context.Context, filters AuditLogListFilters, pr *dtos.PageableRequest) (*dtos.DataResponse[models.AuditLog], error)
	Count(ctx context.Context, filters AuditLogListFilters) (int64, error)
}

type auditLogRepository struct {
	db *db.PostgresDB
}

// ProvideAuditLogRepository wires audit log persistence.
func ProvideAuditLogRepository(db *db.PostgresDB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
	return Create(ctx, r.db, log, DBOp{Operation: "create_audit_log", Resource: "audit_log"}, "Failed to create audit log")
}

func (r *auditLogRepository) List(ctx context.Context, filters AuditLogListFilters, pr *dtos.PageableRequest) (*dtos.DataResponse[models.AuditLog], error) {
	query := r.applyAuditLogFilters(r.db.WithContext(ctx).Model(&models.AuditLog{}), filters)
	return Paginate[models.AuditLog](ctx, query, pr, PaginateOptions{
		OrderBy:             "created_at DESC",
		CountFailureMessage: "Failed to count audit logs",
		FindFailureMessage:  "Failed to list audit logs",
		DBOp:                DBOp{Operation: "list_audit_logs", Resource: "audit_log"},
	})
}

func (r *auditLogRepository) Count(ctx context.Context, filters AuditLogListFilters) (int64, error) {
	query := r.applyAuditLogFilters(r.db.WithContext(ctx).Model(&models.AuditLog{}), filters)
	return CountRows(ctx, query, DBOp{Operation: "count_audit_logs", Resource: "audit_log"}, "Failed to count audit logs")
}

func (r *auditLogRepository) applyAuditLogFilters(query *gorm.DB, filters AuditLogListFilters) *gorm.DB {
	if filters.TenantID != "" {
		query = query.Where("tenant_id = ?", filters.TenantID)
	}
	if filters.Result != "" {
		query = query.Where("result = ?", filters.Result)
	}
	if filters.ActorID != "" {
		query = query.Where("actor_id = ?", filters.ActorID)
	}
	if filters.ResourceType != "" {
		query = query.Where("resource_type = ?", filters.ResourceType)
	}
	if filters.ResourceName != "" {
		query = query.Where("resource_name = ?", filters.ResourceName)
	}
	if len(filters.ActionsIn) > 0 {
		query = query.Where("action IN ?", filters.ActionsIn)
	} else if action := strings.TrimSpace(filters.Action); action != "" {
		if strings.HasSuffix(action, ".") {
			query = query.Where("action LIKE ?", action+"%")
		} else {
			query = query.Where("action = ?", action)
		}
	}
	if filters.From != nil {
		query = query.Where("created_at >= ?", *filters.From)
	}
	if filters.To != nil {
		query = query.Where("created_at <= ?", *filters.To)
	}
	return query
}
