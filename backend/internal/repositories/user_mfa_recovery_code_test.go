package repositories

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/stretchr/testify/require"
)

func TestUserMFARecoveryCodeRepository_ReplaceFindMarkUsed(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserMFARecoveryCodeRepository(pg)
	ctx := testCtx()
	user := seedUser(t, pg, "recovery@test.com")

	rows := []*models.UserMFARecoveryCode{
		{BaseModel: models.NewBaseModel(), UserID: user.ID, CodeHash: "hash1"},
		{BaseModel: models.NewBaseModel(), UserID: user.ID, CodeHash: "hash2"},
	}
	require.NoError(t, repo.ReplaceAllForUser(ctx, user.ID, rows))

	unused, err := repo.FindUnusedByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, unused, 2)

	require.NoError(t, repo.MarkUsed(ctx, unused[0].ID))
	unused, err = repo.FindUnusedByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, unused, 1)

	replacement := []*models.UserMFARecoveryCode{
		{BaseModel: models.NewBaseModel(), UserID: user.ID, CodeHash: "new-hash"},
	}
	require.NoError(t, repo.ReplaceAllForUser(ctx, user.ID, replacement))
	unused, err = repo.FindUnusedByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, unused, 1)
	require.Equal(t, "new-hash", unused[0].CodeHash)
}

func TestUserMFARecoveryCodeRepository_FindUnusedByUserID_Empty(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserMFARecoveryCodeRepository(pg)
	user := seedUser(t, pg, "empty@test.com")

	rows, err := repo.FindUnusedByUserID(testCtx(), user.ID)
	require.NoError(t, err)
	require.Empty(t, rows)
}

func TestUserMFARecoveryCodeRepository_ReplaceAllForUser_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideUserMFARecoveryCodeRepository(pg)
	rows := []*models.UserMFARecoveryCode{
		{BaseModel: models.NewBaseModel(), UserID: "u", CodeHash: "h"},
	}

	err := repo.ReplaceAllForUser(testCtx(), "u", rows)
	requireDatabaseOrClosedErr(t, err)
}

func TestUserMFARecoveryCodeRepository_FindUnusedByUserID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideUserMFARecoveryCodeRepository(pg)

	_, err := repo.FindUnusedByUserID(testCtx(), "u")
	requireDatabaseErr(t, err)
}
