package services

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/crypto"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type oidcFederationProvider struct {
	spec    constants.IdentityProviderSpec
	cfg     *config.Config
	tipRepo repositories.TenantIdentityProviderRepository

	oidcOnce    sync.Once
	oidcProv    *oidc.Provider
	oidcProvErr error
}

func newOIDCFederationProvider(
	cfg *config.Config,
	tipRepo repositories.TenantIdentityProviderRepository,
	spec constants.IdentityProviderSpec,
) FederationIdentityProvider {
	return &oidcFederationProvider{spec: spec, cfg: cfg, tipRepo: tipRepo}
}

func (p *oidcFederationProvider) ID() string {
	return p.spec.ID
}

func (p *oidcFederationProvider) DisplayName() string {
	return p.spec.DisplayName
}

func (p *oidcFederationProvider) OAuthConfiguredForTenant(ctx context.Context, tenantID string) (bool, error) {
	return p.tipRepo.IsProviderConfigured(ctx, tenantID, p.ID())
}

func (p *oidcFederationProvider) loadTenantOAuth(ctx context.Context, tenantID string) (*oauth2.Config, string, error) {
	tip, err := p.tipRepo.GetByTenantAndProvider(ctx, tenantID, p.ID())
	if err != nil {
		return nil, "", err
	}
	clientID := strings.TrimSpace(tip.OAuthClientID)
	encSecret := strings.TrimSpace(tip.OAuthClientSecretEncrypted)
	if clientID == "" || encSecret == "" {
		return nil, "", errors.ForbiddenError(
			fmt.Sprintf("%s sign-in is not configured for this organization", p.spec.DisplayName),
			nil,
		)
	}
	clientSecret, err := crypto.DecryptMFASecret(federationEncryptionKey(p.cfg), encSecret)
	if err != nil {
		return nil, "", errors.InternalError(
			fmt.Sprintf("Failed to decrypt %s OAuth credentials", p.spec.DisplayName),
			err,
		)
	}

	provider, err := p.lazyOIDC(ctx)
	if err != nil {
		return nil, "", errors.ExternalServiceError("OIDC provider init failed", err)
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  FederationCallbackURL(p.cfg, p.ID()),
		Endpoint:     provider.Endpoint(),
		Scopes:       p.spec.Scopes,
	}, clientID, nil
}

func (p *oidcFederationProvider) lazyOIDC(ctx context.Context) (*oidc.Provider, error) {
	p.oidcOnce.Do(func() {
		p.oidcProv, p.oidcProvErr = oidc.NewProvider(ctx, p.spec.IssuerURL)
	})
	return p.oidcProv, p.oidcProvErr
}

func (p *oidcFederationProvider) AuthorizeRedirectURL(ctx context.Context, tenantID, state, nonce string) (string, error) {
	oauthCfg, _, err := p.loadTenantOAuth(ctx, tenantID)
	if err != nil {
		return "", err
	}
	opts := []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("nonce", nonce)}
	for k, v := range p.spec.AuthURLParams {
		opts = append(opts, oauth2.SetAuthURLParam(k, v))
	}
	return oauthCfg.AuthCodeURL(state, opts...), nil
}

func (p *oidcFederationProvider) ExchangeAuthorizationCode(ctx context.Context, tenantID, code, expectedNonce string) (*domains.OIDCUserClaims, error) {
	oauthCfg, clientID, err := p.loadTenantOAuth(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	tok, err := oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, errors.ExternalServiceError(
			fmt.Sprintf("%s token exchange failed", p.spec.DisplayName),
			err,
		)
	}

	rawIDToken, _ := tok.Extra("id_token").(string)
	if rawIDToken == "" {
		return nil, errors.ExternalServiceError(
			fmt.Sprintf("%s did not return an ID token", p.spec.DisplayName),
			nil,
		)
	}

	provider, err := p.lazyOIDC(ctx)
	if err != nil {
		return nil, errors.ExternalServiceError("OIDC provider init failed", err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, errors.UnauthorizedError(
			fmt.Sprintf("Invalid %s ID token", p.spec.DisplayName),
			err,
		)
	}

	var raw domains.OIDCIDTokenClaims
	if err := idToken.Claims(&raw); err != nil {
		return nil, errors.UnauthorizedError(
			fmt.Sprintf("Invalid %s token claims", p.spec.DisplayName),
			err,
		)
	}
	if raw.Nonce != expectedNonce {
		return nil, errors.UnauthorizedError("Invalid OAuth nonce", nil)
	}
	if p.spec.RequireVerifiedEmail && (!raw.EmailVerified || strings.TrimSpace(raw.Email) == "") {
		return nil, errors.ForbiddenError(
			fmt.Sprintf("A verified email is required to sign in with %s", p.spec.DisplayName),
			nil,
		)
	}
	if strings.TrimSpace(raw.Sub) == "" {
		return nil, errors.ExternalServiceError(
			fmt.Sprintf("%s token missing subject", p.spec.DisplayName),
			nil,
		)
	}

	return &domains.OIDCUserClaims{
		Sub:           raw.Sub,
		Email:         raw.Email,
		EmailVerified: raw.EmailVerified,
		GivenName:     raw.GivenName,
		FamilyName:    raw.FamilyName,
		Name:          raw.Name,
	}, nil
}
