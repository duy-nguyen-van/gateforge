package handlers

import (
	"net/http"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func newMFAHandler() *MFAHandler {
	return ProvideMFAHandler(&stubAuthMFAService{}, &stubAuthUserService{}, &stubAuthSessionService{createSID: "sess-mfa"}, handlerTestConfig(), validator.New())
}

func TestMFAHandler_TOTPSetup_success(t *testing.T) {
	h := newMFAHandler()
	c, rec := newJSONContext(http.MethodPost, "/mfa/totp/setup", "")
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.TOTPSetup(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "otpauth://totp/test")
}

func TestMFAHandler_TOTPSetup_unauthorized(t *testing.T) {
	h := newMFAHandler()
	c, rec := newJSONContext(http.MethodPost, "/mfa/totp/setup", "")

	err := h.TOTPSetup(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMFAHandler_TOTPVerifyEnrollment_success(t *testing.T) {
	h := newMFAHandler()
	c, rec := newJSONContext(http.MethodPost, "/mfa/totp/verify", `{"code":"123456"}`)
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.TOTPVerifyEnrollment(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestMFAHandler_TOTPVerifyEnrollment_validationError(t *testing.T) {
	h := newMFAHandler()
	c, rec := newJSONContext(http.MethodPost, "/mfa/totp/verify", `{"code":"12"}`)
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.TOTPVerifyEnrollment(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMFAHandler_RecoveryCodes_success(t *testing.T) {
	h := newMFAHandler()
	c, rec := newJSONContext(http.MethodPost, "/mfa/recovery-codes", "")
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.RecoveryCodes(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "code-1")
}

func TestMFAHandler_ChallengeVerify_success(t *testing.T) {
	h := newMFAHandler()
	c, rec := newJSONContext(http.MethodPost, "/mfa/challenge/verify", `{"mfa_ticket":"ticket-1","code":"123456"}`)

	err := h.ChallengeVerify(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "mfa-access")
	cookies := rec.Result().Cookies()
	require.NotEmpty(t, cookies)
	require.Equal(t, constants.SessionCookieName, cookies[0].Name)
}

func TestMFAHandler_ChallengeVerify_validationError(t *testing.T) {
	h := newMFAHandler()
	c, rec := newJSONContext(http.MethodPost, "/mfa/challenge/verify", `{}`)

	err := h.ChallengeVerify(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMFAHandler_ChallengeVerify_serviceError(t *testing.T) {
	h := ProvideMFAHandler(
		&stubAuthMFAService{verifyChallengeErr: errors.UnauthorizedError("bad ticket", nil)},
		&stubAuthUserService{},
		&stubAuthSessionService{},
		handlerTestConfig(),
		validator.New(),
	)
	c, rec := newJSONContext(http.MethodPost, "/mfa/challenge/verify", `{"mfa_ticket":"bad","code":"123456"}`)

	err := h.ChallengeVerify(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
