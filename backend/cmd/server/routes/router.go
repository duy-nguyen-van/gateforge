package routes

import (
	"errors"
	"net/http"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/cache"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/handlers"
	middlewares "github.com/gateforge-iam/gateforge-iam/internal/middlewares"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
	"github.com/gateforge-iam/gateforge-iam/internal/request"
	"github.com/gateforge-iam/gateforge-iam/internal/static"

	appErrors "github.com/gateforge-iam/gateforge-iam/internal/errors"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func Router(
	authHandler *handlers.AuthHandler,
	healthHandler *handlers.HealthHandler,
	oidcHandler *handlers.OIDCHandler,
	tenantIdentityAdmin *handlers.TenantIdentityAdminHandler,
	adminHandler *handlers.AdminHandler,
	webauthnHandler *handlers.WebauthnHandler,
	mfaHandler *handlers.MFAHandler,
	memberInviteHandler *handlers.MemberInviteHandler,
	tokenService *auth.TokenService,
	userRepo repositories.UserRepository,
	cfg *config.Config,
	shared cache.Cache,
	totp repositories.UserMFATOTPRepository,
	webauthn repositories.WebauthnCredentialRepository,
) *echo.Echo {
	r := echo.New()

	if monitoring.IsOTelEnabled(*cfg) && cfg.OTelTracesEnabled {
		serviceName := cfg.OTelServiceName
		if serviceName == "" {
			serviceName = cfg.AppName
		}
		if serviceName == "" {
			serviceName = "gateforge-iam"
		}
		r.Use(otelecho.Middleware(serviceName))
	}

	r.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	r.Use(sentryCaptureMiddleware(cfg))
	registerGlobalMiddleware(r, cfg, shared)

	if cfg.AppEnv != config.EnvironmentProduction {
		r.GET("/swagger/*", echoSwagger.WrapHandler, middlewares.BasicAuthMiddleware(*cfg))
	}

	registerOIDCRoutes(r, oidcHandler)

	v1 := r.Group("api/v1")
	registerPublicV1Routes(v1, healthHandler, authHandler, webauthnHandler, mfaHandler, tenantIdentityAdmin, memberInviteHandler, cfg, shared)

	authJWT := middlewares.JWTBearerAuth(tokenService)
	adminAuth := middlewares.PlatformAdminAuth(userRepo)
	adminMFA := middlewares.RequireAdminMFA(totp, webauthn)
	registerAuthenticatedV1Routes(v1, authHandler, webauthnHandler, mfaHandler, authJWT)
	registerAdminV1Routes(v1, adminHandler, authJWT, adminAuth, adminMFA)

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
				monitoring.SetScopeData(scope, "method", c.Request().Method)
				monitoring.SetScopeData(scope, "path", c.Request().URL.Path)
				monitoring.SetScopeData(scope, "query", c.QueryParams())
				monitoring.SetScopeData(scope, "headers", monitoring.RedactedHeaders(c.Request().Header))
				if !monitoring.SkipSensitiveRequestPath(c.Request().URL.Path) {
					monitoring.SetScopeData(scope, "body", c.Get("log_body"))
				}
				if cid, ok := request.CorrelationIDFromContext(c.Request().Context()); ok && cid != "" {
					scope.SetTag("correlation_id", cid)
				}
				scope.SetTag("environment", cfg.AppEnv.String())
				scope.SetTag("service", cfg.AppName)
				scope.SetTag("handler", c.Path())
				if orgID := c.Get("organization_id"); orgID != nil {
					scope.SetTag("organization_id", orgID.(string))
				}
				if errors.As(err, &httpErr) {
					scope.SetTag("error_type", "http_error")
					monitoring.SetScopeData(scope, "http_code", httpErr.Code)
				} else {
					scope.SetTag("error_type", "internal_error")
				}
				hub.CaptureException(err)
			})
			return err
		}
	}
}

func registerGlobalMiddleware(r *echo.Echo, cfg *config.Config, shared cache.Cache) {
	r.Use(middlewares.LogBodyMiddleware)
	r.Use(middleware.RequestID())
	r.Use(middlewares.RequestContext(cfg.AppName))
	r.Use(middlewares.AuditContext())
	r.Use(appErrors.RecoveryMiddleware(cfg))
	r.Use(appErrors.ErrorMiddleware())
	r.Use(middlewares.Security())
	r.Use(middlewares.CORS(cfg))
	r.Use(middlewares.CSRF(cfg))
	r.Use(middlewares.ExposeCSRFToken())
	r.Use(middlewares.DefaultRateLimit(*cfg, shared))
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
	r.POST("/introspect", oidcHandler.Introspect)
	r.GET("/userinfo", oidcHandler.UserInfo)
}

func registerPublicV1Routes(
	v1 *echo.Group,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	webauthnHandler *handlers.WebauthnHandler,
	mfaHandler *handlers.MFAHandler,
	tenantIdentityAdmin *handlers.TenantIdentityAdminHandler,
	memberInviteHandler *handlers.MemberInviteHandler,
	cfg *config.Config,
	shared cache.Cache,
) {
	authLimit := middlewares.AuthRateLimit(*cfg, shared)
	publicGroup := v1.Group("")
	publicGroup.GET("/", healthHandler.HealthCheck)
	publicGroup.GET("/health/database", healthHandler.DatabaseHealthCheck)
	publicGroup.GET("/health/metrics", healthHandler.DatabaseMetrics)
	publicGroup.GET("/health/ready", healthHandler.Ready)

	publicGroup.GET("/invites/preview", memberInviteHandler.Preview, authLimit)
	publicGroup.POST("/invites/accept", memberInviteHandler.Accept, authLimit)
	publicGroup.POST("/register", authHandler.Register, authLimit)
	publicGroup.POST("/forgot-password", authHandler.ForgotPassword, authLimit)
	publicGroup.POST("/reset-password", authHandler.ResetPassword, authLimit)
	publicGroup.POST("/login", authHandler.Login, authLimit)
	publicGroup.POST("/login/session", authHandler.ExchangeSession, authLimit)
	publicGroup.GET("/federation/providers", authHandler.ListFederationProviders)
	publicGroup.POST("/refresh", authHandler.Refresh)

	publicGroup.POST("/webauthn/login/start", webauthnHandler.LoginStart, authLimit)
	publicGroup.POST("/webauthn/login/finish", webauthnHandler.LoginFinish, authLimit)
	publicGroup.POST("/mfa/challenge/verify", mfaHandler.ChallengeVerify, authLimit)

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
	adminMFA echo.MiddlewareFunc,
) {
	adminGroup := v1.Group("/admin", authJWT, adminAuth, adminMFA)
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
	adminGroup.GET("/invites", adminHandler.ListInvites)
	adminGroup.GET("/tenants/:tenantId/invites", adminHandler.ListTenantInvites)
	adminGroup.POST("/tenants/:tenantId/invites/:inviteId/resend", adminHandler.ResendTenantInvite)
}
