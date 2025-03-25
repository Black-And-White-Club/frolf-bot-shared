package tempofrolfbot

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Tracer defines the interface for tracing.
type Tracer interface {
	InitTracing(ctx context.Context, opts TracingOptions) (shutdown func(), err error)
	TraceHandler(h message.HandlerFunc) message.HandlerFunc
	InjectTraceContext(ctx context.Context, msg *message.Message)
	SpanContextFromContext(ctx context.Context) trace.SpanContext
	StartSpan(ctx context.Context, name string, msg *message.Message) (context.Context, trace.Span)
}

var shutdownOnce sync.Once

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
	tracerProvider *sdktrace.TracerProvider
	tracer         trace.Tracer // Store the tracer for use in middleware
}

// NewTracer creates and configures a new Tempo tracer
func NewTracer(
	serviceName string,
	tempoEndpoint string,
	insecure bool,
	sampleRate float64,
	serviceVersion string,
	environment string,
) (Tracer, func(), error) {
	tracer := &TempoTracer{}

	opts := TracingOptions{
		ServiceName:    serviceName,
		TempoEndpoint:  tempoEndpoint,
		Insecure:       insecure,
		SampleRate:     sampleRate,
		ServiceVersion: serviceVersion,
		Environment:    environment,
	}

	// Initialize tracing with the provided options
	shutdownFunc, err := tracer.InitTracing(context.Background(), opts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize tracing: %w", err)
	}

	return tracer, shutdownFunc, nil
}

// InitTracing initializes OpenTelemetry tracing with Tempo.
func (t *TempoTracer) InitTracing(ctx context.Context, opts TracingOptions) (shutdownFunc func(), err error) {
	// Create a gRPC client connection to Tempo.
	conn, err := grpc.NewClient(
		opts.TempoEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Corrected placement
		// grpc.WithBlock(),                    // Consider blocking until connected (optional)
		// grpc.WithTimeout(5*time.Second),     //  Consider a connection timeout (optional)
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client connection to Tempo: %w", err)
	}

	// Create an OTLP trace exporter.
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Create a resource with the service name and version.
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(opts.ServiceName),
			semconv.ServiceVersionKey.String(opts.ServiceVersion),
			// Add other resource attributes if needed (environment, etc.)
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Adjust Sampling Rate Dynamically
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(opts.SampleRate))

	// Create a trace provider with the exporter and resource.
	t.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter), // Use batch processing for efficiency.
	)

	// Set the global trace provider.
	otel.SetTracerProvider(t.tracerProvider)
	t.tracer = t.tracerProvider.Tracer(opts.ServiceName) // Create the tracer *once*

	// Set the global propagator (W3C Trace Context).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Ensure cleanup happens only once
	shutdownFunc = func() {
		shutdownOnce.Do(func() {
			if err := t.tracerProvider.Shutdown(context.Background()); err != nil {
				// Log the error, but don't panic.  Shutdown failures are not always fatal.
				fmt.Printf("Error shutting down tracer provider: %v\n", err)
			}
		})
	}

	return shutdownFunc, nil
}

// WatermillTraceMiddleware creates a Watermill middleware for OpenTelemetry tracing.
// Now a method on TempoTracer
func (t *TempoTracer) TraceHandler(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		ctx := otel.GetTextMapPropagator().Extract(msg.Context(), messageMetadataCarrier{msg})

		// Use message topic as span name if available
		spanName := msg.Metadata.Get("topic")
		if spanName == "" {
			spanName = "watermill.handler"
		}

		// Extract JetStream metadata if available
		ackWait := msg.Metadata.Get("jetstream.ack_wait")
		deliveryAttempt := msg.Metadata.Get("jetstream.delivery_attempt")
		streamName := msg.Metadata.Get("jetstream.stream")
		consumerName := msg.Metadata.Get("jetstream.consumer")

		// Convert delivery attempt to int64
		var deliveryAttemptValue int64
		if deliveryAttempt != "" {
			if val, err := strconv.ParseInt(deliveryAttempt, 10, 64); err == nil {
				deliveryAttemptValue = val
			}
		}

		// Start a span with additional attributes
		ctx, span := t.tracer.Start(ctx, spanName, trace.WithAttributes(
			attribute.String("watermill.message_id", msg.UUID),
			attribute.String("watermill.correlation_id", msg.Metadata.Get("correlation_id")),
			attribute.String("watermill.topic", msg.Metadata.Get("topic")),
			attribute.String("watermill.middleware", "watermill.handler"),
			attribute.Bool("watermill.dead_letter", msg.Metadata.Get("dead_letter") == "true"),
			attribute.Int64("watermill.retry_count", deliveryAttemptValue),

			// JetStream attributes
			attribute.String("jetstream.stream", streamName),
			attribute.String("jetstream.consumer", consumerName),
			attribute.Int64("jetstream.delivery_attempt", deliveryAttemptValue),
			attribute.Int64("jetstream.ack_wait", parseInt64(ackWait)),
			attribute.Bool("jetstream.redelivered", msg.Metadata.Get("jetstream.redelivered") == "true"),
		))
		defer span.End()

		// Inject updated context into the message
		msg.SetContext(ctx)

		// Call the next handler
		msgs, err := h(msg)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		return msgs, err
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
