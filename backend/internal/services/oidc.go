package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/cache"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/crypto"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"go.opentelemetry.io/otel/attribute"
)

// OIDCService implements authorization code + PKCE and OIDC token issuance.
type OIDCService interface {
	Authorize(ctx context.Context, userID string, q *dtos.AuthorizeQuery) (successRedirect string, err *domains.OAuthRedirectError)
	AuthorizationCodeToken(ctx context.Context, tenantID, clientID, clientSecret string, form url.Values) (*domains.OIDCTokenResponse, *domains.OAuthTokenError)
	UserInfo(ctx context.Context, accessToken string) (map[string]any, *domains.OAuthTokenError)
	// Introspect validates an OIDC access or refresh token (RFC 7662). Requires confidential client auth.
	Introspect(ctx context.Context, clientID, clientSecret, token, tokenTypeHint string) (*dtos.TokenIntrospectionResponse, *domains.OAuthTokenError)
	OpenIDIssuer() string
}

type oidcService struct {
	cfg            *config.Config
	oidc           *auth.OIDCSigner
	clients        repositories.ClientRepository
	authCodes      repositories.AuthorizationCodeRepository
	users          repositories.UserRepository
	refreshTokens  repositories.RefreshTokenRepository
	membershipRepo repositories.TenantMembershipRepository
	audit          AuditService
	clientCache    cache.Cache
}

// ProvideOIDCService wires Phase 2 OIDC.
func ProvideOIDCService(
	cfg *config.Config,
	oidc *auth.OIDCSigner,
	clients repositories.ClientRepository,
	authCodes repositories.AuthorizationCodeRepository,
	users repositories.UserRepository,
	refreshTokens repositories.RefreshTokenRepository,
	membershipRepo repositories.TenantMembershipRepository,
	audit AuditService,
	clientCache cache.Cache,
) OIDCService {
	return &oidcService{
		cfg:            cfg,
		oidc:           oidc,
		clients:        clients,
		authCodes:      authCodes,
		users:          users,
		refreshTokens:  refreshTokens,
		membershipRepo: membershipRepo,
		audit:          audit,
		clientCache:    clientCache,
	}
}

func (s *oidcService) OpenIDIssuer() string {
	if s.cfg.AppBaseURL != "" {
		return s.cfg.AppBaseURL
	}
	return "http://localhost:3000"
}

func oauthFinishErr(code, description string) error {
	if code == "" {
		return nil
	}
	return fmt.Errorf("%s: %s", code, description)
}

// Authorize validates the OAuth2 request and returns a redirect URL with an authorization code.
// userID is resolved by the handler from the browser session cookie (OIDC login at POST /oidc/login).
func (s *oidcService) Authorize(ctx context.Context, userID string, q *dtos.AuthorizeQuery) (redirect string, oauthErr *domains.OAuthRedirectError) {
	attrs := []attribute.KeyValue{attribute.String("user_id", userID)}
	if q != nil {
		attrs = append(attrs, attribute.String("client_id", q.ClientID))
	}
	ctx, span := monitoring.StartSpan(ctx, oidcTracer, "OIDCService.Authorize", attrs...)
	defer func() {
		code, desc := "", ""
		if oauthErr != nil {
			code, desc = oauthErr.Code, oauthErr.Description
		}
		monitoring.Finish(ctx, span, "OIDCService.Authorize", oauthFinishErr(code, desc))
	}()

	if q == nil {
		return "", &domains.OAuthRedirectError{Code: constants.OAuthInvalidRequest, Description: "missing authorization request"}
	}

	state := q.State
	responseType := q.ResponseType
	clientID := q.ClientID
	redirectURI := q.RedirectURI
	rawScope := q.Scope
	nonce := q.Nonce
	challenge := q.CodeChallenge
	method := q.CodeChallengeMethod

	client, err := s.clients.GetByClientID(ctx, clientID)
	if err != nil {
		return "", &domains.OAuthRedirectError{Code: constants.OAuthInvalidRequest, Description: "unknown client_id", State: state}
	}
	tenantID := client.TenantID

	ok, err := s.membershipRepo.ExistsActive(ctx, userID, tenantID)
	if err != nil || !ok {
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:       constants.AuditActionOIDCAuthorize,
			Result:       constants.AuditResultDenied,
			ActorType:    constants.AuditActorTypeUser,
			ActorID:      userID,
			TenantID:     tenantID,
			ResourceType: constants.AuditResourceTypeClient,
			ResourceName: clientID,
		})
		return "", &domains.OAuthRedirectError{Code: constants.OAuthUnauthorizedClient, Description: "user not authorized for this client", State: state}
	}

	if !redirectAllowed(client, redirectURI) {
		return "", &domains.OAuthRedirectError{Code: constants.OAuthInvalidRequest, Description: "invalid redirect_uri", State: state}
	}

	if responseType != "code" {
		return oauthRedirectErr(redirectURI, state, constants.OAuthUnsupportedResponseType, "only code is supported")
	}

	if !grantAllowed(client, "authorization_code") {
		return oauthRedirectErr(redirectURI, state, constants.OAuthUnauthorizedClient, "authorization_code grant not allowed")
	}

	if rawScope == "" {
		rawScope = "openid"
	}
	if scopeErr := validateScopes(client, rawScope); scopeErr != "" {
		return oauthRedirectErr(redirectURI, state, constants.OAuthInvalidScope, scopeErr)
	}

	if redirectErr := authorizePKCEError(client, challenge, method, redirectURI, state); redirectErr != nil {
		return "", redirectErr
	}

	return s.issueAuthorizationCode(ctx, userID, tenantID, client, clientID, redirectURI, rawScope, state, nonce, challenge, method)
}

func authorizePKCEError(client *models.Client, challenge, method, redirectURI, state string) *domains.OAuthRedirectError {
	publicClient := client.IsPublic || strings.TrimSpace(client.ClientSecret) == ""
	if publicClient {
		if challenge == "" || method == "" {
			_, err := oauthRedirectErr(redirectURI, state, constants.OAuthInvalidRequest, "PKCE code_challenge and code_challenge_method are required for public clients")
			return err
		}
		if method != "S256" {
			_, err := oauthRedirectErr(redirectURI, state, constants.OAuthInvalidRequest, "only S256 code_challenge_method is supported")
			return err
		}
		return nil
	}
	if challenge != "" && method != "" && method != "S256" {
		_, err := oauthRedirectErr(redirectURI, state, constants.OAuthInvalidRequest, "only S256 code_challenge_method is supported")
		return err
	}
	return nil
}

func (s *oidcService) issueAuthorizationCode(
	ctx context.Context,
	userID, tenantID string,
	client *models.Client,
	clientID, redirectURI, rawScope, state, nonce, challenge, method string,
) (string, *domains.OAuthRedirectError) {
	codeRaw, _, err := auth.NewOpaqueRefreshToken()
	if err != nil {
		return "", &domains.OAuthRedirectError{Code: constants.OAuthServerError, Description: "failed to issue code", State: state}
	}

	recordID := client.ID
	row := &models.AuthorizationCode{
		Code:                codeRaw,
		TenantID:            tenantID,
		OAuthClientID:       client.ClientID,
		UserID:              userID,
		Scope:               rawScope,
		RedirectURI:         redirectURI,
		CodeChallenge:       challenge,
		CodeChallengeMethod: method,
		Nonce:               nonce,
		ExpiresAt:           time.Now().UTC().Add(s.cfg.OIDCAuthCodeTTL),
		ClientRecordID:      &recordID,
	}

	if err := s.authCodes.Create(ctx, row); err != nil {
		return "", &domains.OAuthRedirectError{Code: constants.OAuthServerError, Description: "failed to persist code", State: state}
	}

	u, err := url.Parse(redirectURI)
	if err != nil {
		return "", &domains.OAuthRedirectError{Code: constants.OAuthInvalidRequest, Description: "invalid redirect_uri", State: state}
	}
	q2 := u.Query()
	q2.Set("code", codeRaw)
	if state != "" {
		q2.Set("state", state)
	}
	u.RawQuery = q2.Encode()
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionOIDCAuthorize,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      userID,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeClient,
		ResourceName: clientID,
		NewValue:     map[string]any{"scope": rawScope},
	})
	return u.String(), nil
}

func oauthRedirectErr(redirectURI, state, code, description string) (string, *domains.OAuthRedirectError) {
	if redirectURI == "" {
		return "", &domains.OAuthRedirectError{Code: code, Description: description, State: state}
	}
	u, err := url.Parse(redirectURI)
	if err != nil {
		return "", &domains.OAuthRedirectError{Code: code, Description: description, State: state}
	}
	q := u.Query()
	q.Set("error", code)
	if description != "" {
		q.Set("error_description", description)
	}
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	return "", &domains.OAuthRedirectError{RedirectTo: u.String(), Code: code, Description: description, State: state}
}

func redirectAllowed(c *models.Client, redirectURI string) bool {
	if redirectURI == "" {
		return false
	}
	for _, u := range c.RedirectUris {
		if strings.TrimSpace(u) == redirectURI {
			return true
		}
	}
	return false
}

func grantAllowed(c *models.Client, grant string) bool {
	for _, g := range c.GrantTypes {
		if strings.TrimSpace(g) == grant {
			return true
		}
	}
	return false
}

func validateScopes(c *models.Client, requested string) string {
	if len(c.Scopes) == 0 {
		return ""
	}
	allowed := make(map[string]struct{}, len(c.Scopes))
	for _, s := range c.Scopes {
		allowed[strings.TrimSpace(s)] = struct{}{}
	}
	for _, p := range strings.Fields(requested) {
		if _, ok := allowed[p]; !ok {
			return "scope not allowed: " + p
		}
	}
	return ""
}

func verifyPKCE(verifier, challenge, method string) bool {
	if verifier == "" || challenge == "" {
		return false
	}
	if method != "S256" {
		return false
	}
	sum := sha256.Sum256([]byte(verifier))
	enc := base64.RawURLEncoding.EncodeToString(sum[:])
	return enc == challenge
}

// AuthorizationCodeToken exchanges an authorization code for tokens (RFC 6749 + PKCE).
func (s *oidcService) AuthorizationCodeToken(ctx context.Context, _ string, formClientID, clientSecret string, form url.Values) (resp *domains.OIDCTokenResponse, oauthErr *domains.OAuthTokenError) {
	grant := form.Get("grant_type")
	ctx, span := monitoring.StartSpan(ctx, oidcTracer, "OIDCService.Token",
		attribute.String("client_id", formClientID),
		attribute.String("grant_type", grant))
	defer func() {
		code, desc := "", ""
		result := "success"
		if oauthErr != nil {
			code, desc = oauthErr.Code, oauthErr.Description
			result = oauthErr.Code
		}
		monitoring.AddOIDCToken(ctx, grant, result)
		monitoring.Finish(ctx, span, "OIDCService.Token", oauthFinishErr(code, desc))
	}()

	switch grant {
	case "authorization_code":
		return s.authorizationCodeGrant(ctx, formClientID, clientSecret, form)
	case "refresh_token":
		return s.refreshTokenGrant(ctx, formClientID, clientSecret, form)
	default:
		return nil, &domains.OAuthTokenError{Code: constants.OAuthUnsupportedGrantType, Description: "unsupported grant_type"}
	}
}

func (s *oidcService) authorizationCodeGrant(ctx context.Context, formClientID, clientSecret string, form url.Values) (*domains.OIDCTokenResponse, *domains.OAuthTokenError) {
	code := strings.TrimSpace(form.Get("code"))
	if code == "" {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidRequest, Description: "code is required"}
	}

	row, err := s.authCodes.Consume(ctx, code)
	if err != nil {
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:    constants.AuditActionOIDCTokenIssue,
			Result:    constants.AuditResultFailure,
			ActorType: constants.AuditActorTypeOAuthClient,
			ActorID:   formClientID,
		})
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "invalid or expired authorization code"}
	}

	redirectURI := strings.TrimSpace(form.Get("redirect_uri"))
	if redirectURI == "" {
		redirectURI = row.RedirectURI
	}
	if redirectURI == "" {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidRequest, Description: "redirect_uri is required"}
	}
	if row.RedirectURI != redirectURI {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "redirect_uri does not match"}
	}

	client, err := s.clientByID(ctx, row.OAuthClientID)
	if err != nil {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "client not found"}
	}

	if tokenErr := s.validateAuthorizationCodeClient(ctx, client, formClientID, clientSecret, form, row); tokenErr != nil {
		return nil, tokenErr
	}

	verifier := form.Get("code_verifier")
	if row.CodeChallenge != "" {
		if !verifyPKCE(verifier, row.CodeChallenge, row.CodeChallengeMethod) {
			return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "invalid code_verifier"}
		}
	}

	return s.issueTokensForAuthorizationCode(ctx, row, client)
}

func (s *oidcService) validateAuthorizationCodeClient(
	ctx context.Context,
	client *models.Client,
	formClientID, clientSecret string,
	form url.Values,
	row *models.AuthorizationCode,
) *domains.OAuthTokenError {
	publicClient := client.IsPublic || strings.TrimSpace(client.ClientSecret) == ""
	effectiveClientID := formClientID
	if effectiveClientID == "" {
		effectiveClientID = form.Get("client_id")
	}
	if effectiveClientID == "" {
		return &domains.OAuthTokenError{Code: constants.OAuthInvalidRequest, Description: "client_id is required"}
	}
	if effectiveClientID != client.ClientID {
		return &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "client_id does not match authorization code"}
	}
	if !publicClient && (client.ClientSecret == "" || !s.clientSecretMatches(ctx, client, clientSecret)) {
		return &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "invalid client credentials"}
	}
	if row.OAuthClientID != client.ClientID {
		return &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "client_id does not match authorization code"}
	}
	return nil
}

func (s *oidcService) issueTokensForAuthorizationCode(
	ctx context.Context,
	row *models.AuthorizationCode,
	client *models.Client,
) (*domains.OIDCTokenResponse, *domains.OAuthTokenError) {
	u, err := s.users.GetOneByID(ctx, row.UserID)
	if err != nil {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to load user"}
	}
	if u.Status != constants.UserStatusActive {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "user account is not active"}
	}

	scope := row.Scope
	audience := client.ClientID

	access, exp, err := s.oidc.SignAccessTokenOIDC(u.ID, audience, scope, client.ClientID)
	if err != nil {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to issue access token"}
	}

	expiresIn := int64(time.Until(exp).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}

	out := &domains.OIDCTokenResponse{
		AccessToken: access,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		Scope:       scope,
	}

	if tokenErr := s.appendOpenIDToken(out, u, client, row, audience, access, scope); tokenErr != nil {
		return nil, tokenErr
	}

	refreshToken, tokenErr := s.persistRefreshToken(ctx, u, client, row)
	if tokenErr != nil {
		return nil, tokenErr
	}
	out.RefreshToken = refreshToken

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionOIDCTokenIssue,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeOAuthClient,
		ActorID:      client.ClientID,
		TenantID:     row.TenantID,
		ResourceType: constants.AuditResourceTypeClient,
		ResourceID:   client.ID,
		ResourceName: client.ClientID,
		NewValue:     map[string]any{"user_id": u.ID, "scope": scope},
	})
	return out, nil
}

func (s *oidcService) appendOpenIDToken(
	out *domains.OIDCTokenResponse,
	u *models.User,
	client *models.Client,
	row *models.AuthorizationCode,
	audience, access, scope string,
) *domains.OAuthTokenError {
	if !scopeIncludes(scope, "openid") {
		return nil
	}

	profile := &auth.OIDCUserClaims{
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Name:          strings.TrimSpace(u.FirstName + " " + u.LastName),
		GivenName:     u.FirstName,
		FamilyName:    u.LastName,
	}
	idt, err := s.oidc.SignIDToken(u.ID, audience, row.Nonce, access, profile)
	if err != nil {
		return &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to issue id_token"}
	}
	out.IDToken = idt
	return nil
}

func (s *oidcService) persistRefreshToken(
	ctx context.Context,
	u *models.User,
	client *models.Client,
	row *models.AuthorizationCode,
) (string, *domains.OAuthTokenError) {
	opaqueRefreshToken, refreshTokenHash, err := auth.NewOpaqueRefreshToken()
	if err != nil {
		return "", &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to issue refresh token"}
	}
	recID := client.ID
	rt := &models.RefreshToken{
		TenantID:       row.TenantID,
		UserID:         u.ID,
		OAuthClientID:  client.ClientID,
		TokenHash:      refreshTokenHash,
		Scope:          row.Scope,
		Revoked:        false,
		ExpiresAt:      time.Now().UTC().Add(s.cfg.JWTRefreshTTL),
		ClientRecordID: &recID,
	}
	if err := s.refreshTokens.Create(ctx, rt); err != nil {
		return "", &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to persist refresh token"}
	}
	return opaqueRefreshToken, nil
}

func scopeIncludes(scope, needle string) bool {
	for _, p := range strings.Fields(scope) {
		if p == needle {
			return true
		}
	}
	return false
}

// UserInfo returns OIDC standard claims for a valid RS256 access token.
func (s *oidcService) UserInfo(ctx context.Context, accessToken string) (resp map[string]any, oauthErr *domains.OAuthTokenError) {
	ctx, span := monitoring.StartSpan(ctx, oidcTracer, "OIDCService.UserInfo")
	defer func() {
		code, desc := "", ""
		if oauthErr != nil {
			code, desc = oauthErr.Code, oauthErr.Description
		}
		monitoring.Finish(ctx, span, "OIDCService.UserInfo", oauthFinishErr(code, desc))
	}()

	claims, err := s.oidc.ParseAccessTokenOIDC(accessToken)
	if err != nil {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidToken, Description: "invalid or expired access token"}
	}

	u, err := s.users.GetOneByID(ctx, claims.Subject)
	if err != nil {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidToken, Description: "subject not found"}
	}

	sc := claims.Scope
	out := map[string]any{
		"sub": u.ID,
	}
	if scopeIncludes(sc, "email") {
		out["email"] = u.Email
		out["email_verified"] = u.EmailVerified
	}
	if scopeIncludes(sc, "profile") {
		out["name"] = strings.TrimSpace(u.FirstName + " " + u.LastName)
		out["given_name"] = u.FirstName
		out["family_name"] = u.LastName
	}
	return out, nil
}

// Introspect implements RFC 7662 token introspection for OIDC access and refresh tokens.
// Only confidential clients may call this endpoint.
func (s *oidcService) Introspect(ctx context.Context, clientID, clientSecret, token, tokenTypeHint string) (resp *dtos.TokenIntrospectionResponse, oauthErr *domains.OAuthTokenError) {
	ctx, span := monitoring.StartSpan(ctx, oidcTracer, "OIDCService.Introspect",
		attribute.String("client_id", clientID))
	defer func() {
		code, desc := "", ""
		if oauthErr != nil {
			code, desc = oauthErr.Code, oauthErr.Description
		}
		monitoring.Finish(ctx, span, "OIDCService.Introspect", oauthFinishErr(code, desc))
	}()

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidRequest, Description: "token is required"}
	}
	if err := s.authenticateConfidentialClient(ctx, clientID, clientSecret); err != nil {
		return nil, err
	}

	hint := strings.ToLower(strings.TrimSpace(tokenTypeHint))
	switch hint {
	case "access_token":
		if resp := s.introspectAccessToken(token); resp != nil {
			return resp, nil
		}
		if resp := s.introspectRefreshToken(ctx, token); resp != nil {
			return resp, nil
		}
	case "refresh_token":
		if resp := s.introspectRefreshToken(ctx, token); resp != nil {
			return resp, nil
		}
		if resp := s.introspectAccessToken(token); resp != nil {
			return resp, nil
		}
	default:
		if looksLikeJWT(token) {
			if resp := s.introspectAccessToken(token); resp != nil {
				return resp, nil
			}
			if resp := s.introspectRefreshToken(ctx, token); resp != nil {
				return resp, nil
			}
		} else {
			if resp := s.introspectRefreshToken(ctx, token); resp != nil {
				return resp, nil
			}
			if resp := s.introspectAccessToken(token); resp != nil {
				return resp, nil
			}
		}
	}
	return &dtos.TokenIntrospectionResponse{Active: false}, nil
}

func (s *oidcService) authenticateTokenClient(ctx context.Context, clientID, clientSecret string) (*models.Client, *domains.OAuthTokenError) {
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "client authentication required"}
	}
	client, err := s.clientByID(ctx, clientID)
	if err != nil {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "invalid client credentials"}
	}
	publicClient := client.IsPublic || strings.TrimSpace(client.ClientSecret) == ""
	if publicClient {
		if strings.TrimSpace(clientSecret) != "" {
			return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "invalid client credentials"}
		}
		return client, nil
	}
	if !s.clientSecretMatches(ctx, client, clientSecret) {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "invalid client credentials"}
	}
	return client, nil
}

func (s *oidcService) authenticateConfidentialClient(ctx context.Context, clientID, clientSecret string) *domains.OAuthTokenError {
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		return &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "client authentication required"}
	}
	client, err := s.clientByID(ctx, clientID)
	if err != nil {
		return &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "invalid client credentials"}
	}
	publicClient := client.IsPublic || strings.TrimSpace(client.ClientSecret) == ""
	if publicClient {
		return &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "confidential client required"}
	}
	if !s.clientSecretMatches(ctx, client, clientSecret) {
		return &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "invalid client credentials"}
	}
	return nil
}

func (s *oidcService) refreshTokenGrant(ctx context.Context, formClientID, clientSecret string, form url.Values) (*domains.OIDCTokenResponse, *domains.OAuthTokenError) {
	raw := strings.TrimSpace(form.Get("refresh_token"))
	if raw == "" {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidRequest, Description: "refresh_token is required"}
	}
	if strings.TrimSpace(formClientID) == "" {
		formClientID = strings.TrimSpace(form.Get("client_id"))
	}
	client, authErr := s.authenticateTokenClient(ctx, formClientID, clientSecret)
	if authErr != nil {
		return nil, authErr
	}
	if !clientAllowsRefresh(client) {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthUnauthorizedClient, Description: "client is not allowed to use refresh_token"}
	}

	opaque, hash, err := auth.NewOpaqueRefreshToken()
	if err != nil {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to issue refresh token"}
	}
	recID := client.ID
	next := &models.RefreshToken{
		TenantID:       client.TenantID,
		OAuthClientID:  client.ClientID,
		TokenHash:      hash,
		Revoked:        false,
		ExpiresAt:      time.Now().UTC().Add(s.cfg.JWTRefreshTTL),
		ClientRecordID: &recID,
	}
	status, err := s.refreshTokens.Rotate(ctx, auth.HashOpaqueToken(raw), next)
	if err != nil {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to rotate refresh token"}
	}
	if status == repositories.RefreshRotationReuse {
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:       constants.AuditActionOIDCRefreshReuse,
			Result:       constants.AuditResultFailure,
			ActorType:    constants.AuditActorTypeOAuthClient,
			ActorID:      client.ClientID,
			TenantID:     client.TenantID,
			ResourceType: constants.AuditResourceTypeClient,
			ResourceID:   client.ID,
		})
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "refresh token reuse detected"}
	}
	if status != repositories.RefreshRotationOK {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "invalid or expired refresh token"}
	}
	if next.OAuthClientID != client.ClientID {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "refresh token was not issued to this client"}
	}

	u, err := s.users.GetOneByID(ctx, next.UserID)
	if err != nil || u.Status != constants.UserStatusActive {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "user account is not active"}
	}
	scope := next.Scope
	access, exp, err := s.oidc.SignAccessTokenOIDC(u.ID, client.ClientID, scope, client.ClientID)
	if err != nil {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to issue access token"}
	}
	expiresIn := int64(time.Until(exp).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}
	out := &domains.OIDCTokenResponse{
		AccessToken:  access,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		Scope:        scope,
		RefreshToken: opaque,
	}
	if scopeIncludes(scope, "openid") {
		profile := &auth.OIDCUserClaims{
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			Name:          strings.TrimSpace(u.FirstName + " " + u.LastName),
			GivenName:     u.FirstName,
			FamilyName:    u.LastName,
		}
		idt, idErr := s.oidc.SignIDToken(u.ID, client.ClientID, "", access, profile)
		if idErr != nil {
			return nil, &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to issue id_token"}
		}
		out.IDToken = idt
	}
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionOIDCTokenIssue,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeOAuthClient,
		ActorID:      client.ClientID,
		TenantID:     client.TenantID,
		ResourceType: constants.AuditResourceTypeClient,
		ResourceID:   client.ID,
		NewValue:     map[string]any{"user_id": u.ID, "grant_type": "refresh_token"},
	})
	return out, nil
}

func clientAllowsRefresh(client *models.Client) bool {
	if client == nil || len(client.GrantTypes) == 0 {
		return true
	}
	for _, grant := range client.GrantTypes {
		if grant == "refresh_token" || grant == "authorization_code" {
			return true
		}
	}
	return false
}

func (s *oidcService) clientSecretMatches(ctx context.Context, client *models.Client, presented string) bool {
	ok, upgrade := crypto.VerifyClientSecret(s.cfg.ClientSecretPepper, client.ClientSecret, presented)
	if !ok {
		return false
	}
	if upgrade != "" {
		if err := s.clients.UpdateSecretHash(ctx, client.ID, upgrade); err == nil {
			client.ClientSecret = upgrade
			s.storeClientCache(ctx, client)
		}
	}
	return true
}

func (s *oidcService) clientByID(ctx context.Context, clientID string) (*models.Client, error) {
	if s.clientCache != nil {
		if raw, err := s.clientCache.Get(ctx, oauthClientCacheKey(clientID)); err == nil && raw != "" {
			var cached models.Client
			if json.Unmarshal([]byte(raw), &cached) == nil && cached.ClientID != "" {
				return &cached, nil
			}
		}
	}
	client, err := s.clients.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	s.storeClientCache(ctx, client)
	return client, nil
}

func (s *oidcService) storeClientCache(ctx context.Context, client *models.Client) {
	if s.clientCache == nil || client == nil {
		return
	}
	raw, err := json.Marshal(client) //nolint:gosec // G117: cache stores the HMAC client-secret hash, not a plaintext secret
	if err != nil {
		return
	}
	_ = s.clientCache.Set(ctx, oauthClientCacheKey(client.ClientID), string(raw), 30*time.Second)
}

func oauthClientCacheKey(clientID string) string {
	return "iam:oauth-client:" + clientID
}

// ForgetCachedClient drops a short-lived OAuth client cache entry.
func (s *oidcService) ForgetCachedClient(ctx context.Context, clientID string) {
	if s.clientCache == nil || clientID == "" {
		return
	}
	_ = s.clientCache.Delete(ctx, oauthClientCacheKey(clientID))
}

func (s *oidcService) introspectAccessToken(token string) *dtos.TokenIntrospectionResponse {
	claims, err := s.oidc.ParseAccessTokenOIDC(token)
	if err != nil {
		return nil
	}
	resp := &dtos.TokenIntrospectionResponse{
		Active:    true,
		Scope:     claims.Scope,
		ClientID:  claims.ClientID,
		TokenType: "access_token",
		Sub:       claims.Subject,
		Iss:       claims.Issuer,
	}
	if claims.ExpiresAt != nil {
		resp.Exp = claims.ExpiresAt.Unix()
	}
	if claims.IssuedAt != nil {
		resp.Iat = claims.IssuedAt.Unix()
	}
	if len(claims.Audience) > 0 {
		resp.Aud = claims.Audience[0]
	}
	if claims.ID != "" {
		resp.JTI = claims.ID
	}
	return resp
}

func (s *oidcService) introspectRefreshToken(ctx context.Context, token string) *dtos.TokenIntrospectionResponse {
	hash := auth.HashOpaqueToken(token)
	rt, err := s.refreshTokens.FindValidByTokenHash(ctx, hash)
	if err != nil {
		return nil
	}
	now := time.Now().UTC()
	if rt.Revoked || !rt.ExpiresAt.After(now) {
		return nil
	}
	return &dtos.TokenIntrospectionResponse{
		Active:    true,
		ClientID:  rt.OAuthClientID,
		TokenType: "refresh_token",
		Sub:       rt.UserID,
		Exp:       rt.ExpiresAt.Unix(),
	}
}

func looksLikeJWT(token string) bool {
	parts := strings.Split(token, ".")
	return len(parts) == 3 && parts[0] != "" && parts[1] != "" && parts[2] != ""
}
