package dtos

import (
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

// AdminStatsResponse is returned by GET /admin/stats.
type AdminStatsResponse struct {
	TotalUsers        int64   `json:"total_users"`
	MFAEnabledCount   int64   `json:"mfa_enabled_count"`
	MFAEnabledPercent float64 `json:"mfa_enabled_percent"`
	ActiveSessions    int64   `json:"active_sessions"`
}

// AdminUserResponse is a user row for the admin console.
type AdminUserResponse struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Status     string    `json:"status"`
	TenantID   string    `json:"tenant_id"`
	MFAEnabled bool      `json:"mfa_enabled"`
	CreatedAt  time.Time `json:"created_at"`
}

// AdminUserMembershipResponse is a tenant membership on the user detail view.
type AdminUserMembershipResponse struct {
	TenantID   string `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
	Role       string `json:"role"`
	Status     string `json:"status"`
}

// AdminUserDetailResponse is returned by GET /admin/users/:userId.
type AdminUserDetailResponse struct {
	ID              string                        `json:"id"`
	Email           string                        `json:"email"`
	FirstName       string                        `json:"first_name"`
	LastName        string                        `json:"last_name"`
	Status          string                        `json:"status"`
	MFAEnabled      bool                          `json:"mfa_enabled"`
	IsPlatformAdmin bool                          `json:"is_platform_admin"`
	PasskeyCount    int                           `json:"passkey_count"`
	ActiveSessions  int64                         `json:"active_sessions"`
	Memberships     []AdminUserMembershipResponse `json:"memberships"`
	CreatedAt       time.Time                     `json:"created_at"`
}

// AdminClientUsageResponse is returned by GET /admin/clients/:clientId/usage.
type AdminClientUsageResponse struct {
	ClientID            string     `json:"client_id"`
	TotalRefreshTokens  int64      `json:"total_refresh_tokens"`
	ActiveRefreshTokens int64      `json:"active_refresh_tokens"`
	LastTokenIssuedAt   *time.Time `json:"last_token_issued_at,omitempty"`
	AuthorizeEvents30d  int64      `json:"authorize_events_30d"`
	TokenIssueEvents30d int64      `json:"token_issue_events_30d"`
}

// AdminTenantResponse is a tenant row for the admin console.
type AdminTenantResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	UserCount int64     `json:"user_count"`
	CreatedAt time.Time `json:"created_at"`
}

// AdminClientResponse is an OAuth client row for the admin console.
type AdminClientResponse struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	ClientID        string    `json:"client_id"`
	Name            string    `json:"name"`
	IsPublic        bool      `json:"is_public"`
	GrantTypes      []string  `json:"grant_types"`
	RedirectUris    []string  `json:"redirect_uris"`
	Scopes          []string  `json:"scopes"`
	ClientSecretSet bool      `json:"client_secret_set"`
	CreatedAt       time.Time `json:"created_at"`
}

// AdminCreateClientRequest registers a new OAuth client for a tenant.
type AdminCreateClientRequest struct {
	TenantID     string   `json:"tenant_id" validate:"required,uuid"`
	ClientID     string   `json:"client_id,omitempty" validate:"omitempty,min=1,max=255"`
	Name         string   `json:"name" validate:"required,min=1,max=255"`
	IsPublic     bool     `json:"is_public"`
	RedirectUris []string `json:"redirect_uris" validate:"required,min=1,dive,url"`
	GrantTypes   []string `json:"grant_types,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

// AdminUpdateClientRequest updates mutable OAuth client fields.
type AdminUpdateClientRequest struct {
	Name         *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	RedirectUris []string `json:"redirect_uris,omitempty" validate:"omitempty,min=1,dive,url"`
	GrantTypes   []string `json:"grant_types,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
	IsPublic     *bool    `json:"is_public,omitempty"`
	ClientSecret *string  `json:"client_secret,omitempty"`
}

// AdminCreateClientResponse is returned once on client creation; includes the plaintext secret for confidential clients.
type AdminCreateClientResponse struct {
	AdminClientResponse
	ClientSecret string `json:"client_secret,omitempty"`
}

// AdminIdentityProviderResponse describes a federation provider for a tenant.
type AdminIdentityProviderResponse struct {
	Provider             string `json:"provider"`
	Name                 string `json:"name"`
	Enabled              bool   `json:"enabled"`
	Configured           bool   `json:"configured"`
	TenantID             string `json:"tenant_id"`
	OAuthClientID        string `json:"oauth_client_id,omitempty"`
	OAuthClientSecretSet bool   `json:"oauth_client_secret_set"`
	RedirectURI          string `json:"redirect_uri,omitempty"`
	SetupConsoleURL      string `json:"setup_console_url,omitempty"`
}

// PatchIdentityProviderRequest configures an upstream OAuth/OIDC provider for a tenant.
type PatchIdentityProviderRequest struct {
	Enabled           *bool  `json:"enabled"`
	OAuthClientID     string `json:"oauth_client_id,omitempty"`
	OAuthClientSecret string `json:"oauth_client_secret,omitempty"`
}

// PublicFederationProviderResponse is exposed on the login page.
type PublicFederationProviderResponse struct {
	Provider string `json:"provider"`
	Name     string `json:"name"`
}

// AdminCreateTenantRequest creates a new tenant.
type AdminCreateTenantRequest struct {
	Name   string `json:"name" validate:"required,min=1,max=255"`
	Domain string `json:"domain,omitempty" validate:"omitempty,max=255"`
}

// AdminUpdateTenantRequest updates tenant fields.
type AdminUpdateTenantRequest struct {
	Name   *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Domain *string `json:"domain,omitempty" validate:"omitempty,max=255"`
}

// AdminTenantMemberResponse is a tenant member row for the admin console.
type AdminTenantMemberResponse struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	JoinedAt  time.Time `json:"joined_at"`
}

// AdminAddMemberRequest adds an existing user to a tenant.
type AdminAddMemberRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role,omitempty" validate:"omitempty,oneof=member admin"`
}

func NewAdminUserResponse(user *models.User, tenantID string) *AdminUserResponse {
	mfaEnabled := user.UserMFATOTP != nil && user.UserMFATOTP.Enabled
	return &AdminUserResponse{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Status:     string(user.Status),
		TenantID:   tenantID,
		MFAEnabled: mfaEnabled,
		CreatedAt:  user.CreatedAt,
	}
}

func NewAdminTenantResponse(tenant *models.Tenant, userCount int64) *AdminTenantResponse {
	return &AdminTenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Domain:    tenant.Domain,
		UserCount: userCount,
		CreatedAt: tenant.CreatedAt,
	}
}

func NewAdminTenantMemberResponse(m *models.TenantMembership) *AdminTenantMemberResponse {
	resp := &AdminTenantMemberResponse{
		UserID:   m.UserID,
		Role:     string(m.Role),
		Status:   string(m.Status),
		JoinedAt: m.CreatedAt,
	}
	if m.User != nil {
		resp.Email = m.User.Email
		resp.FirstName = m.User.FirstName
		resp.LastName = m.User.LastName
	}
	return resp
}

func NewAdminClientResponse(client *models.Client) *AdminClientResponse {
	return &AdminClientResponse{
		ID:              client.ID,
		TenantID:        client.TenantID,
		ClientID:        client.ClientID,
		Name:            client.Name,
		IsPublic:        client.IsPublic,
		GrantTypes:      []string(client.GrantTypes),
		RedirectUris:    []string(client.RedirectUris),
		Scopes:          []string(client.Scopes),
		ClientSecretSet: client.ClientSecret != "",
		CreatedAt:       client.CreatedAt,
	}
}
