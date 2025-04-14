package tracingfrolfbot

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// NoOpTracer implements the Tracer interface using OpenTelemetry's noop tracer
type NoOpTracer struct {
	tracer trace.Tracer
}

// NewNoOpTracer creates a new NoOpTracer
func NewNoOpTracer() *NoOpTracer {
	provider := noop.NewTracerProvider()
	return &NoOpTracer{
		tracer: provider.Tracer("noop"),
	}
}

// InitTracing does nothing in the no-op implementation
func (n *NoOpTracer) InitTracing(ctx context.Context, opts TracingOptions) (shutdown func(), err error) {
	return func() {}, nil
}

// TraceHandler returns the handler unchanged
func (n *NoOpTracer) TraceHandler(h message.HandlerFunc) message.HandlerFunc {
	return h
}

// InjectTraceContext does nothing
func (n *NoOpTracer) InjectTraceContext(ctx context.Context, msg *message.Message) {}

// SpanContextFromContext returns an empty span context
func (n *NoOpTracer) SpanContextFromContext(ctx context.Context) trace.SpanContext {
	return trace.SpanContextFromContext(ctx)
}

// StartSpan correctly returns a no-op span
func (n *NoOpTracer) StartSpan(ctx context.Context, name string, msg *message.Message) (context.Context, trace.Span) {
	return n.tracer.Start(ctx, name) // ✅ No longer panics
}
