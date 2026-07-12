package middlewares

import (
	"net/http"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CSRF returns a configured CSRF middleware.
func CSRF(cfg *config.Config) echo.MiddlewareFunc {
	//nolint:gosec // G101: cookie name is not a secret
	return middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper: func(c echo.Context) bool {
			// SPA dashboard calls attach Authorization: Bearer; CSRF protects cookie-only flows.
			if strings.HasPrefix(c.Request().Header.Get("Authorization"), "Bearer ") {
				return true
			}
			p := c.Request().URL.Path
			// Stateless JSON login/register/refresh use Bearer tokens; skip CSRF for these API routes.
			// Note: OIDC browser login is a separate endpoint and should keep CSRF protection.
			if p == "/api/v1/register" || p == "/api/v1/login" || p == "/api/v1/login/session" || p == "/api/v1/refresh" || p == "/api/v1/logout" {
				return true
			}
			if strings.HasPrefix(p, "/api/v1/webauthn/") || strings.HasPrefix(p, "/api/v1/mfa/") {
				return true
			}
			// OIDC/OAuth2 token and discovery endpoints (form POST / JSON GET).
			if p == "/token" || p == "/authorize" || p == "/userinfo" || strings.HasPrefix(p, "/.well-known/") {
				return true
			}
			if strings.HasPrefix(p, "/oidc/federation/") {
				return true
			}
			return false
		},
		TokenLookup:    "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookiePath:     "/",
		CookieSameSite: http.SameSiteLaxMode,
		CookieSecure:   cfg.AppEnv == config.EnvironmentProduction,
		CookieHTTPOnly: true,
	})
}

// ExposeCSRFToken adds the current CSRF token to the response header for clients.
func ExposeCSRFToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if token, ok := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string); ok && token != "" {
				c.Response().Header().Set("X-CSRF-Token", token)
			}
			return next(c)
		}
	}
}
