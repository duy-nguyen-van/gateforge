package gateforge

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// PKCE holds a generated Proof Key for Code Exchange pair.
type PKCE struct {
	Verifier  string
	Challenge string
	Method    string // always "S256"
}

// PKCEGenerate creates a new S256 PKCE verifier/challenge pair.
func PKCEGenerate() (PKCE, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return PKCE{}, fmt.Errorf("pkce: generate verifier: %w", err)
	}
	verifier := base64.RawURLEncoding.EncodeToString(buf)
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	return PKCE{
		Verifier:  verifier,
		Challenge: challenge,
		Method:    "S256",
	}, nil
}

// AuthorizeParams are query parameters for an OIDC authorization request.
type AuthorizeParams struct {
	ClientID            string
	RedirectURI         string
	Scope               string
	State               string
	Nonce               string
	CodeChallenge       string
	CodeChallengeMethod string
	ResponseType        string
	// Extra are additional query parameters (e.g. login_hint, prompt).
	Extra map[string]string
}

// BuildAuthorizeURL constructs an OIDC/OAuth2 authorize URL.
// issuerOrAuthorizeURL may be an issuer base (appends /authorize) or a full
// authorize endpoint URL.
func BuildAuthorizeURL(issuerOrAuthorizeURL string, params AuthorizeParams) (string, error) {
	if strings.TrimSpace(issuerOrAuthorizeURL) == "" {
		return "", fmt.Errorf("auth: issuer or authorize URL is required")
	}
	if params.ClientID == "" {
		return "", fmt.Errorf("auth: client_id is required")
	}
	if params.RedirectURI == "" {
		return "", fmt.Errorf("auth: redirect_uri is required")
	}

	base := strings.TrimRight(issuerOrAuthorizeURL, "/")
	authorizeURL := base
	if !strings.HasSuffix(base, "/authorize") {
		authorizeURL = base + "/authorize"
	}

	u, err := url.Parse(authorizeURL)
	if err != nil {
		return "", fmt.Errorf("auth: parse authorize URL: %w", err)
	}

	q := u.Query()
	q.Set("client_id", params.ClientID)
	q.Set("redirect_uri", params.RedirectURI)
	responseType := params.ResponseType
	if responseType == "" {
		responseType = "code"
	}
	q.Set("response_type", responseType)
	if params.Scope != "" {
		q.Set("scope", params.Scope)
	}
	if params.State != "" {
		q.Set("state", params.State)
	}
	if params.Nonce != "" {
		q.Set("nonce", params.Nonce)
	}
	if params.CodeChallenge != "" {
		q.Set("code_challenge", params.CodeChallenge)
		method := params.CodeChallengeMethod
		if method == "" {
			method = "S256"
		}
		q.Set("code_challenge_method", method)
	}
	for k, v := range params.Extra {
		if k == "" {
			continue
		}
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
