package middlewares

import (
	"net/http"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"

	"github.com/labstack/echo/v4"
)

// JWTBearerAuth validates Authorization: Bearer <JWT> and sets user_id and tenant_id.
func JWTBearerAuth(ts *auth.TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
			}
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token required"})
			}
			userID, tenantID, err := ts.ParseAccessToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			}
			c.Set(auth.EchoContextUserIDKey, userID)
			if tenantID != "" {
				c.Set(auth.EchoContextTenantIDKey, tenantID)
			}
			return next(c)
		}
	}
}
