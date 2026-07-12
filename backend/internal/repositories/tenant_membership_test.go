package repositories

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/stretchr/testify/require"
)

func TestTenantMembershipRepository_CreateAndGetActive(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantMembershipRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")
	m := seedMembership(t, pg, user.ID, tenant.ID)

	got, err := repo.GetActive(ctx, user.ID, tenant.ID)
	require.NoError(t, err)
	require.Equal(t, m.ID, got.ID)
	require.NotNil(t, got.Tenant)

	exists, err := repo.ExistsActive(ctx, user.ID, tenant.ID)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestTenantMembershipRepository_GetActive_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantMembershipRepository(pg)

	_, err := repo.GetActive(testCtx(), "u", "t")
	requireNotFound(t, err)
}

func TestTenantMembershipRepository_ExistsActive_False(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantMembershipRepository(pg)

	exists, err := repo.ExistsActive(testCtx(), "u", "t")
	require.NoError(t, err)
	require.False(t, exists)
}

func TestTenantMembershipRepository_ListAndCount(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantMembershipRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")
	seedMembership(t, pg, user.ID, tenant.ID)

	byUser, err := repo.ListByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, byUser, 1)

	byTenant, err := repo.ListByTenantID(ctx, tenant.ID)
	require.NoError(t, err)
	require.Len(t, byTenant, 1)

	count, err := repo.CountByTenantID(ctx, tenant.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
}

func TestTenantMembershipRepository_ListPaginated(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantMembershipRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")
	seedMembership(t, pg, user.ID, tenant.ID)

	byUser, err := repo.ListByUserIDPaginated(ctx, user.ID, &dtos.PageableRequest{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, byUser.Data, 1)
	require.NotNil(t, byUser.Data[0].Tenant)

	byTenant, err := repo.ListByTenantIDPaginated(ctx, tenant.ID, &dtos.PageableRequest{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, byTenant.Data, 1)
	require.NotNil(t, byTenant.Data[0].User)
}

func TestTenantMembershipRepository_Delete(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantMembershipRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")
	seedMembership(t, pg, user.ID, tenant.ID)

	require.NoError(t, repo.Delete(ctx, user.ID, tenant.ID))
	_, err := repo.GetActive(ctx, user.ID, tenant.ID)
	requireNotFound(t, err)
}

func TestTenantMembershipRepository_Delete_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantMembershipRepository(pg)

	err := repo.Delete(testCtx(), "u", "t")
	requireNotFound(t, err)
}

func TestTenantMembershipRepository_Create_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideTenantMembershipRepository(pg)
	m := &models.TenantMembership{
		BaseModel: models.NewBaseModel(),
		UserID:    "u",
		TenantID:  "t",
		Role:      constants.TenantMembershipRoleMember,
		Status:    constants.TenantMembershipStatusActive,
	}

	err := repo.Create(testCtx(), m)
	requireDatabaseErr(t, err)
}
