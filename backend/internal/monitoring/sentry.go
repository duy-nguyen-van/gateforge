package monitoring

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitSentry(cfg config.Config) {
	if strings.TrimSpace(cfg.SentryDSN) == "" {
		return
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.AppEnv.String(),
		Release:          cfg.AppName + "@" + cfg.AppVersion,
		Debug:            cfg.AppEnv == config.EnvironmentDevelopment,
		AttachStacktrace: true,
		EnableTracing:    true,
		EnableLogs:       true,
		TracesSampleRate: 1.0,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if event.Request != nil {
				if event.Request.Headers != nil {
					delete(event.Request.Headers, "Authorization")
					delete(event.Request.Headers, "Cookie")
				}
			}
			return event
		},
	})

	if err != nil {
		log.Fatalf("sentry.Init failed: %v", err)
	}

	if logger.Log != nil {
		logger.Log = logger.Log.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, NewSentryCore(context.Background(), nil))
		}))
		logger.Sugar = logger.Log.Sugar()
	}
}

// Flush buffered events before shutdown
func FlushSentry() {
	sentry.Flush(2 * time.Second)
}

// GetSentryHub returns a Sentry hub from the context if available, otherwise falls back to the current hub.
// This makes error reporting more robust by ensuring we always try to report errors when possible.
func GetSentryHub(ctx context.Context) *sentry.Hub {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	return hub
}

// SentryWriter implements io.Writer to write logs to Sentry
type SentryWriter struct {
	ctx context.Context
}

// NewSentryWriter creates a new writer that writes to Sentry
func NewSentryWriter(ctx context.Context) *SentryWriter {
	if ctx == nil {
		ctx = context.Background()
	}
	return &SentryWriter{ctx: ctx}
}

// Write implements io.Writer
func (w *SentryWriter) Write(p []byte) (n int, err error) {
	sentryLogger := sentry.NewLogger(w.ctx)
	stdLog := log.New(sentryLogger, "", log.LstdFlags)
	stdLog.Print(string(p))
	return len(p), nil
}
