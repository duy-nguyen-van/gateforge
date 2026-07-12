package services

import (
	"context"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/stretchr/testify/require"
)

type adminTenantTestRepo struct {
	tenants    map[string]*models.Tenant
	userCounts map[string]int64
}

func newAdminTenantTestRepo() *adminTenantTestRepo {
	return &adminTenantTestRepo{tenants: map[string]*models.Tenant{}, userCounts: map[string]int64{}}
}

func (r *adminTenantTestRepo) List(_ context.Context, pr *dtos.PageableRequest) (*dtos.DataResponse[models.Tenant], error) {
	rows := make([]models.Tenant, 0, len(r.tenants))
	for _, t := range r.tenants {
		rows = append(rows, *t)
	}
	return &dtos.DataResponse[models.Tenant]{
		Data: rows,
		Pageable: &dtos.Pageable{
			Page:     pr.Page,
			PageSize: pr.PageSize,
			Total:    int64(len(rows)),
		},
	}, nil
}

func (r *adminTenantTestRepo) GetByID(_ context.Context, id string) (*models.Tenant, error) {
	t, ok := r.tenants[id]
	if !ok {
		return nil, errors.NotFoundError("Tenant", nil)
	}
	return t, nil
}

func (r *adminTenantTestRepo) GetByDomain(_ context.Context, domain string) (*models.Tenant, error) {
	for _, t := range r.tenants {
		if t.Domain == domain {
			return t, nil
		}
	}
	return nil, errors.NotFoundError("Tenant", nil)
}

func (r *adminTenantTestRepo) Create(_ context.Context, tenant *models.Tenant) error {
	r.tenants[tenant.ID] = tenant
	return nil
}

func (r *adminTenantTestRepo) Update(_ context.Context, id string, patch repositories.TenantPatch) (*models.Tenant, error) {
	t, ok := r.tenants[id]
	if !ok {
		return nil, errors.NotFoundError("Tenant", nil)
	}
	if patch.Name != nil {
		t.Name = *patch.Name
	}
	if patch.Domain != nil {
		t.Domain = *patch.Domain
	}
	return t, nil
}

func (r *adminTenantTestRepo) Delete(_ context.Context, id string) error {
	if _, ok := r.tenants[id]; !ok {
		return errors.NotFoundError("Tenant", nil)
	}
	delete(r.tenants, id)
	return nil
}

func (r *adminTenantTestRepo) DomainTaken(_ context.Context, domain, excludeID string) (bool, error) {
	for id, t := range r.tenants {
		if t.Domain == domain && domain != "" && id != excludeID {
			return true, nil
		}
	}
	return false, nil
}

func (r *adminTenantTestRepo) CountUsersByTenantID(_ context.Context, tenantID string) (int64, error) {
	return r.userCounts[tenantID], nil
}

type adminTenantMembershipStub struct {
	byTenant map[string][]models.TenantMembership
}

func (s *adminTenantMembershipStub) Create(context.Context, *models.TenantMembership) error {
	return nil
}
func (s *adminTenantMembershipStub) GetActive(context.Context, string, string) (*models.TenantMembership, error) {
	return nil, errors.NotFoundError("Tenant membership", nil)
}
func (s *adminTenantMembershipStub) ExistsActive(context.Context, string, string) (bool, error) {
	return false, nil
}
func (s *adminTenantMembershipStub) ListByUserID(context.Context, string) ([]models.TenantMembership, error) {
	return nil, nil
}
func (s *adminTenantMembershipStub) ListByUserIDPaginated(context.Context, string, *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	return nil, nil
}
func (s *adminTenantMembershipStub) ListByTenantID(_ context.Context, tenantID string) ([]models.TenantMembership, error) {
	return s.byTenant[tenantID], nil
}
func (s *adminTenantMembershipStub) ListByTenantIDPaginated(_ context.Context, tenantID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	rows := s.byTenant[tenantID]
	return &dtos.DataResponse[models.TenantMembership]{
		Data: rows,
		Pageable: &dtos.Pageable{
			Page:     pr.Page,
			PageSize: pr.PageSize,
			Total:    int64(len(rows)),
		},
	}, nil
}
func (s *adminTenantMembershipStub) CountByTenantID(context.Context, string) (int64, error) {
	return 0, nil
}
func (s *adminTenantMembershipStub) Delete(context.Context, string, string) error { return nil }

type adminTenantAuditStub struct {
	last domains.AuditRecordParams
}

func (a *adminTenantAuditStub) Record(_ context.Context, p domains.AuditRecordParams) {
	a.last = p
}

func newAdminTenantTestService(tenants *adminTenantTestRepo, memberships *adminTenantMembershipStub) AdminService {
	cfg := &config.Config{DefaultTenantID: "00000000-0000-0000-0000-000000000001"}
	return ProvideAdminService(
		cfg,
		nil, tenants, nil, nil, nil, nil, nil,
		memberships,
		nil, nil, nil,
		&adminTenantAuditStub{},
		nil, nil,
	)
}

func TestAdminService_CreateTenant(t *testing.T) {
	repo := newAdminTenantTestRepo()
	svc := newAdminTenantTestService(repo, &adminTenantMembershipStub{byTenant: map[string][]models.TenantMembership{}})

	tenant, err := svc.CreateTenant(context.Background(), &dtos.AdminCreateTenantRequest{
		Name:   "Acme Corp",
		Domain: "acme",
	})
	require.NoError(t, err)
	require.Equal(t, "Acme Corp", tenant.Name)
	require.Equal(t, "acme", tenant.Domain)
	require.NotEmpty(t, tenant.ID)
}

func TestAdminService_CreateTenant_DuplicateDomain(t *testing.T) {
	repo := newAdminTenantTestRepo()
	repo.tenants["existing"] = &models.Tenant{BaseModel: models.BaseModel{ID: "existing"}, Name: "Other", Domain: "acme"}
	svc := newAdminTenantTestService(repo, &adminTenantMembershipStub{byTenant: map[string][]models.TenantMembership{}})

	_, err := svc.CreateTenant(context.Background(), &dtos.AdminCreateTenantRequest{
		Name:   "Acme Corp",
		Domain: "acme",
	})
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.Equal(t, errors.ErrorTypeConflict, appErr.Type)
}

func TestAdminService_DeleteTenant_BlocksDefault(t *testing.T) {
	repo := newAdminTenantTestRepo()
	defaultID := "00000000-0000-0000-0000-000000000001"
	repo.tenants[defaultID] = &models.Tenant{BaseModel: models.BaseModel{ID: defaultID}, Name: "Default"}
	svc := newAdminTenantTestService(repo, &adminTenantMembershipStub{byTenant: map[string][]models.TenantMembership{}})

	err := svc.DeleteTenant(context.Background(), defaultID)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.Equal(t, errors.ErrorTypeValidation, appErr.Type)
}

func TestAdminService_UpdateTenant(t *testing.T) {
	repo := newAdminTenantTestRepo()
	id := "tenant-1"
	repo.tenants[id] = &models.Tenant{BaseModel: models.BaseModel{ID: id}, Name: "Old", Domain: "old"}
	svc := newAdminTenantTestService(repo, &adminTenantMembershipStub{byTenant: map[string][]models.TenantMembership{}})

	name := "New Name"
	updated, err := svc.UpdateTenant(context.Background(), id, &dtos.AdminUpdateTenantRequest{Name: &name})
	require.NoError(t, err)
	require.Equal(t, "New Name", updated.Name)
}

func TestAdminService_ListTenantMembers(t *testing.T) {
	repo := newAdminTenantTestRepo()
	tenantID := "tenant-1"
	repo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Acme"}
	memberships := &adminTenantMembershipStub{
		byTenant: map[string][]models.TenantMembership{
			tenantID: {{
				BaseModel: models.BaseModel{ID: "m1", CreatedAt: time.Now()},
				UserID:    "u1",
				TenantID:  tenantID,
				Role:      constants.TenantMembershipRoleMember,
				Status:    constants.TenantMembershipStatusActive,
				User:      &models.User{BaseModel: models.BaseModel{ID: "u1"}, Email: "user@example.com", FirstName: "Jane", LastName: "Doe"},
			}},
		},
	}
	svc := newAdminTenantTestService(repo, memberships)

	members, pageable, err := svc.ListTenantMembers(context.Background(), tenantID, dtos.NewPageableRequest())
	require.NoError(t, err)
	require.Len(t, members, 1)
	require.Equal(t, "user@example.com", members[0].Email)
	require.Equal(t, int64(1), pageable.Total)
}

func TestAdminService_AddMemberByEmail_TenantNotFound(t *testing.T) {
	repo := newAdminTenantTestRepo()
	svc := newAdminTenantTestService(repo, &adminTenantMembershipStub{byTenant: map[string][]models.TenantMembership{}})

	err := svc.AddMemberByEmail(context.Background(), "missing-tenant", "user@example.com", "member")
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.Equal(t, errors.ErrorTypeNotFound, appErr.Type)
}
