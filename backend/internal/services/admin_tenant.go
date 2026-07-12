package services

import (
	"context"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
)

func (s *adminService) GetTenantByID(ctx context.Context, tenantID string) (*dtos.AdminTenantResponse, error) {
	tenant, err := s.tenants.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	userCount, err := s.tenants.CountUsersByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return dtos.NewAdminTenantResponse(tenant, userCount), nil
}

func (s *adminService) CreateTenant(ctx context.Context, req *dtos.AdminCreateTenantRequest) (*dtos.AdminTenantResponse, error) {
	if req == nil {
		return nil, errors.ValidationError("Request body is required", nil)
	}
	name := strings.TrimSpace(req.Name)
	domain := strings.TrimSpace(req.Domain)
	if name == "" {
		return nil, errors.ValidationError("Tenant name is required", nil)
	}
	if domain != "" {
		taken, err := s.tenants.DomainTaken(ctx, domain, "")
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, errors.ConflictError("Tenant domain is already in use", nil)
		}
	}

	tenant := &models.Tenant{
		BaseModel: models.NewBaseModel(),
		Name:      name,
		Domain:    domain,
	}
	if err := s.tenants.Create(ctx, tenant); err != nil {
		return nil, err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminTenantCreate,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		TenantID:     tenant.ID,
		ResourceType: constants.AuditResourceTypeTenant,
		ResourceID:   tenant.ID,
		ResourceName: tenant.Name,
		NewValue:     map[string]any{"name": tenant.Name, "domain": tenant.Domain},
	})

	return dtos.NewAdminTenantResponse(tenant, 0), nil
}

func (s *adminService) UpdateTenant(ctx context.Context, tenantID string, req *dtos.AdminUpdateTenantRequest) (*dtos.AdminTenantResponse, error) {
	if req == nil {
		return nil, errors.ValidationError("Request body is required", nil)
	}
	if req.Name == nil && req.Domain == nil {
		return nil, errors.ValidationError("At least one field must be provided", nil)
	}

	existing, err := s.tenants.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	patch := repositories.TenantPatch{}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.ValidationError("Tenant name cannot be empty", nil)
		}
		patch.Name = &name
	}
	if req.Domain != nil {
		domain := strings.TrimSpace(*req.Domain)
		if domain != "" {
			taken, err := s.tenants.DomainTaken(ctx, domain, tenantID)
			if err != nil {
				return nil, err
			}
			if taken {
				return nil, errors.ConflictError("Tenant domain is already in use", nil)
			}
		}
		patch.Domain = &domain
	}

	updated, err := s.tenants.Update(ctx, tenantID, patch)
	if err != nil {
		return nil, err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminTenantUpdate,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeTenant,
		ResourceID:   tenantID,
		ResourceName: updated.Name,
		OldValue:     map[string]any{"name": existing.Name, "domain": existing.Domain},
		NewValue:     map[string]any{"name": updated.Name, "domain": updated.Domain},
	})

	userCount, err := s.tenants.CountUsersByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return dtos.NewAdminTenantResponse(updated, userCount), nil
}

func (s *adminService) DeleteTenant(ctx context.Context, tenantID string) error {
	if tenantID == s.cfg.DefaultTenantID {
		return errors.ValidationError("The default tenant cannot be deleted", nil)
	}
	existing, err := s.tenants.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}
	if err := s.tenants.Delete(ctx, tenantID); err != nil {
		return err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminTenantDelete,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeTenant,
		ResourceID:   tenantID,
		ResourceName: existing.Name,
		OldValue:     map[string]any{"name": existing.Name, "domain": existing.Domain},
	})
	return nil
}

func (s *adminService) ListTenantMembers(ctx context.Context, tenantID string, pr *dtos.PageableRequest) ([]*dtos.AdminTenantMemberResponse, *dtos.Pageable, error) {
	if _, err := s.tenants.GetByID(ctx, tenantID); err != nil {
		return nil, nil, err
	}
	result, err := s.memberships.ListByTenantIDPaginated(ctx, tenantID, pr)
	if err != nil {
		return nil, nil, err
	}
	out := make([]*dtos.AdminTenantMemberResponse, 0, len(result.Data))
	for i := range result.Data {
		out = append(out, dtos.NewAdminTenantMemberResponse(&result.Data[i]))
	}
	return out, result.Pageable, nil
}
