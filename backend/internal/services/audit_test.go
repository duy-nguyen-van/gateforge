package services

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"github.com/stretchr/testify/require"
)

type stubAuditLogRepo struct {
	created []*models.AuditLog
	err     error
	list    *dtos.DataResponse[models.AuditLog]
	listErr error
}

func (s *stubAuditLogRepo) Create(_ context.Context, log *models.AuditLog) error {
	if s.err != nil {
		return s.err
	}
	s.created = append(s.created, log)
	return nil
}

func (s *stubAuditLogRepo) List(_ context.Context, _ repositories.AuditLogListFilters, _ *dtos.PageableRequest) (*dtos.DataResponse[models.AuditLog], error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return s.list, nil
}

func (s *stubAuditLogRepo) Count(_ context.Context, _ repositories.AuditLogListFilters) (int64, error) {
	return 0, nil
}

type noopAuditService struct{}

func (noopAuditService) Record(context.Context, domains.AuditRecordParams) {}

func TestAuditService_RecordExplicitFields(t *testing.T) {
	repo := &stubAuditLogRepo{}
	svc := ProvideAuditService(repo)

	svc.Record(context.Background(), domains.AuditRecordParams{
		Action:       constants.AuditActionAdminTenantCreate,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      "actor-1",
		TenantID:     "tenant-1",
		ResourceType: constants.AuditResourceTypeTenant,
		ResourceID:   "tenant-1",
		ResourceName: "Acme",
		OldValue:     map[string]any{"name": "Old"},
		NewValue:     map[string]any{"name": "New"},
	})

	require.Len(t, repo.created, 1)
	row := repo.created[0]
	require.NotNil(t, row.ResourceType)
	require.NotNil(t, row.OldValue)
	require.NotNil(t, row.NewValue)
}

func TestAuditService_RecordMergesContext(t *testing.T) {
	repo := &stubAuditLogRepo{}
	svc := ProvideAuditService(repo)

	ctx := request.NewAuditContextContext(context.Background(), request.AuditContext{
		IPAddress:     "127.0.0.1",
		UserAgent:     "test-agent",
		RequestID:     "req-1",
		CorrelationID: "corr-1",
		ActorType:     string(constants.AuditActorTypeUser),
		ActorID:       "actor-from-ctx",
		TenantID:      "tenant-from-ctx",
	})

	svc.Record(ctx, domains.AuditRecordParams{
		Action: constants.AuditActionAuthLogin,
		Result: constants.AuditResultSuccess,
	})

	require.Len(t, repo.created, 1)
	row := repo.created[0]
	require.Equal(t, constants.AuditActionAuthLogin, row.Action)
	require.Equal(t, string(constants.AuditResultSuccess), row.Result)
	require.NotNil(t, row.TenantID)
	require.Equal(t, "tenant-from-ctx", *row.TenantID)
	require.NotNil(t, row.ActorID)
	require.Equal(t, "actor-from-ctx", *row.ActorID)
	require.NotNil(t, row.IPAddress)
	require.Equal(t, "127.0.0.1", *row.IPAddress)
}

func TestAuditService_RecordDoesNotPropagateRepoError(t *testing.T) {
	repo := &stubAuditLogRepo{err: context.Canceled}
	svc := ProvideAuditService(repo)
	require.NotPanics(t, func() {
		svc.Record(context.Background(), domains.AuditRecordParams{
			Action:    constants.AuditActionAuthLogin,
			Result:    constants.AuditResultFailure,
			ActorType: constants.AuditActorTypeUser,
		})
	})
}

func TestAdminService_ListAuditLogs(t *testing.T) {
	tenantID := "00000000-0000-4000-8000-000000000099"
	repo := &stubAuditLogRepo{
		list: &dtos.DataResponse[models.AuditLog]{
			Data: []models.AuditLog{
				{
					BaseModel: models.NewBaseModel(),
					Action:    constants.AuditActionAdminMemberAdd,
					Result:    string(constants.AuditResultSuccess),
					ActorType: string(constants.AuditActorTypeUser),
					TenantID:  &tenantID,
				},
			},
			Pageable: &dtos.Pageable{Page: 1, PageSize: 20, Total: 1},
		},
	}
	svc := ProvideAdminService(&config.Config{}, nil, nil, nil, nil, nil, nil, nil, nil, repo, nil, nil, noopAuditService{}, nil, nil)

	rows, pageable, err := svc.ListAuditLogs(context.Background(), dtos.AdminAuditLogListParams{
		TenantID: tenantID,
	}, dtos.NewPageableRequest())
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, constants.AuditActionAdminMemberAdd, rows[0].Action)
	require.Equal(t, int64(1), pageable.Total)
}
