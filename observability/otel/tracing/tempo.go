package tracingfrolfbot

import (
	"context"
	"strconv"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// OTEL Semantic Convention attribute keys for messaging (v1.26.0+)
// https://opentelemetry.io/docs/specs/semconv/messaging/messaging-spans/
const (
	// MessagingSystemJetStream identifies NATS JetStream as the messaging system.
	MessagingSystemJetStream = "nats_jetstream"

	// Custom attributes for Watermill/Grafana queryability
	AttrCorrelationID = "messaging.correlation_id"
	AttrDeadLetter    = "messaging.dead_letter"
	AttrRetryCount    = "messaging.delivery_attempt"
	AttrHandlerName   = "messaging.handler.name"
)

// MessagingOperationType values per OTEL semconv
var (
	MessagingOperationProcess = semconv.MessagingOperationTypeKey.String("process")
)

// Tracer defines the interface for tracing.
type Tracer interface {
	TraceHandler(h message.HandlerFunc) message.HandlerFunc
	InjectTraceContext(ctx context.Context, msg *message.Message)
	SpanContextFromContext(ctx context.Context) trace.SpanContext
	StartSpan(ctx context.Context, name string, msg *message.Message) (context.Context, trace.Span)
	StartPublishSpan(ctx context.Context, topic string, msg *message.Message) (context.Context, trace.Span)
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

// NewTempoTracer creates a new TempoTracer with the given tracer.
func NewTempoTracer(tracer trace.Tracer) *TempoTracer {
	return &TempoTracer{tracer: tracer}
}

func TraceHandler(tracer trace.Tracer) message.HandlerMiddleware {
	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			ctx := otel.GetTextMapPropagator().Extract(msg.Context(), messageMetadataCarrier{msg})

			topic := msg.Metadata.Get("topic")
			handlerName := msg.Metadata.Get("handler_name")

			// Span name follows OTEL convention: "<destination> <operation>"
			spanName := topic + " process"
			if topic == "" {
				spanName = "message process"
			}

			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					// OTEL Semantic Conventions for Messaging
					semconv.MessagingSystemKey.String(MessagingSystemJetStream),
					MessagingOperationProcess,
					semconv.MessagingDestinationName(topic),
					semconv.MessagingMessageID(msg.UUID),
					semconv.MessagingMessageBodySize(len(msg.Payload)),

					// Custom attributes for queryability
					attribute.String(AttrCorrelationID, msg.Metadata.Get("correlation_id")),
					attribute.String(AttrHandlerName, handlerName),
					attribute.Bool(AttrDeadLetter, msg.Metadata.Get("dead_letter") == "true"),
					attribute.Int64(AttrRetryCount, parseInt64(msg.Metadata.Get("jetstream.delivery_attempt"))),
				),
			)
			defer span.End()

			msg.SetContext(ctx)
			msgs, err := h(msg)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
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
	topic := msg.Metadata.Get("topic")
	return t.tracer.Start(ctx, name,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			semconv.MessagingDestinationName(topic),
			semconv.MessagingMessageID(msg.UUID),
			attribute.String(AttrCorrelationID, msg.Metadata.Get(middleware.CorrelationIDMetadataKey)),
		),
	)
}

// InjectTraceContext injects the trace context into the message metadata before publishing.
func (t *TempoTracer) InjectTraceContext(ctx context.Context, msg *message.Message) {
	otel.GetTextMapPropagator().Inject(ctx, messageMetadataCarrier{msg})
}

// StartPublishSpan creates a producer span for publishing a message.
// Call this before publishing to create a linked span, then call InjectTraceContext.
func (t *TempoTracer) StartPublishSpan(ctx context.Context, topic string, msg *message.Message) (context.Context, trace.Span) {
	spanName := topic + " publish"
	return t.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			semconv.MessagingSystemKey.String(MessagingSystemJetStream),
			semconv.MessagingOperationTypePublish,
			semconv.MessagingDestinationName(topic),
			semconv.MessagingMessageID(msg.UUID),
			semconv.MessagingMessageBodySize(len(msg.Payload)),
			attribute.String(AttrCorrelationID, msg.Metadata.Get("correlation_id")),
		),
	)
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
