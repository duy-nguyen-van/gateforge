package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gateforge-iam/gateforge-iam/docs"
	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/cache"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/handlers"
	"github.com/gateforge-iam/gateforge-iam/internal/integration/email"
	"github.com/gateforge-iam/gateforge-iam/internal/integration/storage"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
	"github.com/gateforge-iam/gateforge-iam/internal/services"
	"github.com/gateforge-iam/gateforge-iam/internal/utils/i18n"

	"github.com/gateforge-iam/gateforge-iam/cmd/server/routes"

	"github.com/gateforge-iam/gateforge-iam/internal/db"

	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	oidcHandler *handlers.OIDCHandler,
	tenantIdentityAdmin *handlers.TenantIdentityAdminHandler,
	adminHandler *handlers.AdminHandler,
	webauthnHandler *handlers.WebauthnHandler,
	mfaHandler *handlers.MFAHandler,
	tokenService *auth.TokenService,
	userRepo repositories.UserRepository,
	cfg *config.Config,
	db *db.PostgresDB,
) *http.Server {
	handler := routes.Router(authHandler, healthHandler, oidcHandler, tenantIdentityAdmin, adminHandler, webauthnHandler, mfaHandler, tokenService, userRepo, cfg)

	srv := &http.Server{
		Addr:              cfg.AppHTTPServer,
		Handler:           handler,
		ReadHeaderTimeout: time.Duration(cfg.AppRequestTimeout) * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			logger.Sugar.Infof("Starting HTTP server at %s", srv.Addr)
			go func() {
				err := srv.Serve(ln)
				if err != nil {
					logger.Sugar.Panicf("HTTP server error: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Sugar.Info("Shutting down HTTP server...")

			// Graceful shutdown with timeout
			shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			// Shutdown HTTP server
			if err := srv.Shutdown(shutdownCtx); err != nil {
				logger.Sugar.Errorf("HTTP server shutdown error: %v", err)
				return err
			}

			// Close database connections
			if err := db.Close(); err != nil {
				logger.Sugar.Errorf("Database shutdown error: %v", err)
				return err
			}

			logger.Sugar.Info("Server shutdown completed")
			return nil
		},
	})

	return srv
}

// @title Golang Boilerplate API
// @version 1.0
// @description This is a backend API for Golang Boilerplate
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.basic  BasicAuth
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Bearer Token Authentication. Use "Bearer {token}" as the value.
func main() {
	// Ensure Swagger spec is registered and optionally override fields at runtime
	docs.SwaggerInfo.BasePath = "/api/v1"
	cfg, err := config.Load()
	if err != nil {
		// Use standard log for fatal errors before logger is initialized
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	// Initialize global logger before any middleware uses it
	logger.Init(cfg.LogLevel, cfg.AppEnv.String())
	monitoring.InitNewRelic(*cfg)
	monitoring.InitSentry(*cfg)

	// Ensure all events are flushed before the program exits
	defer monitoring.FlushSentry()

	// Set application timezone from environment variable, default to UTC if not specified
	timezone := cfg.Timezone
	if timezone == "" {
		timezone = "UTC"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logger.Sugar.Warnf("Invalid timezone %s, falling back to UTC", timezone)
		loc = time.UTC
	}
	time.Local = loc

	fx.New(
		fx.Supply(cfg),
		fx.Provide(
			NewHTTPServer,
			ProvideGormPostgres,
			ProvideValidator,
			auth.ProvideTokenService,
			auth.ProvideOIDCSigner,
			auth.ProvideWebAuthn,
			auth.ProvideEphemeralStore,
			cache.ProvideCache,
			email.ProvideEmailSender,
			storage.ProvideStorageAdapter,
			repositories.ProvideUserRepository,
			repositories.ProvideTenantMembershipRepository,
			repositories.ProvideFederatedIdentityRepository,
			repositories.ProvideTenantRepository,
			repositories.ProvideTenantIdentityProviderRepository,
			repositories.ProvideRefreshTokenRepository,
			repositories.ProvideSessionRepository,
			repositories.ProvideClientRepository,
			repositories.ProvideAuthorizationCodeRepository,
			repositories.ProvideWebauthnCredentialRepository,
			repositories.ProvideUserMFATOTPRepository,
			repositories.ProvideUserMFARecoveryCodeRepository,
			repositories.ProvideAuditLogRepository,
			services.ProvideAuditService,
			services.ProvideEmailService,
			services.ProvideMFAService,
			services.ProvideWebauthnService,
			services.ProvideTenantContextService,
			services.ProvideUserService,
			services.ProvideFederationService,
			services.ProvideSessionService,
			services.ProvideOIDCService,
			services.ProvideAdminService,
			services.ProvidePlatformAdminBootstrap,
			handlers.ProvideHealthHandler,
			handlers.ProvideAuthHandler,
			handlers.ProvideWebauthnHandler,
			handlers.ProvideMFAHandler,
			handlers.ProvideTenantIdentityAdminHandler,
			handlers.ProvideAdminHandler,
			handlers.ProvideOIDCHandler,
		),
		fx.Invoke(runPlatformAdminBootstrap),
		fx.Invoke(func(*http.Server) {}),
		fx.Invoke(func() { i18n.Init() }),
	).Run()
}

func runPlatformAdminBootstrap(bootstrap services.PlatformAdminBootstrap) error {
	return bootstrap.Run(context.Background())
}

func ProvideValidator() *validator.Validate {
	return validator.New()
}

func ProvideGormPostgres(cfg *config.Config) *db.PostgresDB {
	appDB := &db.PostgresDB{}
	err := appDB.NewPostgresDB(cfg)
	if err != nil {
		logger.Sugar.Fatalf("Connecting to Database: %v", err)
	}
	return appDB
}
