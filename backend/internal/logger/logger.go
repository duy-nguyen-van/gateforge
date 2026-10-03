package logger

import (
	"context"
	"os"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/pkg/correlationid"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger
var Sugar *zap.SugaredLogger

// Init initializes the logger with the specified level and environment
func Init(level string, environment string) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Normalize environment to lowercase for comparison
	env := strings.ToLower(environment)
	isProduction := env == "production"

	// Configure encoder based on environment
	var encoderConfig zapcore.EncoderConfig
	if isProduction {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create encoder
	var encoder zapcore.Encoder
	if isProduction {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// Create logger with caller and stack trace
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	Sugar = Log.Sugar()
}

// From returns a logger with correlation_id, trace_id, and span_id from ctx.
// The empty-key context field lets otelzap attach the active trace to OTLP log records.
func From(ctx context.Context) *zap.Logger {
	l := Log
	if l == nil {
		l = zap.NewNop()
	}
	if ctx == nil {
		return l
	}

	fields := make([]zap.Field, 0, 4)
	if cid, ok := correlationid.FromContext(ctx); ok && cid != "" {
		fields = append(fields, zap.String("correlation_id", cid))
	}
	sc := trace.SpanFromContext(ctx).SpanContext()
	if sc.IsValid() {
		fields = append(fields, zap.String("trace_id", sc.TraceID().String()))
		fields = append(fields, zap.String("span_id", sc.SpanID().String()))
	}
	fields = append(fields, zap.Any("", ctx))
	return l.With(fields...)
}
