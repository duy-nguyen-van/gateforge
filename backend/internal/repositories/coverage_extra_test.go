package repositories

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

func TestClientRepository_GetByTenantAndClientID_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideClientRepository(pg)

	_, err := repo.GetByTenantAndClientID(testCtx(), "missing-tenant", "missing-client")
	requireNotFound(t, err)
}

func TestClientRepository_GetByClientID_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideClientRepository(pg)

	_, err := repo.GetByClientID(testCtx(), "missing-client")
	requireNotFound(t, err)
}

func TestClientRepository_GetByTenantAndClientID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideClientRepository(pg)

	_, err := repo.GetByTenantAndClientID(testCtx(), "t", "c")
	requireDatabaseErr(t, err)
}

func TestClientRepository_GetByClientID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideClientRepository(pg)

	_, err := repo.GetByClientID(testCtx(), "c")
	requireDatabaseErr(t, err)
}

func TestClientRepository_Update_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideClientRepository(pg)
	name := "X"

	_, err := repo.Update(testCtx(), "00000000-0000-7000-8000-000000000001", ClientPatch{Name: &name})
	requireDatabaseErr(t, err)
}

func TestClientRepository_Delete_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideClientRepository(pg)

	err := repo.Delete(testCtx(), "00000000-0000-7000-8000-000000000001")
	requireDatabaseErr(t, err)
}

func TestClientRepository_ClientIDTaken_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideClientRepository(pg)

	_, err := repo.ClientIDTaken(testCtx(), "t", "c", "")
	requireDatabaseErr(t, err)
}

func TestAuthorizationCodeRepository_DeleteByCode_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideAuthorizationCodeRepository(pg)

	err := repo.DeleteByCode(testCtx(), "code")
	requireDatabaseErr(t, err)
}

func TestAuthorizationCodeRepository_TakeByCode_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideAuthorizationCodeRepository(pg)

	_, err := repo.TakeByCode(testCtx(), "code")
	requireDatabaseErr(t, err)
}

func TestFederatedIdentityRepository_GetByProviderSubject_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideFederatedIdentityRepository(pg)

	_, err := repo.GetByProviderSubject(testCtx(), "google", "sub-1")
	requireDatabaseErr(t, err)
}

func TestRefreshTokenRepository_RevokeByID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)

	err := repo.RevokeByID(testCtx(), "id")
	requireDatabaseErr(t, err)
}

func TestRefreshTokenRepository_RevokeAllValidForUser_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)

	err := repo.RevokeAllValidForUser(testCtx(), "user")
	requireDatabaseErr(t, err)
}

func TestRefreshTokenRepository_RevokeAndCreate_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)
	rt := &models.RefreshToken{
		HardDeleteModel: models.HardDeleteModel{ID: models.NewBaseModel().ID},
		TenantID:        "t",
		UserID:          "u",
		OAuthClientID:   "app",
		TokenHash:       "hash",
		ExpiresAt:       futureTime(),
	}

	err := repo.RevokeAndCreate(testCtx(), "old-id", rt)
	requireDatabaseErr(t, err)
}

func TestRefreshTokenRepository_UsageByClientRecordID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideRefreshTokenRepository(pg)

	_, err := repo.UsageByClientRecordID(testCtx(), "client-id")
	requireDatabaseErr(t, err)
}

func TestSessionRepository_DeleteByID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideSessionRepository(pg)

	err := repo.DeleteByID(testCtx(), "session-id")
	requireDatabaseErr(t, err)
}

func TestSessionRepository_DeleteAllByUserID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideSessionRepository(pg)

	err := repo.DeleteAllByUserID(testCtx(), "user-id")
	requireDatabaseErr(t, err)
}

func TestSessionRepository_CountActive_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideSessionRepository(pg)

	_, err := repo.CountActive(testCtx())
	requireDatabaseErr(t, err)
}

func TestSessionRepository_CountActiveByUserID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideSessionRepository(pg)

	_, err := repo.CountActiveByUserID(testCtx(), "user-id")
	requireDatabaseErr(t, err)
}

func TestTenantRepository_Delete_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideTenantRepository(pg)

	err := repo.Delete(testCtx(), "tenant-id")
	requireDatabaseErr(t, err)
}

func TestTenantRepository_DomainTaken_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideTenantRepository(pg)

	_, err := repo.DomainTaken(testCtx(), "example.com", "")
	requireDatabaseErr(t, err)
}

func TestTenantRepository_CountUsersByTenantID_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideTenantRepository(pg)

	_, err := repo.CountUsersByTenantID(testCtx(), "tenant-id")
	requireDatabaseErr(t, err)
}

func TestPaginateWithFind_CountErrorOnly(t *testing.T) {
	openPG := newTestDB(t)
	closedPG := closedTestDB(t)
	ctx := testCtx()

	_, err := PaginateWithFind[models.Tenant](ctx, closedPG.WithContext(ctx).Model(&models.Tenant{}), openPG.WithContext(ctx).Model(&models.Tenant{}), nil, PaginateOptions{
		CountFailureMessage: "Failed to count tenants",
		FindFailureMessage:  "Failed to list tenants",
		DBOp:                DBOp{Operation: "list_tenants", Resource: "tenant"},
	})
	requireDatabaseErr(t, err)
}
