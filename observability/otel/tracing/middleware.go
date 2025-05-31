package tracingfrolfbot

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// InjectTraceContext propagates tracing info into a context.
func InjectTraceContext(ctx context.Context, tracer trace.Tracer) context.Context {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ctx
	}
	return trace.ContextWithSpan(ctx, span)
}
