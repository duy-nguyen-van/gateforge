package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func newWebauthnHandler(webauthn *stubWebauthnService, mfa *stubAuthMFAService) *WebauthnHandler {
	if webauthn == nil {
		webauthn = &stubWebauthnService{}
	}
	if mfa == nil {
		mfa = &stubAuthMFAService{}
	}
	return ProvideWebauthnHandler(webauthn, &stubAuthUserService{}, &stubAuthSessionService{createSID: "sess-wa"}, mfa, handlerTestConfig(), validator.New())
}

func TestWebauthnHandler_ListCredentials_success(t *testing.T) {
	h := newWebauthnHandler(nil, nil)
	c, rec := newJSONContext(http.MethodGet, "/webauthn/credentials?page=1", "")
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.ListCredentials(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "MacBook")
}

func TestWebauthnHandler_ListCredentials_unauthorized(t *testing.T) {
	h := newWebauthnHandler(nil, nil)
	c, rec := newJSONContext(http.MethodGet, "/webauthn/credentials", "")

	err := h.ListCredentials(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestWebauthnHandler_RegisterStart_success(t *testing.T) {
	h := newWebauthnHandler(nil, nil)
	c, rec := newJSONContext(http.MethodPost, "/webauthn/register/start", `{"device_name":"Phone"}`)
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.RegisterStart(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "reg-token")
}

func TestWebauthnHandler_RegisterFinish_success(t *testing.T) {
	h := newWebauthnHandler(nil, nil)
	body := `{"session_token":"tok","credential":{"id":"cred"}}`
	c, rec := newJSONContext(http.MethodPost, "/webauthn/register/finish", body)
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.RegisterFinish(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestWebauthnHandler_RegisterFinish_missingCredential(t *testing.T) {
	h := newWebauthnHandler(nil, nil)
	c, rec := newJSONContext(http.MethodPost, "/webauthn/register/finish", `{"session_token":"tok"}`)
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.RegisterFinish(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestWebauthnHandler_LoginStart_success(t *testing.T) {
	h := newWebauthnHandler(nil, nil)
	c, rec := newJSONContext(http.MethodPost, "/webauthn/login/start", `{"email":"user@example.com"}`)

	err := h.LoginStart(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "login-token")
}

func TestWebauthnHandler_LoginFinish_success(t *testing.T) {
	h := newWebauthnHandler(nil, nil)
	body := `{"email":"user@example.com","session_token":"tok","credential":{"id":"cred"},"remember_me":true}`
	c, rec := newJSONContext(http.MethodPost, "/webauthn/login/finish", body)

	err := h.LoginFinish(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "access-token")
	cookies := rec.Result().Cookies()
	require.NotEmpty(t, cookies)
	require.Equal(t, constants.SessionCookieName, cookies[0].Name)
}

func TestWebauthnHandler_LoginFinish_mfaRequired(t *testing.T) {
	h := newWebauthnHandler(nil, &stubAuthMFAService{hasMFA: true, ticket: "wa-ticket"})
	body := `{"email":"user@example.com","session_token":"tok","credential":{"id":"cred"}}`
	c, rec := newJSONContext(http.MethodPost, "/webauthn/login/finish", body)

	err := h.LoginFinish(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"mfa_required":true`)
	require.Contains(t, rec.Body.String(), "wa-ticket")
}

func TestWebauthnHandler_LoginFinish_tenantSelection(t *testing.T) {
	userSvc := &stubAuthUserService{selection: &dtos.TenantSelectionResponse{SelectionRequired: true, SelectionToken: "sel"}}
	h := ProvideWebauthnHandler(&stubWebauthnService{}, userSvc, &stubAuthSessionService{}, &stubAuthMFAService{}, handlerTestConfig(), validator.New())
	body := `{"email":"user@example.com","session_token":"tok","credential":{"id":"cred"}}`
	c, rec := newJSONContext(http.MethodPost, "/webauthn/login/finish", body)

	err := h.LoginFinish(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"selection_required":true`)
}

func TestWebauthnHandler_LoginFinish_serviceError(t *testing.T) {
	h := newWebauthnHandler(&stubWebauthnService{finishErr: errors.UnauthorizedError("bad assertion", nil)}, nil)
	body := `{"email":"user@example.com","session_token":"tok","credential":{"id":"cred"}}`
	c, rec := newJSONContext(http.MethodPost, "/webauthn/login/finish", body)

	err := h.LoginFinish(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestWebauthnHandler_RegisterFinish_serviceError(t *testing.T) {
	h := newWebauthnHandler(&stubWebauthnService{finishErr: errors.ValidationError("bad attestation", nil)}, nil)
	body := `{"session_token":"tok","credential":` + string(mustRawJSON(map[string]string{"id": "cred"})) + `}`
	c, rec := newJSONContext(http.MethodPost, "/webauthn/register/finish", body)
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.RegisterFinish(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func mustRawJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func TestWebauthnHandler_ListCredentials_mapsResponse(t *testing.T) {
	cred := models.WebauthnCredential{BaseModel: models.NewBaseModel(), DeviceName: "YubiKey"}
	h := newWebauthnHandler(&stubWebauthnService{credentials: []models.WebauthnCredential{cred}}, nil)
	c, rec := newJSONContext(http.MethodGet, "/webauthn/credentials", "")
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.ListCredentials(c)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), "YubiKey")
}
