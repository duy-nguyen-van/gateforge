package domains

import "strings"

// OAuthRedirectError is used by the authorize handler: RedirectTo is a full URL for 302 (success or OAuth error); if empty, respond with JSON (handler maps Code to HTTP status).
type OAuthRedirectError struct {
	RedirectTo  string
	Code        string
	Description string
	State       string
}

// OAuthTokenError is RFC 6749 token endpoint error (400 + JSON body).
type OAuthTokenError struct {
	Code        string
	Description string
}

func (e *OAuthTokenError) Error() string {
	if e.Description != "" {
		return e.Description
	}
	return e.Code
}

// OIDCIDTokenClaims holds standard claims parsed from an upstream OIDC ID token.
type OIDCIDTokenClaims struct {
	Nonce         string `json:"nonce"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Sub           string `json:"sub"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Name          string `json:"name"`
}

// OIDCUserClaims is a normalized identity after upstream OIDC ID token verification.
type OIDCUserClaims struct {
	Sub           string
	Email         string
	EmailVerified bool
	GivenName     string
	FamilyName    string
	Name          string
}

// SplitDisplayName derives first/last name from standard OIDC profile-style claims.
func (c *OIDCUserClaims) SplitDisplayName() (first, last string) {
	if c == nil {
		return "", ""
	}
	first = strings.TrimSpace(c.GivenName)
	last = strings.TrimSpace(c.FamilyName)
	if first == "" && last == "" && strings.TrimSpace(c.Name) != "" {
		parts := strings.SplitN(strings.TrimSpace(c.Name), " ", 2)
		first = parts[0]
		if len(parts) > 1 {
			last = parts[1]
		}
	}
	return first, last
}

// OIDCTokenResponse is the successful token endpoint JSON.
type OIDCTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	IDToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}
