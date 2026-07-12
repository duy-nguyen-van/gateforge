package dtos

import (
	"net/url"
	"strings"
)

// AuthorizeQuery holds OAuth 2.0 / OIDC authorization request query parameters (RFC 6749, OIDC Core).
type AuthorizeQuery struct {
	State               string
	ResponseType        string
	ClientID            string
	RedirectURI         string
	Scope               string
	Nonce               string
	CodeChallenge       string
	CodeChallengeMethod string
}

// NewAuthorizeQueryFromURLValues parses GET /authorize query parameters.
func NewAuthorizeQueryFromURLValues(v url.Values) AuthorizeQuery {
	return AuthorizeQuery{
		State:               v.Get("state"),
		ResponseType:        v.Get("response_type"),
		ClientID:            v.Get("client_id"),
		RedirectURI:         v.Get("redirect_uri"),
		Scope:               strings.TrimSpace(v.Get("scope")),
		Nonce:               v.Get("nonce"),
		CodeChallenge:       v.Get("code_challenge"),
		CodeChallengeMethod: v.Get("code_challenge_method"),
	}
}

// OpenIDConfigurationResponse is the response for the OpenID Configuration endpoint.
type OpenIDConfigurationResponse struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	JWKSURI                           string   `json:"jwks_uri"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	ScopesSupported                   []string `json:"scopes_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
}
