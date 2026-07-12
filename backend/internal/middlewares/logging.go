package middlewares

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type requestLogContext struct {
	userID         string
	userLoginInfo  string
	organizationID string
	body           any
	queryJSON      string
	pathParamsJSON string
	headersJSON    string
	correlationID  string
	languageCode   string
}

// RequestLogging provides structured logs to Sentry and zap
// and mirrors the rich request logging previously configured in the router.
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
			ctx := buildRequestLogContext(c)
			logRequestToSentry(c, values, cfg, ctx)
			logRequestToZap(c, values, cfg, ctx)
			return nil
		},
	})
}

func buildRequestLogContext(c echo.Context) requestLogContext {
	userLoginInfo := ""
	if len(c.Request().Header["Authorization"]) > 0 {
		tmp := strings.Split(c.Request().Header["Authorization"][0], ".")
		if len(tmp) == 3 {
			sDesc, _ := base64.RawStdEncoding.DecodeString(tmp[1])
			userLoginInfo = string(sDesc)
		}
	}

	userID := ""
	if u := c.Get(auth.EchoContextUserIDKey); u != nil {
		if s, ok := u.(string); ok {
			userID = s
		}
	}

	organizationID := ""
	if org := c.Get("organization_id"); org != nil {
		organizationID = org.(string)
	}

	body := c.Get("log_body")
	if body == nil {
		body = ""
	}

	query := echo.Map{}
	_ = (&echo.DefaultBinder{}).BindQueryParams(c, &query)
	jsonQueryStr, _ := json.Marshal(query)

	param := echo.Map{}
	_ = (&echo.DefaultBinder{}).BindPathParams(c, &param)
	jsonParamStr, _ := json.Marshal(param)

	headers := make(map[string]string)
	for k, v := range c.Request().Header {
		if strings.EqualFold(k, "Authorization") ||
			strings.EqualFold(k, "Cookie") ||
			strings.EqualFold(k, "Set-Cookie") {
			continue
		}
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	jsonHeadersStr, _ := json.Marshal(headers)

	reqCtx := c.Request().Context()
	correlationID, _ := request.CorrelationIDFromContext(reqCtx)
	languageCode, _ := request.LanguageCodeFromContext(reqCtx)

	return requestLogContext{
		userID:         userID,
		userLoginInfo:  userLoginInfo,
		organizationID: organizationID,
		body:           body,
		queryJSON:      string(jsonQueryStr),
		pathParamsJSON: string(jsonParamStr),
		headersJSON:    string(jsonHeadersStr),
		correlationID:  correlationID,
		languageCode:   languageCode,
	}
}

func logRequestToSentry(c echo.Context, values middleware.RequestLoggerValues, cfg *config.Config, ctx requestLogContext) {
	sentryLogger := sentry.NewLogger(c.Request().Context())
	stdLogger := log.New(sentryLogger, "", log.LstdFlags)
	logMsg := fmt.Sprintf(
		"Request: %s %s (status=%d, latency=%v, request_id=%s, correlation_id=%s, language=%s, remote_ip=%s, user_agent=%s, user_id=%s, org_id=%s, payload=%s, query=%s, path_params=%s, headers=%s, environment=%s, service=%s, version=%s, timestamp=%d, hostname=%s, protocol=%s)",
		values.Method,
		values.URI,
		values.Status,
		values.Latency,
		values.RequestID,
		ctx.correlationID,
		ctx.languageCode,
		values.RemoteIP,
		values.UserAgent,
		ctx.userID,
		ctx.organizationID,
		ctx.body,
		ctx.queryJSON,
		ctx.pathParamsJSON,
		ctx.headersJSON,
		cfg.AppEnv.String(),
		cfg.AppName,
		cfg.AppVersion,
		time.Now().UnixMilli(),
		c.Request().Host,
		c.Request().Proto,
	)
	stdLogger.Println(logMsg)
}

func logRequestToZap(c echo.Context, values middleware.RequestLoggerValues, cfg *config.Config, ctx requestLogContext) {
	logger.Log.Info("Request: "+values.Method+" "+values.URI,
		zap.String("uri", values.URI),
		zap.String("method", values.Method),
		zap.Int("status", values.Status),
		zap.Duration("latency", values.Latency),
		zap.String("request_id", values.RequestID),
		zap.String("remote_ip", values.RemoteIP),
		zap.String("correlation_id", ctx.correlationID),
		zap.String("language_code", ctx.languageCode),
		zap.String("user_agent", values.UserAgent),
		zap.String("user_id", ctx.userID),
		zap.String("user_login", ctx.userLoginInfo),
		zap.String("org_id", ctx.organizationID),
		zap.String("payload", fmt.Sprintf("%v", ctx.body)),
		zap.String("query", ctx.queryJSON),
		zap.String("path_params", ctx.pathParamsJSON),
		zap.String("headers", ctx.headersJSON),
		zap.String("environment", cfg.AppEnv.String()),
		zap.String("service", cfg.AppName),
		zap.String("version", cfg.AppVersion),
		zap.Int64("timestamp", time.Now().UnixMilli()),
		zap.String("hostname", c.Request().Host),
		zap.String("protocol", c.Request().Proto),
	)
}

func LogBodyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		data, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		c.Request().Body = io.NopCloser(bytes.NewReader(data))
		c.Set("log_body", string(data))
		return next(c)
	}
}
