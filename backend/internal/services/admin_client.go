package services

import (
	"context"
	"strings"

	"github.com/lib/pq"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
)

const seededDevClientID = "oidc-dev"

var (
	defaultClientGrantTypes = []string{"authorization_code"}
	defaultClientScopes     = []string{"openid", "email", "profile"}
)

func (s *adminService) GetClientByID(ctx context.Context, clientID string) (*dtos.AdminClientResponse, error) {
	client, err := s.clients.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return dtos.NewAdminClientResponse(client), nil
}

func (s *adminService) CreateClient(ctx context.Context, req *dtos.AdminCreateClientRequest) (*dtos.AdminCreateClientResponse, error) {
	if req == nil {
		return nil, errors.ValidationError("Request body is required", nil)
	}

	tenantID := strings.TrimSpace(req.TenantID)
	if _, err := s.tenants.GetByID(ctx, tenantID); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.ValidationError("Client name is required", nil)
	}

	clientID := strings.TrimSpace(req.ClientID)
	if clientID == "" {
		generated, err := generateClientIdentifier()
		if err != nil {
			return nil, err
		}
		clientID = generated
	}

	taken, err := s.clients.ClientIDTaken(ctx, tenantID, clientID, "")
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, errors.ConflictError("OAuth client ID is already in use for this tenant", nil)
	}

	redirectUris, err := normalizeRedirectUris(req.RedirectUris)
	if err != nil {
		return nil, err
	}

	grantTypes := normalizeGrantTypes(req.GrantTypes)
	scopes := normalizeScopes(req.Scopes)

	var plaintextSecret string
	clientSecret := ""
	if !req.IsPublic {
		secret, err := generateClientSecret()
		if err != nil {
			return nil, err
		}
		plaintextSecret = secret
		clientSecret = secret
	}

	client := &models.Client{
		BaseModel:    models.NewBaseModel(),
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Name:         name,
		RedirectUris: pq.StringArray(redirectUris),
		GrantTypes:   pq.StringArray(grantTypes),
		Scopes:       pq.StringArray(scopes),
		IsPublic:     req.IsPublic,
	}
	if err := s.clients.Create(ctx, client); err != nil {
		return nil, err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminClientCreate,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeClient,
		ResourceID:   client.ID,
		ResourceName: client.ClientID,
		NewValue: map[string]any{
			"client_id":     client.ClientID,
			"name":          client.Name,
			"is_public":     client.IsPublic,
			"redirect_uris": redirectUris,
			"grant_types":   grantTypes,
			"scopes":        scopes,
		},
	})

	resp := &dtos.AdminCreateClientResponse{
		AdminClientResponse: *dtos.NewAdminClientResponse(client),
		ClientSecret:        plaintextSecret,
	}
	return resp, nil
}

func (s *adminService) UpdateClient(ctx context.Context, clientID string, req *dtos.AdminUpdateClientRequest) (*dtos.AdminClientResponse, error) {
	if req == nil {
		return nil, errors.ValidationError("Request body is required", nil)
	}
	if req.Name == nil && req.RedirectUris == nil && req.GrantTypes == nil && req.Scopes == nil && req.IsPublic == nil && req.ClientSecret == nil {
		return nil, errors.ValidationError("At least one field must be provided", nil)
	}

	existing, err := s.clients.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	patch := repositories.ClientPatch{}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.ValidationError("Client name cannot be empty", nil)
		}
		patch.Name = &name
	}
	if req.RedirectUris != nil {
		redirectUris, err := normalizeRedirectUris(req.RedirectUris)
		if err != nil {
			return nil, err
		}
		patch.RedirectUris = &redirectUris
	}
	if req.GrantTypes != nil {
		grantTypes := normalizeGrantTypes(req.GrantTypes)
		patch.GrantTypes = &grantTypes
	}
	if req.Scopes != nil {
		scopes := normalizeScopes(req.Scopes)
		patch.Scopes = &scopes
	}
	if req.IsPublic != nil {
		patch.IsPublic = req.IsPublic
		if *req.IsPublic {
			empty := ""
			patch.ClientSecret = &empty
		} else if strings.TrimSpace(existing.ClientSecret) == "" && req.ClientSecret == nil {
			secret, err := generateClientSecret()
			if err != nil {
				return nil, err
			}
			patch.ClientSecret = &secret
		}
	}
	if req.ClientSecret != nil {
		secret := strings.TrimSpace(*req.ClientSecret)
		if secret == "" {
			return nil, errors.ValidationError("Client secret cannot be empty", nil)
		}
		if existing.IsPublic {
			return nil, errors.ValidationError("Public clients cannot have a client secret", nil)
		}
		patch.ClientSecret = &secret
	}

	updated, err := s.clients.Update(ctx, clientID, patch)
	if err != nil {
		return nil, err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminClientUpdate,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		TenantID:     existing.TenantID,
		ResourceType: constants.AuditResourceTypeClient,
		ResourceID:   clientID,
		ResourceName: existing.ClientID,
		OldValue: map[string]any{
			"name":          existing.Name,
			"is_public":     existing.IsPublic,
			"redirect_uris": []string(existing.RedirectUris),
			"grant_types":   []string(existing.GrantTypes),
			"scopes":        []string(existing.Scopes),
		},
		NewValue: map[string]any{
			"name":          updated.Name,
			"is_public":     updated.IsPublic,
			"redirect_uris": []string(updated.RedirectUris),
			"grant_types":   []string(updated.GrantTypes),
			"scopes":        []string(updated.Scopes),
		},
	})

	return dtos.NewAdminClientResponse(updated), nil
}

func (s *adminService) DeleteClient(ctx context.Context, clientID string) error {
	existing, err := s.clients.GetByID(ctx, clientID)
	if err != nil {
		return err
	}
	if existing.TenantID == s.cfg.DefaultTenantID && existing.ClientID == seededDevClientID {
		return errors.ValidationError("The default development OAuth client cannot be deleted", nil)
	}
	if err := s.clients.Delete(ctx, clientID); err != nil {
		return err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminClientDelete,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		TenantID:     existing.TenantID,
		ResourceType: constants.AuditResourceTypeClient,
		ResourceID:   clientID,
		ResourceName: existing.ClientID,
		OldValue: map[string]any{
			"client_id": existing.ClientID,
			"name":      existing.Name,
		},
	})
	return nil
}

func generateClientIdentifier() (string, error) {
	raw, _, err := auth.NewOpaqueRefreshToken()
	if err != nil {
		return "", errors.InternalError("Failed to generate client ID", err)
	}
	return raw, nil
}

func generateClientSecret() (string, error) {
	raw, _, err := auth.NewOpaqueRefreshToken()
	if err != nil {
		return "", errors.InternalError("Failed to generate client secret", err)
	}
	return raw, nil
}

func normalizeRedirectUris(uris []string) ([]string, error) {
	if len(uris) == 0 {
		return nil, errors.ValidationError("At least one redirect URI is required", nil)
	}
	out := make([]string, 0, len(uris))
	seen := make(map[string]struct{}, len(uris))
	for _, uri := range uris {
		uri = strings.TrimSpace(uri)
		if uri == "" {
			continue
		}
		if _, ok := seen[uri]; ok {
			continue
		}
		seen[uri] = struct{}{}
		out = append(out, uri)
	}
	if len(out) == 0 {
		return nil, errors.ValidationError("At least one redirect URI is required", nil)
	}
	return out, nil
}

func normalizeGrantTypes(grantTypes []string) []string {
	if len(grantTypes) == 0 {
		return append([]string(nil), defaultClientGrantTypes...)
	}
	out := make([]string, 0, len(grantTypes))
	seen := make(map[string]struct{}, len(grantTypes))
	for _, gt := range grantTypes {
		gt = strings.TrimSpace(gt)
		if gt == "" {
			continue
		}
		if _, ok := seen[gt]; ok {
			continue
		}
		seen[gt] = struct{}{}
		out = append(out, gt)
	}
	if len(out) == 0 {
		return append([]string(nil), defaultClientGrantTypes...)
	}
	return out
}

func normalizeScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return append([]string(nil), defaultClientScopes...)
	}
	out := make([]string, 0, len(scopes))
	seen := make(map[string]struct{}, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		out = append(out, scope)
	}
	if len(out) == 0 {
		return append([]string(nil), defaultClientScopes...)
	}
	return out
}
