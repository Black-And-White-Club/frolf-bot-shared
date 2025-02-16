package attr

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
)

// --- Attribute Constructors ---

// CorrelationIDFromMsg creates an slog.Attr for the correlation ID from message metadata.
func CorrelationIDFromMsg(msg *message.Message) slog.Attr {
	return slog.String("correlation_id", msg.Metadata.Get("correlation_id"))
}

// UserID creates an slog.Attr for a user ID.
func UserID(userID string) slog.Attr {
	return slog.String("user_id", userID)
}

// MessageID creates an slog.Attr for a Watermill message ID.
func MessageID(msg *message.Message) slog.Attr {
	return slog.String("message_id", msg.UUID)
}

// Error creates an slog.Attr for an error.
func Error(err error) slog.Attr {
	return slog.Any("error", err)
}

// EventName creates an slog.Attr for an event name.
func EventName(event string) slog.Attr {
	return slog.String("event", event)
}

// Topic creates an slog.Attr for a topic.
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

// Group creates a grouped attribute.  VERY useful.
func Group(key string, attrs ...slog.Attr) slog.Attr {
	convertedAttrs := make([]any, len(attrs))
	for i, attr := range attrs {
		convertedAttrs[i] = attr
	}
	return slog.Group(key, convertedAttrs...)
}
