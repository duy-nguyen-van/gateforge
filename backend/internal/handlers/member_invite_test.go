package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type stubMemberInviteService struct {
	preview    *dtos.MemberInvitePreview
	previewErr error
	accept     *services.MemberInviteAcceptResult
	acceptErr  error
}

func (s *stubMemberInviteService) Preview(context.Context, string) (*dtos.MemberInvitePreview, error) {
	return s.preview, s.previewErr
}

func (s *stubMemberInviteService) Accept(context.Context, *dtos.AcceptMemberInviteRequest) (*services.MemberInviteAcceptResult, error) {
	return s.accept, s.acceptErr
}

func newMemberInviteHandler(invites services.MemberInviteService, sessions services.SessionService) *MemberInviteHandler {
	return ProvideMemberInviteHandler(invites, sessions, handlerTestConfig(), validator.New())
}

func TestMemberInviteHandler_Preview(t *testing.T) {
	h := newMemberInviteHandler(&stubMemberInviteService{
		preview: &dtos.MemberInvitePreview{Email: "ada@example.com", OrganizationName: "Acme", Role: "member"},
	}, &stubAuthSessionService{})
	c, rec := newJSONContext(http.MethodGet, "/api/v1/invites/preview?token=raw", "")

	require.NoError(t, h.Preview(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "ada@example.com")
	require.Contains(t, rec.Body.String(), "Acme")
}

func TestMemberInviteHandler_Preview_Error(t *testing.T) {
	h := newMemberInviteHandler(&stubMemberInviteService{
		previewErr: errors.ForbiddenError("Invite is invalid or expired", nil),
	}, &stubAuthSessionService{})
	c, rec := newJSONContext(http.MethodGet, "/api/v1/invites/preview?token=bad", "")

	require.NoError(t, h.Preview(c))
	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestMemberInviteHandler_Accept_Session(t *testing.T) {
	h := newMemberInviteHandler(&stubMemberInviteService{
		accept: &services.MemberInviteAcceptResult{
			UserID: testUserID,
			Login:  &dtos.LoginResponse{AccessToken: "access-token", ActiveTenantID: testTenantID},
		},
	}, &stubAuthSessionService{createSID: "invite-session"})
	c, rec := newJSONContext(http.MethodPost, "/api/v1/invites/accept", `{"token":"raw-token","password":"long-secure-passphrase","first_name":"Ada"}`)

	require.NoError(t, h.Accept(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "access-token")
	require.Len(t, rec.Result().Cookies(), 1)
	require.Equal(t, constants.SessionCookieName, rec.Result().Cookies()[0].Name)
	require.Equal(t, "invite-session", rec.Result().Cookies()[0].Value)
}

func TestMemberInviteHandler_Accept_MFA(t *testing.T) {
	h := newMemberInviteHandler(&stubMemberInviteService{
		accept: &services.MemberInviteAcceptResult{
			UserID: testUserID,
			MFA:    &dtos.MFALoginChallengeResponse{MfaRequired: true, MfaTicket: "ticket-1", ExpiresIn: 600},
		},
	}, &stubAuthSessionService{})
	c, rec := newJSONContext(http.MethodPost, "/api/v1/invites/accept", `{"token":"raw-token","password":"long-secure-passphrase","first_name":"Ada"}`)

	require.NoError(t, h.Accept(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"mfa_required":true`)
	require.Empty(t, rec.Result().Cookies())
}

func TestMemberInviteHandler_Accept_Errors(t *testing.T) {
	t.Run("invalid body", func(t *testing.T) {
		h := newMemberInviteHandler(&stubMemberInviteService{}, &stubAuthSessionService{})
		c, rec := newJSONContext(http.MethodPost, "/api/v1/invites/accept", `{`)
		require.NoError(t, h.Accept(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation", func(t *testing.T) {
		h := newMemberInviteHandler(&stubMemberInviteService{}, &stubAuthSessionService{})
		c, rec := newJSONContext(http.MethodPost, "/api/v1/invites/accept", `{"password":"short"}`)
		require.NoError(t, h.Accept(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		h := newMemberInviteHandler(&stubMemberInviteService{
			acceptErr: errors.ForbiddenError("Invite is invalid or expired", nil),
		}, &stubAuthSessionService{})
		c, rec := newJSONContext(http.MethodPost, "/api/v1/invites/accept", `{"token":"raw-token","password":"long-secure-passphrase"}`)
		require.NoError(t, h.Accept(c))
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("empty result", func(t *testing.T) {
		h := newMemberInviteHandler(&stubMemberInviteService{
			accept: &services.MemberInviteAcceptResult{},
		}, &stubAuthSessionService{})
		c, rec := newJSONContext(http.MethodPost, "/api/v1/invites/accept", `{"token":"raw-token","password":"long-secure-passphrase"}`)
		require.NoError(t, h.Accept(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("session error", func(t *testing.T) {
		h := newMemberInviteHandler(&stubMemberInviteService{
			accept: &services.MemberInviteAcceptResult{
				UserID: testUserID,
				Login:  &dtos.LoginResponse{AccessToken: "access-token", ActiveTenantID: testTenantID},
			},
		}, &stubAuthSessionService{createErr: errors.InternalError("session failed", nil)})
		h.cfg = &config.Config{AppEnv: config.EnvironmentProduction}
		c, rec := newJSONContext(http.MethodPost, "/api/v1/invites/accept", `{"token":"raw-token","password":"long-secure-passphrase"}`)
		require.NoError(t, h.Accept(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
