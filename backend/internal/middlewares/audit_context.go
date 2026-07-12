package middlewares

import (
	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"github.com/labstack/echo/v4"
)

// AuditContext enriches request context with metadata used when persisting audit logs.
func AuditContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			ac := request.AuditContext{
				IPAddress: c.RealIP(),
				UserAgent: c.Request().UserAgent(),
			}
			if rid, ok := c.Get(echo.HeaderXRequestID).(string); ok {
				ac.RequestID = rid
			}
			if cid, ok := request.CorrelationIDFromContext(ctx); ok {
				ac.CorrelationID = cid
			}
			if u := c.Get(auth.EchoContextUserIDKey); u != nil {
				if s, ok := u.(string); ok && s != "" {
					ac.ActorType = string(constants.AuditActorTypeUser)
					ac.ActorID = s
				}
			}
			if t := c.Get(auth.EchoContextTenantIDKey); t != nil {
				if s, ok := t.(string); ok {
					ac.TenantID = s
				}
			}

			c.SetRequest(c.Request().WithContext(request.NewAuditContextContext(ctx, ac)))
			return next(c)
		}
	}
}
