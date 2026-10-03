package logger

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/pkg/correlationid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name        string
		level       string
		environment string
		wantLevel   zapcore.Level
	}{
		{"debug development", "debug", "development", zapcore.DebugLevel},
		{"info development", "info", "development", zapcore.InfoLevel},
		{"warn development", "warn", "development", zapcore.WarnLevel},
		{"error development", "error", "development", zapcore.ErrorLevel},
		{"unknown level defaults to info", "trace", "development", zapcore.InfoLevel},
		{"production json", "info", "production", zapcore.InfoLevel},
		{"production case insensitive", "warn", "PRODUCTION", zapcore.WarnLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.level, tt.environment)
			require.NotNil(t, Log)
			require.NotNil(t, Sugar)

			require.Equal(t, tt.wantLevel, Log.Level())
		})
	}
}

func TestInit_ProductionUsesJSONEncoder(t *testing.T) {
	Init("info", "production")
	require.NotNil(t, Log)
	// Smoke: logger should accept structured fields without panic.
	require.NotPanics(t, func() {
		Log.Info("production log test", zap.String("key", "value"))
	})
}

func TestInit_DevelopmentUsesConsoleEncoder(t *testing.T) {
	Init("debug", "development")
	require.NotNil(t, Log)
	require.NotPanics(t, func() {
		Sugar.Debugf("development log %s", "test")
	})
}

func TestFrom_AttachesCorrelationAndTrace(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	prev := Log
	Log = zap.New(core)
	t.Cleanup(func() { Log = prev })

	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	ctx := correlationid.NewContext(context.Background(), "corr-ops-1")
	ctx, span := tp.Tracer("logger-test").Start(ctx, "op")
	defer span.End()

	From(ctx).Info("hello")

	require.Equal(t, 1, logs.Len())
	fields := logs.All()[0].ContextMap()
	assert.Equal(t, "corr-ops-1", fields["correlation_id"])
	sc := span.SpanContext()
	assert.Equal(t, sc.TraceID().String(), fields["trace_id"])
	assert.Equal(t, sc.SpanID().String(), fields["span_id"])
}

func TestFrom_NilLoggerAndContext(t *testing.T) {
	prev := Log
	Log = nil
	t.Cleanup(func() { Log = prev })

	From(nil).Info("noop") //nolint:staticcheck // SA1012: nil context is the branch under test
	From(context.Background()).Info("noop")
}

func TestFrom_OmitsInvalidTrace(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	prev := Log
	Log = zap.New(core)
	t.Cleanup(func() { Log = prev })

	From(context.Background()).Info("plain")
	require.Equal(t, 1, logs.Len())
	fields := logs.All()[0].ContextMap()
	_, hasTrace := fields["trace_id"]
	_, hasCorr := fields["correlation_id"]
	assert.False(t, hasTrace)
	assert.False(t, hasCorr)
}
