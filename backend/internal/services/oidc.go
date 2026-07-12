package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"go.uber.org/zap"
)

// OIDCService implements authorization code + PKCE and OIDC token issuance.
type OIDCService interface {
	Authorize(ctx context.Context, userID string, q *dtos.AuthorizeQuery) (successRedirect string, err *domains.OAuthRedirectError)
	AuthorizationCodeToken(ctx context.Context, tenantID, clientID, clientSecret string, form url.Values) (*domains.OIDCTokenResponse, *domains.OAuthTokenError)
	UserInfo(ctx context.Context, accessToken string) (map[string]any, *domains.OAuthTokenError)
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
	}
}

func (s *oidcService) OpenIDIssuer() string {
	if s.cfg.AppBaseURL != "" {
		return s.cfg.AppBaseURL
	}
	return "http://localhost:3000"
}

// Authorize validates the OAuth2 request and returns a redirect URL with an authorization code.
// userID is resolved by the handler from the browser session cookie (OIDC login at POST /oidc/login).
func (s *oidcService) Authorize(ctx context.Context, userID string, q *dtos.AuthorizeQuery) (string, *domains.OAuthRedirectError) {
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
		logger.Log.Error("oidc authorize: issue code failed",
			zap.String("operation", "oidc_authorize"),
			zap.String("client_id", clientID),
			zap.Error(err))
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
		logger.Log.Error("oidc authorize: persist authorization code failed",
			zap.String("operation", "oidc_authorize"),
			zap.String("client_id", clientID),
			zap.String("user_id", userID),
			zap.Error(err))
		return "", &domains.OAuthRedirectError{Code: constants.OAuthServerError, Description: "failed to persist code", State: state}
	}

	u, err := url.Parse(redirectURI)
	if err != nil {
		logger.Log.Error("oidc authorize: parse redirect_uri failed",
			zap.String("operation", "oidc_authorize"),
			zap.String("client_id", clientID),
			zap.Error(err))
		return "", &domains.OAuthRedirectError{Code: constants.OAuthInvalidRequest, Description: "invalid redirect_uri", State: state}
	}
	q2 := u.Query()
	q2.Set("code", codeRaw)
	if state != "" {
		q2.Set("state", state)
	}
	u.RawQuery = q2.Encode()
	logger.Log.Info("oidc authorization code issued",
		zap.String("operation", "oidc_authorize"),
		zap.String("client_id", clientID),
		zap.String("user_id", userID))
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
func (s *oidcService) AuthorizationCodeToken(ctx context.Context, _ string, formClientID, clientSecret string, form url.Values) (*domains.OIDCTokenResponse, *domains.OAuthTokenError) {
	if form.Get("grant_type") != "authorization_code" {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthUnsupportedGrantType, Description: "only authorization_code is supported"}
	}
	code := strings.TrimSpace(form.Get("code"))
	if code == "" {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidRequest, Description: "code is required"}
	}

	row, err := s.authCodes.TakeByCode(ctx, code)
	if err != nil {
		logger.Log.Error("oidc token: take authorization code failed",
			zap.String("operation", "oidc_token"),
			zap.Error(err))
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:    constants.AuditActionOIDCTokenIssue,
			Result:    constants.AuditResultFailure,
			ActorType: constants.AuditActorTypeOAuthClient,
			ActorID:   formClientID,
		})
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "invalid or expired authorization code"}
	}

	tenantID := row.TenantID
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

	client, err := s.clients.GetByClientID(ctx, row.OAuthClientID)
	if err != nil {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidClient, Description: "client not found"}
	}

	if tokenErr := validateAuthorizationCodeClient(client, formClientID, clientSecret, form, row); tokenErr != nil {
		return nil, tokenErr
	}

	verifier := form.Get("code_verifier")
	if row.CodeChallenge != "" {
		if !verifyPKCE(verifier, row.CodeChallenge, row.CodeChallengeMethod) {
			return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "invalid code_verifier"}
		}
	}

	if err := s.authCodes.DeleteByCode(ctx, code); err != nil {
		logger.Log.Error("oidc token: delete authorization code failed",
			zap.String("operation", "oidc_token"),
			zap.String("tenant_id", tenantID),
			zap.Error(err))
		return nil, &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to finalize authorization code"}
	}

	return s.issueTokensForAuthorizationCode(ctx, row, client)
}

func validateAuthorizationCodeClient(
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
	if !publicClient && (client.ClientSecret == "" || clientSecret != client.ClientSecret) {
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
		logger.Log.Error("oidc token: load user failed",
			zap.String("operation", "oidc_token"),
			zap.String("user_id", row.UserID),
			zap.Error(err))
		return nil, &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to load user"}
	}
	if u.Status != constants.UserStatusActive {
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "user account is not active"}
	}

	scope := row.Scope
	audience := client.ClientID

	access, exp, err := s.oidc.SignAccessTokenOIDC(u.ID, audience, scope, client.ClientID)
	if err != nil {
		logger.Log.Error("oidc token: sign access token failed",
			zap.String("operation", "oidc_token"),
			zap.String("user_id", u.ID),
			zap.String("client_id", client.ClientID),
			zap.Error(err))
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

	logger.Log.Info("oidc tokens issued",
		zap.String("operation", "oidc_token"),
		zap.String("user_id", u.ID),
		zap.String("client_id", client.ClientID))
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
		logger.Log.Error("oidc token: sign id_token failed",
			zap.String("operation", "oidc_token"),
			zap.String("user_id", u.ID),
			zap.String("client_id", client.ClientID),
			zap.Error(err))
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
		logger.Log.Error("oidc token: new refresh token failed",
			zap.String("operation", "oidc_token"),
			zap.String("user_id", u.ID),
			zap.String("client_id", client.ClientID),
			zap.Error(err))
		return "", &domains.OAuthTokenError{Code: constants.OAuthServerError, Description: "failed to issue refresh token"}
	}
	recID := client.ID
	rt := &models.RefreshToken{
		TenantID:       row.TenantID,
		UserID:         u.ID,
		OAuthClientID:  client.ClientID,
		TokenHash:      refreshTokenHash,
		Revoked:        false,
		ExpiresAt:      time.Now().UTC().Add(s.cfg.JWTRefreshTTL),
		ClientRecordID: &recID,
	}
	if err := s.refreshTokens.Create(ctx, rt); err != nil {
		logger.Log.Error("oidc token: persist refresh token failed",
			zap.String("operation", "oidc_token"),
			zap.String("user_id", u.ID),
			zap.String("client_id", client.ClientID),
			zap.Error(err))
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
func (s *oidcService) UserInfo(ctx context.Context, accessToken string) (map[string]any, *domains.OAuthTokenError) {
	claims, err := s.oidc.ParseAccessTokenOIDC(accessToken)
	if err != nil {
		logger.Log.Warn("oidc userinfo: invalid access token",
			zap.String("operation", "oidc_userinfo"),
			zap.Error(err))
		return nil, &domains.OAuthTokenError{Code: constants.OAuthInvalidToken, Description: "invalid or expired access token"}
	}

	u, err := s.users.GetOneByID(ctx, claims.Subject)
	if err != nil {
		logger.Log.Warn("oidc userinfo: subject not found",
			zap.String("operation", "oidc_userinfo"),
			zap.String("user_id", claims.Subject),
			zap.Error(err))
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
