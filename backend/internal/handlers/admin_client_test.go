package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type stubAdminClientService struct {
	stubAdminAuditService
}

func (stubAdminClientService) GetClientByID(_ context.Context, clientID string) (*dtos.AdminClientResponse, error) {
	return &dtos.AdminClientResponse{ID: clientID, ClientID: "my-app", Name: "My App"}, nil
}

func (stubAdminClientService) CreateClient(_ context.Context, req *dtos.AdminCreateClientRequest) (*dtos.AdminCreateClientResponse, error) {
	return &dtos.AdminCreateClientResponse{
		AdminClientResponse: dtos.AdminClientResponse{
			ClientID: req.ClientID,
			Name:     req.Name,
			IsPublic: req.IsPublic,
		},
		ClientSecret: "generated-secret",
	}, nil
}

func (stubAdminClientService) UpdateClient(_ context.Context, clientID string, req *dtos.AdminUpdateClientRequest) (*dtos.AdminClientResponse, error) {
	resp := &dtos.AdminClientResponse{ID: clientID, ClientID: "my-app"}
	if req.Name != nil {
		resp.Name = *req.Name
	}
	return resp, nil
}

func (stubAdminClientService) DeleteClient(context.Context, string) error { return nil }

func TestAdminHandler_CreateClient(t *testing.T) {
	h := ProvideAdminHandler(stubAdminClientService{}, validator.New())
	e := echo.New()
	body := `{"tenant_id":"00000000-0000-0000-0000-000000000001","client_id":"my-app","name":"My App","is_public":false,"redirect_uris":["http://localhost/callback"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/clients", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.CreateClient(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Contains(t, rec.Body.String(), "My App")
	require.Contains(t, rec.Body.String(), "generated-secret")
}

func TestAdminHandler_GetClient(t *testing.T) {
	h := ProvideAdminHandler(stubAdminClientService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/clients/client-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("clientId")
	c.SetParamValues("client-1")

	err := h.GetClient(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "My App")
}

func TestAdminHandler_UpdateClient(t *testing.T) {
	h := ProvideAdminHandler(stubAdminClientService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/clients/client-1", strings.NewReader(`{"name":"Updated"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("clientId")
	c.SetParamValues("client-1")

	err := h.UpdateClient(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Updated")
}

func TestAdminHandler_DeleteClient(t *testing.T) {
	h := ProvideAdminHandler(stubAdminClientService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/clients/client-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("clientId")
	c.SetParamValues("client-1")

	err := h.DeleteClient(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
}
