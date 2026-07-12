package handlers

import (
	"net/http"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// MFAHandler exposes MFA enrollment and post-login challenge endpoints.
type MFAHandler struct {
	BaseHandler
	mfaService     services.MFAService
	userService    services.UserService
	sessionService services.SessionService
	cfg            *config.Config
	validator      *validator.Validate
}

func ProvideMFAHandler(
	mfaService services.MFAService,
	userService services.UserService,
	sessionService services.SessionService,
	cfg *config.Config,
	validator *validator.Validate,
) *MFAHandler {
	return &MFAHandler{
		BaseHandler:    *NewBaseHandler(),
		mfaService:     mfaService,
		userService:    userService,
		sessionService: sessionService,
		cfg:            cfg,
		validator:      validator,
	}
}

// TOTPSetup godoc
// @Summary Start TOTP enrollment
// @Description Generates a TOTP secret, persists it (encrypted) as pending, and returns the raw secret plus otpauth URI for QR scanning. Call POST /mfa/totp/verify with a valid code to enable MFA.
// @Tags MFA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.MFATOTPSetupResponse}
// @Failure 401 {object} object{meta=dtos.Meta}
// @Failure 500 {object} object{meta=dtos.Meta}
// @Router /mfa/totp/setup [post]
func (h *MFAHandler) TOTPSetup(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	u, err := h.userService.GetOneByID(c.Request().Context(), uid)
	if err != nil {
		return h.HandleError(c, err)
	}
	secret, uri, err := h.mfaService.SetupTOTP(c.Request().Context(), uid, u.Email)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "TOTP setup started", dtos.MFATOTPSetupResponse{
		Secret:     secret,
		OtpauthURI: uri,
	}, nil)
}

// TOTPVerifyEnrollment godoc
// @Summary Confirm TOTP enrollment
// @Description Validates a 6-digit code against the pending secret and enables TOTP for the user.
// @Tags MFA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dtos.MFATOTPVerifyRequest true "6-digit authenticator code"
// @Success 200 {object} object{meta=dtos.Meta,data=object}
// @Failure 400 {object} object{meta=dtos.Meta}
// @Failure 401 {object} object{meta=dtos.Meta}
// @Router /mfa/totp/verify [post]
func (h *MFAHandler) TOTPVerifyEnrollment(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	var req dtos.MFATOTPVerifyRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		return h.validationErr(c, err)
	}
	if err := h.mfaService.VerifyTOTPEnrollment(c.Request().Context(), uid, req.Code); err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "TOTP enabled", map[string]any{}, nil)
}

// RecoveryCodes godoc
// @Summary Generate MFA recovery codes
// @Description Replaces any existing recovery codes. Requires TOTP to be enabled. Plain codes are returned once; only hashes are stored.
// @Tags MFA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.MFARecoveryCodesResponse}
// @Failure 400 {object} object{meta=dtos.Meta}
// @Failure 401 {object} object{meta=dtos.Meta}
// @Failure 403 {object} object{meta=dtos.Meta}
// @Router /mfa/recovery-codes [post]
func (h *MFAHandler) RecoveryCodes(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	codes, err := h.mfaService.RegenerateRecoveryCodes(c.Request().Context(), uid)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Recovery codes generated", dtos.MFARecoveryCodesResponse{Codes: codes}, nil)
}

// ChallengeVerify godoc
// @Summary Complete MFA after password or passkey login
// @Description Consumes the `mfa_ticket` from password or passkey login when MFA is required. `code` is a 6-digit TOTP, or a recovery code (non–6-digit path). On success returns LoginResponse and sets `iam_session` cookie.
// @Tags MFA
// @Accept json
// @Produce json
// @Param body body dtos.MFAChallengeVerifyRequest true "Ticket and code"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.LoginResponse}
// @Failure 400 {object} object{meta=dtos.Meta}
// @Failure 401 {object} object{meta=dtos.Meta}
// @Router /mfa/challenge/verify [post]
func (h *MFAHandler) ChallengeVerify(c echo.Context) error {
	var req dtos.MFAChallengeVerifyRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		return h.validationErr(c, err)
	}
	payload, err := h.mfaService.VerifyLoginChallenge(c.Request().Context(), req.MfaTicket, req.Code)
	if err != nil {
		return h.HandleError(c, err)
	}
	u, err := h.userService.GetOneByID(c.Request().Context(), payload.UserID)
	if err != nil {
		return h.HandleError(c, err)
	}
	out, err := h.userService.IssueTokensForUser(c.Request().Context(), u, payload.TenantID)
	if err != nil {
		return h.HandleError(c, err)
	}
	sid, ttl, err := h.sessionService.Create(c.Request().Context(), u.ID, payload.TenantID, c.RealIP(), c.Request().UserAgent(), payload.RememberMe)
	if err != nil {
		return h.HandleError(c, err)
	}
	c.SetCookie(&http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cfg.AppEnv == config.EnvironmentProduction,
		MaxAge:   int(ttl.Seconds()),
	})
	return h.SuccessResponse(c, "Login successful", out, nil)
}

func (h *MFAHandler) validationErr(c echo.Context, err error) error {
	fieldErrors := errors.ParseValidationErrors(err)
	if len(fieldErrors) > 0 {
		return h.HandleError(c, errors.ValidationErrorWithDetails("Validation failed", err, fieldErrors))
	}
	return h.HandleError(c, errors.ValidationError("Validation failed", err))
}
