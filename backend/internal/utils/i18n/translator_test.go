package i18n

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestT_ResolvesKnownErrorCode(t *testing.T) {
	Init()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(request.NewLanguageCodeContext(req.Context(), "en"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	msg := T(c, constants.InvalidLoginCredentials, nil)
	require.Equal(t, "Invalid email or password", msg)
}

func TestTFromContext_ResolvesKnownErrorCode(t *testing.T) {
	Init()

	ctx := request.NewLanguageCodeContext(context.Background(), "en")
	msg := TFromContext(ctx, constants.InvalidRefreshToken, nil)
	require.Equal(t, "Invalid or expired refresh token", msg)
}

func TestT_ReturnsCodeForUnknownMessage(t *testing.T) {
	Init()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	msg := T(c, "DOES_NOT_EXIST", nil)
	require.Equal(t, "DOES_NOT_EXIST", msg)
}

func TestT_FallsBackToEnglishLocale(t *testing.T) {
	Init()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(request.NewLanguageCodeContext(req.Context(), "en-US"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	msg := T(c, constants.Unauthorized, nil)
	require.Equal(t, "Unauthorized", msg)
}
