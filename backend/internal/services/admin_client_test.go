package services

import (
	"context"
	"testing"

	"github.com/lib/pq"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/stretchr/testify/require"
)

type adminClientTestRepo struct {
	clients map[string]*models.Client
}

func newAdminClientTestRepo() *adminClientTestRepo {
	return &adminClientTestRepo{clients: map[string]*models.Client{}}
}

func (r *adminClientTestRepo) GetByTenantAndClientID(_ context.Context, tenantID, clientID string) (*models.Client, error) {
	for _, c := range r.clients {
		if c.TenantID == tenantID && c.ClientID == clientID {
			return c, nil
		}
	}
	return nil, errors.NotFoundError("OAuth client", nil)
}

func (r *adminClientTestRepo) GetByClientID(_ context.Context, clientID string) (*models.Client, error) {
	for _, c := range r.clients {
		if c.ClientID == clientID {
			return c, nil
		}
	}
	return nil, errors.NotFoundError("OAuth client", nil)
}

func (r *adminClientTestRepo) GetByID(_ context.Context, id string) (*models.Client, error) {
	c, ok := r.clients[id]
	if !ok {
		return nil, errors.NotFoundError("OAuth client", nil)
	}
	return c, nil
}

func (r *adminClientTestRepo) List(_ context.Context, _ string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.Client], error) {
	rows := make([]models.Client, 0, len(r.clients))
	for _, c := range r.clients {
		rows = append(rows, *c)
	}
	return &dtos.DataResponse[models.Client]{
		Data: rows,
		Pageable: &dtos.Pageable{
			Page:     pr.Page,
			PageSize: pr.PageSize,
			Total:    int64(len(rows)),
		},
	}, nil
}

func (r *adminClientTestRepo) Create(_ context.Context, client *models.Client) error {
	r.clients[client.ID] = client
	return nil
}

func (r *adminClientTestRepo) Update(_ context.Context, id string, patch repositories.ClientPatch) (*models.Client, error) {
	c, ok := r.clients[id]
	if !ok {
		return nil, errors.NotFoundError("OAuth client", nil)
	}
	if patch.Name != nil {
		c.Name = *patch.Name
	}
	if patch.RedirectUris != nil {
		c.RedirectUris = pq.StringArray(*patch.RedirectUris)
	}
	if patch.GrantTypes != nil {
		c.GrantTypes = pq.StringArray(*patch.GrantTypes)
	}
	if patch.Scopes != nil {
		c.Scopes = pq.StringArray(*patch.Scopes)
	}
	if patch.IsPublic != nil {
		c.IsPublic = *patch.IsPublic
	}
	if patch.ClientSecret != nil {
		c.ClientSecret = *patch.ClientSecret
	}
	return c, nil
}

func (r *adminClientTestRepo) Delete(_ context.Context, id string) error {
	if _, ok := r.clients[id]; !ok {
		return errors.NotFoundError("OAuth client", nil)
	}
	delete(r.clients, id)
	return nil
}

func (r *adminClientTestRepo) ClientIDTaken(_ context.Context, tenantID, clientID, excludeID string) (bool, error) {
	for id, c := range r.clients {
		if c.TenantID == tenantID && c.ClientID == clientID && id != excludeID {
			return true, nil
		}
	}
	return false, nil
}

func newAdminClientTestService(tenants *adminTenantTestRepo, clients *adminClientTestRepo) AdminService {
	cfg := &config.Config{DefaultTenantID: "00000000-0000-0000-0000-000000000001"}
	return ProvideAdminService(
		cfg,
		nil, tenants, clients, nil, nil, nil, nil,
		&adminTenantMembershipStub{byTenant: map[string][]models.TenantMembership{}},
		nil, nil, nil,
		&adminTenantAuditStub{},
		nil, nil,
	)
}

func TestAdminService_CreateClient_Public(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "00000000-0000-0000-0000-000000000002"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Acme"}
	clientRepo := newAdminClientTestRepo()
	svc := newAdminClientTestService(tenantRepo, clientRepo)

	client, err := svc.CreateClient(context.Background(), &dtos.AdminCreateClientRequest{
		TenantID:     tenantID,
		ClientID:     "my-app",
		Name:         "My App",
		IsPublic:     true,
		RedirectUris: []string{"http://localhost:5173/callback"},
	})
	require.NoError(t, err)
	require.Equal(t, "my-app", client.ClientID)
	require.Equal(t, "My App", client.Name)
	require.True(t, client.IsPublic)
	require.Empty(t, client.ClientSecret)
	require.False(t, client.ClientSecretSet)
	require.Equal(t, defaultClientGrantTypes, client.GrantTypes)
	require.Equal(t, defaultClientScopes, client.Scopes)
}

func TestAdminService_CreateClient_Confidential(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "00000000-0000-0000-0000-000000000002"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Acme"}
	clientRepo := newAdminClientTestRepo()
	svc := newAdminClientTestService(tenantRepo, clientRepo)

	client, err := svc.CreateClient(context.Background(), &dtos.AdminCreateClientRequest{
		TenantID:     tenantID,
		Name:         "Backend App",
		RedirectUris: []string{"http://localhost:3000/callback"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, client.ClientID)
	require.NotEmpty(t, client.ClientSecret)
	require.True(t, client.ClientSecretSet)
	require.False(t, client.IsPublic)
}

func TestAdminService_CreateClient_DuplicateClientID(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "00000000-0000-0000-0000-000000000002"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Acme"}
	clientRepo := newAdminClientTestRepo()
	clientRepo.clients["existing"] = &models.Client{
		BaseModel: models.BaseModel{ID: "existing"},
		TenantID:  tenantID,
		ClientID:  "my-app",
	}
	svc := newAdminClientTestService(tenantRepo, clientRepo)

	_, err := svc.CreateClient(context.Background(), &dtos.AdminCreateClientRequest{
		TenantID:     tenantID,
		ClientID:     "my-app",
		Name:         "Duplicate",
		IsPublic:     true,
		RedirectUris: []string{"http://localhost:5173/callback"},
	})
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.Equal(t, errors.ErrorTypeConflict, appErr.Type)
}

func TestAdminService_UpdateClient(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	clientRepo := newAdminClientTestRepo()
	id := "client-1"
	clientRepo.clients[id] = &models.Client{
		BaseModel:    models.BaseModel{ID: id},
		TenantID:     "tenant-1",
		ClientID:     "my-app",
		Name:         "Old Name",
		RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes:   pq.StringArray{"authorization_code"},
		Scopes:       pq.StringArray{"openid"},
		IsPublic:     true,
	}
	svc := newAdminClientTestService(tenantRepo, clientRepo)

	name := "New Name"
	updated, err := svc.UpdateClient(context.Background(), id, &dtos.AdminUpdateClientRequest{
		Name: &name,
	})
	require.NoError(t, err)
	require.Equal(t, "New Name", updated.Name)
}

func TestAdminService_UpdateClient_RotateSecret(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	clientRepo := newAdminClientTestRepo()
	id := "client-1"
	clientRepo.clients[id] = &models.Client{
		BaseModel:    models.BaseModel{ID: id},
		TenantID:     "tenant-1",
		ClientID:     "backend",
		Name:         "Backend",
		ClientSecret: "old-secret",
		IsPublic:     false,
	}
	svc := newAdminClientTestService(tenantRepo, clientRepo)

	secret := "new-secret"
	updated, err := svc.UpdateClient(context.Background(), id, &dtos.AdminUpdateClientRequest{
		ClientSecret: &secret,
	})
	require.NoError(t, err)
	require.True(t, updated.ClientSecretSet)
	require.Equal(t, "new-secret", clientRepo.clients[id].ClientSecret)
}

func TestAdminService_DeleteClient_BlocksDevClient(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	clientRepo := newAdminClientTestRepo()
	defaultTenant := "00000000-0000-0000-0000-000000000001"
	id := "dev-client"
	clientRepo.clients[id] = &models.Client{
		BaseModel: models.BaseModel{ID: id},
		TenantID:  defaultTenant,
		ClientID:  seededDevClientID,
		Name:      "Dev",
		IsPublic:  true,
	}
	svc := newAdminClientTestService(tenantRepo, clientRepo)

	err := svc.DeleteClient(context.Background(), id)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.Equal(t, errors.ErrorTypeValidation, appErr.Type)
}

func TestAdminService_DeleteClient(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	clientRepo := newAdminClientTestRepo()
	id := "client-1"
	clientRepo.clients[id] = &models.Client{
		BaseModel: models.BaseModel{ID: id},
		TenantID:  "tenant-1",
		ClientID:  "my-app",
	}
	svc := newAdminClientTestService(tenantRepo, clientRepo)

	err := svc.DeleteClient(context.Background(), id)
	require.NoError(t, err)
	_, ok := clientRepo.clients[id]
	require.False(t, ok)
}

func TestAdminService_GetClientByID(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	clientRepo := newAdminClientTestRepo()
	id := "client-1"
	clientRepo.clients[id] = &models.Client{
		BaseModel: models.BaseModel{ID: id},
		TenantID:  "tenant-1",
		ClientID:  "my-app",
		Name:      "My App",
	}
	svc := newAdminClientTestService(tenantRepo, clientRepo)

	client, err := svc.GetClientByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, "my-app", client.ClientID)
}
