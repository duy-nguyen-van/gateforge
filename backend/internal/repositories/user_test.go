package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

func TestUserRepository_CreateWithPasswordHash(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)
	ctx := testCtx()

	user := &models.User{
		BaseModel: models.NewBaseModel(),
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane@test.com",
		Status:    constants.UserStatusActive,
	}
	require.NoError(t, repo.CreateWithPasswordHash(ctx, user, "hash"))

	got, err := repo.GetOneByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "jane@test.com", got.Email)

	var pc models.PasswordCredential
	require.NoError(t, pg.WithContext(ctx).Where("user_id = ?", user.ID).First(&pc).Error)
	require.Equal(t, "hash", pc.PasswordHash)
}

func TestUserRepository_CreateWithPasswordHash_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideUserRepository(pg)
	user := &models.User{BaseModel: models.NewBaseModel(), Email: "x@test.com", Status: constants.UserStatusActive}

	err := repo.CreateWithPasswordHash(testCtx(), user, "hash")
	requireDatabaseOrClosedErr(t, err)
}

func TestUserRepository_CreateUserOnly(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)
	ctx := testCtx()

	user := &models.User{
		BaseModel: models.NewBaseModel(),
		Email:     "new@test.com",
		Status:    constants.UserStatusActive,
	}
	require.NoError(t, repo.CreateUserOnly(ctx, user))

	got, err := repo.GetOneByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "new@test.com", got.Email)
}

func TestUserRepository_GetOneByID_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)

	_, err := repo.GetOneByID(testCtx(), "00000000-0000-7000-8000-000000000001")
	requireNotFound(t, err)
}

func TestUserRepository_GetOneByID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideUserRepository(pg)

	_, err := repo.GetOneByID(testCtx(), "00000000-0000-7000-8000-000000000001")
	requireDatabaseErr(t, err)
}

func TestUserRepository_GetByEmailLower(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)
	ctx := testCtx()
	user := seedUserWithPassword(t, pg, "FindMe@Test.com", "hash")

	got, err := repo.GetByEmailLower(ctx, user.EmailLower)
	require.NoError(t, err)
	require.Equal(t, user.ID, got.ID)
	require.NotNil(t, got.PasswordCredential)
}

func TestUserRepository_GetByEmailLower_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)

	_, err := repo.GetByEmailLower(testCtx(), "missing@test.com")
	requireNotFound(t, err)
}

func TestUserRepository_CountAndPlatformAdmins(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)
	ctx := testCtx()

	u1 := seedUser(t, pg, "a@test.com")
	_ = seedUser(t, pg, "b@test.com")
	require.NoError(t, repo.SetPlatformAdmin(ctx, u1.ID, true))

	count, err := repo.Count(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(2), count)

	adminCount, err := repo.CountPlatformAdmins(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(1), adminCount)
}

func TestUserRepository_UpdateStatus(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)
	ctx := testCtx()
	user := seedUser(t, pg, "status@test.com")

	require.NoError(t, repo.UpdateStatus(ctx, user.ID, constants.UserStatusDisabled))
	got, err := repo.GetOneByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, constants.UserStatusDisabled, got.Status)
}

func TestUserRepository_UpdateStatus_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)

	err := repo.UpdateStatus(testCtx(), "00000000-0000-7000-8000-000000000099", constants.UserStatusDisabled)
	requireNotFound(t, err)
}

func TestUserRepository_UpdateProfile(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)
	ctx := testCtx()
	user := seedUser(t, pg, "profile@test.com")

	first := "Updated"
	got, err := repo.UpdateProfile(ctx, user.ID, UserProfilePatch{FirstName: &first})
	require.NoError(t, err)
	require.Equal(t, "Updated", got.FirstName)
}

func TestUserRepository_UpdateProfile_EmptyPatch(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)
	ctx := testCtx()
	user := seedUser(t, pg, "noop@test.com")

	got, err := repo.UpdateProfile(ctx, user.ID, UserProfilePatch{})
	require.NoError(t, err)
	require.Equal(t, user.ID, got.ID)
}

func TestUserRepository_UpdateProfile_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)
	first := "X"

	_, err := repo.UpdateProfile(testCtx(), "00000000-0000-7000-8000-000000000099", UserProfilePatch{FirstName: &first})
	requireNotFound(t, err)
}

func TestUserRepository_SetPlatformAdmin_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)

	err := repo.SetPlatformAdmin(testCtx(), "00000000-0000-7000-8000-000000000099", true)
	requireNotFound(t, err)
}

func TestUserRepository_List(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideUserRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	u1 := seedUser(t, pg, "u1@test.com")
	_ = seedUser(t, pg, "u2@test.com")
	seedMembership(t, pg, u1.ID, tenant.ID)

	resp, err := repo.List(ctx, tenant.ID, "", &dtos.PageableRequest{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, u1.ID, resp.Data[0].ID)

	respAll, err := repo.List(ctx, "", "", &dtos.PageableRequest{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, respAll.Data, 2)
}

func TestUserRepository_List_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideUserRepository(pg)

	_, err := repo.List(testCtx(), "", "", nil)
	requireDatabaseErr(t, err)
	appErr := errors.GetAppError(err)
	require.Equal(t, "list_users", appErr.Operation)
}
