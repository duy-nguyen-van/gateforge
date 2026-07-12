package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTenantIdentityProviderRepository_UpdateProviderCreate(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantIdentityProviderRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")

	enabled := true
	tip, err := repo.UpdateProvider(ctx, tenant.ID, "google", TenantIdentityProviderPatch{
		Enabled:                    &enabled,
		OAuthClientID:              "client-id",
		OAuthClientSecretEncrypted: "enc-secret",
	})
	require.NoError(t, err)
	require.True(t, tip.Enabled)
	require.Equal(t, "client-id", tip.OAuthClientID)

	got, err := repo.GetByTenantAndProvider(ctx, tenant.ID, "google")
	require.NoError(t, err)
	require.Equal(t, tip.ID, got.ID)
}

func TestTenantIdentityProviderRepository_IsProviderEnabled(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantIdentityProviderRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")

	enabled, err := repo.IsProviderEnabled(ctx, tenant.ID, "google")
	require.NoError(t, err)
	require.False(t, enabled)

	require.NoError(t, repo.SetProviderEnabled(ctx, tenant.ID, "google", true))
	enabled, err = repo.IsProviderEnabled(ctx, tenant.ID, "google")
	require.NoError(t, err)
	require.True(t, enabled)
}

func TestTenantIdentityProviderRepository_IsProviderConfigured(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantIdentityProviderRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")

	configured, err := repo.IsProviderConfigured(ctx, tenant.ID, "google")
	require.NoError(t, err)
	require.False(t, configured)

	enabled := true
	_, err = repo.UpdateProvider(ctx, tenant.ID, "google", TenantIdentityProviderPatch{
		Enabled:                    &enabled,
		OAuthClientID:              "cid",
		OAuthClientSecretEncrypted: "secret",
	})
	require.NoError(t, err)

	configured, err = repo.IsProviderConfigured(ctx, tenant.ID, "google")
	require.NoError(t, err)
	require.True(t, configured)
}

func TestTenantIdentityProviderRepository_GetByTenantAndProvider_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantIdentityProviderRepository(pg)
	tenant := seedTenant(t, pg, "Org", "org.test")

	_, err := repo.GetByTenantAndProvider(testCtx(), tenant.ID, "google")
	requireNotFound(t, err)
}

func TestTenantIdentityProviderRepository_ListByTenant(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantIdentityProviderRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	enabled := true
	_, err := repo.UpdateProvider(ctx, tenant.ID, "google", TenantIdentityProviderPatch{Enabled: &enabled})
	require.NoError(t, err)

	rows, err := repo.ListByTenant(ctx, tenant.ID)
	require.NoError(t, err)
	require.Len(t, rows, 1)
}

func TestTenantIdentityProviderRepository_UpdateProvider_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideTenantIdentityProviderRepository(pg)
	enabled := true

	_, err := repo.UpdateProvider(testCtx(), "t", "google", TenantIdentityProviderPatch{Enabled: &enabled})
	requireDatabaseOrClosedErr(t, err)
}

func TestTenantIdentityProviderRepository_ListByTenant_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideTenantIdentityProviderRepository(pg)

	_, err := repo.ListByTenant(testCtx(), "t")
	requireDatabaseErr(t, err)
}

func TestFederationRepoNotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideTenantIdentityProviderRepository(pg)
	tenant := seedTenant(t, pg, "Org", "org.test")

	enabled, err := repo.IsProviderEnabled(testCtx(), tenant.ID, "missing")
	require.NoError(t, err)
	require.False(t, enabled)
}
