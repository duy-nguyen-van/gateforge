package repositories

import (
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/stretchr/testify/require"
)

func TestAuditLogRepository_CreateListCount(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideAuditLogRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	tenantID := tenant.ID
	actorID := "actor-1"
	resourceType := "user"
	resourceName := "jane@test.com"

	log := &models.AuditLog{
		BaseModel:    models.NewBaseModel(),
		TenantID:     &tenantID,
		Action:       "user.login",
		Result:       "success",
		ActorType:    "user",
		ActorID:      &actorID,
		ResourceType: &resourceType,
		ResourceName: &resourceName,
	}
	require.NoError(t, repo.Create(ctx, log))

	filters := AuditLogListFilters{
		TenantID:     tenantID,
		Action:       "user.login",
		Result:       "success",
		ActorID:      actorID,
		ResourceType: resourceType,
		ResourceName: resourceName,
	}
	resp, err := repo.List(ctx, filters, &dtos.PageableRequest{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)

	count, err := repo.Count(ctx, filters)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
}

func TestAuditLogRepository_ListFilters(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideAuditLogRepository(pg)
	ctx := testCtx()
	tenant := seedTenant(t, pg, "Org", "org.test")
	tenantID := tenant.ID

	for _, action := range []string{"user.create", "user.update", "tenant.create"} {
		a := action
		require.NoError(t, repo.Create(ctx, &models.AuditLog{
			BaseModel: models.NewBaseModel(),
			TenantID:  &tenantID,
			Action:    a,
			Result:    "success",
			ActorType: "system",
		}))
	}

	from := time.Now().UTC().Add(-time.Hour)
	to := time.Now().UTC().Add(time.Hour)

	resp, err := repo.List(ctx, AuditLogListFilters{
		TenantID:  tenantID,
		ActionsIn: []string{"user.create", "user.update"},
		From:      &from,
		To:        &to,
	}, &dtos.PageableRequest{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, resp.Data, 2)

	prefix, err := repo.List(ctx, AuditLogListFilters{
		TenantID: tenantID,
		Action:   "user.",
	}, &dtos.PageableRequest{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, prefix.Data, 2)
}

func TestAuditLogRepository_Create_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideAuditLogRepository(pg)
	log := &models.AuditLog{
		BaseModel: models.NewBaseModel(),
		Action:    "test",
		Result:    "success",
		ActorType: "system",
	}

	err := repo.Create(testCtx(), log)
	requireDatabaseErr(t, err)
}

func TestAuditLogRepository_Count_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideAuditLogRepository(pg)

	_, err := repo.Count(testCtx(), AuditLogListFilters{})
	requireDatabaseErr(t, err)
}
