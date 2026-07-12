package repositories

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/stretchr/testify/require"
)

func TestAuthorizationCodeRepository_CreateTakeDelete(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideAuthorizationCodeRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")

	row := &models.AuthorizationCode{
		BaseModel:     models.NewBaseModel(),
		Code:          "auth-code-123",
		TenantID:      tenant.ID,
		OAuthClientID: "app",
		UserID:        user.ID,
		ExpiresAt:     futureTime(),
	}
	require.NoError(t, repo.Create(ctx, row))

	got, err := repo.TakeByCode(ctx, "auth-code-123")
	require.NoError(t, err)
	require.Equal(t, user.ID, got.UserID)

	require.NoError(t, repo.DeleteByCode(ctx, "auth-code-123"))
	_, err = repo.TakeByCode(ctx, "auth-code-123")
	requireNotFound(t, err)
}

func TestAuthorizationCodeRepository_TakeByCode_Expired(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideAuthorizationCodeRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	user := seedUser(t, pg, "u@test.com")

	row := &models.AuthorizationCode{
		BaseModel:     models.NewBaseModel(),
		Code:          "expired-code",
		TenantID:      tenant.ID,
		OAuthClientID: "app",
		UserID:        user.ID,
		ExpiresAt:     pastTime(),
	}
	require.NoError(t, repo.Create(ctx, row))

	_, err := repo.TakeByCode(ctx, "expired-code")
	requireNotFound(t, err)
}

func TestAuthorizationCodeRepository_TakeByCode_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideAuthorizationCodeRepository(pg)

	_, err := repo.TakeByCode(testCtx(), "missing")
	requireNotFound(t, err)
}

func TestAuthorizationCodeRepository_Create_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideAuthorizationCodeRepository(pg)
	row := &models.AuthorizationCode{
		BaseModel:     models.NewBaseModel(),
		Code:          "x",
		TenantID:      "t",
		OAuthClientID: "app",
		UserID:        "u",
		ExpiresAt:     futureTime(),
	}

	err := repo.Create(testCtx(), row)
	requireDatabaseErr(t, err)
}
