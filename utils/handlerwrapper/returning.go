package handlerwrapper

import (
	"context"
	"log/slog"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	"github.com/Black-And-White-Club/frolf-bot-shared/utils"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// DiscordMetadataCarrier identifies payloads that carry a Discord message ID
type DiscordMetadataCarrier interface {
	GetEventMessageID() string
}

// ReturningMetrics is minimal for module implementation
type ReturningMetrics interface {
	RecordAttempt(ctx context.Context, handler string)
	RecordSuccess(ctx context.Context, handler string)
	RecordFailure(ctx context.Context, handler string)
	RecordDuration(ctx context.Context, handler string, d time.Duration)
}

// Result represents a domain event outcome
type Result struct {
	Topic    string
	Payload  interface{}
	Metadata map[string]string
}

// WrapTransformingTyped wraps pure domain handlers that return []Result.
func WrapTransformingTyped[T any](
	handlerName string,
	logger *slog.Logger,
	tracer trace.Tracer,
	helpers utils.Helpers,
	metrics ReturningMetrics,
	handler func(ctx context.Context, payload *T) ([]Result, error),
) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		// 1. Correlation ID Propagation
		// We extract it from the message and manually inject it into the context.
		// This is the fix for "correlation_id=unknown" in logs.
		corrID := middleware.MessageCorrelationID(msg)

		// Ensure the message metadata itself is updated
		middleware.SetCorrelationID(corrID, msg)

		// Inject into Context using the key Watermill/Slog look for
		ctx := context.WithValue(msg.Context(), middleware.CorrelationIDMetadataKey, corrID)

		// 2. Tracing Setup
		ctx, span := tracer.Start(
			ctx,
			handlerName,
			trace.WithAttributes(
				attribute.String("message.id", msg.UUID),
				attribute.String("correlation_id", corrID),
			),
		)
		defer span.End()

		// 3. Metadata Extraction (The "Discord" Bridge)
		if dID := msg.Metadata.Get("discord_message_id"); dID != "" {
			ctx = context.WithValue(ctx, "discord_message_id", dID)
		}
		if channelID := msg.Metadata.Get("channel_id"); channelID != "" {
			ctx = context.WithValue(ctx, "channel_id", channelID)
		}
		if messageID := msg.Metadata.Get("message_id"); messageID != "" {
			ctx = context.WithValue(ctx, "message_id", messageID)
		}
		// Propagate response token (used by Discord handlers to carry RSVP choice)
		if resp := msg.Metadata.Get("response"); resp != "" {
			ctx = context.WithValue(ctx, "response", resp)
		}
		if subAt := msg.Metadata.Get("submitted_at"); subAt != "" {
			if t, err := time.Parse(time.RFC3339, subAt); err == nil {
				ctx = context.WithValue(ctx, "submitted_at", t)
			}
		}

		// 4. Metrics & Logs
		start := time.Now()
		if metrics != nil {
			metrics.RecordAttempt(ctx, handlerName)
			defer func() {
				metrics.RecordDuration(ctx, handlerName, time.Since(start))
			}()
		}

		// Logger now pulls correlationID from ctx automatically
		logger.InfoContext(ctx, "handler started",
			attr.String("handler", handlerName),
		)

		// 5. Unmarshal Payload
		payload := new(T)
		if err := helpers.UnmarshalPayload(msg, payload); err != nil {
			if metrics != nil {
				metrics.RecordFailure(ctx, handlerName)
			}
			logger.ErrorContext(ctx, "payload unmarshal failed", attr.Error(err))
			span.RecordError(err)
			return nil, err
		}

		// 6. Execute Pure Domain Logic
		results, err := handler(ctx, payload)
		if err != nil {
			if metrics != nil {
				metrics.RecordFailure(ctx, handlerName)
			}
			logger.ErrorContext(ctx, "handler failed", attr.Error(err))
			span.RecordError(err)
			return nil, err
		}

		// 7. Transform Results -> Outgoing Messages
		var outMessages []*message.Message
		for _, res := range results {
			outMsg, err := helpers.CreateResultMessage(msg, res.Payload, res.Topic)
			if err != nil {
				if metrics != nil {
					metrics.RecordFailure(ctx, handlerName)
				}
				logger.ErrorContext(ctx, "failed to create result message",
					attr.String("topic", res.Topic),
					attr.Error(err),
				)
				span.RecordError(err)
				return nil, err
			}

			// Ensure the topic is set in metadata so the Watermill router/NATS publisher
			// knows where to send this message when the handler output topic is dynamic.
			if res.Topic != "" {
				outMsg.Metadata.Set("topic", res.Topic)
			}

			// Apply Explicit Result Metadata
			if res.Metadata != nil {
				for k, v := range res.Metadata {
					outMsg.Metadata.Set(k, v)
				}
			}

			// Downstream Propagation (Context -> Outbound Metadata)
			if dID, ok := ctx.Value("discord_message_id").(string); ok && dID != "" {
				if outMsg.Metadata.Get("discord_message_id") == "" {
					outMsg.Metadata.Set("discord_message_id", dID)
				}
			}

			// Type-safe fallback
			if carrier, ok := res.Payload.(DiscordMetadataCarrier); ok {
				if id := carrier.GetEventMessageID(); id != "" {
					if outMsg.Metadata.Get("discord_message_id") == "" {
						outMsg.Metadata.Set("discord_message_id", id)
					}
				}
			}

			outMessages = append(outMessages, outMsg)
		}

		// 8. Success Finalization
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
