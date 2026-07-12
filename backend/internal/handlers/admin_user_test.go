package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type stubAdminUserService struct {
	stubAdminAuditService
	getUserID string
}

func (s *stubAdminUserService) GetUserByID(_ context.Context, userID string) (*dtos.AdminUserDetailResponse, error) {
	if userID == "missing" {
		return nil, errors.NotFoundError("User", nil)
	}
	return &dtos.AdminUserDetailResponse{ID: userID, Email: "user@example.com", Status: "active"}, nil
}

func (s *stubAdminUserService) DisableUser(_ context.Context, actorUserID, targetUserID string) error {
	if targetUserID == "missing" {
		return errors.NotFoundError("User", nil)
	}
	s.getUserID = targetUserID
	return nil
}

func TestAdminHandler_GetUser(t *testing.T) {
	h := ProvideAdminHandler(&stubAdminUserService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/user-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues("user-1")

	err := h.GetUser(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "user@example.com")
}

func TestAdminHandler_DisableUser(t *testing.T) {
	svc := &stubAdminUserService{}
	h := ProvideAdminHandler(svc, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/target-1/disable", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues("target-1")
	c.Set(auth.EchoContextUserIDKey, "actor-1")

	err := h.DisableUser(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, "target-1", svc.getUserID)
}
