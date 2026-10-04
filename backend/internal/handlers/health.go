package handlers

import (
	"net/http"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/cache"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"

	"github.com/labstack/echo/v4"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	BaseHandler
	cfg   *config.Config
	db    *db.PostgresDB
	cache cache.Cache
}

// NewHealthHandler creates a new health handler
func ProvideHealthHandler(cfg *config.Config, database *db.PostgresDB, shared cache.Cache) *HealthHandler {
	return &HealthHandler{
		BaseHandler: *NewBaseHandler(),
		cfg:         cfg,
		db:          database,
		cache:       shared,
	}
}

// HealthCheck godoc
// @Summary Health Check
// @Description Check if the service is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.HealthResponse}
// @Router / [get]
func (h *HealthHandler) HealthCheck(c echo.Context) error {
	healthResponse := dtos.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Version:   h.cfg.AppVersion,
		Service:   h.cfg.AppName,
	}

	response := h.SuccessResponse(c, "Service is healthy", healthResponse, nil)
	return response
}

// Ready godoc
// @Summary Readiness check
// @Description Ping Postgres and Redis
// @Tags Health
// @Produce json
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.HealthResponse}
// @Failure 500 {object} object{meta=dtos.Meta}
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c echo.Context) error {
	ctx := c.Request().Context()
	if h.db == nil || h.db.DB == nil || h.cache == nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}
	sqlDB, err := h.db.DB.DB()
	if err != nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}
	if err := h.cache.Ping(ctx); err != nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}
	return h.SuccessResponse(c, "Service is ready", dtos.HealthResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC(),
		Version:   h.cfg.AppVersion,
		Service:   h.cfg.AppName,
	}, nil)
}

// DatabaseHealthCheck godoc
// @Summary Database Health Check
// @Description Check if the database connection is healthy
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} object{meta=dtos.Meta,data=object}
// @Router /health/database [get]
func (h *HealthHandler) DatabaseHealthCheck(c echo.Context) error {
	if h.db == nil {
		return h.InternalErrorResponse(c, "Database not initialized", nil)
	}

	// Use fast, cached health checks to avoid excessive load on the database
	healthStatus := h.db.FastHealthCheck()
	metrics := h.db.GetMetrics()

	response := map[string]interface{}{
		"database_health":    healthStatus,
		"connection_metrics": metrics,
		"timestamp":          time.Now().UTC(),
	}

	if healthStatus.IsHealthy {
		return h.SuccessResponse(c, "Database is healthy", response, nil)
	}
	return h.InternalErrorResponse(c, "Database is unhealthy", nil)
}

// DatabaseMetrics godoc
// @Summary Database Connection Metrics
// @Description Get detailed database connection metrics
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} object{meta=dtos.Meta,data=object}
// @Router /health/metrics [get]
func (h *HealthHandler) DatabaseMetrics(c echo.Context) error {
	if h.db == nil {
		return h.InternalErrorResponse(c, "Database not initialized", nil)
	}

	metrics := h.db.GetMetrics()
	healthStatus := h.db.HealthCheck()

	response := map[string]interface{}{
		"connection_metrics": metrics,
		"health_status":      healthStatus,
		"timestamp":          time.Now().UTC(),
		"configuration": map[string]interface{}{
			"max_open_connections": h.cfg.DatabaseMaxOpenConns,
			"max_idle_connections": h.cfg.DatabaseMaxIdleConns,
			"conn_max_lifetime":    h.cfg.DatabaseConnMaxLifetime.String(),
			"conn_max_idle_time":   h.cfg.DatabaseConnMaxIdleTime.String(),
			"connect_timeout":      h.cfg.DatabaseConnectTimeout.String(),
			"query_timeout":        h.cfg.DatabaseQueryTimeout.String(),
		},
	}

	return h.SuccessResponse(c, "Database metrics retrieved successfully", response, nil)
}
