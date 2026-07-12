package dtos

// RegisterRequest is the body for POST /register.
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email" example:"user@example.com"`
	Password  string `json:"password" validate:"required,min=8,max=128" example:"secretpassword"`
	FirstName string `json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,max=100"`
	TenantID  string `json:"tenant_id,omitempty" validate:"omitempty,uuid"`
}

// LoginRequest is the body for POST /login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	// TenantID optional explicit tenant context.
	TenantID string `json:"tenant_id,omitempty" validate:"omitempty,uuid"`
	// ReturnTo when set (OIDC browser flow): creates a session cookie and redirects here with 302 — no access_token in the response body.
	ReturnTo string `json:"return_to,omitempty" validate:"omitempty"`
	// RememberMe extends the iam_session cookie lifetime (SSO_SESSION_REMEMBER_TTL). API login and POST /oidc/login both honor this.
	RememberMe bool `json:"remember_me,omitempty" form:"remember_me"`
}

// LoginResponse is returned after successful login or token refresh (OAuth2-style token bundle).
type LoginResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int64  `json:"expires_in"`         // access token lifetime (seconds)
	RefreshExpiresIn int64  `json:"refresh_expires_in"` // refresh token lifetime (seconds)
	ActiveTenantID   string `json:"active_tenant_id,omitempty"`
}

// TenantSummary is a tenant the user can access.
type TenantSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Role   string `json:"role"`
}

// TenantSelectionResponse is returned when auth succeeded but tenant must be chosen.
type TenantSelectionResponse struct {
	SelectionRequired bool            `json:"selection_required"`
	Tenants           []TenantSummary `json:"tenants"`
	SelectionToken    string          `json:"selection_token"`
	ExpiresIn         int64           `json:"expires_in"`
}

// TenantSelectRequest is the body for POST /tenants/select.
type TenantSelectRequest struct {
	SelectionToken string `json:"selection_token" validate:"required"`
	TenantID       string `json:"tenant_id" validate:"required,uuid"`
	RememberMe     bool   `json:"remember_me,omitempty"`
}

// TenantSwitchRequest is the body for POST /tenants/switch.
type TenantSwitchRequest struct {
	TenantID string `json:"tenant_id" validate:"required,uuid"`
}

// RefreshTokenRequest is the body for POST /refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
