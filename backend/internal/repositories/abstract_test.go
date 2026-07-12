package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

func TestNormalizePageableDefaults(t *testing.T) {
	pr := NormalizePageable(nil, 0, 0)
	if pr.Page != 1 {
		t.Fatalf("page = %d, want 1", pr.Page)
	}
	if pr.PageSize != constants.DefaultPageSize {
		t.Fatalf("page_size = %d, want %d", pr.PageSize, constants.DefaultPageSize)
	}
}

func TestNormalizePageableZeroPageSize(t *testing.T) {
	pr := NormalizePageable(&dtos.PageableRequest{Page: 1, PageSize: 0}, 0, 0)
	if pr.PageSize != constants.DefaultPageSize {
		t.Fatalf("page_size = %d, want %d", pr.PageSize, constants.DefaultPageSize)
	}
}

func TestNormalizePageableClamp(t *testing.T) {
	pr := NormalizePageable(&dtos.PageableRequest{Page: 0, PageSize: 500}, 20, constants.MaxPageSize)
	if pr.PageSize != constants.MaxPageSize {
		t.Fatalf("page_size = %d, want %d", pr.PageSize, constants.MaxPageSize)
	}
}

func TestNormalizePageablePreservesValidInput(t *testing.T) {
	pr := NormalizePageable(&dtos.PageableRequest{Page: 3, PageSize: 25}, 20, constants.MaxPageSize)
	if pr.Page != 3 || pr.PageSize != 25 {
		t.Fatalf("got page=%d page_size=%d, want page=3 page_size=25", pr.Page, pr.PageSize)
	}
}

func TestCreate_Success(t *testing.T) {
	pg := newTestDB(t)
	ctx := testCtx()
	tenant := &models.Tenant{BaseModel: models.NewBaseModel(), Name: "Acme", Domain: "acme.test"}

	err := Create(ctx, pg, tenant, DBOp{Operation: "create_tenant", Resource: "tenant"}, "Failed to create tenant")
	require.NoError(t, err)
	require.NotEmpty(t, tenant.ID)

	var count int64
	require.NoError(t, pg.WithContext(ctx).Model(&models.Tenant{}).Count(&count).Error)
	require.Equal(t, int64(1), count)
}

func TestCreate_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	ctx := testCtx()
	tenant := &models.Tenant{BaseModel: models.NewBaseModel(), Name: "Acme", Domain: "acme.test"}

	err := Create(ctx, pg, tenant, DBOp{Operation: "create_tenant", Resource: "tenant"}, "Failed to create tenant")
	requireDatabaseErr(t, err)
	appErr := errors.GetAppError(err)
	require.Equal(t, "create_tenant", appErr.Operation)
	require.Equal(t, "tenant", appErr.Resource)
}

func TestPaginate_Success(t *testing.T) {
	pg := newTestDB(t)
	ctx := testCtx()
	seedTenant(t, pg, "A", "a.test")
	seedTenant(t, pg, "B", "b.test")

	query := pg.WithContext(ctx).Model(&models.Tenant{})
	resp, err := Paginate[models.Tenant](ctx, query, &dtos.PageableRequest{Page: 1, PageSize: 10}, PaginateOptions{
		OrderBy:             "created_at DESC",
		CountFailureMessage: "Failed to count tenants",
		FindFailureMessage:  "Failed to list tenants",
		DBOp:                DBOp{Operation: "list_tenants", Resource: "tenant"},
	})
	require.NoError(t, err)
	require.Len(t, resp.Data, 2)
	require.Equal(t, int64(2), resp.Pageable.Total)
}

func TestPaginate_CountError(t *testing.T) {
	pg := closedTestDB(t)
	ctx := testCtx()
	query := pg.WithContext(ctx).Model(&models.Tenant{})

	_, err := Paginate[models.Tenant](ctx, query, nil, PaginateOptions{
		CountFailureMessage: "Failed to count tenants",
		FindFailureMessage:  "Failed to list tenants",
		DBOp:                DBOp{Operation: "list_tenants", Resource: "tenant"},
	})
	requireDatabaseErr(t, err)
}

func TestPaginateWithFind_FindError(t *testing.T) {
	openPG := newTestDB(t)
	closedPG := closedTestDB(t)
	ctx := testCtx()
	seedTenant(t, openPG, "A", "a.test")

	countQuery := openPG.WithContext(ctx).Model(&models.Tenant{})
	findQuery := closedPG.WithContext(ctx).Model(&models.Tenant{})

	_, err := PaginateWithFind[models.Tenant](ctx, countQuery, findQuery, nil, PaginateOptions{
		OrderBy:             "created_at DESC",
		CountFailureMessage: "Failed to count tenants",
		FindFailureMessage:  "Failed to list tenants",
		DBOp:                DBOp{Operation: "list_tenants", Resource: "tenant"},
	})
	requireDatabaseErr(t, err)
}

func TestCountRows_Success(t *testing.T) {
	pg := newTestDB(t)
	ctx := testCtx()
	seedUser(t, pg, "one@test.com")
	seedUser(t, pg, "two@test.com")

	query := pg.WithContext(ctx).Model(&models.User{})
	total, err := CountRows(ctx, query, DBOp{Operation: "count_users", Resource: "user"}, "Failed to count users")
	require.NoError(t, err)
	require.Equal(t, int64(2), total)
}

func TestCountRows_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	ctx := testCtx()
	query := pg.WithContext(ctx).Model(&models.User{})

	_, err := CountRows(ctx, query, DBOp{Operation: "count_users", Resource: "user"}, "Failed to count users")
	requireDatabaseErr(t, err)
}
