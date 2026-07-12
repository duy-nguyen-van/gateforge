package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

func TestClientRepository_CRUD(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideClientRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	client := seedClient(t, pg, tenant.ID, "app-client")

	got, err := repo.GetByID(ctx, client.ID)
	require.NoError(t, err)
	require.Equal(t, "app-client", got.ClientID)

	byTenant, err := repo.GetByTenantAndClientID(ctx, tenant.ID, "app-client")
	require.NoError(t, err)
	require.Equal(t, client.ID, byTenant.ID)

	byClientID, err := repo.GetByClientID(ctx, "app-client")
	require.NoError(t, err)
	require.Equal(t, client.ID, byClientID.ID)

	name := "Renamed"
	redirects := []string{"https://app/cb"}
	updated, err := repo.Update(ctx, client.ID, ClientPatch{Name: &name, RedirectUris: &redirects})
	require.NoError(t, err)
	require.Equal(t, "Renamed", updated.Name)
	require.Equal(t, redirects, []string(updated.RedirectUris))

	require.NoError(t, repo.Delete(ctx, client.ID))
	_, err = repo.GetByID(ctx, client.ID)
	requireNotFound(t, err)
}

func TestClientRepository_GetByID_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideClientRepository(pg)

	_, err := repo.GetByID(testCtx(), "00000000-0000-7000-8000-000000000001")
	requireNotFound(t, err)
}

func TestClientRepository_Update_EmptyPatch(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideClientRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	client := seedClient(t, pg, tenant.ID, "client-a")

	got, err := repo.Update(ctx, client.ID, ClientPatch{})
	require.NoError(t, err)
	require.Equal(t, client.Name, got.Name)
}

func TestClientRepository_Update_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideClientRepository(pg)
	name := "X"

	_, err := repo.Update(testCtx(), "00000000-0000-7000-8000-000000000099", ClientPatch{Name: &name})
	requireNotFound(t, err)
}

func TestClientRepository_Delete_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideClientRepository(pg)

	err := repo.Delete(testCtx(), "00000000-0000-7000-8000-000000000099")
	requireNotFound(t, err)
}

func TestClientRepository_List(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideClientRepository(pg)
	ctx := testCtx()
	t1 := seedTenant(t, pg, "T1", "t1.test")
	t2 := seedTenant(t, pg, "T2", "t2.test")
	seedClient(t, pg, t1.ID, "c1")
	seedClient(t, pg, t2.ID, "c2")

	resp, err := repo.List(ctx, t1.ID, &dtos.PageableRequest{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
}

func TestClientRepository_ClientIDTaken(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideClientRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	client := seedClient(t, pg, tenant.ID, "dup-id")

	taken, err := repo.ClientIDTaken(ctx, tenant.ID, "dup-id", "")
	require.NoError(t, err)
	require.True(t, taken)

	free, err := repo.ClientIDTaken(ctx, tenant.ID, "dup-id", client.ID)
	require.NoError(t, err)
	require.False(t, free)

	empty, err := repo.ClientIDTaken(ctx, tenant.ID, "", "")
	require.NoError(t, err)
	require.False(t, empty)
}

func TestClientRepository_Create_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideClientRepository(pg)
	client := &models.Client{BaseModel: models.NewBaseModel(), TenantID: "t", ClientID: "c"}

	err := repo.Create(testCtx(), client)
	requireDatabaseErr(t, err)
}
