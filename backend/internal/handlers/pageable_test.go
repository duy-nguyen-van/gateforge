package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestPageableFromQuery_defaults(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	pr := pageableFromQuery(c)
	require.Equal(t, 1, pr.Page)
	require.Equal(t, constants.DefaultPageSize, pr.PageSize)
}

func TestPageableFromQuery_customValues(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items?page=3&page_size=50", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	pr := pageableFromQuery(c)
	require.Equal(t, 3, pr.Page)
	require.Equal(t, 50, pr.PageSize)
}

func TestPageableFromQuery_clampsPageSize(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items?page_size=500", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	pr := pageableFromQuery(c)
	require.Equal(t, constants.MaxPageSize, pr.PageSize)
}

func TestPageableFromQuery_ignoresInvalidValues(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items?page=0&page_size=-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	pr := pageableFromQuery(c)
	require.Equal(t, 1, pr.Page)
	require.Equal(t, constants.DefaultPageSize, pr.PageSize)
}
