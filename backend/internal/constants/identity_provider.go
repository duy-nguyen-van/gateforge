package constants

import (
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

// External identity provider IDs for federation.
const (
	IdentityProviderGoogle = "google"
)

// IdentityProviderSpec describes a built-in upstream OIDC IdP.
type IdentityProviderSpec struct {
	ID                   string
	DisplayName          string
	IssuerURL            string
	Scopes               []string
	AuthURLParams        map[string]string
	RequireVerifiedEmail bool
	SetupConsoleURL      string
}

// SupportedIdentityProviders is the code-defined catalog of federation IdPs.
var SupportedIdentityProviders = []IdentityProviderSpec{
	{
		ID:                   IdentityProviderGoogle,
		DisplayName:          "Google",
		IssuerURL:            "https://accounts.google.com",
		Scopes:               []string{oidc.ScopeOpenID, "email", "profile"},
		AuthURLParams:        map[string]string{"prompt": "select_account"},
		RequireVerifiedEmail: true,
		SetupConsoleURL:      "https://console.cloud.google.com/apis/credentials",
	},
}

// IdentityProviderByID returns catalog metadata for a provider id.
func IdentityProviderByID(id string) (IdentityProviderSpec, bool) {
	key := strings.ToLower(strings.TrimSpace(id))
	for _, spec := range SupportedIdentityProviders {
		if strings.EqualFold(spec.ID, key) {
			return spec, true
		}
	}
	return IdentityProviderSpec{}, false
}

// IsSupportedIdentityProvider reports whether id is in the catalog.
func IsSupportedIdentityProvider(id string) bool {
	_, ok := IdentityProviderByID(id)
	return ok
}
