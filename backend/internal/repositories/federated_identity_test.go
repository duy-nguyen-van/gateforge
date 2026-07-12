package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

func TestFederatedIdentityRepository_CreateAndGet(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideFederatedIdentityRepository(pg)
	ctx := testCtx()
	user := seedUserWithPassword(t, pg, "fed@test.com", "hash")

	row := &models.FederatedIdentity{
		BaseModel:   models.NewBaseModel(),
		UserID:      user.ID,
		Provider:    "google",
		Subject:     "sub-123",
		EmailAtLink: "fed@test.com",
	}
	require.NoError(t, repo.Create(ctx, row))

	got, err := repo.GetByProviderSubject(ctx, "google", "sub-123")
	require.NoError(t, err)
	require.Equal(t, user.ID, got.UserID)
	require.NotNil(t, got.User)
	require.NotNil(t, got.User.PasswordCredential)
}

func TestFederatedIdentityRepository_GetByProviderSubject_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideFederatedIdentityRepository(pg)

	_, err := repo.GetByProviderSubject(testCtx(), "google", "missing")
	requireNotFound(t, err)
}

func TestFederatedIdentityRepository_ListByUserID(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideFederatedIdentityRepository(pg)
	ctx := testCtx()
	user := seedUser(t, pg, "u@test.com")

	for _, sub := range []string{"s1", "s2"} {
		require.NoError(t, repo.Create(ctx, &models.FederatedIdentity{
			BaseModel: models.NewBaseModel(),
			UserID:    user.ID,
			Provider:  "google",
			Subject:   sub,
		}))
	}

	rows, err := repo.ListByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, rows, 2)
}

func TestFederatedIdentityRepository_Create_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideFederatedIdentityRepository(pg)
	row := &models.FederatedIdentity{
		BaseModel: models.NewBaseModel(),
		UserID:    "u",
		Provider:  "google",
		Subject:   "s",
	}

	err := repo.Create(testCtx(), row)
	requireDatabaseErr(t, err)
}

func TestFederatedIdentityRepository_ListByUserID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideFederatedIdentityRepository(pg)

	_, err := repo.ListByUserID(testCtx(), "u")
	requireDatabaseErr(t, err)
}
