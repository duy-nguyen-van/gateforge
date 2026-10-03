package monitoring

import (
	"context"

	"github.com/gateforge-iam/gateforge-iam/internal/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// StartSpan starts a child span linked to any span in ctx.
// When OpenTelemetry is disabled, a no-op span is returned.
func StartSpan(ctx context.Context, tracerName, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, spanName)
	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}
	return ctx, span
}

// EndSpan ends a span and records err when non-nil.
func EndSpan(span trace.Span, err error) {
	if span == nil {
		return
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}

// Finish ends a span and logs the operation outcome via logger.From.
func Finish(ctx context.Context, span trace.Span, spanName string, err error) {
	EndSpan(span, err)
	logObserve(ctx, spanName, err)
}

// Observe starts a child span, runs fn, and ends the span with the returned error.
func Observe[T any](ctx context.Context, tracerName, spanName string, attrs []attribute.KeyValue, fn func(context.Context) (T, error)) (T, error) {
	ctx, span := StartSpan(ctx, tracerName, spanName, attrs...)
	out, err := fn(ctx)
	Finish(ctx, span, spanName, err)
	return out, err
}

// Observe2 is Observe for functions that return two values plus an error.
func Observe2[A, B any](ctx context.Context, tracerName, spanName string, attrs []attribute.KeyValue, fn func(context.Context) (A, B, error)) (A, B, error) {
	type pair struct {
		a A
		b B
	}
	out, err := Observe(ctx, tracerName, spanName, attrs, func(ctx context.Context) (pair, error) {
		a, b, err := fn(ctx)
		return pair{a: a, b: b}, err
	})
	return out.a, out.b, err
}

// Observe3 is Observe for functions that return three values plus an error.
func Observe3[A, B, C any](ctx context.Context, tracerName, spanName string, attrs []attribute.KeyValue, fn func(context.Context) (A, B, C, error)) (A, B, C, error) {
	type triple struct {
		a A
		b B
		c C
	}
	out, err := Observe(ctx, tracerName, spanName, attrs, func(ctx context.Context) (triple, error) {
		a, b, c, err := fn(ctx)
		return triple{a: a, b: b, c: c}, err
	})
	return out.a, out.b, out.c, err
}

// ObserveErr is Observe for functions that return only an error.
func ObserveErr(ctx context.Context, tracerName, spanName string, attrs []attribute.KeyValue, fn func(context.Context) error) error {
	_, err := Observe(ctx, tracerName, spanName, attrs, func(ctx context.Context) (struct{}, error) {
		return struct{}{}, fn(ctx)
	})
	return err
}

func logObserve(ctx context.Context, spanName string, err error) {
	log := logger.From(ctx)
	if err != nil {
		log.Error(spanName+" failed", zap.Error(err))
		return
	}
	log.Info(spanName)
}
