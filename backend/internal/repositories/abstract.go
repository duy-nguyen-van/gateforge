package repositories

import (
	"context"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DBOp identifies a repository operation for AppError metadata.
type DBOp struct {
	Operation string
	Resource  string
}

// PaginateOptions configures shared list/count helpers.
type PaginateOptions struct {
	OrderBy             string
	DefaultPageSize     int
	MaxPageSize         int // 0 = do not clamp
	CountFailureMessage string
	FindFailureMessage  string
	DBOp
}

// NormalizePageable applies list defaults (page 1, configurable page size).
func NormalizePageable(pr *dtos.PageableRequest, defaultPageSize, maxPageSize int) *dtos.PageableRequest {
	if pr == nil {
		pr = dtos.NewPageableRequest()
	}
	if pr.Page <= 0 {
		pr.Page = 1
	}
	if defaultPageSize <= 0 {
		defaultPageSize = constants.DefaultPageSize
	}
	if pr.PageSize <= 0 {
		pr.PageSize = defaultPageSize
	}
	if maxPageSize > 0 {
		dtos.ClampPageSize(pr, maxPageSize)
	}
	return pr
}

// Create inserts one row with standard database error wrapping.
// Primary keys are assigned by model BeforeCreate hooks (BaseModel / HardDeleteModel).
func Create[T any](ctx context.Context, db *db.PostgresDB, entity *T, op DBOp, message string) error {
	if err := db.WithContext(ctx).Create(entity).Error; err != nil {
		return errors.DatabaseError(message, err).
			WithOperation(op.Operation).
			WithResource(op.Resource)
	}
	return nil
}

// Paginate counts and fetches rows from a pre-scoped GORM query.
func Paginate[T any](ctx context.Context, query *gorm.DB, pr *dtos.PageableRequest, opts PaginateOptions) (*dtos.DataResponse[T], error) {
	return PaginateWithFind[T](ctx, query, query, pr, opts)
}

// PaginateWithFind counts rows from countQuery and fetches from findQuery (e.g. when preloads differ).
func PaginateWithFind[T any](ctx context.Context, countQuery, findQuery *gorm.DB, pr *dtos.PageableRequest, opts PaginateOptions) (*dtos.DataResponse[T], error) {
	pr = NormalizePageable(pr, opts.DefaultPageSize, opts.MaxPageSize)

	var total int64
	if err := countQuery.Session(&gorm.Session{}).
		Clauses(clause.OrderBy{}).
		Limit(-1).Offset(-1).
		Count(&total).Error; err != nil {
		return nil, errors.DatabaseError(opts.CountFailureMessage, err).
			WithOperation(opts.Operation).
			WithResource(opts.Resource)
	}

	findQuery = findQuery.Session(&gorm.Session{})
	if opts.OrderBy != "" {
		findQuery = findQuery.Order(opts.OrderBy)
	}

	var rows []T
	if err := findQuery.
		Limit(pr.GetLimit()).
		Offset(pr.GetOffset()).
		Find(&rows).Error; err != nil {
		return nil, errors.DatabaseError(opts.FindFailureMessage, err).
			WithOperation(opts.Operation).
			WithResource(opts.Resource)
	}

	return &dtos.DataResponse[T]{
		Data: rows,
		Pageable: &dtos.Pageable{
			Page:     pr.Page,
			PageSize: pr.PageSize,
			Total:    total,
		},
	}, nil
}

// CountRows returns the row count for a pre-scoped GORM query.
func CountRows(ctx context.Context, query *gorm.DB, op DBOp, message string) (int64, error) {
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return 0, errors.DatabaseError(message, err).
			WithOperation(op.Operation).
			WithResource(op.Resource)
	}
	return total, nil
}
