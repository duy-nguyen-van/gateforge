package services

import (
	"context"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
)

// TenantResolveInput carries hints for resolving the active tenant.
type TenantResolveInput struct {
	Host          string
	TenantIDParam string
	OAuthClientID string
	UserID        string
}

// TenantResolveResult is the outcome of tenant resolution.
type TenantResolveResult struct {
	TenantID       string
	RequiresPicker bool
}

// TenantContextService resolves active tenant from request context and validates membership.
type TenantContextService interface {
	Resolve(ctx context.Context, in TenantResolveInput) (*TenantResolveResult, error)
	ValidateMembership(ctx context.Context, userID, tenantID string) error
}

type tenantContextService struct {
	cfg         *config.Config
	clients     repositories.ClientRepository
	tenants     repositories.TenantRepository
	memberships repositories.TenantMembershipRepository
}

// ProvideTenantContextService wires tenant resolution.
func ProvideTenantContextService(
	cfg *config.Config,
	clients repositories.ClientRepository,
	tenants repositories.TenantRepository,
	memberships repositories.TenantMembershipRepository,
) TenantContextService {
	return &tenantContextService{
		cfg:         cfg,
		clients:     clients,
		tenants:     tenants,
		memberships: memberships,
	}
}

func (s *tenantContextService) Resolve(ctx context.Context, in TenantResolveInput) (*TenantResolveResult, error) {
	tenantID := s.resolveTenantID(ctx, in)
	if tenantID != "" {
		if in.UserID != "" {
			if err := s.ValidateMembership(ctx, in.UserID, tenantID); err != nil {
				return nil, err
			}
		}
		return &TenantResolveResult{TenantID: tenantID}, nil
	}

	if in.UserID == "" {
		return nil, errors.ValidationError("tenant context required", nil)
	}

	memberships, err := s.memberships.ListByUserID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}
	if len(memberships) == 0 {
		return nil, errors.ForbiddenError("No tenant access", nil)
	}
	if len(memberships) == 1 {
		return &TenantResolveResult{TenantID: memberships[0].TenantID}, nil
	}

	// Multiple memberships and no context — caller should show tenant picker.
	if s.cfg.DefaultTenantID != "" {
		for _, m := range memberships {
			if m.TenantID == s.cfg.DefaultTenantID {
				return &TenantResolveResult{TenantID: m.TenantID}, nil
			}
		}
	}

	return &TenantResolveResult{RequiresPicker: true}, nil
}

func (s *tenantContextService) resolveTenantID(ctx context.Context, in TenantResolveInput) string {
	if clientID := strings.TrimSpace(in.OAuthClientID); clientID != "" {
		client, err := s.clients.GetByClientID(ctx, clientID)
		if err == nil && client != nil {
			return client.TenantID
		}
	}

	if tid := strings.TrimSpace(in.TenantIDParam); tid != "" {
		return tid
	}

	if host := strings.TrimSpace(in.Host); host != "" {
		domain := extractDomainFromHost(host)
		if domain != "" {
			tenant, err := s.tenants.GetByDomain(ctx, domain)
			if err == nil && tenant != nil {
				return tenant.ID
			}
		}
	}

	if in.UserID != "" {
		memberships, err := s.memberships.ListByUserID(ctx, in.UserID)
		if err == nil && len(memberships) == 1 {
			return memberships[0].TenantID
		}
	}

	return ""
}

func (s *tenantContextService) ValidateMembership(ctx context.Context, userID, tenantID string) error {
	ok, err := s.memberships.ExistsActive(ctx, userID, tenantID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.ForbiddenError("No access to this tenant", nil)
	}
	return nil
}

func extractDomainFromHost(host string) string {
	host = strings.ToLower(strings.TrimSpace(host))
	if idx := strings.Index(host, ":"); idx >= 0 {
		host = host[:idx]
	}
	parts := strings.Split(host, ".")
	if len(parts) >= 3 && parts[0] != "www" {
		return parts[0]
	}
	if len(parts) == 2 {
		return host
	}
	return ""
}
