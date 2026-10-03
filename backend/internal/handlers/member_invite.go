package handlers

import (
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// MemberInviteHandler accepts an organization invite and creates the account when needed.
type MemberInviteHandler struct {
	BaseHandler
	invites   services.MemberInviteService
	sessions  services.SessionService
	cfg       *config.Config
	validator *validator.Validate
}

// ProvideMemberInviteHandler wires public invite preview and accept routes.
func ProvideMemberInviteHandler(
	invites services.MemberInviteService,
	sessions services.SessionService,
	cfg *config.Config,
	validator *validator.Validate,
) *MemberInviteHandler {
	return &MemberInviteHandler{
		BaseHandler: *NewBaseHandler(),
		invites:     invites,
		sessions:    sessions,
		cfg:         cfg,
		validator:   validator,
	}
}

// Preview godoc
// @Summary Preview an organization invite
// @Tags Auth
// @Produce json
// @Param token query string true "Invite token"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.MemberInvitePreview}
// @Router /invites/preview [get]
func (h *MemberInviteHandler) Preview(c echo.Context) error {
	preview, err := h.invites.Preview(c.Request().Context(), c.QueryParam("token"))
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Invite retrieved", preview, nil)
}

// Accept godoc
// @Summary Accept an organization invite
// @Description Creates the user when they do not exist, adds the membership, and returns a session. If the existing account has MFA enabled, `data` is MFALoginChallengeResponse instead of LoginResponse.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dtos.AcceptMemberInviteRequest true "Invite acceptance"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.LoginResponse} "When MFA is off: LoginResponse and iam_session cookie. When MFA is on: MFALoginChallengeResponse."
// @Router /invites/accept [post]
func (h *MemberInviteHandler) Accept(c echo.Context) error {
	var req dtos.AcceptMemberInviteRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}
	out, err := h.invites.Accept(c.Request().Context(), &req)
	if err != nil {
		return h.HandleError(c, err)
	}
	if out == nil || (out.MFA == nil && out.Login == nil) {
		return h.HandleError(c, errors.InternalError("Invite acceptance did not produce a session", nil))
	}
	if out.MFA != nil {
		return h.SuccessResponse(c, "MFA required", out.MFA, nil)
	}
	ctx := c.Request().Context()
	sid, ttl, err := h.sessions.Create(ctx, out.UserID, out.Login.ActiveTenantID, c.RealIP(), c.Request().UserAgent(), false)
	if err != nil {
		return h.HandleError(c, err)
	}
	c.SetCookie(newSessionCookie(h.cfg.AppEnv == config.EnvironmentProduction, sid, int(ttl.Seconds())))
	return h.SuccessResponse(c, "Signed in", out.Login, nil)
}
