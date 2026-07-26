package middlewares

import (
	"net/http"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"

	"github.com/labstack/echo/v4"
)

// JWTBearerAuth validates Authorization: Bearer <JWT> and sets user_id and tenant_id.
func JWTBearerAuth(ts *auth.TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return metaError(c, http.StatusUnauthorized, constants.Unauthorized, "Authorization header required")
			}
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return metaError(c, http.StatusUnauthorized, constants.Unauthorized, "Invalid authorization header format")
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				return metaError(c, http.StatusUnauthorized, constants.Unauthorized, "Token required")
			}
			userID, tenantID, err := ts.ParseAccessToken(token)
			if err != nil {
				return metaError(c, http.StatusUnauthorized, constants.Unauthorized, "Invalid or expired token")
			}
			c.Set(auth.EchoContextUserIDKey, userID)
			if tenantID != "" {
				c.Set(auth.EchoContextTenantIDKey, tenantID)
			}
			return next(c)
		}
	}
}

func metaError(c echo.Context, status int, code, message string) error {
	return c.JSON(status, dtos.BaseResponse[any]{
		Meta: dtos.Meta{
			ErrorCode: code,
			Message:   message,
			Code:      status,
		},
		Data: nil,
	})
}
