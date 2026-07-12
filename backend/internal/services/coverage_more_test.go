package services

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func TestAdminService_GetTenantByID_WithUserCount(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "tenant-count"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Counted"}
	tenantRepo.userCounts[tenantID] = 42
	svc := &adminService{tenants: tenantRepo}

	tenant, err := svc.GetTenantByID(context.Background(), tenantID)
	require.NoError(t, err)
	require.Equal(t, int64(42), tenant.UserCount)
}

func TestAdminTenantTestRepo_CountUsersByTenantID(t *testing.T) {
	repo := newAdminTenantTestRepo()
	repo.userCounts["t1"] = 5
	n, err := repo.CountUsersByTenantID(context.Background(), "t1")
	require.NoError(t, err)
	require.Equal(t, int64(5), n)
}

func TestAdminService_AddMemberByEmail_AlreadyMember(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "tenant-dup"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}}
	users := newUserTestRepo()
	u := users.seed("dup@example.com", "secret")
	memberships := &stubMembershipRepo{active: map[string]map[string]bool{u.ID: {tenantID: true}}}
	svc := &adminService{tenants: tenantRepo, users: users, memberships: memberships, audit: &auditCapture{}}

	require.NoError(t, svc.AddMemberByEmail(context.Background(), tenantID, "dup@example.com", ""))
	require.Len(t, memberships.byUser[u.ID], 0)
}

func TestAdminService_CreateTenant_Validation(t *testing.T) {
	svc := ProvideAdminService(testConfig(), nil, newAdminTenantTestRepo(), nil, nil, nil, nil, nil, nil, nil, nil, nil, &auditCapture{}, nil, nil)
	_, err := svc.CreateTenant(context.Background(), nil)
	require.Error(t, err)
	_, err = svc.CreateTenant(context.Background(), &dtos.AdminCreateTenantRequest{Name: "  "})
	require.Error(t, err)
}

func TestAdminService_UpdateTenant_Validation(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "t1"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "A"}
	svc := &adminService{tenants: tenantRepo}
	_, err := svc.UpdateTenant(context.Background(), tenantID, nil)
	require.Error(t, err)
	empty := ""
	_, err = svc.UpdateTenant(context.Background(), tenantID, &dtos.AdminUpdateTenantRequest{Name: &empty})
	require.Error(t, err)
}

func TestAdminService_ForceLogoutUser_NotFound(t *testing.T) {
	svc, _, _, _, _ := newAdminUserTestService(&adminUserTestRepo{users: map[string]*models.User{}})
	err := svc.ForceLogoutUser(context.Background(), "actor", "missing")
	require.Error(t, err)
}

func TestUserService_Login_Flow(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("flow@example.com", "secret123")
	tenantID := "tenant-flow"
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			u.ID: {{UserID: u.ID, TenantID: tenantID}},
		},
		active: map[string]map[string]bool{u.ID: {tenantID: true}},
	}
	cfg := testConfig()
	tokenSvc, _ := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	svc := ProvideUserService(users, memberships, newRefreshTokenTestRepo(), ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships), cfg, tokenSvc, &auditCapture{})

	resp, sel, err := svc.Login(context.Background(), &dtos.LoginRequest{Email: "flow@example.com", Password: "secret123", TenantID: tenantID}, "")
	require.NoError(t, err)
	require.Nil(t, sel)
	require.NotEmpty(t, resp.AccessToken)
}

func TestUserService_AuthenticateUser_NilPasswordCredential(t *testing.T) {
	users := newUserTestRepo()
	u := &models.User{BaseModel: models.NewBaseModel(), Email: "x@example.com", EmailLower: "x@example.com", Status: constants.UserStatusActive}
	users.users[u.ID] = u
	users.byEmail[u.EmailLower] = u
	svc := newUserTestService(t, users, &stubMembershipRepo{active: map[string]map[string]bool{}})
	_, err := svc.AuthenticateUser(context.Background(), &dtos.LoginRequest{Email: "x@example.com", Password: "p"})
	require.Error(t, err)
}
