package attr

import (
	"context"
	"log/slog"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

// CorrelationIDFromMsg creates a slog.Attr for the correlation ID from message metadata.
func CorrelationIDFromMsg(msg *message.Message) slog.Attr {
	return slog.String("correlation_id", msg.Metadata.Get("correlation_id"))
}

// UserID creates a slog.Attr for a user ID.
func UserID(userID sharedtypes.DiscordID) slog.Attr {
	return slog.String("user_id", string(userID))
}

// MessageID creates a slog.Attr for a Watermill message ID.
func MessageID(msg *message.Message) slog.Attr {
	return slog.String("message_id", msg.UUID)
}

// Error creates a slog.Attr for an error.
func Error(err error) slog.Attr {
	return slog.Any("error", err)
}

// EventName creates a slog.Attr for an event name.
func EventName(event string) slog.Attr {
	return slog.String("event", event)
}

// Topic creates a slog.Attr for a topic.
func Topic(topic string) slog.Attr {
	return slog.String("topic", topic)
}

// String creates a generic string attribute.
func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

// Int creates a generic int attribute.
func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

// RoundID creates a slog.Attr for a RoundID.
func RoundID(key string, value sharedtypes.RoundID) slog.Attr {
	return slog.String(key, value.String())
}

// UUIDValue creates a slog.Attr for a UUID.
func UUIDValue(key string, value uuid.UUID) slog.Attr {
	return slog.String(key, value.String())
}

// StringUUID safely converts a string to a UUID for logging.
func StringUUID(key string, uuidStr string) slog.Attr {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return slog.String(key, "<invalid-uuid>")
	}
	return slog.String(key, id.String())
}

// Int64 creates a generic int64 attribute.
func Int64(key string, value int64) slog.Attr {
	return slog.Int64(key, value)
}

// Uint64 creates a generic uint64 attribute.
func Uint64(key string, value uint64) slog.Attr {
	return slog.Uint64(key, value)
}

// Float64 creates a generic float64 attribute.
func Float64(key string, value float64) slog.Attr {
	return slog.Float64(key, value)
}

// Bool creates a generic bool attribute.
func Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}

// Time creates a generic time attribute.
func Time(key string, value time.Time) slog.Attr {
	return slog.Time(key, value)
}

// Duration creates a generic duration attribute.
func Duration(key string, value time.Duration) slog.Attr {
	return slog.Duration(key, value)
}

// Any creates a generic interface attribute.
func Any(key string, value interface{}) slog.Attr {
	return slog.Any(key, value)
}

// Group creates a grouped attribute.
func Group(key string, attrs ...slog.Attr) slog.Attr {
	slogAttrs := make([]any, len(attrs))
	for i, attr := range attrs {
		slogAttrs[i] = attr
	}
	return slog.Group(key, slogAttrs...)
}

// ExtractCorrelationID pulls the correlation ID from context, or sets a default.
func ExtractCorrelationID(ctx context.Context) slog.Attr {
	if ctx == nil {
		return slog.String("correlation_id", "unknown") // Safe fallback
	}

	// Assume correlation_id is stored in context under this key
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		return slog.String("correlation_id", correlationID)
	}

	return slog.String("correlation_id", "unknown") // Default if missing
}

// TraceID extracts trace_id from context for Grafana Trace ↔ Logs correlation.
// Returns empty attribute if no valid span context exists.
func TraceID(ctx context.Context) slog.Attr {
	if ctx == nil {
		return slog.String("trace_id", "")
	}
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return slog.String("trace_id", "")
	}
	return slog.String("trace_id", spanCtx.TraceID().String())
}

// SpanID extracts span_id from context for Grafana Trace ↔ Logs correlation.
// Returns empty attribute if no valid span context exists.
func SpanID(ctx context.Context) slog.Attr {
	if ctx == nil {
		return slog.String("span_id", "")
	}
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return slog.String("span_id", "")
	}
	return slog.String("span_id", spanCtx.SpanID().String())
}

// TraceContext extracts both trace_id and span_id as a group for structured logging.
// Use this when you need both values together for Grafana correlation.
func TraceContext(ctx context.Context) []slog.Attr {
	if ctx == nil {
		return nil
	}
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return nil
	}
	return []slog.Attr{
		slog.String("trace_id", spanCtx.TraceID().String()),
		slog.String("span_id", spanCtx.SpanID().String()),
	}
}

// ConvertAttrsToAny converts a slice of slog.Attr to a slice of any.
func ConvertAttrsToAny(attrs []slog.Attr) []any {
	anyAttrs := make([]any, len(attrs))
	for i, attr := range attrs {
		anyAttrs[i] = attr
	}
	return anyAttrs
}
