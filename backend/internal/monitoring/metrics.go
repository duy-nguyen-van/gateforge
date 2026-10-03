package monitoring

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const meterName = "gateforge-iam"

var (
	metricsOnce          sync.Once
	oidcTokenTotal       metric.Int64Counter
	authLockoutTotal     metric.Int64Counter
	rateLimitDeniedTotal metric.Int64Counter
	auditDroppedTotal    metric.Int64Counter
	retentionRowsDeleted metric.Int64Counter
	authLoginTotal       metric.Int64Counter
)

func instruments() {
	metricsOnce.Do(func() {
		m := otel.Meter(meterName)
		oidcTokenTotal, _ = m.Int64Counter("oidc_token_total")
		authLockoutTotal, _ = m.Int64Counter("auth_lockout_total")
		rateLimitDeniedTotal, _ = m.Int64Counter("rate_limit_denied_total")
		auditDroppedTotal, _ = m.Int64Counter("audit_dropped_total")
		retentionRowsDeleted, _ = m.Int64Counter("retention_rows_deleted")
		authLoginTotal, _ = m.Int64Counter("auth_login_total")
	})
}

func add(ctx context.Context, c metric.Int64Counter, n int64, attrs ...attribute.KeyValue) {
	if c == nil {
		return
	}
	c.Add(ctx, n, metric.WithAttributes(attrs...))
}

// AddOIDCToken records a token endpoint outcome. grant and result are low-cardinality.
func AddOIDCToken(ctx context.Context, grant, result string) {
	instruments()
	add(ctx, oidcTokenTotal, 1, attribute.String("grant", grant), attribute.String("result", result))
}

// AddAuthLogin records a password login outcome.
func AddAuthLogin(ctx context.Context, result string) {
	instruments()
	add(ctx, authLoginTotal, 1, attribute.String("result", result))
}

// AddAuthLockout records an account lockout.
func AddAuthLockout(ctx context.Context) {
	instruments()
	add(ctx, authLockoutTotal, 1)
}

// AddRateLimitDenied records a denied request. routeGroup is auth, public, or default.
func AddRateLimitDenied(ctx context.Context, routeGroup string) {
	instruments()
	add(ctx, rateLimitDeniedTotal, 1, attribute.String("route_group", routeGroup))
}

// AddAuditDropped records an audit event that could not be queued.
func AddAuditDropped(ctx context.Context) {
	instruments()
	add(ctx, auditDroppedTotal, 1)
}

// AddRetentionRows records rows removed by the retention job.
func AddRetentionRows(ctx context.Context, table string, n int64) {
	if n <= 0 {
		return
	}
	instruments()
	add(ctx, retentionRowsDeleted, n, attribute.String("table", table))
}
