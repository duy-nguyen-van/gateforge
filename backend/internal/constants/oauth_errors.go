package constants

// OAuth / OIDC error values for JSON and redirect query parameters (RFC 6749, OpenID Connect).
// Names follow the app constants style; values MUST stay exactly as registered for clients.

const (
	OAuthInvalidRequest          = "invalid_request"
	OAuthInvalidClient           = "invalid_client"
	OAuthInvalidGrant            = "invalid_grant"
	OAuthInvalidScope            = "invalid_scope"
	OAuthInvalidToken            = "invalid_token"
	OAuthUnsupportedGrantType    = "unsupported_grant_type"
	OAuthUnsupportedResponseType = "unsupported_response_type"
	OAuthUnauthorizedClient      = "unauthorized_client"
	OAuthServerError             = "server_error"
)
