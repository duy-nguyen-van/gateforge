package middlewares

import (
	"net/http"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/labstack/echo/v4"
)

// RequireAdminMFA runs after PlatformAdminAuth. Admin APIs stay closed until the
// caller has a verified TOTP factor or at least one passkey.
func RequireAdminMFA(totp repositories.UserMFATOTPRepository, creds repositories.WebauthnCredentialRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get(auth.EchoContextUserIDKey).(string)
			if !ok || userID == "" {
				return metaError(c, http.StatusUnauthorized, constants.Unauthorized, "User not authenticated")
			}
			if totp == nil || creds == nil {
				return metaError(c, http.StatusForbidden, constants.AdminMFARequired, "Admin MFA enrollment is required")
			}
			ctx := c.Request().Context()
			row, err := totp.GetActiveByUserID(ctx, userID)
			if err != nil {
				return metaError(c, http.StatusForbidden, constants.AdminMFARequired, "Admin MFA enrollment is required")
			}
			if row != nil {
				return next(c)
			}
			list, err := creds.ListByUserID(ctx, userID)
			if err != nil || len(list) == 0 {
				return metaError(c, http.StatusForbidden, constants.AdminMFARequired, "Admin MFA enrollment is required")
			}
			return next(c)
		}
	}
}
