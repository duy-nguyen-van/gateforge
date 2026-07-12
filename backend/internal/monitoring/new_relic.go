package monitoring

import (
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"

	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitNewRelic(config config.Config) *newrelic.Application {
	appName := config.NewRelicAppName
	license := config.NewRelicLicense

	// Skip NewRelic initialization if license is not provided
	if license == "" {
		logger.Sugar.Warn("NewRelic license not provided, skipping NewRelic initialization")
		return nil
	}

	// Use default app name if not provided
	if appName == "" {
		appName = "github.com/gateforge-iam/gateforge-iam"
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(license),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		logger.Sugar.Warnf("Failed to initialize NewRelic: %v", err)
		return nil
	}

	// Add NewRelic core to zap logger
	if logger.Log != nil {
		logger.Log = logger.Log.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, NewNRCore(app))
		}))
		logger.Sugar = logger.Log.Sugar()
	}

	logger.Sugar.Info("NewRelic initialized successfully")

	return app
}
