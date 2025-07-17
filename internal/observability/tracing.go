// Package observability provides a wrapper around the OpenTelemetry tracing
package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// StartTracer is a wrapper around oteltrace.Tracer.Start that uses the global tracer.
func StartTracer(ctx context.Context, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	tracer := GetTracer()

	if tracer == nil {
		// return the global tracer
		return otel.GetTracerProvider().Tracer("").Start(ctx, spanName, opts...) //nolint:spancheck // spanName != span
	}

	return tracer.Start(ctx, spanName, opts...) //nolint:spancheck // spanName != span
}
