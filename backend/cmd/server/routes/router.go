package routes

import (
	"errors"
	"net/http"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/handlers"
	middlewares "github.com/gateforge-iam/gateforge-iam/internal/middlewares"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
	"github.com/gateforge-iam/gateforge-iam/internal/static"

	appErrors "github.com/gateforge-iam/gateforge-iam/internal/errors"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Router(
	authHandler *handlers.AuthHandler,
	healthHandler *handlers.HealthHandler,
	oidcHandler *handlers.OIDCHandler,
	tenantIdentityAdmin *handlers.TenantIdentityAdminHandler,
	adminHandler *handlers.AdminHandler,
	webauthnHandler *handlers.WebauthnHandler,
	mfaHandler *handlers.MFAHandler,
	tokenService *auth.TokenService,
	userRepo repositories.UserRepository,
	cfg *config.Config,
) *echo.Echo {
	r := echo.New()

	r.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	r.Use(sentryCaptureMiddleware(cfg))
	registerGlobalMiddleware(r, cfg)

	if cfg.AppEnv != config.EnvironmentProduction {
		r.GET("/swagger/*", echoSwagger.WrapHandler, middlewares.BasicAuthMiddleware(*cfg))
	}

	registerOIDCRoutes(r, oidcHandler)

	v1 := r.Group("api/v1")
	registerPublicV1Routes(v1, healthHandler, authHandler, webauthnHandler, mfaHandler, tenantIdentityAdmin)

	authJWT := middlewares.JWTBearerAuth(tokenService)
	adminAuth := middlewares.PlatformAdminAuth(userRepo)
	registerAuthenticatedV1Routes(v1, authHandler, webauthnHandler, mfaHandler, authJWT)
	registerAdminV1Routes(v1, adminHandler, authJWT, adminAuth)

	if cfg.ServeEmbeddedFrontend {
		distPath := cfg.FrontendDistPath
		if distPath == "" && cfg.AppEnv.IsDevelopment() {
			distPath = static.ResolveDevDistPath()
		}
		static.Register(r, static.Dist, distPath)
	}

	return r
}

func sentryCaptureMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			var httpErr *echo.HTTPError
			if errors.As(err, &httpErr) && (httpErr.Code == http.StatusNotFound || httpErr.Code == http.StatusMethodNotAllowed) {
				return err
			}

			hub := sentryecho.GetHubFromContext(c)
			if hub == nil {
				return err
			}

			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("method", c.Request().Method)
				scope.SetExtra("path", c.Request().URL.Path)
				scope.SetExtra("query", c.QueryParams())
				scope.SetExtra("headers", c.Request().Header)
				scope.SetExtra("body", c.Get("log_body"))
				scope.SetTag("environment", cfg.AppEnv.String())
				scope.SetTag("service", cfg.AppName)
				scope.SetTag("handler", c.Path())
				if orgID := c.Get("organization_id"); orgID != nil {
					scope.SetTag("organization_id", orgID.(string))
				}
				if errors.As(err, &httpErr) {
					scope.SetTag("error_type", "http_error")
					scope.SetExtra("http_code", httpErr.Code)
				} else {
					scope.SetTag("error_type", "internal_error")
				}
				hub.CaptureException(err)
			})
			return err
		}
	}
}

func registerGlobalMiddleware(r *echo.Echo, cfg *config.Config) {
	r.Use(middlewares.LogBodyMiddleware)
	r.Use(middleware.RequestID())
	r.Use(middlewares.RequestContext(cfg.AppName))
	r.Use(middlewares.AuditContext())
	r.Use(appErrors.RecoveryMiddleware(cfg))
	r.Use(appErrors.ErrorMiddleware())
	r.Use(middlewares.Security())
	r.Use(middlewares.CORS())
	r.Use(middlewares.CSRF(cfg))
	r.Use(middlewares.ExposeCSRFToken())
	r.Use(middlewares.DefaultRateLimit())
	r.Use(middlewares.RequestLogging(cfg))
}

func registerOIDCRoutes(r *echo.Echo, oidcHandler *handlers.OIDCHandler) {
	r.GET("/.well-known/jwks.json", oidcHandler.JWKS)
	r.GET("/.well-known/openid-configuration", oidcHandler.OpenIDConfiguration)
	r.GET("/authorize", oidcHandler.Authorize)
	r.POST("/oidc/login", oidcHandler.Login)
	r.GET("/oidc/federation/:provider/start", oidcHandler.FederationOAuthStart)
	r.GET("/oidc/federation/:provider/callback", oidcHandler.FederationOAuthCallback)
	r.POST("/token", oidcHandler.Token)
	r.GET("/userinfo", oidcHandler.UserInfo)
}

func registerPublicV1Routes(
	v1 *echo.Group,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	webauthnHandler *handlers.WebauthnHandler,
	mfaHandler *handlers.MFAHandler,
	tenantIdentityAdmin *handlers.TenantIdentityAdminHandler,
) {
	publicGroup := v1.Group("")
	publicGroup.GET("/", healthHandler.HealthCheck)
	publicGroup.GET("/health/database", healthHandler.DatabaseHealthCheck)
	publicGroup.GET("/health/metrics", healthHandler.DatabaseMetrics)

	publicGroup.POST("/register", authHandler.Register, middlewares.AuthRateLimit())
	publicGroup.POST("/login", authHandler.Login, middlewares.AuthRateLimit())
	publicGroup.POST("/login/session", authHandler.ExchangeSession, middlewares.AuthRateLimit())
	publicGroup.GET("/federation/providers", authHandler.ListFederationProviders)
	publicGroup.POST("/refresh", authHandler.Refresh)

	publicGroup.POST("/webauthn/login/start", webauthnHandler.LoginStart, middlewares.AuthRateLimit())
	publicGroup.POST("/webauthn/login/finish", webauthnHandler.LoginFinish, middlewares.AuthRateLimit())
	publicGroup.POST("/mfa/challenge/verify", mfaHandler.ChallengeVerify, middlewares.AuthRateLimit())

	publicGroup.PATCH("/internal/tenants/:tenantId/identity-providers/:provider", tenantIdentityAdmin.PatchIdentityProvider)
}

func registerAuthenticatedV1Routes(
	v1 *echo.Group,
	authHandler *handlers.AuthHandler,
	webauthnHandler *handlers.WebauthnHandler,
	mfaHandler *handlers.MFAHandler,
	authJWT echo.MiddlewareFunc,
) {
	v1.POST("/logout", authHandler.Logout, authJWT)
	v1.GET("/me", authHandler.Me, authJWT)
	v1.PATCH("/me", authHandler.UpdateMe, authJWT)
	v1.GET("/me/tenants", authHandler.ListMyTenants, authJWT)
	v1.POST("/tenants/select", authHandler.SelectTenant)
	v1.POST("/tenants/switch", authHandler.SwitchTenant, authJWT)
	v1.GET("/webauthn/credentials", webauthnHandler.ListCredentials, authJWT)
	v1.POST("/webauthn/register/start", webauthnHandler.RegisterStart, authJWT)
	v1.POST("/webauthn/register/finish", webauthnHandler.RegisterFinish, authJWT)
	v1.POST("/mfa/totp/setup", mfaHandler.TOTPSetup, authJWT)
	v1.POST("/mfa/totp/verify", mfaHandler.TOTPVerifyEnrollment, authJWT)
	v1.POST("/mfa/recovery-codes", mfaHandler.RecoveryCodes, authJWT)
}

func registerAdminV1Routes(
	v1 *echo.Group,
	adminHandler *handlers.AdminHandler,
	authJWT echo.MiddlewareFunc,
	adminAuth echo.MiddlewareFunc,
) {
	adminGroup := v1.Group("/admin", authJWT, adminAuth)
	adminGroup.GET("/stats", adminHandler.GetStats)
	adminGroup.GET("/users", adminHandler.ListUsers)
	adminGroup.GET("/users/:userId", adminHandler.GetUser)
	adminGroup.POST("/users/:userId/disable", adminHandler.DisableUser)
	adminGroup.POST("/users/:userId/force-logout", adminHandler.ForceLogoutUser)
	adminGroup.POST("/users/:userId/reset-passkey", adminHandler.ResetUserPasskeys)
	adminGroup.POST("/users/:userId/reset-mfa", adminHandler.ResetUserMFA)
	adminGroup.GET("/tenants", adminHandler.ListTenants)
	adminGroup.POST("/tenants", adminHandler.CreateTenant)
	adminGroup.GET("/tenants/:tenantId", adminHandler.GetTenant)
	adminGroup.PATCH("/tenants/:tenantId", adminHandler.UpdateTenant)
	adminGroup.DELETE("/tenants/:tenantId", adminHandler.DeleteTenant)
	adminGroup.GET("/clients", adminHandler.ListClients)
	adminGroup.GET("/clients/:clientId/usage", adminHandler.GetClientUsage)
	adminGroup.POST("/clients", adminHandler.CreateClient)
	adminGroup.GET("/clients/:clientId", adminHandler.GetClient)
	adminGroup.PATCH("/clients/:clientId", adminHandler.UpdateClient)
	adminGroup.DELETE("/clients/:clientId", adminHandler.DeleteClient)
	adminGroup.GET("/audit-logs", adminHandler.ListAuditLogs)
	adminGroup.GET("/login-history", adminHandler.ListLoginHistory)
	adminGroup.GET("/tenants/:tenantId/identity-providers", adminHandler.ListIdentityProviders)
	adminGroup.PATCH("/tenants/:tenantId/identity-providers/:provider", adminHandler.PatchIdentityProvider)
	adminGroup.GET("/tenants/:tenantId/members", adminHandler.ListTenantMembers)
	adminGroup.POST("/tenants/:tenantId/members", adminHandler.AddMember)
	adminGroup.DELETE("/tenants/:tenantId/members/:userId", adminHandler.RemoveMember)
}
