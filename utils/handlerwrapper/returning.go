package handlerwrapper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	"github.com/Black-And-White-Club/frolf-bot-shared/utils"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Context keys for type-safe context value access.
type contextKey string

const (
	CtxKeyDiscordMessageID contextKey = "discord_message_id"
	CtxKeyChannelID        contextKey = "channel_id"
	CtxKeyMessageID        contextKey = "message_id"
	CtxKeyResponse         contextKey = "response"
	CtxKeySubmittedAt      contextKey = "submitted_at"
	CtxKeyReplyTo          contextKey = "reply_to"
)

// DiscordMetadataCarrier identifies payloads that carry a Discord message ID.
type DiscordMetadataCarrier interface {
	GetEventMessageID() string
}

// ReturningMetrics defines metrics for handler instrumentation.
type ReturningMetrics interface {
	RecordAttempt(ctx context.Context, handler string)
	RecordSuccess(ctx context.Context, handler string)
	RecordFailure(ctx context.Context, handler string)
	RecordDuration(ctx context.Context, handler string, d time.Duration)
}

// Result represents a domain event outcome with explicit routing.
type Result struct {
	Topic    string            // Required: target topic for this message
	Payload  any               // Required: event payload
	Metadata map[string]string // Optional: additional message metadata
}

// Validate ensures the result has required fields.
func (r Result) Validate() error {
	if r.Topic == "" {
		return errors.New("result topic is required")
	}
	if r.Payload == nil {
		return errors.New("result payload is required")
	}
	return nil
}

// WrapTransformingTyped wraps domain handlers that return []Result.
// Each Result is validated and transformed into a Watermill message with explicit topic routing.
func WrapTransformingTyped[T any](
	handlerName string,
	logger *slog.Logger,
	tracer trace.Tracer,
	helpers utils.Helpers,
	metrics ReturningMetrics,
	handler func(ctx context.Context, payload *T) ([]Result, error),
) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		// Correlation ID propagation
		corrID := middleware.MessageCorrelationID(msg)
		middleware.SetCorrelationID(corrID, msg)
		ctx := context.WithValue(msg.Context(), middleware.CorrelationIDMetadataKey, corrID)

		// Tracing
		ctx, span := tracer.Start(ctx, handlerName,
			trace.WithAttributes(
				attribute.String("message.id", msg.UUID),
				attribute.String("correlation_id", corrID),
			),
		)
		defer span.End()

		// Extract and propagate metadata from incoming message
		ctx = extractMetadataToContext(ctx, msg)

		// Metrics setup
		start := time.Now()
		if metrics != nil {
			metrics.RecordAttempt(ctx, handlerName)
			defer func() {
				metrics.RecordDuration(ctx, handlerName, time.Since(start))
			}()
		}

		logger.InfoContext(ctx, "handler started", attr.String("handler", handlerName))

		// Unmarshal payload
		payload := new(T)
		if err := helpers.UnmarshalPayload(msg, payload); err != nil {
			recordFailure(ctx, metrics, handlerName, logger, span, "unmarshal failed", err)
			return nil, err
		}

		// Execute handler
		results, err := handler(ctx, payload)
		if err != nil {
			recordFailure(ctx, metrics, handlerName, logger, span, "handler failed", err)
			return nil, err
		}

		// Transform results to messages
		outMessages, err := transformResults(ctx, msg, results, helpers, logger)
		if err != nil {
			recordFailure(ctx, metrics, handlerName, logger, span, "transform failed", err)
			return nil, err
		}

		// Success
		if metrics != nil {
			metrics.RecordSuccess(ctx, handlerName)
		}
		logger.InfoContext(ctx, "handler completed successfully",
			attr.String("handler", handlerName),
			attr.Int("results_count", len(results)),
		)

		return outMessages, nil
	}
}

// extractMetadataToContext extracts known metadata keys from message to context.
// Values are stored under both typed keys (for type-safe access) and string keys (for backward compatibility).
func extractMetadataToContext(ctx context.Context, msg *message.Message) context.Context {
	if v := msg.Metadata.Get("discord_message_id"); v != "" {
		ctx = context.WithValue(ctx, CtxKeyDiscordMessageID, v)
		ctx = context.WithValue(ctx, "discord_message_id", v) // backward compat
	}
	if v := msg.Metadata.Get("channel_id"); v != "" {
		ctx = context.WithValue(ctx, CtxKeyChannelID, v)
		ctx = context.WithValue(ctx, "channel_id", v) // backward compat
	}
	if v := msg.Metadata.Get("message_id"); v != "" {
		ctx = context.WithValue(ctx, CtxKeyMessageID, v)
		ctx = context.WithValue(ctx, "message_id", v) // backward compat
	}
	if v := msg.Metadata.Get("response"); v != "" {
		ctx = context.WithValue(ctx, CtxKeyResponse, v)
		ctx = context.WithValue(ctx, "response", v) // backward compat
	}
	if v := msg.Metadata.Get("submitted_at"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			ctx = context.WithValue(ctx, CtxKeySubmittedAt, t)
			ctx = context.WithValue(ctx, "submitted_at", t) // backward compat
		}
	}
	// Try standard Watermill/NATS reply keys
	if v := msg.Metadata.Get("reply_to"); v != "" {
		ctx = context.WithValue(ctx, CtxKeyReplyTo, v)
	} else if v := msg.Metadata.Get("reply"); v != "" {
		ctx = context.WithValue(ctx, CtxKeyReplyTo, v)
	}
	return ctx
}

// transformResults converts handler results to Watermill messages.
func transformResults(ctx context.Context, origMsg *message.Message, results []Result, helpers utils.Helpers, logger *slog.Logger) ([]*message.Message, error) {
	out := make([]*message.Message, 0, len(results))

	for i, res := range results {
		if err := res.Validate(); err != nil {
			return nil, fmt.Errorf("result[%d]: %w", i, err)
		}

		outMsg, err := helpers.CreateResultMessage(origMsg, res.Payload, res.Topic)
		if err != nil {
			return nil, fmt.Errorf("result[%d] create message: %w", i, err)
		}

		// Ensure topic is set for dynamic routing
		outMsg.Metadata.Set("topic", res.Topic)

		// Apply explicit metadata from result
		for k, v := range res.Metadata {
			outMsg.Metadata.Set(k, v)
		}

		// Propagate discord_message_id if not already set
		applyDiscordMetadata(ctx, outMsg, res.Payload)

		logger.DebugContext(ctx, "created result message metadata",
			attr.String("correlation_id", middleware.MessageCorrelationID(outMsg)),
			attr.String("topic", res.Topic),
		)

		out = append(out, outMsg)
	}

	return out, nil
}

// applyDiscordMetadata sets discord_message_id from context or payload carrier.
func applyDiscordMetadata(ctx context.Context, msg *message.Message, payload any) {
	// Skip if already set
	if msg.Metadata.Get("discord_message_id") != "" {
		return
	}

	// Try context first
	if v, ok := ctx.Value(CtxKeyDiscordMessageID).(string); ok && v != "" {
		msg.Metadata.Set("discord_message_id", v)
		return
	}

	// Try payload carrier interface
	if carrier, ok := payload.(DiscordMetadataCarrier); ok {
		if id := carrier.GetEventMessageID(); id != "" {
			msg.Metadata.Set("discord_message_id", id)
		}
	}
}

// recordFailure handles failure logging and metrics.
func recordFailure(ctx context.Context, metrics ReturningMetrics, handler string, logger *slog.Logger, span trace.Span, msg string, err error) {
	if metrics != nil {
		metrics.RecordFailure(ctx, handler)
	}
	logger.ErrorContext(ctx, msg, attr.Error(err))
	span.RecordError(err)
}
