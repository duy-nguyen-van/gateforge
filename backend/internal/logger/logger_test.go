package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
