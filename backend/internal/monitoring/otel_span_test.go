package monitoring

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestStartSpanAndEndSpan(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	t.Run("with attributes", func(t *testing.T) {
		ctx, span := StartSpan(context.Background(), "test-tracer", "op", attribute.String("k", "v"))
		require.NotNil(t, span)
		assert.NotNil(t, ctx)
		EndSpan(span, nil)
	})

	t.Run("records error", func(t *testing.T) {
		_, span := StartSpan(context.Background(), "test-tracer", "fail")
		EndSpan(span, errors.New("boom"))
	})

	t.Run("nil span is no-op", func(t *testing.T) {
		EndSpan(nil, errors.New("ignored"))
	})
}

func TestFinish_RecordsErrorAndLogs(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	ctx, span := StartSpan(context.Background(), "test-tracer", "OIDCService.Token")
	Finish(ctx, span, "OIDCService.Token", errors.New("invalid_grant"))

	spans := recorder.Ended()
	require.Len(t, spans, 1)
	assert.Equal(t, codes.Error, spans[0].Status().Code)
}

func TestObserve_MarksSpanOnError(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	err := ObserveErr(context.Background(), "test-tracer", "op.fail", nil, func(ctx context.Context) error {
		return errors.New("boom")
	})
	require.EqualError(t, err, "boom")

	spans := recorder.Ended()
	require.Len(t, spans, 1)
	assert.Equal(t, "op.fail", spans[0].Name())
	assert.Equal(t, codes.Error, spans[0].Status().Code)
}

func TestObserve_SuccessReturnsValue(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	out, err := Observe(context.Background(), "test-tracer", "op.ok",
		[]attribute.KeyValue{attribute.String("k", "v")},
		func(ctx context.Context) (int, error) {
			return 7, nil
		},
	)
	require.NoError(t, err)
	assert.Equal(t, 7, out)

	spans := recorder.Ended()
	require.Len(t, spans, 1)
	assert.Equal(t, "op.ok", spans[0].Name())
	assert.Equal(t, codes.Unset, spans[0].Status().Code)
}

func TestObserve2And3(t *testing.T) {
	a, b, err := Observe2(context.Background(), "test-tracer", "op.two", nil,
		func(ctx context.Context) (string, int, error) {
			return "ok", 2, nil
		})
	require.NoError(t, err)
	assert.Equal(t, "ok", a)
	assert.Equal(t, 2, b)

	x, y, z, err := Observe3(context.Background(), "test-tracer", "op.three", nil,
		func(ctx context.Context) (int, int, int, error) {
			return 1, 2, 3, nil
		})
	require.NoError(t, err)
	assert.Equal(t, 1, x)
	assert.Equal(t, 2, y)
	assert.Equal(t, 3, z)
}
