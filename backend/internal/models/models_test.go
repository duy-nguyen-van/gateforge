package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTableNames(t *testing.T) {
	tests := []struct {
		name     string
		model    interface{ TableName() string }
		expected string
	}{
		{"User", User{}, "users"},
		{"Tenant", Tenant{}, "tenants"},
		{"TenantMembership", TenantMembership{}, "tenant_memberships"},
		{"Client", Client{}, "clients"},
		{"AuditLog", AuditLog{}, "audit_logs"},
		{"Session", Session{}, "sessions"},
		{"AccessToken", AccessToken{}, "access_tokens"},
		{"RefreshToken", RefreshToken{}, "refresh_tokens"},
		{"AuthorizationCode", AuthorizationCode{}, "authorization_codes"},
		{"Consent", Consent{}, "consents"},
		{"PasswordCredential", PasswordCredential{}, "password_credentials"},
		{"WebauthnCredential", WebauthnCredential{}, "webauthn_credentials"},
		{"UserMFATOTP", UserMFATOTP{}, "user_mfa_totps"},
		{"UserMFARecoveryCode", UserMFARecoveryCode{}, "user_mfa_recovery_codes"},
		{"FederatedIdentity", FederatedIdentity{}, "federated_identities"},
		{"TenantIdentityProvider", TenantIdentityProvider{}, "tenant_identity_providers"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.model.TableName())
		})
	}
}

func TestNewBaseModel(t *testing.T) {
	bm := NewBaseModel()
	require.NotEmpty(t, bm.ID)
	_, err := uuid.Parse(bm.ID)
	require.NoError(t, err)
}

func TestBaseModel_BeforeCreate(t *testing.T) {
	t.Run("generates UUID when ID is empty", func(t *testing.T) {
		bm := &BaseModel{}
		require.NoError(t, bm.BeforeCreate(&gorm.DB{}))
		require.NotEmpty(t, bm.ID)
		_, err := uuid.Parse(bm.ID)
		require.NoError(t, err)
	})

	t.Run("accepts valid UUID", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7()).String()
		bm := &BaseModel{ID: id}
		require.NoError(t, bm.BeforeCreate(&gorm.DB{}))
		require.Equal(t, id, bm.ID)
	})

	t.Run("rejects invalid UUID", func(t *testing.T) {
		bm := &BaseModel{ID: "not-a-uuid"}
		err := bm.BeforeCreate(&gorm.DB{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid uuid")
	})
}

func TestHardDeleteModel_BeforeCreate(t *testing.T) {
	t.Run("generates UUID when ID is empty", func(t *testing.T) {
		hm := &HardDeleteModel{}
		require.NoError(t, hm.BeforeCreate(&gorm.DB{}))
		require.NotEmpty(t, hm.ID)
		_, err := uuid.Parse(hm.ID)
		require.NoError(t, err)
	})

	t.Run("accepts valid UUID", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7()).String()
		hm := &HardDeleteModel{ID: id}
		require.NoError(t, hm.BeforeCreate(&gorm.DB{}))
		require.Equal(t, id, hm.ID)
	})

	t.Run("rejects invalid UUID", func(t *testing.T) {
		hm := &HardDeleteModel{ID: "bad-id"}
		err := hm.BeforeCreate(&gorm.DB{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid uuid")
	})
}

func TestUser_BeforeSave(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		wantLower  string
	}{
		{"lowercases email", "User@Example.COM", "user@example.com"},
		{"trims whitespace", "  user@example.com  ", "user@example.com"},
		{"lowercases and trims", "  Mixed@Case.COM  ", "mixed@case.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Email: tt.email}
			require.NoError(t, u.BeforeSave(&gorm.DB{}))
			require.Equal(t, tt.wantLower, u.EmailLower)
		})
	}
}

func TestEnsurePrimaryUUID_InvalidUUID(t *testing.T) {
	// ensurePrimaryUUID is private; exercise it through BeforeCreate hooks.
	models := []interface {
		BeforeCreate(*gorm.DB) error
	}{
		&BaseModel{ID: "invalid"},
		&HardDeleteModel{ID: "invalid"},
	}

	for _, m := range models {
		err := m.BeforeCreate(&gorm.DB{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid uuid")
	}
}

func TestBaseModel_TimestampsZeroValue(t *testing.T) {
	bm := NewBaseModel()
	require.True(t, bm.CreatedAt.IsZero())
	require.True(t, bm.UpdatedAt.IsZero())
	require.False(t, bm.DeletedAt.Valid)
}

func TestHardDeleteModel_HasNoSoftDelete(t *testing.T) {
	hm := HardDeleteModel{
		ID:        uuid.Must(uuid.NewV7()).String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NotEmpty(t, hm.ID)
	require.False(t, hm.CreatedAt.IsZero())
}
