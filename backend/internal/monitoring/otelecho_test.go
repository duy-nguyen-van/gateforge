package monitoring

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestOtelEcho_SetsHTTPResponseStatusCode(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		_ = tp.Shutdown(t.Context())
		otel.SetTracerProvider(prev)
	})

	e := echo.New()
	e.Use(otelecho.Middleware("gateforge-iam"))
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	spans := recorder.Ended()
	require.NotEmpty(t, spans)
	span := spans[len(spans)-1]

	var status int64
	found := false
	for _, attr := range span.Attributes() {
		switch string(attr.Key) {
		case "http.response.status_code", "http.status_code":
			status = attr.Value.AsInt64()
			found = true
		}
	}
	require.True(t, found, "expected http response status attribute on server span")
	assert.Equal(t, int64(http.StatusOK), status)
}
