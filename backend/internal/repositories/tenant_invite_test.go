package repositories

import (
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func TestTenantInviteRepository_CreateGetAndSave(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantInviteRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Acme", "acme.test")
	user := seedUser(t, pg, "accepted@test.com")

	invite := &models.TenantInvite{
		BaseModel:  models.NewBaseModel(),
		Email:      "New@Example.com",
		EmailLower: "new@example.com",
		TenantID:   tenant.ID,
		Role:       "member",
		TokenHash:  "hash-1",
		Status:     constants.TenantInviteStatusPending,
		ExpiresAt:  time.Now().UTC().Add(time.Hour),
	}
	require.NoError(t, repo.Create(ctx, invite))

	got, err := repo.GetByTokenHash(ctx, "hash-1")
	require.NoError(t, err)
	require.Equal(t, "new@example.com", got.EmailLower)
	require.NotNil(t, got.Tenant)
	require.Equal(t, "Acme", got.Tenant.Name)

	pending, err := repo.FindPendingByEmailAndTenant(ctx, "new@example.com", tenant.ID)
	require.NoError(t, err)
	require.Equal(t, invite.ID, pending.ID)

	invite.Status = constants.TenantInviteStatusAccepted
	invite.Role = "admin"
	invite.AcceptedUserID = &user.ID
	require.NoError(t, repo.Save(ctx, invite))

	saved, err := repo.GetByTokenHash(ctx, "hash-1")
	require.NoError(t, err)
	require.Equal(t, constants.TenantInviteStatusAccepted, saved.Status)
	require.Equal(t, "admin", saved.Role)

	missing, err := repo.FindPendingByEmailAndTenant(ctx, "new@example.com", tenant.ID)
	require.NoError(t, err)
	require.Nil(t, missing)
}

func TestTenantInviteRepository_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantInviteRepository(pg)

	_, err := repo.GetByTokenHash(testCtx(), "missing")
	requireNotFound(t, err)

	pending, err := repo.FindPendingByEmailAndTenant(testCtx(), "nobody@test.com", "00000000-0000-7000-8000-000000000099")
	require.NoError(t, err)
	require.Nil(t, pending)
}

func TestTenantInviteRepository_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideTenantInviteRepository(pg)
	ctx := testCtx()
	invite := &models.TenantInvite{BaseModel: models.NewBaseModel(), Email: "a@test.com", EmailLower: "a@test.com", Role: "member", TokenHash: "h", Status: constants.TenantInviteStatusPending, ExpiresAt: time.Now().UTC()}

	requireDatabaseErr(t, repo.Create(ctx, invite))
	requireDatabaseErr(t, repo.Save(ctx, invite))
	_, err := repo.GetByTokenHash(ctx, "h")
	requireDatabaseErr(t, err)
	_, err = repo.FindPendingByEmailAndTenant(ctx, "a@test.com", invite.ID)
	requireDatabaseErr(t, err)
}

func TestPendingInviteFresh(t *testing.T) {
	now := time.Now().UTC()
	require.False(t, PendingInviteFresh(nil, now))
	require.False(t, PendingInviteFresh(&models.TenantInvite{Status: constants.TenantInviteStatusAccepted, ExpiresAt: now.Add(time.Hour)}, now))
	require.False(t, PendingInviteFresh(&models.TenantInvite{Status: constants.TenantInviteStatusPending, ExpiresAt: now.Add(-time.Second)}, now))
	require.True(t, PendingInviteFresh(&models.TenantInvite{Status: constants.TenantInviteStatusPending, ExpiresAt: now.Add(time.Minute)}, now))
}
