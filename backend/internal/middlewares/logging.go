package middlewares

import (
	"bytes"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"
	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// RequestLogging provides structured logs to Sentry and zap.
// Access lines record method, URI, status, latency, and user_id — never bodies, JWT payloads, or auth headers.
func RequestLogging(cfg *config.Config) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogRequestID: true,
		LogUserAgent: true,
		LogMethod:    true,
		LogRemoteIP:  true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			userID := ""
			if u := c.Get(auth.EchoContextUserIDKey); u != nil {
				if s, ok := u.(string); ok {
					userID = s
				}
			}

			ctx := c.Request().Context()
			correlationID, _ := request.CorrelationIDFromContext(ctx)

			sentryLogger := sentry.NewLogger(ctx)
			stdLogger := log.New(sentryLogger, "", log.LstdFlags)
			stdLogger.Printf(
				"Request: %s %s (status=%d, latency=%v, request_id=%s, correlation_id=%s, user_id=%s, service=%s)",
				values.Method,
				values.URI,
				values.Status,
				values.Latency,
				values.RequestID,
				correlationID,
				userID,
				cfg.AppName,
			)

			logger.From(ctx).Info("Request: "+values.Method+" "+values.URI,
				zap.String("uri", values.URI),
				zap.String("method", values.Method),
				zap.Int("status", values.Status),
				zap.Duration("latency", values.Latency),
				zap.String("request_id", values.RequestID),
				zap.String("remote_ip", values.RemoteIP),
				zap.String("user_id", userID),
				zap.String("environment", cfg.AppEnv.String()),
				zap.String("service", cfg.AppName),
				zap.String("version", cfg.AppVersion),
				zap.Int64("timestamp", time.Now().UnixMilli()),
			)

			return nil
		},
	})
}

func LogBodyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if monitoring.SkipSensitiveRequestPath(c.Request().URL.Path) {
			return next(c)
		}
		ct := c.Request().Header.Get(echo.HeaderContentType)
		if strings.HasPrefix(ct, "multipart/form-data") {
			return next(c)
		}
		data, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		c.Request().Body = io.NopCloser(bytes.NewReader(data))
		c.Set("log_body", string(data))
		return next(c)
	}
}
