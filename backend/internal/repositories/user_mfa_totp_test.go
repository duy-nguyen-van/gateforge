package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

func TestUserMFATOTPRepository_UpsertAndActivate(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserMFATOTPRepository(pg)
	ctx := testCtx()
	user := seedUser(t, pg, "mfa@test.com")

	row := &models.UserMFATOTP{
		BaseModel:       models.NewBaseModel(),
		UserID:          user.ID,
		SecretEncrypted: "enc",
		Enabled:         false,
	}
	require.NoError(t, repo.UpsertPending(ctx, row))

	got, err := repo.GetByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.False(t, got.Enabled)

	active, err := repo.GetActiveByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Nil(t, active)

	require.NoError(t, repo.MarkVerifiedAndEnabled(ctx, user.ID))
	active, err = repo.GetActiveByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, active)
	require.True(t, active.Enabled)
}

func TestUserMFATOTPRepository_UpsertReplacesExisting(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserMFATOTPRepository(pg)
	ctx := testCtx()
	user := seedUser(t, pg, "mfa2@test.com")

	require.NoError(t, repo.UpsertPending(ctx, &models.UserMFATOTP{
		BaseModel: models.NewBaseModel(), UserID: user.ID, SecretEncrypted: "old",
	}))
	require.NoError(t, repo.UpsertPending(ctx, &models.UserMFATOTP{
		BaseModel: models.NewBaseModel(), UserID: user.ID, SecretEncrypted: "new",
	}))

	got, err := repo.GetByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "new", got.SecretEncrypted)
}

func TestUserMFATOTPRepository_Disable(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserMFATOTPRepository(pg)
	ctx := testCtx()
	user := seedUser(t, pg, "mfa3@test.com")

	require.NoError(t, repo.UpsertPending(ctx, &models.UserMFATOTP{
		BaseModel: models.NewBaseModel(), UserID: user.ID, SecretEncrypted: "enc",
	}))
	require.NoError(t, repo.Disable(ctx, user.ID))

	_, err := repo.GetByUserID(ctx, user.ID)
	requireNotFound(t, err)
}

func TestUserMFATOTPRepository_GetByUserID_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserMFATOTPRepository(pg)

	_, err := repo.GetByUserID(testCtx(), "missing")
	requireNotFound(t, err)
}

func TestUserMFATOTPRepository_CountEnabled(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserMFATOTPRepository(pg)
	ctx := testCtx()
	user := seedUser(t, pg, "mfa4@test.com")

	require.NoError(t, repo.UpsertPending(ctx, &models.UserMFATOTP{
		BaseModel: models.NewBaseModel(), UserID: user.ID, SecretEncrypted: "enc",
	}))
	require.NoError(t, repo.MarkVerifiedAndEnabled(ctx, user.ID))

	count, err := repo.CountEnabled(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
}

func TestUserMFATOTPRepository_UpsertPending_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideUserMFATOTPRepository(pg)
	row := &models.UserMFATOTP{
		BaseModel: models.NewBaseModel(), UserID: "u", SecretEncrypted: "enc",
	}

	err := repo.UpsertPending(testCtx(), row)
	requireDatabaseOrClosedErr(t, err)
}
