package handlers

import (
	"net/http"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// WebauthnHandler exposes WebAuthn passkey endpoints.
type WebauthnHandler struct {
	BaseHandler
	webauthnService services.WebauthnService
	userService     services.UserService
	sessionService  services.SessionService
	mfaService      services.MFAService
	cfg             *config.Config
	validator       *validator.Validate
}

func ProvideWebauthnHandler(
	webauthnService services.WebauthnService,
	userService services.UserService,
	sessionService services.SessionService,
	mfaService services.MFAService,
	cfg *config.Config,
	validator *validator.Validate,
) *WebauthnHandler {
	return &WebauthnHandler{
		BaseHandler:     *NewBaseHandler(),
		webauthnService: webauthnService,
		userService:     userService,
		sessionService:  sessionService,
		mfaService:      mfaService,
		cfg:             cfg,
		validator:       validator,
	}
}

// ListCredentials godoc
// @Summary List registered passkeys
// @Description Returns passkeys registered for the current user (device label and created date only).
// @Tags WebAuthn
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} object{meta=dtos.Meta,data=[]dtos.WebauthnCredentialResponse}
// @Failure 401 {object} object{meta=dtos.Meta}
// @Router /webauthn/credentials [get]
func (h *WebauthnHandler) ListCredentials(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	rows, pageable, err := h.webauthnService.ListCredentials(c.Request().Context(), uid, pageableFromQuery(c))
	if err != nil {
		return h.HandleError(c, err)
	}
	out := make([]dtos.WebauthnCredentialResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, dtos.WebauthnCredentialResponse{
			ID:         row.ID,
			DeviceName: row.DeviceName,
			CreatedAt:  row.CreatedAt,
		})
	}
	return h.SuccessResponse(c, "Passkeys retrieved", out, pageable)
}

// RegisterStart godoc
// @Summary Begin passkey registration
// @Description Returns WebAuthn PublicKeyCredentialCreationOptions as `data.options` and a `session_token` to send with register/finish. Requires a logged-in user (Bearer access token).
// @Tags WebAuthn
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dtos.WebauthnRegisterStartRequest false "Optional device label"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.WebauthnRegisterStartResponse}
// @Failure 400 {object} object{meta=dtos.Meta}
// @Failure 401 {object} object{meta=dtos.Meta}
// @Router /webauthn/register/start [post]
func (h *WebauthnHandler) RegisterStart(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	var req dtos.WebauthnRegisterStartRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		return h.validationErr(c, err)
	}
	raw, token, err := h.webauthnService.RegisterStart(c.Request().Context(), uid, req.DeviceName)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "WebAuthn registration started", dtos.WebauthnRegisterStartResponse{
		Options:      raw,
		SessionToken: token,
	}, nil)
}

// RegisterFinish godoc
// @Summary Finish passkey registration
// @Description Verifies the attestation and stores the passkey for the current user. Body includes `session_token` from register/start and the PublicKeyCredential JSON from the browser.
// @Tags WebAuthn
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dtos.WebauthnRegisterFinishRequest true "Credential response"
// @Success 200 {object} object{meta=dtos.Meta,data=object}
// @Failure 400 {object} object{meta=dtos.Meta}
// @Failure 401 {object} object{meta=dtos.Meta}
// @Router /webauthn/register/finish [post]
func (h *WebauthnHandler) RegisterFinish(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	var req dtos.WebauthnRegisterFinishRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		return h.validationErr(c, err)
	}
	if len(req.Credential) == 0 {
		return h.HandleError(c, errors.ValidationError("credential is required", nil))
	}
	if err := h.webauthnService.RegisterFinish(c.Request().Context(), uid, req.SessionToken, req.Credential); err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Passkey registered", map[string]any{}, nil)
}

// LoginStart godoc
// @Summary Begin passkey login
// @Description Returns WebAuthn PublicKeyCredentialRequestOptions as `data.options` and a `session_token` for login/finish. Public endpoint (rate limited).
// @Tags WebAuthn
// @Accept json
// @Produce json
// @Param body body dtos.WebauthnLoginStartRequest true "User lookup"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.WebauthnLoginStartResponse}
// @Failure 400 {object} object{meta=dtos.Meta}
// @Router /webauthn/login/start [post]
func (h *WebauthnHandler) LoginStart(c echo.Context) error {
	var req dtos.WebauthnLoginStartRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		return h.validationErr(c, err)
	}
	raw, token, err := h.webauthnService.LoginStart(c.Request().Context(), req.Email)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "WebAuthn login started", dtos.WebauthnLoginStartResponse{
		Options:      raw,
		SessionToken: token,
	}, nil)
}

// LoginFinish godoc
// @Summary Finish passkey login
// @Description Verifies the assertion. On success without MFA: returns `data` as LoginResponse and sets `iam_session` cookie. When MFA is enabled for the user, returns MFALoginChallengeResponse in `data` instead (then call POST /mfa/challenge/verify).
// @Tags WebAuthn
// @Accept json
// @Produce json
// @Param body body dtos.WebauthnLoginFinishRequest true "Assertion response"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.LoginResponse} "When MFA is off: LoginResponse + iam_session. When MFA is on: `data` is MFALoginChallengeResponse (see description)."
// @Failure 400 {object} object{meta=dtos.Meta}
// @Failure 401 {object} object{meta=dtos.Meta}
// @Router /webauthn/login/finish [post]
func (h *WebauthnHandler) LoginFinish(c echo.Context) error {
	var req dtos.WebauthnLoginFinishRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		return h.validationErr(c, err)
	}
	if len(req.Credential) == 0 {
		return h.HandleError(c, errors.ValidationError("credential is required", nil))
	}
	u, err := h.webauthnService.LoginFinish(c.Request().Context(), req.Email, req.SessionToken, req.Credential)
	if err != nil {
		return h.HandleError(c, err)
	}

	tenantParam := req.TenantID
	if tenantParam == "" {
		tenantParam = h.cfg.DefaultTenantID
	}

	hasMFA, err := h.mfaService.HasActiveMFA(c.Request().Context(), u.ID)
	if err != nil {
		return h.HandleError(c, err)
	}
	if hasMFA {
		ticket, exp, err := h.mfaService.CreateLoginTicket(c.Request().Context(), auth.MFAPendingPayload{
			UserID:     u.ID,
			TenantID:   tenantParam,
			RememberMe: req.RememberMe,
			ReturnTo:   req.ReturnTo,
		})
		if err != nil {
			return h.HandleError(c, err)
		}
		return h.SuccessResponse(c, "MFA required", dtos.MFALoginChallengeResponse{
			MfaRequired: true,
			MfaTicket:   ticket,
			ExpiresIn:   exp,
		}, nil)
	}

	loginOut, selectionOut, err := h.userService.CompleteAuth(c.Request().Context(), u, services.TenantResolveInput{
		Host:          c.Request().Host,
		TenantIDParam: req.TenantID,
		UserID:        u.ID,
	})
	if err != nil {
		return h.HandleError(c, err)
	}
	if selectionOut != nil {
		return h.SuccessResponse(c, "Tenant selection required", selectionOut, nil)
	}

	sid, ttl, err := h.sessionService.Create(c.Request().Context(), u.ID, loginOut.ActiveTenantID, c.RealIP(), c.Request().UserAgent(), req.RememberMe)
	if err != nil {
		return h.HandleError(c, err)
	}
	h.setSessionCookie(c, sid, ttl)
	return h.SuccessResponse(c, "Login successful", loginOut, nil)
}

func (h *WebauthnHandler) setSessionCookie(c echo.Context, sid string, ttl time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cfg.AppEnv == config.EnvironmentProduction,
		MaxAge:   int(ttl.Seconds()),
	})
}

func (h *WebauthnHandler) validationErr(c echo.Context, err error) error {
	fieldErrors := errors.ParseValidationErrors(err)
	if len(fieldErrors) > 0 {
		return h.HandleError(c, errors.ValidationErrorWithDetails("Validation failed", err, fieldErrors))
	}
	return h.HandleError(c, errors.ValidationError("Validation failed", err))
}
