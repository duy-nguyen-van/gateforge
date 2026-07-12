package middlewares

import (
	"net/http"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/labstack/echo/v4"
)

// PlatformAdminAuth requires JWTBearerAuth upstream and verifies the user is a platform admin (DB flag).
func PlatformAdminAuth(users repositories.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get(auth.EchoContextUserIDKey).(string)
			if !ok || userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
			}

			user, err := users.GetOneByID(c.Request().Context(), userID)
			if err != nil {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
			}
			if !user.IsPlatformAdmin {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
			}
			return next(c)
		}
	}
}
