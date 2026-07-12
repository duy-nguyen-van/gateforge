package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestAdminHandler_ListAuditLogs(t *testing.T) {
	h := ProvideAdminHandler(stubAdminAuditService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/audit-logs?page=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ListAuditLogs(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), constants.AuditActionAuthLogin)
}
