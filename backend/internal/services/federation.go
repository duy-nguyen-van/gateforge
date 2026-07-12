package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/cache"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// FederationIdentityProvider implements OAuth2/OIDC against one upstream IdP (Google, Microsoft, …).
type FederationIdentityProvider interface {
	ID() string
	DisplayName() string
	OAuthConfiguredForTenant(ctx context.Context, tenantID string) (bool, error)
	AuthorizeRedirectURL(ctx context.Context, tenantID, state, nonce string) (string, error)
	ExchangeAuthorizationCode(ctx context.Context, tenantID, code, expectedNonce string) (*domains.OIDCUserClaims, error)
}

const (
	federationStateKeyPrefix = "oidc_federation_state:"
	federationStateTTL       = 10 * time.Minute
)

// FederationService coordinates OAuth state, tenant checks, and user provisioning for registered IdPs.
type FederationService interface {
	BuildAuthorizeRedirectURL(ctx context.Context, providerID, returnTo, tenantID string) (redirectURL string, err error)
	CompleteOAuthLogin(ctx context.Context, providerID, code, state string) (user *models.User, tenantID string, returnTo string, err error)
	ListAvailableProviders(ctx context.Context, tenantID string) ([]ProviderAvailability, error)
}

// ProviderAvailability describes a sign-in method exposed to the login page.
type ProviderAvailability struct {
	Provider string
	Name     string
}

type oauthStatePayload struct {
	ReturnTo   string `json:"return_to"`
	TenantID   string `json:"tenant_id"`
	Nonce      string `json:"nonce"`
	ProviderID string `json:"provider"`
}

type federationService struct {
	cfg            *config.Config
	cache          cache.Cache
	userRepo       repositories.UserRepository
	membershipRepo repositories.TenantMembershipRepository
	fedRepo        repositories.FederatedIdentityRepository
	tipRepo        repositories.TenantIdentityProviderRepository
	pg             *db.PostgresDB
	providers      map[string]FederationIdentityProvider
	audit          AuditService
}

// ProvideFederationService wires federation and registers built-in IdP implementations.
func ProvideFederationService(
	cfg *config.Config,
	c cache.Cache,
	userRepo repositories.UserRepository,
	membershipRepo repositories.TenantMembershipRepository,
	fedRepo repositories.FederatedIdentityRepository,
	tipRepo repositories.TenantIdentityProviderRepository,
	pg *db.PostgresDB,
	audit AuditService,
) FederationService {
	providers := make(map[string]FederationIdentityProvider, len(constants.SupportedIdentityProviders))
	for _, spec := range constants.SupportedIdentityProviders {
		providers[spec.ID] = newOIDCFederationProvider(cfg, tipRepo, spec)
	}
	return &federationService{
		cfg:            cfg,
		cache:          c,
		userRepo:       userRepo,
		membershipRepo: membershipRepo,
		fedRepo:        fedRepo,
		tipRepo:        tipRepo,
		pg:             pg,
		providers:      providers,
		audit:          audit,
	}
}

func (s *federationService) providerByID(id string) (FederationIdentityProvider, error) {
	key := strings.ToLower(strings.TrimSpace(id))
	if key == "" {
		return nil, errors.ValidationError("Identity provider is required", nil)
	}
	p, ok := s.providers[key]
	if !ok {
		return nil, errors.ValidationError("Unsupported identity provider", nil)
	}
	return p, nil
}

func federationStateCacheKey(providerID, state string) string {
	return federationStateKeyPrefix + strings.ToLower(strings.TrimSpace(providerID)) + ":" + state
}

func (s *federationService) BuildAuthorizeRedirectURL(ctx context.Context, providerID, returnTo, tenantID string) (string, error) {
	p, err := s.providerByID(providerID)
	if err != nil {
		return "", err
	}
	configured, err := p.OAuthConfiguredForTenant(ctx, tenantID)
	if err != nil {
		return "", err
	}
	if !configured {
		return "", errors.ForbiddenError("This sign-in method is not configured for this organization", nil)
	}

	enabled, err := s.tipRepo.IsProviderEnabled(ctx, tenantID, p.ID())
	if err != nil {
		return "", err
	}
	if !enabled {
		return "", errors.ForbiddenError("This sign-in method is disabled for this organization", nil)
	}

	state, err := federationRandomHex(32)
	if err != nil {
		return "", errors.InternalError("Failed to generate OAuth state", err)
	}
	nonce, err := federationRandomHex(32)
	if err != nil {
		return "", errors.InternalError("Failed to generate OAuth nonce", err)
	}

	payload, err := json.Marshal(oauthStatePayload{
		ReturnTo:   returnTo,
		TenantID:   tenantID,
		Nonce:      nonce,
		ProviderID: p.ID(),
	})
	if err != nil {
		return "", errors.InternalError("Failed to encode OAuth state", err)
	}

	key := federationStateCacheKey(p.ID(), state)
	if err := s.cache.Set(ctx, key, string(payload), federationStateTTL); err != nil {
		return "", err
	}

	url, err := p.AuthorizeRedirectURL(ctx, tenantID, state, nonce)
	if err != nil {
		return "", err
	}
	logger.Log.Info("federation authorize URL built",
		zap.String("operation", "federation_oauth_authorize"),
		zap.String("provider", p.ID()),
		zap.String("tenant_id", tenantID))
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionFederationStart,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeSystem,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeIdentityProvider,
		ResourceName: p.ID(),
	})
	return url, nil
}

func (s *federationService) CompleteOAuthLogin(ctx context.Context, providerID, code, state string) (*models.User, string, string, error) {
	p, err := s.providerByID(providerID)
	if err != nil {
		return nil, "", "", err
	}
	if strings.TrimSpace(code) == "" || strings.TrimSpace(state) == "" {
		return nil, "", "", errors.ValidationError("Missing code or state", nil)
	}

	key := federationStateCacheKey(p.ID(), state)
	raw, err := s.cache.Get(ctx, key)
	if err != nil {
		if federationIsNotFoundAppError(err) {
			return nil, "", "", errors.ValidationError("Invalid or expired OAuth state", nil)
		}
		return nil, "", "", err
	}
	_ = s.cache.Delete(ctx, key)

	var payload oauthStatePayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, "", "", errors.ValidationError("Invalid OAuth state payload", err)
	}
	if !strings.EqualFold(strings.TrimSpace(payload.ProviderID), p.ID()) {
		return nil, "", "", errors.ValidationError("OAuth state does not match provider", nil)
	}

	enabled, err := s.tipRepo.IsProviderEnabled(ctx, payload.TenantID, p.ID())
	if err != nil {
		return nil, "", "", err
	}
	if !enabled {
		return nil, "", "", errors.ForbiddenError("This sign-in method is disabled for this organization", nil)
	}

	claims, err := p.ExchangeAuthorizationCode(ctx, payload.TenantID, code, payload.Nonce)
	if err != nil {
		logger.Log.Error("federation token exchange or verify failed",
			zap.String("operation", "federation_oauth_callback"),
			zap.String("provider", p.ID()),
			zap.String("tenant_id", payload.TenantID),
			zap.Error(err))
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:       constants.AuditActionFederationLogin,
			Result:       constants.AuditResultFailure,
			ActorType:    constants.AuditActorTypeSystem,
			TenantID:     payload.TenantID,
			ResourceType: constants.AuditResourceTypeIdentityProvider,
			ResourceName: p.ID(),
		})
		return nil, "", "", err
	}

	u, err := s.federationResolveOrProvisionUser(ctx, payload.TenantID, p.ID(), claims)
	if err != nil {
		logger.Log.Error("federation resolve or provision user failed",
			zap.String("operation", "federation_oauth_callback"),
			zap.String("provider", p.ID()),
			zap.String("tenant_id", payload.TenantID),
			zap.Error(err))
		return nil, "", "", err
	}
	logger.Log.Info("federation login completed",
		zap.String("operation", "federation_oauth_callback"),
		zap.String("provider", p.ID()),
		zap.String("user_id", u.ID),
		zap.String("tenant_id", payload.TenantID))
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionFederationLogin,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      u.ID,
		TenantID:     payload.TenantID,
		ResourceType: constants.AuditResourceTypeIdentityProvider,
		ResourceName: p.ID(),
		ResourceID:   u.ID,
	})
	return u, payload.TenantID, payload.ReturnTo, nil
}

func (s *federationService) federationResolveOrProvisionUser(ctx context.Context, tenantID, providerID string, claims *domains.OIDCUserClaims) (*models.User, error) {
	sub := strings.TrimSpace(claims.Sub)
	email := strings.TrimSpace(claims.Email)
	emailLower := strings.ToLower(email)

	fi, err := s.fedRepo.GetByProviderSubject(ctx, providerID, sub)
	if err == nil && fi != nil && fi.User != nil {
		if fi.User.Status != constants.UserStatusActive {
			return nil, errors.ForbiddenError("Account is not active", nil).
				WithOperation("federation_login").
				WithResource("user")
		}
		if err := s.ensureMembership(ctx, fi.User.ID, tenantID); err != nil {
			return nil, err
		}
		return fi.User, nil
	}
	if err != nil && !federationIsNotFoundAppError(err) {
		return nil, err
	}

	existing, err := s.userRepo.GetByEmailLower(ctx, emailLower)
	if err == nil && existing != nil {
		if existing.Status != constants.UserStatusActive {
			return nil, errors.ForbiddenError("Account is not active", nil).
				WithOperation("federation_login").
				WithResource("user")
		}
		link := &models.FederatedIdentity{
			UserID:      existing.ID,
			Provider:    providerID,
			Subject:     sub,
			EmailAtLink: email,
		}
		if err := s.fedRepo.Create(ctx, link); err != nil {
			return nil, err
		}
		if err := s.ensureMembership(ctx, existing.ID, tenantID); err != nil {
			return nil, err
		}
		return existing, nil
	}
	if err != nil && !federationIsNotFoundAppError(err) {
		return nil, err
	}

	firstName, lastName := claims.SplitDisplayName()
	user := &models.User{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		EmailVerified: claims.EmailVerified,
		Status:        constants.UserStatusActive,
	}
	link := &models.FederatedIdentity{
		Provider:    providerID,
		Subject:     sub,
		EmailAtLink: email,
	}
	membership := &models.TenantMembership{
		TenantID: tenantID,
		Role:     constants.TenantMembershipRoleMember,
		Status:   constants.TenantMembershipStatusActive,
	}

	err = s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return errors.DatabaseError("Failed to create federated user", err).
				WithOperation("federation_login").
				WithResource("user")
		}
		link.UserID = user.ID
		if err := tx.Create(link).Error; err != nil {
			return errors.DatabaseError("Failed to link federated identity", err).
				WithOperation("federation_login").
				WithResource("federated_identity")
		}
		membership.UserID = user.ID
		if err := tx.Create(membership).Error; err != nil {
			return errors.DatabaseError("Failed to create tenant membership", err).
				WithOperation("federation_login").
				WithResource("tenant_membership")
		}
		return nil
	})
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			return nil, appErr
		}
		return nil, err
	}
	return user, nil
}

func (s *federationService) ensureMembership(ctx context.Context, userID, tenantID string) error {
	ok, err := s.membershipRepo.ExistsActive(ctx, userID, tenantID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return s.membershipRepo.Create(ctx, &models.TenantMembership{
		BaseModel: models.NewBaseModel(),
		UserID:    userID,
		TenantID:  tenantID,
		Role:      constants.TenantMembershipRoleMember,
		Status:    constants.TenantMembershipStatusActive,
	})
}

func federationIsNotFoundAppError(err error) bool {
	appErr := errors.GetAppError(err)
	return appErr != nil && appErr.Type == errors.ErrorTypeNotFound
}

func federationRandomHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *federationService) ListAvailableProviders(ctx context.Context, tenantID string) ([]ProviderAvailability, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return nil, errors.ValidationError("tenant_id is required", nil)
	}

	var out []ProviderAvailability
	for _, p := range s.providers {
		configured, err := p.OAuthConfiguredForTenant(ctx, tenantID)
		if err != nil {
			return nil, err
		}
		if !configured {
			continue
		}
		enabled, err := s.tipRepo.IsProviderEnabled(ctx, tenantID, p.ID())
		if err != nil {
			return nil, err
		}
		if !enabled {
			continue
		}
		out = append(out, ProviderAvailability{Provider: p.ID(), Name: p.DisplayName()})
	}
	return out, nil
}
