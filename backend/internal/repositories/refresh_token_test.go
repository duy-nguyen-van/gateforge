package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

func TestRefreshTokenRepository_CreateAndFindValid(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")
	client := seedClient(t, pg, tenant.ID, "app")
	clientID := client.ID

	rt := &models.RefreshToken{
		HardDeleteModel: models.HardDeleteModel{ID: models.NewBaseModel().ID},
		TenantID:        tenant.ID,
		UserID:          user.ID,
		OAuthClientID:   client.ClientID,
		TokenHash:       "hash-valid",
		ExpiresAt:       futureTime(),
		ClientRecordID:  &clientID,
	}
	require.NoError(t, repo.Create(ctx, rt))

	got, err := repo.FindValidByTokenHash(ctx, "hash-valid")
	require.NoError(t, err)
	require.Equal(t, rt.ID, got.ID)
}

func TestRefreshTokenRepository_FindValidByTokenHash_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)

	_, err := repo.FindValidByTokenHash(testCtx(), "missing")
	requireNotFound(t, err)
}

func TestRefreshTokenRepository_RevokeByID(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")

	rt := &models.RefreshToken{
		HardDeleteModel: models.HardDeleteModel{ID: models.NewBaseModel().ID},
		TenantID:        tenant.ID,
		UserID:          user.ID,
		OAuthClientID:   "app",
		TokenHash:       "hash-revoke",
		ExpiresAt:       futureTime(),
	}
	require.NoError(t, repo.Create(ctx, rt))
	require.NoError(t, repo.RevokeByID(ctx, rt.ID))

	_, err := repo.FindValidByTokenHash(ctx, "hash-revoke")
	requireNotFound(t, err)
}

func TestRefreshTokenRepository_RevokeAllValidForUser(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")

	for _, hash := range []string{"h1", "h2"} {
		rt := &models.RefreshToken{
			HardDeleteModel: models.HardDeleteModel{ID: models.NewBaseModel().ID},
			TenantID:        tenant.ID,
			UserID:          user.ID,
			OAuthClientID:   "app",
			TokenHash:       hash,
			ExpiresAt:       futureTime(),
		}
		require.NoError(t, repo.Create(ctx, rt))
	}

	require.NoError(t, repo.RevokeAllValidForUser(ctx, user.ID))
	_, err := repo.FindValidByTokenHash(ctx, "h1")
	requireNotFound(t, err)
}

func TestRefreshTokenRepository_RevokeAndCreate(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")

	old := &models.RefreshToken{
		HardDeleteModel: models.HardDeleteModel{ID: models.NewBaseModel().ID},
		TenantID:        tenant.ID,
		UserID:          user.ID,
		OAuthClientID:   "app",
		TokenHash:       "old-hash",
		ExpiresAt:       futureTime(),
	}
	require.NoError(t, repo.Create(ctx, old))

	newRT := &models.RefreshToken{
		HardDeleteModel: models.HardDeleteModel{ID: models.NewBaseModel().ID},
		TenantID:        tenant.ID,
		UserID:          user.ID,
		OAuthClientID:   "app",
		TokenHash:       "new-hash",
		ExpiresAt:       futureTime(),
	}
	require.NoError(t, repo.RevokeAndCreate(ctx, old.ID, newRT))

	_, err := repo.FindValidByTokenHash(ctx, "old-hash")
	requireNotFound(t, err)
	got, err := repo.FindValidByTokenHash(ctx, "new-hash")
	require.NoError(t, err)
	require.Equal(t, newRT.ID, got.ID)
}

func TestRefreshTokenRepository_UsageByClientRecordID(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")
	client := seedClient(t, pg, tenant.ID, "app")
	clientID := client.ID

	rt := &models.RefreshToken{
		HardDeleteModel: models.HardDeleteModel{ID: models.NewBaseModel().ID},
		TenantID:        tenant.ID,
		UserID:          user.ID,
		OAuthClientID:   client.ClientID,
		TokenHash:       "usage-hash",
		ExpiresAt:       futureTime(),
		ClientRecordID:  &clientID,
	}
	require.NoError(t, repo.Create(ctx, rt))

	usage, err := repo.UsageByClientRecordID(ctx, client.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), usage.TotalIssued)
	require.Equal(t, int64(1), usage.ActiveCount)
	require.NotNil(t, usage.LastIssuedAt)
}

func TestRefreshTokenRepository_UsageByClientRecordID_Empty(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)

	usage, err := repo.UsageByClientRecordID(testCtx(), "00000000-0000-7000-8000-000000000001")
	require.NoError(t, err)
	require.Equal(t, int64(0), usage.TotalIssued)
	require.Nil(t, usage.LastIssuedAt)
}

func TestRefreshTokenRepository_Create_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)
	rt := &models.RefreshToken{
		HardDeleteModel: models.HardDeleteModel{ID: models.NewBaseModel().ID},
		TenantID:        "t",
		UserID:          "u",
		OAuthClientID:   "app",
		TokenHash:       "x",
		ExpiresAt:       futureTime(),
	}

	err := repo.Create(testCtx(), rt)
	requireDatabaseErr(t, err)
}
