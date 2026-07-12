package monitoring

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SentryCore is a zap core that sends logs to Sentry
type SentryCore struct {
	ctx    context.Context
	levels []zapcore.Level
}

// NewSentryCore creates a new zap core for Sentry
func NewSentryCore(ctx context.Context, levels []zapcore.Level) *SentryCore {
	if levels == nil {
		levels = []zapcore.Level{
			zapcore.ErrorLevel,
			zapcore.FatalLevel,
			zapcore.PanicLevel,
		}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return &SentryCore{
		ctx:    ctx,
		levels: levels,
	}
}

// Enabled returns whether the core should process this entry
func (c *SentryCore) Enabled(level zapcore.Level) bool {
	for _, l := range c.levels {
		if l == level {
			return true
		}
	}
	return false
}

// With adds structured context to the core
func (c *SentryCore) With(fields []zap.Field) zapcore.Core {
	return c
}

// Check determines whether the supplied entry should be logged
func (c *SentryCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

// Write sends the log entry to Sentry
func (c *SentryCore) Write(entry zapcore.Entry, fields []zap.Field) error {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentryZapLevel(entry.Level))

		for _, field := range fields {
			if field.Key == "error" {
				if err, ok := field.Interface.(error); ok {
					scope.SetExtra("error_details", err.Error())
					hub.CaptureException(err)
					continue
				}
			}
			scope.SetExtra(field.Key, field.Interface)
		}

		scope.SetTag("logger", "zap")
		scope.SetTag("log_level", entry.Level.String())
		scope.SetExtra("timestamp", entry.Time.Format(time.RFC3339))

		if entry.Caller.Defined {
			scope.SetExtra("caller_file", entry.Caller.File)
			scope.SetExtra("caller_line", entry.Caller.Line)
			scope.SetExtra("caller_function", entry.Caller.Function)
		}
	})

	msg := entry.Message
	for _, field := range fields {
		if field.Key != "error" {
			msg = msg + " " + field.Key + "=" + sentryZapFormatValue(field.Interface)
		}
	}

	sentryLogger := sentry.NewLogger(c.ctx)
	stdLogger := log.New(sentryLogger, "", log.LstdFlags)
	stdLogger.Println(msg)

	return nil
}

// Sync flushes any buffered logs
func (c *SentryCore) Sync() error {
	return nil
}

func sentryZapLevel(level zapcore.Level) sentry.Level {
	switch level {
	case zapcore.DebugLevel:
		return sentry.LevelDebug
	case zapcore.InfoLevel:
		return sentry.LevelInfo
	case zapcore.WarnLevel:
		return sentry.LevelWarning
	case zapcore.ErrorLevel:
		return sentry.LevelError
	case zapcore.FatalLevel, zapcore.PanicLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelInfo
	}
}

func sentryZapFormatValue(v interface{}) string {
	switch val := v.(type) {
	case error:
		return val.Error()
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}
