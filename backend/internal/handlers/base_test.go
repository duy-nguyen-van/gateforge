package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestBaseHandler_SuccessResponse(t *testing.T) {
	h := NewBaseHandler()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.SuccessResponse(c, "ok", map[string]string{"key": "value"}, &dtos.Pageable{Page: 1, PageSize: 20, Total: 1})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"key":"value"`)
}

func TestBaseHandler_ValidationErrorResponse(t *testing.T) {
	h := NewBaseHandler()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ValidationErrorResponse(c, "Validation failed", map[string]string{"email": "required"})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), constants.ValidationError)
}

func TestBaseHandler_NotFoundErrorResponse(t *testing.T) {
	h := NewBaseHandler()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.NotFoundErrorResponse(c, "User")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestBaseHandler_UnauthorizedErrorResponse(t *testing.T) {
	h := NewBaseHandler()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.UnauthorizedErrorResponse(c, "Authentication required")
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestBaseHandler_ForbiddenErrorResponse(t *testing.T) {
	h := NewBaseHandler()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ForbiddenErrorResponse(c, "Access denied")
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestBaseHandler_InternalErrorResponse(t *testing.T) {
	h := NewBaseHandler()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.InternalErrorResponse(c, "Something broke", errors.InternalError("db down", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestBaseHandler_HandleError(t *testing.T) {
	h := NewBaseHandler()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.HandleError(c, errors.ValidationError("bad input", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
