package repositories

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/stretchr/testify/require"
)

func TestSessionRepository_CreateAndGetValid(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideSessionRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")
	exp := futureTime()

	s := &models.Session{
		BaseModel: models.NewBaseModel(),
		UserID:    user.ID,
		TenantID:  tenant.ID,
		ExpiresAt: &exp,
	}
	require.NoError(t, repo.Create(ctx, s))

	got, err := repo.GetValidByID(ctx, s.ID)
	require.NoError(t, err)
	require.Equal(t, user.ID, got.UserID)
}

func TestSessionRepository_GetValidByID_Expired(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideSessionRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")
	exp := pastTime()

	s := &models.Session{
		BaseModel: models.NewBaseModel(),
		UserID:    user.ID,
		TenantID:  tenant.ID,
		ExpiresAt: &exp,
	}
	require.NoError(t, repo.Create(ctx, s))

	_, err := repo.GetValidByID(ctx, s.ID)
	requireNotFound(t, err)
}

func TestSessionRepository_GetValidByID_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideSessionRepository(pg)

	_, err := repo.GetValidByID(testCtx(), "00000000-0000-7000-8000-000000000001")
	requireNotFound(t, err)
}

func TestSessionRepository_DeleteByID(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideSessionRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")
	s := &models.Session{BaseModel: models.NewBaseModel(), UserID: user.ID, TenantID: tenant.ID}
	require.NoError(t, repo.Create(ctx, s))

	require.NoError(t, repo.DeleteByID(ctx, s.ID))
	_, err := repo.GetValidByID(ctx, s.ID)
	requireNotFound(t, err)
}

func TestSessionRepository_DeleteAllByUserID(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideSessionRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")
	for range 2 {
		s := &models.Session{BaseModel: models.NewBaseModel(), UserID: user.ID, TenantID: tenant.ID}
		require.NoError(t, repo.Create(ctx, s))
	}

	require.NoError(t, repo.DeleteAllByUserID(ctx, user.ID))
	count, err := repo.CountActiveByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}

func TestSessionRepository_CountActive(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideSessionRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")

	expFuture := futureTime()
	expPast := pastTime()
	require.NoError(t, repo.Create(ctx, &models.Session{
		BaseModel: models.NewBaseModel(), UserID: user.ID, TenantID: tenant.ID, ExpiresAt: &expFuture,
	}))
	require.NoError(t, repo.Create(ctx, &models.Session{
		BaseModel: models.NewBaseModel(), UserID: user.ID, TenantID: tenant.ID, ExpiresAt: &expPast,
	}))

	total, err := repo.CountActive(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)

	byUser, err := repo.CountActiveByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), byUser)
}

func TestSessionRepository_Create_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideSessionRepository(pg)
	s := &models.Session{BaseModel: models.NewBaseModel(), UserID: "u", TenantID: "t"}

	err := repo.Create(testCtx(), s)
	requireDatabaseErr(t, err)
}
