package dtos

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

func TestNewAdminUserResponse(t *testing.T) {
	createdAt := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	userID := uuid.Must(uuid.NewV7()).String()
	tenantID := uuid.Must(uuid.NewV7()).String()

	t.Run("MFA disabled when no TOTP record", func(t *testing.T) {
		user := &models.User{
			BaseModel: models.BaseModel{ID: userID, CreatedAt: createdAt},
			Email:     "user@example.com",
			FirstName: "Jane",
			LastName:  "Doe",
			Status:    constants.UserStatusActive,
		}

		resp := NewAdminUserResponse(user, tenantID)
		require.Equal(t, userID, resp.ID)
		require.Equal(t, "user@example.com", resp.Email)
		require.Equal(t, "Jane", resp.FirstName)
		require.Equal(t, "Doe", resp.LastName)
		require.Equal(t, string(constants.UserStatusActive), resp.Status)
		require.Equal(t, tenantID, resp.TenantID)
		require.False(t, resp.MFAEnabled)
		require.Equal(t, createdAt, resp.CreatedAt)
	})

	t.Run("MFA disabled when TOTP not enabled", func(t *testing.T) {
		user := &models.User{
			BaseModel:   models.BaseModel{ID: userID, CreatedAt: createdAt},
			UserMFATOTP: &models.UserMFATOTP{Enabled: false},
		}
		resp := NewAdminUserResponse(user, tenantID)
		require.False(t, resp.MFAEnabled)
	})

	t.Run("MFA enabled when TOTP enabled", func(t *testing.T) {
		user := &models.User{
			BaseModel:   models.BaseModel{ID: userID, CreatedAt: createdAt},
			UserMFATOTP: &models.UserMFATOTP{Enabled: true},
		}
		resp := NewAdminUserResponse(user, tenantID)
		require.True(t, resp.MFAEnabled)
	})
}

func TestNewAdminTenantResponse(t *testing.T) {
	createdAt := time.Date(2024, 1, 15, 8, 30, 0, 0, time.UTC)
	tenant := &models.Tenant{
		BaseModel: models.BaseModel{ID: uuid.Must(uuid.NewV7()).String(), CreatedAt: createdAt},
		Name:      "Acme Corp",
		Domain:    "acme.example.com",
	}

	resp := NewAdminTenantResponse(tenant, 42)
	require.Equal(t, tenant.ID, resp.ID)
	require.Equal(t, "Acme Corp", resp.Name)
	require.Equal(t, "acme.example.com", resp.Domain)
	require.Equal(t, int64(42), resp.UserCount)
	require.Equal(t, createdAt, resp.CreatedAt)
}

func TestNewAdminTenantMemberResponse(t *testing.T) {
	joinedAt := time.Date(2024, 3, 10, 9, 0, 0, 0, time.UTC)
	membership := &models.TenantMembership{
		BaseModel: models.BaseModel{CreatedAt: joinedAt},
		UserID:    uuid.Must(uuid.NewV7()).String(),
		Role:      constants.TenantMembershipRoleAdmin,
		Status:    constants.TenantMembershipStatusActive,
	}

	t.Run("without preloaded user", func(t *testing.T) {
		resp := NewAdminTenantMemberResponse(membership)
		require.Equal(t, membership.UserID, resp.UserID)
		require.Equal(t, string(constants.TenantMembershipRoleAdmin), resp.Role)
		require.Equal(t, string(constants.TenantMembershipStatusActive), resp.Status)
		require.Equal(t, joinedAt, resp.JoinedAt)
		require.Empty(t, resp.Email)
		require.Empty(t, resp.FirstName)
		require.Empty(t, resp.LastName)
	})

	t.Run("with preloaded user", func(t *testing.T) {
		membership.User = &models.User{
			Email:     "member@example.com",
			FirstName: "Alex",
			LastName:  "Smith",
		}
		resp := NewAdminTenantMemberResponse(membership)
		require.Equal(t, "member@example.com", resp.Email)
		require.Equal(t, "Alex", resp.FirstName)
		require.Equal(t, "Smith", resp.LastName)
	})
}

func TestNewAdminClientResponse(t *testing.T) {
	createdAt := time.Date(2024, 5, 20, 14, 0, 0, 0, time.UTC)
	clientID := uuid.Must(uuid.NewV7()).String()
	tenantID := uuid.Must(uuid.NewV7()).String()

	t.Run("confidential client with secret", func(t *testing.T) {
		client := &models.Client{
			BaseModel:    models.BaseModel{ID: clientID, CreatedAt: createdAt},
			TenantID:     tenantID,
			ClientID:     "my-app",
			Name:         "My App",
			IsPublic:     false,
			GrantTypes:   pq.StringArray{"authorization_code", "refresh_token"},
			RedirectUris: pq.StringArray{"https://app.example/callback"},
			Scopes:       pq.StringArray{"openid", "profile"},
			ClientSecret: "hashed-secret",
		}

		resp := NewAdminClientResponse(client)
		require.Equal(t, clientID, resp.ID)
		require.Equal(t, tenantID, resp.TenantID)
		require.Equal(t, "my-app", resp.ClientID)
		require.Equal(t, "My App", resp.Name)
		require.False(t, resp.IsPublic)
		require.Equal(t, []string{"authorization_code", "refresh_token"}, resp.GrantTypes)
		require.Equal(t, []string{"https://app.example/callback"}, resp.RedirectUris)
		require.Equal(t, []string{"openid", "profile"}, resp.Scopes)
		require.True(t, resp.ClientSecretSet)
		require.Equal(t, createdAt, resp.CreatedAt)
	})

	t.Run("public client without secret", func(t *testing.T) {
		client := &models.Client{
			BaseModel: models.BaseModel{ID: clientID, CreatedAt: createdAt},
			IsPublic:  true,
		}
		resp := NewAdminClientResponse(client)
		require.True(t, resp.IsPublic)
		require.False(t, resp.ClientSecretSet)
	})
}
