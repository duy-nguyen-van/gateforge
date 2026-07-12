package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/db"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	cfg := handlerTestConfig()
	h := ProvideHealthHandler(cfg, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.HealthCheck(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"status":"healthy"`)
	require.Contains(t, rec.Body.String(), cfg.AppVersion)
}

func TestHealthHandler_DatabaseHealthCheck_nilDB(t *testing.T) {
	h := ProvideHealthHandler(handlerTestConfig(), nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health/database", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.DatabaseHealthCheck(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHealthHandler_DatabaseHealthCheck_unhealthy(t *testing.T) {
	h := ProvideHealthHandler(handlerTestConfig(), &db.PostgresDB{})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health/database", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.DatabaseHealthCheck(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.Contains(t, rec.Body.String(), constants.InternalError)
}

func TestHealthHandler_DatabaseMetrics_nilDB(t *testing.T) {
	h := ProvideHealthHandler(handlerTestConfig(), nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health/metrics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.DatabaseMetrics(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHealthHandler_DatabaseHealthCheck_healthy(t *testing.T) {
	h := ProvideHealthHandler(handlerTestConfig(), db.NewHealthyPostgresDBStub())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health/database", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.DatabaseHealthCheck(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"database_health"`)
}

func TestHealthHandler_DatabaseMetrics_success(t *testing.T) {
	cfg := handlerTestConfig()
	cfg.DatabaseMaxOpenConns = 25
	cfg.DatabaseMaxIdleConns = 5
	cfg.DatabaseConnMaxLifetime = time.Hour
	cfg.DatabaseConnMaxIdleTime = 30 * time.Minute
	cfg.DatabaseConnectTimeout = 5 * time.Second
	cfg.DatabaseQueryTimeout = 10 * time.Second

	h := ProvideHealthHandler(cfg, &db.PostgresDB{})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health/metrics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.DatabaseMetrics(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "connection_metrics")
	require.Contains(t, rec.Body.String(), `"max_open_connections":25`)
}
