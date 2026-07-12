package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Environment represents the application environment.
type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentStaging     Environment = "staging"
	EnvironmentProduction  Environment = "production"
	EnvironmentTest        Environment = "test"
)

func (e Environment) String() string {
	return string(e)
}

func (e Environment) IsDevelopment() bool {
	return e == EnvironmentDevelopment || e == EnvironmentTest
}

func (e Environment) IsProduction() bool {
	return e == EnvironmentProduction
}

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	AppEnv        Environment
	AppName       string
	AppVersion    string
	Timezone      string
	AppHTTPServer string
	// AppRequestTimeout is the HTTP server read header timeout in seconds.
	AppRequestTimeout int
	AppBaseURL        string
	// ServeEmbeddedFrontend serves the Vite SPA from go:embed when true (production single-binary deploy).
	ServeEmbeddedFrontend bool
	// FrontendDistPath serves SPA assets from disk when embed FS is empty (local debug). Auto-detected in development when unset.
	FrontendDistPath string
	// OIDCLoginPageURL is where unauthenticated /authorize redirects (e.g. https://app.example.com/login). Defaults to AppBaseURL + "/login".
	OIDCLoginPageURL string

	// Database configuration
	DatabaseHost        string
	DatabasePort        string
	DatabaseUsername    string
	DatabasePassword    string
	DatabaseName        string
	DatabaseEnableDebug bool

	// Database connection management
	DatabaseMaxOpenConns    int
	DatabaseMaxIdleConns    int
	DatabaseConnMaxLifetime time.Duration
	DatabaseConnMaxIdleTime time.Duration
	DatabaseConnectTimeout  time.Duration
	DatabaseQueryTimeout    time.Duration
	DatabaseHealthTimeout   time.Duration
	DatabaseRetryAttempts   int
	DatabaseRetryDelay      time.Duration
	DatabaseSSLMode         string
	DatabaseTimezone        string

	// Cache configuration
	CacheProvider   string
	RedisHost       string
	RedisPort       string
	RedisPassword   string
	RedisDB         int
	PoolSize        int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolTimeout     time.Duration
	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	// Logging configuration
	LogLevel string

	// Email configuration
	EmailProvider   string
	AWSSESRegion    string
	AWSSESAccessKey string
	AWSSESSecretKey string

	// Environment
	Environment string

	// Rate limiting configuration
	DefaultRateLimit  int
	AuthRateLimit     int
	PublicRateLimit   int
	RateLimit         int
	RateLimitDuration time.Duration

	// NewRelic configuration
	NewRelicAppName string
	NewRelicLicense string

	// Sentry configuration
	SentryDSN string

	// JWT (Phase 1 local auth)
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
	// SSOSessionTTL is the browser SSO cookie (iam_session) lifetime when remember_me is false. Zero means use JWTRefreshTTL.
	SSOSessionTTL time.Duration
	// SSOSessionRememberTTL is iam_session lifetime when remember_me is true. Zero means 30 days.
	SSOSessionRememberTTL time.Duration
	NativeOAuthClientID   string // first-party refresh tokens (stored in oauth_client_id column)
	DefaultTenantID       string

	// Phase 5 — WebAuthn passkeys + MFA
	WebauthnRPID          string
	WebauthnRPDisplayName string
	WebauthnRPOrigins     []string // parsed from WEBAUTHN_RP_ORIGINS comma-separated
	WebauthnSessionTTL    time.Duration
	MFAEncryptionKey      string // AES-256 key material (min 32 bytes recommended)
	MFAPendingTicketTTL   time.Duration
	MFARecoveryCodeCount  int

	// Phase 4 — Federation (upstream OAuth/OIDC); credentials are stored per tenant in the database.
	// AdminAPIKey protects internal tenant/provider management routes (X-Admin-API-Key). Empty disables those routes.
	AdminAPIKey string
	// BootstrapAdminEmail and BootstrapAdminPassword create or promote the first platform admin when none exist (first startup only).
	BootstrapAdminEmail    string
	BootstrapAdminPassword string

	// OIDC Phase 2 (RS256, JWKS)
	OIDCRSAPrivateKeyPEM  string
	OIDCRSAPrivateKeyPath string
	OIDCKeyID             string
	OIDCAccessTTL         time.Duration
	OIDCIDTokenTTL        time.Duration
	OIDCAuthCodeTTL       time.Duration

	// Basic Auth configuration
	BasicAuthUsername string
	BasicAuthPassword string

	// HTTP Client configuration
	HTTPClientTimeout            time.Duration
	HTTPClientRetryCount         int
	HTTPClientRetryWaitMin       time.Duration
	HTTPClientRetryWaitMax       time.Duration
	HTTPClientDebug              bool
	HTTPClientTLSInsecureSkipTLS bool

	// Cloud storage configuration
	StorageProvider           string
	GCSBucket                 string
	GCSCredentialsJSONPath    string
	GCSPresignedURLDuration   time.Duration
	GCSPresignedURLExpiration time.Duration
	S3Bucket                  string
	S3Region                  string
	S3AccessKey               string
	S3SecretKey               string
	S3PresignedURLDuration    time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (optional; ignore errors)
	_ = godotenv.Load()

	appEnv := Environment(getEnv("APP_ENV", "development"))
	cfg := baseConfigFromEnv(appEnv)
	applyHTTPAndStorageEnv(cfg)
	return cfg, nil
}

func baseConfigFromEnv(appEnv Environment) *Config {
	return &Config{
		AppEnv:                  appEnv,
		AppName:                 getEnv("APP_NAME", ""),
		AppVersion:              getEnv("APP_VERSION", "1.0.0"),
		Timezone:                getEnv("TIMEZONE", "UTC"),
		AppHTTPServer:           getEnv("APP_HTTP_SERVER", ":3000"),
		AppRequestTimeout:       getEnvAsInt("APP_REQUEST_TIMEOUT", 30),
		AppBaseURL:              getEnv("APP_BASE_URL", ""),
		ServeEmbeddedFrontend:   getEnvAsBool("SERVE_EMBEDDED_FRONTEND", appEnv.IsProduction()),
		FrontendDistPath:        getEnv("FRONTEND_DIST_PATH", ""),
		OIDCLoginPageURL:        getEnv("OIDC_LOGIN_PAGE_URL", ""),
		DatabaseHost:            getEnv("POSTGRES_HOST", "localhost"),
		DatabasePort:            getEnv("POSTGRES_PORT", "5432"),
		DatabaseUsername:        getEnv("POSTGRES_USER", "postgres"),
		DatabasePassword:        getEnv("POSTGRES_PASSWORD", ""),
		DatabaseName:            getEnv("POSTGRES_DB", ""),
		DatabaseEnableDebug:     getEnvAsBool("DATABASE_DEBUG", false),
		DatabaseMaxOpenConns:    getEnvAsInt("DATABASE_MAX_OPEN_CONNS", 25),
		DatabaseMaxIdleConns:    getEnvAsInt("DATABASE_MAX_IDLE_CONNS", 5),
		DatabaseConnMaxLifetime: getEnvAsDuration("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
		DatabaseConnMaxIdleTime: getEnvAsDuration("DATABASE_CONN_MAX_IDLE_TIME", 1*time.Minute),
		DatabaseConnectTimeout:  getEnvAsDuration("DATABASE_CONNECT_TIMEOUT", 30*time.Second),
		DatabaseQueryTimeout:    getEnvAsDuration("DATABASE_QUERY_TIMEOUT", 30*time.Second),
		DatabaseHealthTimeout:   getEnvAsDuration("DATABASE_HEALTH_TIMEOUT", 5*time.Second),
		DatabaseRetryAttempts:   getEnvAsInt("DATABASE_RETRY_ATTEMPTS", 3),
		DatabaseRetryDelay:      getEnvAsDuration("DATABASE_RETRY_DELAY", 1*time.Second),
		DatabaseSSLMode:         getEnv("DATABASE_SSL_MODE", "disable"),
		DatabaseTimezone:        getEnv("DATABASE_TIMEZONE", "UTC"),
		CacheProvider:           getEnv("CACHE_PROVIDER", "redis"),
		RedisHost:               getEnv("REDIS_HOST", "localhost"),
		RedisPort:               getEnv("REDIS_PORT", "6379"),
		RedisPassword:           getEnv("REDIS_PASSWORD", ""),
		RedisDB:                 getEnvAsInt("REDIS_DB", 0),
		PoolSize:                getEnvAsInt("REDIS_POOL_SIZE", 10),
		DialTimeout:             getEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:             getEnvAsDuration("REDIS_READ_TIMEOUT", 5*time.Second),
		WriteTimeout:            getEnvAsDuration("REDIS_WRITE_TIMEOUT", 5*time.Second),
		PoolTimeout:             getEnvAsDuration("REDIS_POOL_TIMEOUT", 5*time.Second),
		MaxRetries:              getEnvAsInt("REDIS_MAX_RETRIES", 3),
		MinRetryBackoff:         getEnvAsDuration("REDIS_MIN_RETRY_BACKOFF", 1*time.Second),
		MaxRetryBackoff:         getEnvAsDuration("REDIS_MAX_RETRY_BACKOFF", 5*time.Second),
		LogLevel:                getEnv("LOG_LEVEL", "info"),
		EmailProvider:           getEnv("EMAIL_PROVIDER", "ses"),
		AWSSESRegion:            getEnv("AWS_SES_REGION", ""),
		AWSSESAccessKey:         getEnv("AWS_SES_ACCESS_KEY", ""),
		AWSSESSecretKey:         getEnv("AWS_SES_SECRET_KEY", ""),
		RateLimit:               getEnvAsInt("RATE_LIMIT", 20),
		RateLimitDuration:       getEnvAsDuration("RATE_LIMIT_DURATION", 1*time.Second),
		Environment:             getEnv("ENVIRONMENT", "development"),
		DefaultRateLimit:        getEnvAsInt("DEFAULT_RATE_LIMIT", 20),
		AuthRateLimit:           getEnvAsInt("AUTH_RATE_LIMIT", 3),
		PublicRateLimit:         getEnvAsInt("PUBLIC_RATE_LIMIT", 100),
		NewRelicAppName:         getEnv("NEWRELIC_APP_NAME", "github.com/gateforge-iam/gateforge-iam"),
		NewRelicLicense:         getEnv("NEWRELIC_LICENSE", ""),
		SentryDSN:               getEnv("SENTRY_DSN", ""),
		JWTSecret:               getEnv("JWT_SECRET", "dev-only-change-me-min-32-chars-long!!"),
		JWTAccessTTL:            getEnvAsDuration("JWT_ACCESS_TTL", 24*time.Hour),
		JWTRefreshTTL:           getEnvAsDuration("JWT_REFRESH_TTL", 168*time.Hour),
		SSOSessionTTL:           getEnvAsDuration("SSO_SESSION_TTL", 0),
		SSOSessionRememberTTL:   getEnvAsDuration("SSO_SESSION_REMEMBER_TTL", 720*time.Hour),
		NativeOAuthClientID:     getEnv("NATIVE_OAUTH_CLIENT_ID", "native"),
		DefaultTenantID:         getEnv("DEFAULT_TENANT_ID", "00000000-0000-0000-0000-000000000001"),
		WebauthnRPID:            getEnv("WEBAUTHN_RP_ID", "localhost"),
		WebauthnRPDisplayName:   getEnv("WEBAUTHN_RP_DISPLAY_NAME", "IAM"),
		WebauthnRPOrigins:       splitCommaTrim(getEnv("WEBAUTHN_RP_ORIGINS", "http://localhost:3000")),
		WebauthnSessionTTL:      getEnvAsDuration("WEBAUTHN_SESSION_TTL", 5*time.Minute),
		MFAEncryptionKey:        getEnv("MFA_ENCRYPTION_KEY", ""),
		MFAPendingTicketTTL:     getEnvAsDuration("MFA_PENDING_TICKET_TTL", 10*time.Minute),
		MFARecoveryCodeCount:    getEnvAsInt("MFA_RECOVERY_CODE_COUNT", 10),
		AdminAPIKey:             getEnv("ADMIN_API_KEY", ""),
		BootstrapAdminEmail:     getEnv("BOOTSTRAP_ADMIN_EMAIL", ""),
		BootstrapAdminPassword:  getEnv("BOOTSTRAP_ADMIN_PASSWORD", ""),
		OIDCRSAPrivateKeyPEM:    getEnv("OIDC_RSA_PRIVATE_KEY_PEM", ""),
		OIDCRSAPrivateKeyPath:   getEnv("OIDC_RSA_PRIVATE_KEY_FILE", ""),
		OIDCKeyID:               getEnv("OIDC_KEY_ID", "default-key-1"),
		OIDCAccessTTL:           getEnvAsDuration("OIDC_ACCESS_TTL", 0),
		OIDCIDTokenTTL:          getEnvAsDuration("OIDC_ID_TOKEN_TTL", 0),
		OIDCAuthCodeTTL:         getEnvAsDuration("OIDC_AUTH_CODE_TTL", 10*time.Minute),
		BasicAuthUsername:       getEnv("BASIC_AUTH_USER", ""),
		BasicAuthPassword:       getEnv("BASIC_AUTH_SECRET", ""),
	}
}

func applyHTTPAndStorageEnv(cfg *Config) {
	cfg.HTTPClientTimeout = getEnvAsDuration("HTTP_CLIENT_TIMEOUT", 30*time.Second)
	cfg.HTTPClientRetryCount = getEnvAsInt("HTTP_CLIENT_RETRY_COUNT", 2)
	cfg.HTTPClientRetryWaitMin = getEnvAsDuration("HTTP_CLIENT_RETRY_WAIT_MIN", 250*time.Millisecond)
	cfg.HTTPClientRetryWaitMax = getEnvAsDuration("HTTP_CLIENT_RETRY_WAIT_MAX", 2*time.Second)
	cfg.HTTPClientDebug = getEnvAsBool("HTTP_CLIENT_DEBUG", false)
	cfg.HTTPClientTLSInsecureSkipTLS = getEnvAsBool("HTTP_CLIENT_TLS_INSECURE_SKIP_TLS", false)
	cfg.StorageProvider = getEnv("STORAGE_PROVIDER", "gcs")
	cfg.GCSBucket = getEnv("GCS_BUCKET", "")
	cfg.GCSCredentialsJSONPath = getEnv("GCS_CREDENTIALS_JSON", "")
	cfg.GCSPresignedURLDuration = getEnvAsDuration("GCS_PRESIGNED_URL_DURATION", 1*time.Hour)
	cfg.GCSPresignedURLExpiration = getEnvAsDuration("GCS_PRESIGNED_URL_EXPIRATION", 1*time.Hour)
	cfg.S3Bucket = getEnv("S3_BUCKET", "")
	cfg.S3Region = getEnv("S3_REGION", "")
	cfg.S3AccessKey = getEnv("S3_ACCESS_KEY", "")
	cfg.S3SecretKey = getEnv("S3_SECRET_KEY", "")
	cfg.S3PresignedURLDuration = getEnvAsDuration("S3_PRESIGNED_URL_DURATION", 1*time.Hour)
}

func splitCommaTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getEnvAsBool gets an environment variable as boolean with a fallback value
func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}

// getEnvAsDuration gets an environment variable as duration with a fallback value
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=%v timezone=%v connect_timeout=%d",
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseUsername,
		c.DatabasePassword,
		c.DatabaseName,
		c.DatabaseSSLMode,
		c.DatabaseTimezone,
		int(c.DatabaseConnectTimeout.Seconds()),
	)
}

func (c *Config) IsDebugMode() bool {
	return c.DatabaseEnableDebug
}

// PopulateFromJSON reads a service account JSON and fills email and private key
func (c *Config) PopulateFromJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return c.PopulateFromJSONBytes(data)
}

// PopulateFromJSONBytes parses service account JSON bytes (same validation as PopulateFromJSON).
func (c *Config) PopulateFromJSONBytes(data []byte) error {
	var payload struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
		ProjectID   string `json:"project_id"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	return nil
}
