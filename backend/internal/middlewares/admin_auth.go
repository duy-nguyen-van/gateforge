package middlewares

import (
	"net/http"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/labstack/echo/v4"
)

// PlatformAdminAuth requires JWTBearerAuth upstream and verifies the user is a platform admin (DB flag).
func PlatformAdminAuth(users repositories.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get(auth.EchoContextUserIDKey).(string)
			if !ok || userID == "" {
				return metaError(c, http.StatusUnauthorized, constants.Unauthorized, "User not authenticated")
			}

			user, err := users.GetOneByID(c.Request().Context(), userID)
			if err != nil {
				return metaError(c, http.StatusForbidden, constants.Forbidden, "Admin access required")
			}
			if !user.IsPlatformAdmin {
				return metaError(c, http.StatusForbidden, constants.Forbidden, "Admin access required")
			}
			return next(c)
		}
	}
}
