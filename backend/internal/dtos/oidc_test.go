package dtos

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAuthorizeQueryFromURLValues(t *testing.T) {
	values := url.Values{
		"state":                 {"opaque-state"},
		"response_type":         {"code"},
		"client_id":             {"my-client"},
		"redirect_uri":          {"https://app.example/callback"},
		"scope":                 {"  openid profile  "},
		"nonce":                 {"nonce-123"},
		"code_challenge":        {"challenge"},
		"code_challenge_method": {"S256"},
	}

	q := NewAuthorizeQueryFromURLValues(values)
	require.Equal(t, "opaque-state", q.State)
	require.Equal(t, "code", q.ResponseType)
	require.Equal(t, "my-client", q.ClientID)
	require.Equal(t, "https://app.example/callback", q.RedirectURI)
	require.Equal(t, "openid profile", q.Scope)
	require.Equal(t, "nonce-123", q.Nonce)
	require.Equal(t, "challenge", q.CodeChallenge)
	require.Equal(t, "S256", q.CodeChallengeMethod)
}

func TestNewAuthorizeQueryFromURLValues_Empty(t *testing.T) {
	q := NewAuthorizeQueryFromURLValues(url.Values{})
	require.Empty(t, q.State)
	require.Empty(t, q.ResponseType)
	require.Empty(t, q.ClientID)
	require.Empty(t, q.RedirectURI)
	require.Empty(t, q.Scope)
	require.Empty(t, q.Nonce)
	require.Empty(t, q.CodeChallenge)
	require.Empty(t, q.CodeChallengeMethod)
}

func TestNewAuthorizeQueryFromURLValues_MissingKeys(t *testing.T) {
	values := url.Values{"client_id": {"only-client"}}
	q := NewAuthorizeQueryFromURLValues(values)
	require.Equal(t, "only-client", q.ClientID)
	require.Empty(t, q.State)
}
