package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type mfaStub struct {
	active *models.UserMFATOTP
}

func (m mfaStub) GetByUserID(context.Context, string) (*models.UserMFATOTP, error) {
	return nil, nil
}
func (m mfaStub) GetActiveByUserID(context.Context, string) (*models.UserMFATOTP, error) {
	return m.active, nil
}
func (m mfaStub) UpsertPending(context.Context, *models.UserMFATOTP) error { return nil }
func (m mfaStub) MarkVerifiedAndEnabled(context.Context, string) error     { return nil }
func (m mfaStub) Disable(context.Context, string) error                    { return nil }
func (m mfaStub) CountEnabled(context.Context) (int64, error)              { return 0, nil }

type passkeyStub struct {
	rows []models.WebauthnCredential
}

func (p passkeyStub) Create(context.Context, *models.WebauthnCredential) error { return nil }
func (p passkeyStub) ListByUserID(context.Context, string) ([]models.WebauthnCredential, error) {
	return p.rows, nil
}
func (p passkeyStub) ListByUserIDPaginated(context.Context, string, *dtos.PageableRequest) (*dtos.DataResponse[models.WebauthnCredential], error) {
	return nil, nil
}
func (p passkeyStub) GetByCredentialID(context.Context, string) (*models.WebauthnCredential, error) {
	return nil, nil
}
func (p passkeyStub) UpdateCredentialJSON(context.Context, string, string, int64) error { return nil }
func (p passkeyStub) DeleteAllByUserID(context.Context, string) (int64, error)          { return 0, nil }

func TestRequireAdminMFA_WithoutFactor(t *testing.T) {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(auth.EchoContextUserIDKey, "user-1")
			return next(c)
		}
	})
	e.Use(RequireAdminMFA(mfaStub{}, passkeyStub{}))
	e.GET("/admin", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), constants.AdminMFARequired)
}

func TestRequireAdminMFA_WithPasskey(t *testing.T) {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(auth.EchoContextUserIDKey, "user-1")
			return next(c)
		}
	})
	e.Use(RequireAdminMFA(mfaStub{}, passkeyStub{rows: []models.WebauthnCredential{{}}}))
	e.GET("/admin", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}
