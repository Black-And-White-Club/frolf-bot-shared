package tracingfrolfbot

import (
	"context"
	"strconv"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Tracer defines the interface for tracing.
type Tracer interface {
	TraceHandler(h message.HandlerFunc) message.HandlerFunc
	InjectTraceContext(ctx context.Context, msg *message.Message)
	SpanContextFromContext(ctx context.Context) trace.SpanContext
	StartSpan(ctx context.Context, name string, msg *message.Message) (context.Context, trace.Span)
}

// TracingOptions defines configuration for OpenTelemetry tracing.
type TracingOptions struct {
	ServiceName    string
	TempoEndpoint  string
	Insecure       bool // Whether to use an insecure connection to Tempo (for testing ONLY!)
	ServiceVersion string
	SampleRate     float64
	Environment    string
}

// TempoTracer is the concrete implementation of the Tracer interface.
type TempoTracer struct {
	tracer trace.Tracer // Store the tracer for use in middleware
}

func TraceHandler(tracer trace.Tracer) message.HandlerMiddleware {
	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			ctx := otel.GetTextMapPropagator().Extract(msg.Context(), messageMetadataCarrier{msg})
			spanName := msg.Metadata.Get("topic")
			if spanName == "" {
				spanName = "watermill.handler"
			}

			ctx, span := tracer.Start(ctx, spanName,
				trace.WithAttributes(
					attribute.String("watermill.message_id", msg.UUID),
					attribute.String("watermill.correlation_id", msg.Metadata.Get("correlation_id")),
					attribute.String("watermill.topic", msg.Metadata.Get("topic")),
					attribute.String("watermill.middleware", "watermill.handler"),
					attribute.Bool("watermill.dead_letter", msg.Metadata.Get("dead_letter") == "true"),
					attribute.Int64("watermill.retry_count", parseInt64(msg.Metadata.Get("jetstream.delivery_attempt"))),
				),
			)
			defer span.End()

			msg.SetContext(ctx)
			msgs, err := h(msg)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			return msgs, err
		}
	}
}

// parseInt64 is a helper function to safely parse int64 values from strings.
func parseInt64(value string) int64 {
	if value == "" {
		return 0
	}
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// StartSpan creates a new span with the given name and includes common attributes.
func (t *TempoTracer) StartSpan(ctx context.Context, name string, msg *message.Message) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name,
		trace.WithAttributes(attribute.String("correlation_id", msg.Metadata.Get(middleware.CorrelationIDMetadataKey))),
	)
}

// InjectTraceContext injects the trace context into the message metadata before publishing.
func (t *TempoTracer) InjectTraceContext(ctx context.Context, msg *message.Message) {
	otel.GetTextMapPropagator().Inject(ctx, messageMetadataCarrier{msg})
}

// SpanContextFromContext extracts the span context from the provided context.
func (t *TempoTracer) SpanContextFromContext(ctx context.Context) trace.SpanContext {
	return trace.SpanFromContext(ctx).SpanContext()
}

// messageMetadataCarrier adapts Watermill's message.Metadata to OpenTelemetry's TextMapCarrier.
type messageMetadataCarrier struct {
	msg *message.Message
}

func (c messageMetadataCarrier) Get(key string) string {
	return c.msg.Metadata.Get(key)
}

func (c messageMetadataCarrier) Set(key, value string) {
	c.msg.Metadata.Set(key, value)
}

func (c messageMetadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c.msg.Metadata))
	for k := range c.msg.Metadata {
		keys = append(keys, k)
	}
	return keys
}
