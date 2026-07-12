package services

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/stretchr/testify/require"
)

type stubMembershipRepo struct {
	byUser map[string][]models.TenantMembership
	active map[string]map[string]bool
}

func (s *stubMembershipRepo) Create(ctx context.Context, m *models.TenantMembership) error {
	s.byUser[m.UserID] = append(s.byUser[m.UserID], *m)
	if s.active[m.UserID] == nil {
		s.active[m.UserID] = map[string]bool{}
	}
	s.active[m.UserID][m.TenantID] = true
	return nil
}

func (s *stubMembershipRepo) GetActive(ctx context.Context, userID, tenantID string) (*models.TenantMembership, error) {
	return nil, nil
}

func (s *stubMembershipRepo) ExistsActive(ctx context.Context, userID, tenantID string) (bool, error) {
	return s.active[userID][tenantID], nil
}

func (s *stubMembershipRepo) ListByUserID(ctx context.Context, userID string) ([]models.TenantMembership, error) {
	return s.byUser[userID], nil
}

func (s *stubMembershipRepo) ListByUserIDPaginated(_ context.Context, userID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	rows := s.byUser[userID]
	page, pageable := dtos.PaginateSlice(rows, pr)
	return &dtos.DataResponse[models.TenantMembership]{Data: page, Pageable: pageable}, nil
}

func (s *stubMembershipRepo) ListByTenantID(ctx context.Context, tenantID string) ([]models.TenantMembership, error) {
	return nil, nil
}

func (s *stubMembershipRepo) ListByTenantIDPaginated(_ context.Context, tenantID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	return &dtos.DataResponse[models.TenantMembership]{Pageable: &dtos.Pageable{}}, nil
}

func (s *stubMembershipRepo) CountByTenantID(ctx context.Context, tenantID string) (int64, error) {
	return 0, nil
}

func (s *stubMembershipRepo) Delete(ctx context.Context, userID, tenantID string) error {
	return nil
}

type stubClientRepo struct {
	byClientID map[string]*models.Client
}

func (s *stubClientRepo) GetByTenantAndClientID(ctx context.Context, tenantID, clientID string) (*models.Client, error) {
	return nil, nil
}

func (s *stubClientRepo) GetByClientID(ctx context.Context, clientID string) (*models.Client, error) {
	c, ok := s.byClientID[clientID]
	if !ok || c == nil {
		return nil, errors.NotFoundError("OAuth client", nil)
	}
	return c, nil
}

func (s *stubClientRepo) GetByID(ctx context.Context, id string) (*models.Client, error) {
	for _, c := range s.byClientID {
		if c != nil && c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}

func (s *stubClientRepo) List(ctx context.Context, tenantID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.Client], error) {
	return nil, nil
}

func (s *stubClientRepo) Create(ctx context.Context, client *models.Client) error { return nil }
func (s *stubClientRepo) Update(ctx context.Context, id string, patch repositories.ClientPatch) (*models.Client, error) {
	return nil, nil
}
func (s *stubClientRepo) Delete(ctx context.Context, id string) error { return nil }
func (s *stubClientRepo) ClientIDTaken(ctx context.Context, tenantID, clientID, excludeID string) (bool, error) {
	return false, nil
}

type stubTenantRepo struct {
	byDomain map[string]*models.Tenant
}

func (s *stubTenantRepo) List(ctx context.Context, pr *dtos.PageableRequest) (*dtos.DataResponse[models.Tenant], error) {
	return nil, nil
}

func (s *stubTenantRepo) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	return nil, nil
}

func (s *stubTenantRepo) GetByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	return s.byDomain[domain], nil
}

func (s *stubTenantRepo) Create(ctx context.Context, tenant *models.Tenant) error { return nil }
func (s *stubTenantRepo) Update(ctx context.Context, id string, patch repositories.TenantPatch) (*models.Tenant, error) {
	return nil, nil
}
func (s *stubTenantRepo) Delete(ctx context.Context, id string) error { return nil }
func (s *stubTenantRepo) DomainTaken(ctx context.Context, domain, excludeID string) (bool, error) {
	return false, nil
}

func (s *stubTenantRepo) CountUsersByTenantID(ctx context.Context, tenantID string) (int64, error) {
	return 0, nil
}

func TestTenantContextService_SingleMembershipAutoSelect(t *testing.T) {
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			"user-1": {{UserID: "user-1", TenantID: "tenant-a", Status: constants.TenantMembershipStatusActive}},
		},
		active: map[string]map[string]bool{"user-1": {"tenant-a": true}},
	}
	svc := ProvideTenantContextService(testConfig(), &stubClientRepo{}, &stubTenantRepo{}, memberships)

	res, err := svc.Resolve(context.Background(), TenantResolveInput{UserID: "user-1"})
	require.NoError(t, err)
	require.False(t, res.RequiresPicker)
	require.Equal(t, "tenant-a", res.TenantID)
}

func TestTenantContextService_MultipleMembershipsRequiresPicker(t *testing.T) {
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			"user-1": {
				{UserID: "user-1", TenantID: "tenant-a", Status: constants.TenantMembershipStatusActive},
				{UserID: "user-1", TenantID: "tenant-b", Status: constants.TenantMembershipStatusActive},
			},
		},
		active: map[string]map[string]bool{"user-1": {"tenant-a": true, "tenant-b": true}},
	}
	svc := ProvideTenantContextService(testConfig(), &stubClientRepo{}, &stubTenantRepo{}, memberships)

	res, err := svc.Resolve(context.Background(), TenantResolveInput{UserID: "user-1"})
	require.NoError(t, err)
	require.True(t, res.RequiresPicker)
}

func TestTenantContextService_ValidateMembership(t *testing.T) {
	memberships := &stubMembershipRepo{active: map[string]map[string]bool{"user-1": {"tenant-a": true}}}
	svc := ProvideTenantContextService(testConfig(), &stubClientRepo{}, &stubTenantRepo{}, memberships)

	require.NoError(t, svc.ValidateMembership(context.Background(), "user-1", "tenant-a"))
	require.Error(t, svc.ValidateMembership(context.Background(), "user-1", "tenant-b"))
}

func TestTenantContextService_ResolveByHostDomain(t *testing.T) {
	tenants := &stubTenantRepo{byDomain: map[string]*models.Tenant{
		"acme": {BaseModel: models.BaseModel{ID: "tenant-acme"}, Name: "Acme"},
	}}
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			"user-1": {{UserID: "user-1", TenantID: "tenant-acme"}},
		},
		active: map[string]map[string]bool{"user-1": {"tenant-acme": true}},
	}
	svc := ProvideTenantContextService(testConfig(), &stubClientRepo{}, tenants, memberships)

	res, err := svc.Resolve(context.Background(), TenantResolveInput{
		Host:   "acme.example.com",
		UserID: "user-1",
	})
	require.NoError(t, err)
	require.Equal(t, "tenant-acme", res.TenantID)
}

func TestTenantContextService_NoMembershipsForbidden(t *testing.T) {
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := ProvideTenantContextService(testConfig(), &stubClientRepo{}, &stubTenantRepo{}, memberships)

	_, err := svc.Resolve(context.Background(), TenantResolveInput{UserID: "user-none"})
	require.Error(t, err)
}

func TestTenantContextService_DefaultTenantPreference(t *testing.T) {
	cfg := testConfig()
	cfg.DefaultTenantID = "tenant-default"
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			"user-1": {
				{UserID: "user-1", TenantID: "tenant-a"},
				{UserID: "user-1", TenantID: "tenant-default"},
			},
		},
		active: map[string]map[string]bool{"user-1": {"tenant-a": true, "tenant-default": true}},
	}
	svc := ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships)

	res, err := svc.Resolve(context.Background(), TenantResolveInput{UserID: "user-1"})
	require.NoError(t, err)
	require.Equal(t, "tenant-default", res.TenantID)
}

func TestTenantContextService_OAuthClientResolvesTenant(t *testing.T) {
	memberships := &stubMembershipRepo{
		active: map[string]map[string]bool{"user-1": {"tenant-oauth": true}},
	}
	clients := &stubClientRepo{
		byClientID: map[string]*models.Client{
			"app-client": {TenantID: "tenant-oauth", ClientID: "app-client"},
		},
	}
	svc := ProvideTenantContextService(testConfig(), clients, &stubTenantRepo{}, memberships)

	res, err := svc.Resolve(context.Background(), TenantResolveInput{
		OAuthClientID: "app-client",
		UserID:        "user-1",
	})
	require.NoError(t, err)
	require.Equal(t, "tenant-oauth", res.TenantID)
}

func TestExtractDomainFromHost(t *testing.T) {
	require.Equal(t, "acme", extractDomainFromHost("acme.example.com"))
	require.Equal(t, "acme", extractDomainFromHost("ACME.EXAMPLE.COM:443"))
	require.Equal(t, "example.com", extractDomainFromHost("example.com"))
	require.Equal(t, "", extractDomainFromHost("localhost"))
}
