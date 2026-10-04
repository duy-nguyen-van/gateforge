# OpenTelemetry

GateForge IAM exports **traces**, **metrics**, and **logs** over OTLP using the same pattern as production Go services: optional exporters, Zap on stdout with trace correlation, business spans on exported service methods, and auto-spans on I/O.

Instrumentation is off when `OTEL_EXPORTER_OTLP_ENDPOINT` is empty (same as an empty Sentry DSN). Sentry stays the error product; do not add Sentry or New Relic child spans as a second trace model.

## Architecture

```text
HTTP (otelecho) → Observe(service method) → GORM / Redis / otelhttp child spans
                     ↓
            logger.From(ctx)  →  correlation_id + trace_id + span_id
                     ↓
         stdout (all levels)  +  OTLP error+  +  Sentry error+
```

Init order in `cmd/server/main.go`:

1. `logger.Init`
2. `InitNewRelic` / `InitSentry`
3. `InitOpenTelemetry` (no-op if endpoint empty)
4. `AttachOTelZapLogger` (error+)
5. Sentry zap core (error+)
6. `defer` OTel `Shutdown` + `FlushSentry`

Middleware order in `cmd/server/routes/router.go`: **otelecho (if enabled) → Sentry →** RequestID / RequestContext / recovery / CSRF / RequestLogging.

## Local collector

```bash
make otel-up
# In cmd/server/.env:
#   OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
#   OTEL_EXPORTER_OTLP_PROTOCOL=grpc
#   OTEL_EXPORTER_OTLP_INSECURE=true
make -C backend up
```

Jaeger UI: http://localhost:16686. Collector stdout for metrics and OTLP logs: `docker compose -f backend/docker-compose.yml logs -f otel-collector`.

```bash
make otel-down
```

## Configuration

| Env | Default | Notes |
|-----|---------|-------|
| `OTEL_SERVICE_NAME` | `APP_NAME` / `gateforge-iam` | Span resource service name |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | empty (off) | e.g. `localhost:4317` |
| `OTEL_EXPORTER_OTLP_PROTOCOL` | `grpc` | `grpc` or `http/protobuf` |
| `OTEL_EXPORTER_OTLP_INSECURE` | `false` | `true` for local collector |
| `OTEL_TRACES_ENABLED` | `true` | Ignored if endpoint empty |
| `OTEL_METRICS_ENABLED` | `true` | Ignored if endpoint empty |
| `OTEL_LOGS_ENABLED` | `true` | Zap error+ via otelzap |

## What is instrumented

| Surface | How |
|---------|-----|
| HTTP server | Official Echo v4 `otelecho.Middleware` when traces are on |
| IAM services | `monitoring.Observe` / `Observe2` / `Observe3` / `ObserveErr`; OIDC uses `StartSpan`/`Finish` because OAuth errors are not `error` |
| Postgres | GORM OpenTelemetry tracing plugin |
| Redis | `redisotel.InstrumentTracing` / `InstrumentMetrics` |
| Federation IdP HTTP | `otelhttp` on the `oauth2`/`oidc` HTTP client |
| Zap | `logger.From(ctx)` on request paths; OTLP + Sentry tee error+ only |

Span attributes are IDs only (`client_id`, `tenant_id`, `user_id`, `provider`). Never email, password, JWT, tokens, or request bodies.

Do **not** wrap `AuditService.Record` — it is called from every other service and would double Info logs.

## Counters

Meter name `gateforge-iam`. Labels stay low-cardinality: no email, token, IP, or user id.

| Counter | Labels |
|---------|--------|
| `oidc_token_total` | `grant`, `result` |
| `auth_login_total` | `result` |
| `auth_lockout_total` | none |
| `rate_limit_denied_total` | `route_group` (`auth`, `public`, `default`) |
| `audit_dropped_total` | none |
| `retention_rows_deleted` | `table` |

Auth and OIDC audit events are queued (default 1024). A full queue increments `audit_dropped_total` and does not fail `/token`. Admin mutations insert synchronously and fail the request if the insert fails.

## Logs and redaction

Access logs are method/URI/status/latency/`correlation_id`/`user_id`. `LogBodyMiddleware` skips auth/OIDC/MFA/WebAuthn paths so passwords and tokens are not captured for Sentry extras. Sentry scopes use `RedactedHeaders` (`Authorization`, `Cookie`, CSRF, admin API key).

Use `logger.From(ctx)` on request paths so stdout lines join Jaeger via `trace_id`. `Observe` already logs `"Service.Method"` / `"… failed"` — do not also `logger.Log.Info` then return the same error.

## Related files

| Path | Role |
|------|------|
| `internal/monitoring/otel.go` | OTLP exporters, providers, shutdown |
| `internal/monitoring/otel_span.go` | `Observe` / `StartSpan` / `Finish` |
| `internal/monitoring/otel_zap.go` | Zap → OTLP error+ bridge |
| `internal/logger/logger.go` | `From(ctx)` |
| `internal/monitoring/sentry_headers.go` | Header/body redaction |
| `deploy/otel/otel-collector-config.yaml` | Local collector pipelines |
| `.cursor/rules/observability.mdc` | Agent rules |
