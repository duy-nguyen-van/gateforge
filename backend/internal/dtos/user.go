package dtos

import (
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

// UserResponse is returned for the current user (e.g. GET /me).
type UserResponse struct {
	ID              string          `json:"id"`
	Email           string          `json:"email"`
	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	EmailVerified   bool            `json:"email_verified"`
	IsPlatformAdmin bool            `json:"is_platform_admin"`
	MFAEnabled      bool            `json:"mfa_enabled"`
	ActiveTenantID  string          `json:"active_tenant_id,omitempty"`
	Tenants         []TenantSummary `json:"tenants,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// UpdateProfileRequest is the body for PATCH /me.
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,max=100"`
}

func NewUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		EmailVerified:   user.EmailVerified,
		IsPlatformAdmin: user.IsPlatformAdmin,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}
