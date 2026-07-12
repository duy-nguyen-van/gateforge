package monitoring

import (
	"testing"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNRCore_NilAppDisabled(t *testing.T) {
	core := NewNRCore(nil)
	require.False(t, core.Enabled(zapcore.InfoLevel))
	require.NoError(t, core.Write(zapcore.Entry{Level: zapcore.InfoLevel, Message: "skip"}, nil))
	require.NoError(t, core.Sync())
}

func TestNRCore_WithCheckWriteSync(t *testing.T) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("gateforge-iam-test"),
		newrelic.ConfigLicense("0000000000000000000000000000000000000000"),
		newrelic.ConfigEnabled(false),
	)
	require.NoError(t, err)

	core := NewNRCore(app)
	require.True(t, core.Enabled(zapcore.InfoLevel))
	require.Equal(t, core, core.With([]zap.Field{zap.String("k", "v")}))

	entry := zapcore.Entry{
		Level:   zapcore.WarnLevel,
		Time:    time.Now(),
		Message: "warn message",
	}
	checked := core.Check(entry, &zapcore.CheckedEntry{})
	require.NotNil(t, checked)

	require.NoError(t, core.Write(entry, []zap.Field{
		zap.Duration("latency", 2*time.Millisecond),
		zap.String("key", "value"),
	}))

	for _, level := range []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.ErrorLevel,
		zapcore.FatalLevel,
		zapcore.PanicLevel,
	} {
		require.NoError(t, core.Write(zapcore.Entry{Level: level, Message: "msg", Time: time.Now()}, nil))
	}

	require.NoError(t, core.Sync())
}

func TestNRCore_DisabledLevelSkipsCheck(t *testing.T) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("gateforge-iam-test"),
		newrelic.ConfigLicense("0000000000000000000000000000000000000000"),
		newrelic.ConfigEnabled(false),
	)
	require.NoError(t, err)

	core := NewNRCore(app)
	entry := zapcore.Entry{Level: zapcore.Level(99), Message: "skip"}
	checked := core.Check(entry, &zapcore.CheckedEntry{})
	require.NotNil(t, checked)
	require.False(t, core.Enabled(zapcore.Level(99)))
}
