package monitoring

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestSentryCore_Enabled(t *testing.T) {
	core := NewSentryCore(context.Background(), nil)
	requireEnabled(t, core, zapcore.ErrorLevel, true)
	requireEnabled(t, core, zapcore.InfoLevel, false)

	custom := NewSentryCore(context.Background(), []zapcore.Level{zapcore.InfoLevel})
	requireEnabled(t, custom, zapcore.InfoLevel, true)
	requireEnabled(t, custom, zapcore.ErrorLevel, false)
}

func TestSentryCore_WithCheckWriteSync(t *testing.T) {
	core := NewSentryCore(context.Background(), nil)
	withCore := core.With([]zap.Field{zap.String("k", "v")})
	require.Equal(t, core, withCore)

	entry := zapcore.Entry{
		Level:   zapcore.ErrorLevel,
		Time:    time.Now(),
		Message: "boom",
		Caller:  zapcore.NewEntryCaller(0, "file.go", 42, true),
	}
	checked := &zapcore.CheckedEntry{}
	result := core.Check(entry, checked)
	require.NotNil(t, result)

	skipped := core.Check(zapcore.Entry{Level: zapcore.InfoLevel}, &zapcore.CheckedEntry{})
	require.NotNil(t, skipped)

	err := core.Write(entry, []zap.Field{
		zap.String("detail", "value"),
		zap.Error(errors.New("write failed")),
	})
	require.NoError(t, err)

	err = core.Write(zapcore.Entry{Level: zapcore.FatalLevel, Message: "fatal"}, nil)
	require.NoError(t, err)

	require.NoError(t, core.Sync())
}

func TestSentryZapLevelAndFormatValue(t *testing.T) {
	require.Equal(t, sentry.LevelDebug, sentryZapLevel(zapcore.DebugLevel))
	require.Equal(t, sentry.LevelInfo, sentryZapLevel(zapcore.InfoLevel))
	require.Equal(t, sentry.LevelWarning, sentryZapLevel(zapcore.WarnLevel))
	require.Equal(t, sentry.LevelError, sentryZapLevel(zapcore.ErrorLevel))
	require.Equal(t, sentry.LevelFatal, sentryZapLevel(zapcore.FatalLevel))
	require.Equal(t, sentry.LevelInfo, sentryZapLevel(zapcore.Level(99)))
	require.NotEmpty(t, sentryZapFormatValue(errors.New("err")))
	require.Equal(t, "plain", sentryZapFormatValue("plain"))
	require.Equal(t, "42", sentryZapFormatValue(42))
}

func requireEnabled(t *testing.T, core *SentryCore, level zapcore.Level, want bool) {
	t.Helper()
	if core.Enabled(level) != want {
		t.Fatalf("Enabled(%v) = %v, want %v", level, !want, want)
	}
}
