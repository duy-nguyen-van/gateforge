package services

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/stretchr/testify/require"
)

type adminIdpTipStub struct {
	byKey map[string]*models.TenantIdentityProvider
}

func (s *adminIdpTipStub) key(tenantID, provider string) string {
	return tenantID + ":" + provider
}

func (s *adminIdpTipStub) IsProviderEnabled(context.Context, string, string) (bool, error) {
	return false, nil
}

func (s *adminIdpTipStub) IsProviderConfigured(context.Context, string, string) (bool, error) {
	return false, nil
}

func (s *adminIdpTipStub) GetByTenantAndProvider(_ context.Context, tenantID, provider string) (*models.TenantIdentityProvider, error) {
	if s.byKey == nil {
		return nil, errors.NotFoundError("TenantIdentityProvider", nil)
	}
	tip, ok := s.byKey[s.key(tenantID, provider)]
	if !ok {
		return nil, errors.NotFoundError("TenantIdentityProvider", nil)
	}
	return tip, nil
}

func (s *adminIdpTipStub) SetProviderEnabled(context.Context, string, string, bool) error {
	return nil
}

func (s *adminIdpTipStub) UpdateProvider(_ context.Context, tenantID, provider string, patch repositories.TenantIdentityProviderPatch) (*models.TenantIdentityProvider, error) {
	if s.byKey == nil {
		s.byKey = map[string]*models.TenantIdentityProvider{}
	}
	key := s.key(tenantID, provider)
	tip := s.byKey[key]
	if tip == nil {
		tip = &models.TenantIdentityProvider{TenantID: tenantID, Provider: provider}
		s.byKey[key] = tip
	}
	if patch.Enabled != nil {
		tip.Enabled = *patch.Enabled
	}
	if patch.OAuthClientID != "" {
		tip.OAuthClientID = patch.OAuthClientID
	}
	if patch.OAuthClientSecretEncrypted != "" {
		tip.OAuthClientSecretEncrypted = patch.OAuthClientSecretEncrypted
	}
	return tip, nil
}

func (s *adminIdpTipStub) ListByTenant(context.Context, string) ([]models.TenantIdentityProvider, error) {
	return nil, nil
}

type adminIdpAuditStub struct{}

func (adminIdpAuditStub) Record(context.Context, domains.AuditRecordParams) {}

func newAdminIdpTestService(tip *adminIdpTipStub) AdminService {
	cfg := &config.Config{
		AppBaseURL:       "http://localhost:3000",
		MFAEncryptionKey: "01234567890123456789012345678901",
	}
	return ProvideAdminService(
		cfg,
		nil, nil, nil, nil, nil, nil,
		tip,
		nil, nil, nil, nil,
		adminIdpAuditStub{},
		nil, nil,
	)
}

func TestConfigureIdentityProvider_UnknownProvider(t *testing.T) {
	svc := newAdminIdpTestService(&adminIdpTipStub{})
	enabled := true
	err := svc.ConfigureIdentityProvider(
		context.Background(),
		"tenant-1",
		"github",
		&dtos.PatchIdentityProviderRequest{Enabled: &enabled},
		constants.AuditActorTypeUser,
	)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeValidation, appErr.Type)
}

func TestConfigureIdentityProvider_EnableWithoutCredentials(t *testing.T) {
	svc := newAdminIdpTestService(&adminIdpTipStub{})
	enabled := true
	err := svc.ConfigureIdentityProvider(
		context.Background(),
		"tenant-1",
		constants.IdentityProviderGoogle,
		&dtos.PatchIdentityProviderRequest{Enabled: &enabled},
		constants.AuditActorTypeUser,
	)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeValidation, appErr.Type)
	require.Contains(t, appErr.Message, "Google")
}
