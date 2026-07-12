package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

func TestTenantRepository_CreateAndGet(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)
	ctx := testCtx()

	tenant := &models.Tenant{BaseModel: models.NewBaseModel(), Name: "Acme", Domain: "acme.test"}
	require.NoError(t, repo.Create(ctx, tenant))

	got, err := repo.GetByID(ctx, tenant.ID)
	require.NoError(t, err)
	require.Equal(t, "Acme", got.Name)

	byDomain, err := repo.GetByDomain(ctx, "acme.test")
	require.NoError(t, err)
	require.Equal(t, tenant.ID, byDomain.ID)
}

func TestTenantRepository_GetByID_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)

	_, err := repo.GetByID(testCtx(), "00000000-0000-7000-8000-000000000001")
	requireNotFound(t, err)
}

func TestTenantRepository_GetByDomain_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)

	_, err := repo.GetByDomain(testCtx(), "missing.test")
	requireNotFound(t, err)
}

func TestTenantRepository_List(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)
	ctx := testCtx()
	seedTenant(t, pg, "A", "a.test")
	seedTenant(t, pg, "B", "b.test")

	resp, err := repo.List(ctx, &dtos.PageableRequest{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, resp.Data, 2)
}

func TestTenantRepository_Update(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Old", "old.test")

	name := "New"
	got, err := repo.Update(ctx, tenant.ID, TenantPatch{Name: &name})
	require.NoError(t, err)
	require.Equal(t, "New", got.Name)
}

func TestTenantRepository_Update_EmptyPatch(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Same", "same.test")

	got, err := repo.Update(ctx, tenant.ID, TenantPatch{})
	require.NoError(t, err)
	require.Equal(t, tenant.Name, got.Name)
}

func TestTenantRepository_Update_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)
	name := "X"

	_, err := repo.Update(testCtx(), "00000000-0000-7000-8000-000000000099", TenantPatch{Name: &name})
	requireNotFound(t, err)
}

func TestTenantRepository_Delete(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Del", "del.test")

	require.NoError(t, repo.Delete(ctx, tenant.ID))
	_, err := repo.GetByID(ctx, tenant.ID)
	requireNotFound(t, err)
}

func TestTenantRepository_Delete_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)

	err := repo.Delete(testCtx(), "00000000-0000-7000-8000-000000000099")
	requireNotFound(t, err)
}

func TestTenantRepository_DomainTaken(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Acme", "acme.test")

	taken, err := repo.DomainTaken(ctx, "acme.test", "")
	require.NoError(t, err)
	require.True(t, taken)

	free, err := repo.DomainTaken(ctx, "acme.test", tenant.ID)
	require.NoError(t, err)
	require.False(t, free)

	empty, err := repo.DomainTaken(ctx, "  ", "")
	require.NoError(t, err)
	require.False(t, empty)
}

func TestTenantRepository_CountUsersByTenantID(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "m@test.com")
	seedMembership(t, pg, user.ID, tenant.ID)

	count, err := repo.CountUsersByTenantID(ctx, tenant.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
}

func TestTenantRepository_Create_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideTenantRepository(pg)
	tenant := &models.Tenant{BaseModel: models.NewBaseModel(), Name: "X", Domain: "x.test"}

	err := repo.Create(testCtx(), tenant)
	requireDatabaseErr(t, err)
}
