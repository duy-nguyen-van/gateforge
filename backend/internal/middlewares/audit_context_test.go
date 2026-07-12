package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestAuditContext(t *testing.T) {
	t.Run("enriches context with actor and tenant", func(t *testing.T) {
		e := echo.New()
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				ctx := request.NewCorrelationIDContext(c.Request().Context(), "corr-abc")
				c.SetRequest(c.Request().WithContext(ctx))
				c.Set(echo.HeaderXRequestID, "req-123")
				c.Set(auth.EchoContextUserIDKey, "user-42")
				c.Set(auth.EchoContextTenantIDKey, "tenant-99")
				return next(c)
			}
		})
		e.Use(AuditContext())
		e.GET("/audit", func(c echo.Context) error {
			ac, ok := request.AuditContextFromContext(c.Request().Context())
			require.True(t, ok)
			require.Equal(t, "req-123", ac.RequestID)
			require.Equal(t, "corr-abc", ac.CorrelationID)
			require.Equal(t, string(constants.AuditActorTypeUser), ac.ActorType)
			require.Equal(t, "user-42", ac.ActorID)
			require.Equal(t, "tenant-99", ac.TenantID)
			require.NotEmpty(t, ac.IPAddress)
			return c.NoContent(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/audit", nil)
		req.Header.Set("User-Agent", "test-agent/1.0")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("works without authenticated user", func(t *testing.T) {
		e := echo.New()
		e.Use(AuditContext())
		e.GET("/public", func(c echo.Context) error {
			ac, ok := request.AuditContextFromContext(c.Request().Context())
			require.True(t, ok)
			require.Empty(t, ac.ActorID)
			require.Empty(t, ac.ActorType)
			return c.NoContent(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})
}
