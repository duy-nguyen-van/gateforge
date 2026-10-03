package middlewares

import (
	"net/http"

	"github.com/gateforge-iam/gateforge-iam/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CORS(cfg *config.Config) echo.MiddlewareFunc {
	origins := []string{"http://localhost:5173", "http://localhost:3000"}
	if cfg != nil && len(cfg.CORSAllowedOrigins) > 0 {
		origins = cfg.CORSAllowedOrigins
	}
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     origins,
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-CSRF-Token"},
		AllowMethods: []string{
			http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch,
			http.MethodPost, http.MethodDelete, http.MethodOptions,
		},
		ExposeHeaders: []string{"X-CSRF-Token"},
	})
}
